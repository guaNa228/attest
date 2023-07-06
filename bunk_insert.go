package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"

	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/parsing"
)

var insertContextByTypes = map[string][]string{
	"groups":      {"id", "created_at", "updated_at", "subcode", "stream", "course"},
	"users":       {"id", "created_at", "updated_at", "name", "login", "password", "role", "teacher_id", "group_id", "email"},
	"streams":     {"id", "created_at", "updated_at", "name", "code", "program"},
	"classes":     {"id", "created_at", "updated_at", "name"},
	"workloads":   {"id", "group_id", "class", "teacher"},
	"attestation": {"id", "semester_activity_id", "student_id", "month", "result", "comment"},
}

func fillTypeData(v interface{}) ([]interface{}, error) {
	switch x := v.(type) {
	case db.Group:
		return []interface{}{x.ID, x.CreatedAt, x.UpdatedAt, x.Subcode, x.Stream, x.Course}, nil
	case db.User:
		return []interface{}{x.ID, x.CreatedAt, x.UpdatedAt, x.Name, x.Login, x.Password, x.Role, x.TeacherID, x.GroupID, x.Email}, nil
	case db.Stream:
		return []interface{}{x.ID, x.CreatedAt, x.UpdatedAt, x.Name, x.Code, x.Program}, nil
	case db.Class:
		return []interface{}{x.ID, x.CreatedAt, x.UpdatedAt, x.Name}, nil
	case db.Workload:
		return []interface{}{x.ID, x.GroupID, x.Class, x.Teacher}, nil
	case userUpdatedGroup:
		return []interface{}{x.ID, x.Group_id}, nil
	case db.Attestation:
		return []interface{}{x.ID, x.SemesterActivityID, x.StudentID, x.Month, nil, nil}, nil
	default:
		return []interface{}{}, fmt.Errorf("unsupported type %t thrown for bunk creation", x)
	}
}

func haveErrorsHappend(errCounter *int) bool {
	return *errCounter > 0
}

func sendFailMessage(logChannel *chan string, process string) {
	*logChannel <- fmt.Sprintf("%s failed due to errors, none of data added, for details see the error log", process)
}

func (apiCfg *apiConfig) parsedBunkInsert(codedToDatabaseFormatParsedData map[string]interface{}, logChan *chan string, errChan *chan error, errCounter *int) {

	*logChan <- "Starting to fill db with parsed data"

	err := apiCfg.DB.ClearStreamsTable(context.Background())

	if err != nil {
		*errChan <- fmt.Errorf("error clearing Streams table: %s", err.Error())
		sendFailMessage(logChan, "Filling db with data")
		return
	}

	*logChan <- "Successfully cleared streams table, which leads to clearing groups and workloads tables"

	err = apiCfg.DB.ClearClassesTable(context.Background())

	if err != nil {
		*errChan <- fmt.Errorf("error clearing Classes table: %s", err.Error())
		sendFailMessage(logChan, "Filling db with data")
		return
	}

	*logChan <- "Successfully cleared classes table"

	streamsClassesUsersWg := sync.WaitGroup{}

	usersData := *codedToDatabaseFormatParsedData["users"].(*[]*db.User)
	streamsClassesUsersWg.Add(1)
	go itemsBunkCreate[db.User](usersData, "users", &streamsClassesUsersWg, errChan, errCounter)

	classesData := *codedToDatabaseFormatParsedData["classes"].(*[]*db.Class)
	streamsClassesUsersWg.Add(1)
	go itemsBunkCreate[db.Class](classesData, "classes", &streamsClassesUsersWg, errChan, errCounter)

	streamsData := *codedToDatabaseFormatParsedData["streams"].(*[]*db.Stream)
	streamsClassesUsersWg.Add(1)
	go itemsBunkCreate[db.Stream](streamsData, "streams", &streamsClassesUsersWg, errChan, errCounter)

	streamsClassesUsersWg.Wait()

	if haveErrorsHappend(errCounter) {
		sendFailMessage(logChan, "Filling db with data")
		return
	}

	*logChan <- "Successfully added teachers, classes and streams"

	groupsWg := sync.WaitGroup{}
	groupsData := *codedToDatabaseFormatParsedData["groups"].(*[]*db.Group)
	groupsWg.Add(1)
	go itemsBunkCreate[db.Group](groupsData, "groups", &groupsWg, errChan, errCounter)

	groupsWg.Wait()

	if haveErrorsHappend(errCounter) {
		sendFailMessage(logChan, "Filling db with data")
		return
	}

	*logChan <- "Successfully added groups"

	workloadsWg := sync.WaitGroup{}
	workloadsData := *codedToDatabaseFormatParsedData["workloads"].(*[]*db.Workload)
	workloadsWg.Add(1)
	go itemsBunkCreate[db.Workload](workloadsData, "workloads", &workloadsWg, errChan, errCounter)

	workloadsWg.Wait()

	if haveErrorsHappend(errCounter) {
		sendFailMessage(logChan, "Filling db with data")
		return
	}

	*logChan <- "Successfully added workloads"
	*logChan <- "Succsessfully filled db with parsed data"
}

