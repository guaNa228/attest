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

type ParsedTeachersEmails struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
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

	teachersChunks := ChunkItems(teachers, int(len(teachers)/11))

	teachersDataChannel := make(chan *ParsedTeachersEmails)
	go readChannelData(&parsingResult, &teachersDataChannel)

	duplicatedTeachers := []*uuid.UUID{}
	duplicateTeachersChannel := make(chan *uuid.UUID)
	go readChannelData(&duplicatedTeachers, &duplicateTeachersChannel)

	teachersWg := sync.WaitGroup{}
	for _, chunk := range teachersChunks {
		teachersWg.Add(1)
		go parseTeachersMailsChunk(chunk, errCh, logCh, &teachersDataChannel, &duplicateTeachersChannel, &teachersWg)
	}

	teachersWg.Wait()

	dataWG := sync.WaitGroup{}
	dataWG.Add(1)
	go func() {
		defer dataWG.Done()
		close(teachersDataChannel)
	}()

	duplicatesWG := sync.WaitGroup{}
	duplicatesWG.Add(1)
	go func() {
		defer duplicatesWG.Done()
		close(duplicateTeachersChannel)
	}()

	dataWG.Wait()
	duplicatesWG.Wait()

	clearResult := []*ParsedTeachersEmails{}

	var isDuplicate bool

	for _, teacherEmail := range parsingResult {
		isDuplicate = false
		for _, duplicateTeacherId := range duplicatedTeachers {
			if teacherEmail.Id == *duplicateTeacherId {
				break
			}
		}
		if !isDuplicate {
			clearResult = append(clearResult, teacherEmail)
		}
	}

	*logCh <- fmt.Sprintf("Found %v emails from %v, which were ordered for search", len(clearResult), len(teachers))

	return &clearResult
}

func parseTeachersMailsChunk(teachersDataChunk []*db.GetTeachersWithUniqueNameRow, errCh *chan error, logCh *chan string, dataCh *chan *ParsedTeachersEmails, maliciousChan *chan *uuid.UUID, wg *sync.WaitGroup) {
	defer wg.Done()

	currentTeacherErrorCounter := 0
	facultyPrefixesLength := len(facultyPrefixes)
	for _, teacher := range teachersDataChunk {
		currentTeacherErrorCounter = 0
		for _, facultyPrefix := range facultyPrefixes {
			err := parseSingleTeacher(teacher.ID, teacher.Name, facultyPrefix, dataCh)
			if err != nil {
				currentTeacherErrorCounter++
			}
		}
		if currentTeacherErrorCounter < facultyPrefixesLength-1 {
			*maliciousChan <- &teacher.ID
			*logCh <- fmt.Sprintf("can't determine %s email, found more than 1 faculty", teacher.Name)
		} else {
			if currentTeacherErrorCounter == facultyPrefixesLength-1 {
				*logCh <- fmt.Sprintf("Found email of %s", teacher.Name)
			} else {
				*logCh <- fmt.Sprintf("%s email is not found", teacher.Name)
			}
		}
	}
}

func parseSingleTeacher(id uuid.UUID, name string, facultyPrefix string, dataChan *chan *ParsedTeachersEmails) error {
	url := fmt.Sprintf("https://%s.spbstu.ru/person/%s", facultyPrefix, translit.ToLatin(strings.Replace(strings.ToLower(name), " ", "_", -1), translit.RussianEmail))

	resp, err := http.Head(url)
	if err != nil {
		return fmt.Errorf("%s can't make a request to page", name)
	}
	defer resp.Body.Close()

	responseCode := resp.StatusCode

	if responseCode == http.StatusNotFound {
		return fmt.Errorf("%s page is 404", name)
	} else {
		email, err := waitForPageLoad(url)
		if err != nil {
			return fmt.Errorf("%s page didn't load", name)
		}

		if !IsValidEmail(email) {
			return fmt.Errorf("%s not valid email", name)
		}

		*dataChan <- &ParsedTeachersEmails{
			Id:    id,
			Email: email,
		}
		return nil
	}
}
