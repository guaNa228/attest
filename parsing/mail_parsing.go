package parsing

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/translit"
)

type parsedTeachersEmails struct {
	id    uuid.UUID
	email string
}

var facultyPrefixes = []string{"icst", "physics", "hmath", "hum", "ic"}

func ParseTeachersMails(apiCfg *db.Queries, logCh *chan string, errCh *chan error) *[]*parsedTeachersEmails {

	*logCh <- "Starting to parse teachers emails(only teachers with unique names and who does not have email yet)"

	parsingResult := []*parsedTeachersEmails{}

	teachers, err := apiCfg.GetTeachersWithUniqueName(context.Background())
	if err != nil {
		*errCh <- fmt.Errorf("failed to get teachers to parse their emails: %s", err)
		return nil
	}

	teachersChunks := ChunkItems(teachers, 10)

	teachersDataChannel := make(chan *parsedTeachersEmails)
	go readTeachersMailsData(&parsingResult, &teachersDataChannel)

	teachersWg := sync.WaitGroup{}
	for _, chunk := range teachersChunks {
		teachersWg.Add(1)
		go parseTeachersMailsChunk(&chunk, errCh, logCh, &teachersDataChannel, &teachersWg)
	}

	teachersWg.Wait()

	*logCh <- fmt.Sprintf("Found %v emails from %v, which were ordered for search", len(parsingResult), len(teachers))

	return &parsingResult
}

func parseTeachersMailsChunk(teachersDataChunk *[]*db.GetTeachersWithUniqueNameRow, errCh *chan error, logCh *chan string, dataCh *chan *parsedTeachersEmails, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, teacher := range *teachersDataChunk {
		for facultyIndex, facultyPrefix := range facultyPrefixes {
			fmt.Println(facultyPrefix, teacher.Name)
			parsedTeacherEmail, err := parseSingleTeacher(teacher.ID, teacher.Name, facultyPrefix)
			if err != nil {
				if facultyIndex == len(facultyPrefixes)-1 {
					*logCh <- err.Error()
				}
			} else {
				*dataCh <- &parsedTeacherEmail
				break
			}
		}
	}
}

func parseSingleTeacher(id uuid.UUID, name string, facultyPrefix string) (parsedTeachersEmails, error) {
	url := fmt.Sprintf("https://%s.spbstu.ru/person/%s", facultyPrefix, translit.ToLatin(strings.Replace(strings.ToLower(name), " ", "_", -1), translit.RussianEmail))

	resp, err := http.Head(url)
	if err != nil {
		return parsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
	}
	defer resp.Body.Close()

	responseCode := resp.StatusCode

	resp.Body.Close()

	if responseCode == http.StatusNotFound {
		return parsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
	} else {
		doc, err := waitForPageLoad(url)
		if err != nil {
			return parsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
		}

		email := doc.Find("li.mail").Children().First().Children().First().Text()

		if !isValidEmail(email) {
			return parsedTeachersEmails{}, fmt.Errorf("%s email parsing failed", name)
		}

		fmt.Println(facultyPrefix)

		return parsedTeachersEmails{
			id:    id,
			email: email,
		}, nil
	}
}