func itemsBunkCreate[T any](items []*T, typeTitle string, outerWg *sync.WaitGroup, errCh *chan error, errCount *int) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=%v dbname=attest password=%v sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")))
	if err != nil {
		*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
		return
	}

	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
		return
	}
	defer conn.Close()

	columnNamesForCurrentTable, ok := insertContextByTypes[typeTitle]
	if !ok {
		err := fmt.Errorf("unsupported type %v thrown for bunk creation", typeTitle)
		*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
		return
	}

	size := 500
	chunkList := parsing.ChunkItems(items, size)

	tx, err := db.Begin()
	if err != nil {
		*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
		return
	}

	defer outerWg.Done()

	var chunkWg sync.WaitGroup
	for _, chunk := range chunkList {
		chunkWg.Add(1)
		go func(chunk []*T) {

			valueStrings := []string{}
			valueArgs := []interface{}{}
			defer chunkWg.Done()
			for itemIndex, item := range chunk {
				values, err := fillTypeData(*item)
				if err != nil {
					*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
					tx.Rollback()
					return
				}

				tempValues := []string{}

				for index, value := range values {
					tempValues = append(tempValues, fmt.Sprintf("$%v", itemIndex*len(columnNamesForCurrentTable)+(index+1)))
					valueArgs = append(valueArgs, value)
				}

				valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(tempValues, ", ")))
			}

			stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s(%s) VALUES %s", typeTitle, strings.Join(columnNamesForCurrentTable, ", "), strings.Join(valueStrings, ", ")))
			if err != nil {
				*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
				tx.Rollback()
				return
			}

			_, err = stmt.ExecContext(context.Background(), valueArgs...)
			if err != nil {
				*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
				tx.Rollback()
				return
			}
		}(chunk)
	}

	chunkWg.Wait()

	if *errCount == 0 {
		err = tx.Commit()
		if err != nil {
			*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
		}
	}
}

func itemsBunkUpdate[T any](items []*T, typeTitle string, fieldToUpdate string, fieldToIdentificate string, outerWg *sync.WaitGroup, errCh *chan error, errCount *int) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=%v dbname=attest password=%v sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")))
	if err != nil {
		*errCh <- fmt.Errorf("error while updating %s table: %s", typeTitle, err.Error())
		return
	}

	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		*errCh <- fmt.Errorf("error while updating %s table: %s", typeTitle, err.Error())
		return
	}
	defer conn.Close()

	size := 500
	chunkList := parsing.ChunkItems(items, size)

	tx, err := db.Begin()
	if err != nil {
		*errCh <- fmt.Errorf("error while updating %s table: %s", typeTitle, err.Error())
		return
	}

	defer outerWg.Done()

	var chunkWg sync.WaitGroup
	for _, chunk := range chunkList {
		chunkWg.Add(1)
		go func(chunk []*T) {

			caseStrings := []string{}
			valueArgs := []interface{}{}
			defer chunkWg.Done()
			for itemIndex, item := range chunk {

				values, err := fillTypeData(*item)
				if err != nil {
					*errCh <- fmt.Errorf("error while updating %s table: %s", typeTitle, err.Error())
					tx.Rollback()
					return
				}

				caseStrings = append(caseStrings, fmt.Sprintf("WHEN %s=$%v THEN $%v", fieldToIdentificate, itemIndex*2+1, itemIndex*2+2))
				valueArgs = append(valueArgs, values...)
			}

			//log.Println(fmt.Sprintf("UPDATE %s SET %s = (CASE %s ELSE %s)", typeTitle, fieldToUpdate, strings.Join(caseStrings, "\n"), fieldToUpdate))

			stmt, err := tx.Prepare(fmt.Sprintf("UPDATE %s SET %s = (CASE %s ELSE %s END)", typeTitle, fieldToUpdate, strings.Join(caseStrings, "\n"), fieldToUpdate))
			if err != nil {
				*errCh <- fmt.Errorf("error while updating table %s table: %s", typeTitle, err.Error())
				tx.Rollback()
				return
			}

			_, err = stmt.ExecContext(context.Background(), valueArgs...)
			if err != nil {
				*errCh <- fmt.Errorf("error while updating %s table: %s", typeTitle, err.Error())
				tx.Rollback()
				return
			}
		}(chunk)
	}

	chunkWg.Wait()

	if *errCount == 0 {
		err = tx.Commit()
		if err != nil {
			*errCh <- fmt.Errorf("error while filling %s table: %s", typeTitle, err.Error())
		}
	}
}
