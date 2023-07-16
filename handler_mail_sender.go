package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/logger"
	"github.com/guaNa228/attest/parsing"
	"gopkg.in/gomail.v2"
)

func (apiCfg *apiConfig) handleEmailSending(w http.ResponseWriter, r *http.Request, user db.User) {
	logChan := make(chan string)
	errorChan := make(chan error)
	usersToUpdateChan := make(chan *uuid.UUID)
	usersWithSuccessfullySentEmails := []*uuid.UUID{}
	var errorCounter int

	readingWg := sync.WaitGroup{}

	go logger.Logger(logChan, GlobalWsConn, true)
	go logger.ErrLogger(errorChan, &errorCounter, GlobalWsConn, true)
	go parsing.ReadChannelDataWithWG(&usersWithSuccessfullySentEmails, &usersToUpdateChan, &readingWg)

	dialer := gomail.Dialer{
		Host:     "smtp.yandex.ru",
		Port:     465,
		Username: "g-u-a-N-a@yandex.ru",
		Password: os.Getenv("MAIL_PASSWORD"),
		SSL:      true,
	}
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	logChan <- "Starting to sending emails to users, which have email and haven't recieved creadentials yet"

	dbGetUsersToSend, err := apiCfg.DB.GetUsersWithEmails(context.Background())
	if err != nil {
		errorChan <- fmt.Errorf("failed to get users to send emails, probably all users with emails already recieved their mails: %s", err)
		return
	}

	sendingWg := sync.WaitGroup{}

	for _, user := range dbGetUsersToSend {
		if user.Email.String == "frolov.go@edu.spbstu.ru" || user.Email.String == "turalchuk_ka@spbstu.ru" {
			sendingWg.Add(1)
			go sendEmail(user, &logChan, &usersToUpdateChan, &sendingWg, &dialer, &readingWg)
		}
	}

	sendingWg.Wait()
	readingWg.Wait()

	dataChannelWg := sync.WaitGroup{}
	dataChannelWg.Add(1)
	go func() {
		defer dataChannelWg.Done()
		close(usersToUpdateChan)
	}()

	logChan <- fmt.Sprintf("Sent email to %v users of %v ordered for sending", len(usersWithSuccessfullySentEmails), len(dbGetUsersToSend))

	logChan <- "Starting to update emails statuses"

	updateWg := sync.WaitGroup{}
	updateWg.Add(1)
	itemsBunkUpdate(usersWithSuccessfullySentEmails, "users", "email_sent", "id", &updateWg, &errorChan, &errorCounter)

	updateWg.Wait()

	logChan <- "Succesfully updated users statuses"

	logWG := sync.WaitGroup{}
	logWG.Add(1)
	go func() {
		defer logWG.Done()
		close(logChan)
	}()

	errorWG := sync.WaitGroup{}
	errorWG.Add(1)
	go func() {
		defer errorWG.Done()
		close(errorChan)
	}()

	logWG.Wait()
	errorWG.Wait()

	if errorCounter > 0 {
		respondWithError(w, 400, "Something went wrong, see the error log")
	} else {
		respondWithJSON(w, 201, struct{}{})
	}

	GlobalWsWg.Done()
}

func sendEmail(recipient db.GetUsersWithEmailsRow, logChan *chan string, dataCh *chan *uuid.UUID, outerWg *sync.WaitGroup, dialer *gomail.Dialer, readingWg *sync.WaitGroup) {
	defer outerWg.Done()
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "g-u-a-N-a@yandex.ru")
	mailer.SetHeader("To", recipient.Email.String)
	mailer.SetHeader("Subject", "Ваши учетные данные на портале ежемесячной аттестации ИКНТ")
	mailer.SetBody("text/plain", fmt.Sprintf("Логин %s\nПароль: %s", recipient.Login, recipient.Password))

	if err := dialer.DialAndSend(mailer); err != nil {
		*logChan <- fmt.Sprintf("Failed to send email to %s: %v", recipient.Email.String, err)
	} else {
		*logChan <- fmt.Sprintf("Successfully sent email to %s", recipient.Email.String)
		uuidToSend := recipient.ID
		readingWg.Add(1)
		*dataCh <- &uuidToSend
	}
}
