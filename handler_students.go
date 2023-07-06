package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/logger"
	"github.com/xuri/excelize/v2"
)

func (apiCfg *apiConfig) uploadStudentsUpload(w http.ResponseWriter, r *http.Request, user db.User) {
	logChan := make(chan string)
	errorChan := make(chan error)

	var errorCounter int

	go logger.Logger(logChan)
	go logger.ErrLogger(errorChan, &errorCounter)

	err := r.ParseMultipartForm(10 << 20) // 10MB limit for file size

	if err != nil {
		errorChan <- fmt.Errorf("file is too big")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		errorChan <- fmt.Errorf("error retrieving file")
		return
	}
	defer file.Close()

	if handler.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		errorChan <- fmt.Errorf("invalid file format: Only Excel files are allowed")
		return
	}

	// Create a temporary file to store the uploaded file
	tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		errorChan <- fmt.Errorf("eror creating temporary file")
		return
	}
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		errorChan <- fmt.Errorf("eror copying file")
		return
	}

	apiCfg.processExcelFile(tempFile, &logChan, &errorChan, &errorCounter)

	respondWithJSON(w, 201, struct{}{})
}

type parsedStudent struct {
	name  string
	email string
}

type userUpdatedGroup struct {
	ID       uuid.UUID
	Group_id uuid.UUID
}

type EmailGroup struct {
	group string
	email string
}

func (apiCfg *apiConfig) processExcelFile(file *os.File, logCh *chan string, errCh *chan error, errCount *int) {
	*logCh <- "Starting to process file"

	groupStudents := make(map[string][]*parsedStudent)

	f, err := excelize.OpenFile(file.Name())
	if err != nil {
		sendFailMessage(logCh, "Processing student's file")
		*errCh <- fmt.Errorf("failed to open file, file seems to be damaged")
		return
	}

	// Get all sheet names
	mainSheetList := f.GetSheetList()[0]

	rows, err := f.GetRows(mainSheetList)
	if err != nil {
		sendFailMessage(logCh, "Processing student's file")
		*errCh <- fmt.Errorf("failed to read rows, file seems to be damaged")
		return
	}

	var groupNumber string
	var students []*parsedStudent

	// Iterate over each row
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		if row[0] != "" {
			if groupNumber != "" {
				studentsToPaste := make([]*parsedStudent, len(students))
				copy(studentsToPaste, students)
				groupStudents[groupNumber] = students
			}

			groupNumber = row[0]

			students = []*parsedStudent{}

		} else {
			studentName := row[1]
			studentEmail := ""
			if len(row) == 3 {
				studentEmail = row[2]
			}

			student := &parsedStudent{name: studentName, email: studentEmail}
			students = append(students, student)
		}
	}

	groupStudents[groupNumber] = students

	uniqueEmails := []*EmailGroup{}
	var shouldAddStudent bool
	for group, students := range groupStudents {
		for _, student := range students {
			if student.email != "" {
				shouldAddStudent = true
				for _, emailGroupInstance := range uniqueEmails {
					if student.email == emailGroupInstance.email {
						*errCh <- fmt.Errorf("found duplicate email %s in groups %s and %s", student.email, group, emailGroupInstance.email)
						shouldAddStudent = false
						break
					}
				}

				if shouldAddStudent {
					newGroupEmail := EmailGroup{
						group: group,
						email: student.email,
					}
					uniqueEmails = append(uniqueEmails, &newGroupEmail)
				}
			}
		}
	}

	if haveErrorsHappend(errCount) {
		sendFailMessage(logCh, "Processing student's file")
		return
	}

	*logCh <- "File readed properly, starting to prepare data for db"

	err = apiCfg.DB.DeleteSemesterUsers(context.Background())
	if err != nil {
		*errCh <- fmt.Errorf("error deleting temporary semester users: %s", err)
		sendFailMessage(logCh, "Processing student's file")
		return
	}

	*logCh <- "Old users with no email deleted"

	usersToAdd := []*db.User{}
	usersToUpdate := []*userUpdatedGroup{}

	for group, students := range groupStudents {
		groupSplitted := strings.Split(group, "/")

		code, subcode := groupSplitted[0], groupSplitted[1]

		groupId, err := apiCfg.DB.GetGroupByFullCode(context.Background(), db.GetGroupByFullCodeParams{
			Code:    code,
			Subcode: subcode,
		})

		if err != nil {
			*errCh <- fmt.Errorf("group %s not found, wrong group number in file or maybe you forgot to parse the timetable, if not, you can add the group manually through the admin panel", group)
			sendFailMessage(logCh, "Processing student's file")
			return
		}

		for _, student := range students {
			if student.email != "" {
				studentId, err := apiCfg.DB.GetUserByEmail(context.Background(), sql.NullString{String: student.email, Valid: true})
				if err != nil {
					if err == sql.ErrNoRows {
						newUserInstance, err := apiCfg.createStudentInstance(student.name, student.email, groupId)
						if err != nil {
							*errCh <- fmt.Errorf("error generating student instance %s", student.name)
							sendFailMessage(logCh, "Processing student's file")
							return
						}
						usersToAdd = append(usersToAdd, &newUserInstance)
					} else {
						*errCh <- errors.New("db broken")
						sendFailMessage(logCh, "Processing student's file")
						return
					}
				} else {
					usersToUpdate = append(usersToUpdate, &userUpdatedGroup{
						Group_id: groupId,
						ID:       studentId,
					})
				}
			} else {
				newUserInstance, err := apiCfg.createStudentInstance(student.name, student.email, groupId)
				if err != nil {
					*errCh <- fmt.Errorf("error generating student instance %s", student.name)
					sendFailMessage(logCh, "Processing student's file")
					return
				}
				usersToAdd = append(usersToAdd, &newUserInstance)
			}
		}
	}

	*logCh <- "Data formed successfully"

	err = apiCfg.DB.RemoveGroupID(context.Background())
	if err != nil {
		*errCh <- fmt.Errorf("error removing group ids from users: %s", err)
		sendFailMessage(logCh, "Processing student's file")
		return
	}

	*logCh <- "Old groups are removed from students, starting db loading"

	actionsWg := sync.WaitGroup{}
	actionsWg.Add(1)
	go itemsBunkCreate[db.User](usersToAdd, "users", &actionsWg, errCh, errCount)

	actionsWg.Add(1)
	go itemsBunkUpdate[userUpdatedGroup](usersToUpdate, "users", "group_id", "id", &actionsWg, errCh, errCount)

	actionsWg.Wait()

	os.Remove(file.Name())

	if *errCount > 0 {
		sendFailMessage(logCh, "Processing student's file")
		return
	}

	*logCh <- "Succesfully filled groups with students, operation passed clear"
}

func (apiCfg *apiConfig) createStudentInstance(name string, email string, group uuid.UUID) (db.User, error) {
	emailToAdd := sql.NullString{}
	if email != "" {
		emailToAdd = sql.NullString{
			String: email,
			Valid:  true,
		}
	}

	uniqueCredentials, err := apiCfg.credentialsByName(name)
	if err != nil {
		return db.User{}, fmt.Errorf("error generating students credentials: %s", err)
	}

	return db.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		TeacherID: sql.NullInt32{},
		Login:     uniqueCredentials.login,
		Password:  uniqueCredentials.password,
		Role:      "student",
		GroupID:   uuid.NullUUID{UUID: group, Valid: true},
		Email:     emailToAdd,
	}, nil
}
