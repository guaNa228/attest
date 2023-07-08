package parsing

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/translit"
)

type ParsedTeachersEmails struct {
	Id    uuid.UUID
	Email string
}

var facultyPrefixes = []string{"icst", "physics", "hmath", "hum", "ic"}

func ParseTeachersMails(apiCfg *db.Queries, logCh *chan string, errCh *chan error) *[]*ParsedTeachersEmails {

	*logCh <- "Starting to parse teachers emails(only teachers with unique names and who does not have email yet)"

	parsingResult := []*ParsedTeachersEmails{}

	teachers, err := apiCfg.GetTeachersWithUniqueName(context.Background())
	if err != nil {
		*errCh <- fmt.Errorf("failed to get teachers to parse their emails: %s", err)
		return nil
	}

	facultyPrefixesLength := len(facultyPrefixes)
	for _, teacher := range teachers {
		for facultyIndex, facultyPrefix := range facultyPrefixes {
			parsedTeacherEmail, err := parseSingleTeacher(teacher.ID, teacher.Name, facultyPrefix)
			if err != nil {
				if facultyIndex == facultyPrefixesLength-1 {
					*logCh <- err.Error()
				}
			} else {
				parsingResult = append(parsingResult, &parsedTeacherEmail)
				*logCh <- fmt.Sprintf("%s email %s found on %s", teacher.Name, parsedTeacherEmail.Email, facultyPrefix)
				break
			}
		}
	}

	*logCh <- fmt.Sprintf("Found %v emails from %v, which were ordered for search", len(parsingResult), len(teachers))

	return &parsingResult
}

// func parseTeachersMailsChunk(teachersDataChunk *[]*db.GetTeachersWithUniqueNameRow, errCh *chan error, logCh *chan string, dataCh *chan *parsedTeachersEmails, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	facultyPrefixesLength := len(facultyPrefixes)
// 	for teacherIndexInChunk, teacher := range *teachersDataChunk {
// 		fmt.Println(teacherIndexInChunk, teacher.Name)
// 		for facultyIndex, facultyPrefix := range facultyPrefixes {
// 			parsedTeacherEmail, err := parseSingleTeacher(teacher.ID, teacher.Name, facultyPrefix)
// 			if err != nil {
// 				if facultyIndex == facultyPrefixesLength-1 {
// 					*logCh <- err.Error()
// 				}
// 			} else {
// 				*dataCh <- &parsedTeacherEmail
// 				*logCh <- fmt.Sprintf("%s email fount on the %v try on %s", teacher.Name, facultyIndex, facultyPrefix)
// 				break
// 			}
// 		}
// 	}
// }

func parseSingleTeacher(id uuid.UUID, name string, facultyPrefix string) (ParsedTeachersEmails, error) {
	url := fmt.Sprintf("https://%s.spbstu.ru/person/%s", facultyPrefix, translit.ToLatin(strings.Replace(strings.ToLower(name), " ", "_", -1), translit.RussianEmail))

	resp, err := http.Head(url)
	if err != nil {
		return ParsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
	}
	defer resp.Body.Close()

	responseCode := resp.StatusCode

	fmt.Println(fmt.Sprintf("https://%s.spbstu.ru/person/%s", facultyPrefix, translit.ToLatin(strings.Replace(strings.ToLower(name), " ", "_", -1), translit.RussianEmail)), responseCode)

	if responseCode == http.StatusNotFound {
		return ParsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
	} else {
		doc, err := waitForPageLoad(url)
		if err != nil {
			return ParsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
		}

		email := doc.Find("li.mail").Children().First().Children().First().Text()

		if !IsValidEmail(email) {
			return ParsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
		}

		return ParsedTeachersEmails{
			Id:    id,
			Email: email,
		}, nil
	}
}
