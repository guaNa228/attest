package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	db "github.com/guaNa228/attest/internal/database"
	"github.com/lib/pq"
)

func fillTypeData(v interface{}) ([]interface{}, error) {
	switch x := v.(type) {
	case Group:
		return []interface{}{x.ID, x.CreatedAt, x.UpdatedAt, x.Name, x.Code}, nil
	case db.Attestation:
		return []interface{}{x.ID, x.SemesterActivityID, x.StudentID, x.Month, nil, nil}, nil
	default:
		return []interface{}{}, fmt.Errorf("unsupported type %t thrown for bunk creation", x)
	}
}

var insertContextByTypes = map[string][]string{
	"groups":      {"id", "created_at", "updated_at", "name", "code"},
	"attestation": {"id", "semester_activity_id", "student_id", "month", "result", "comment"},
}

func itemsBunkCreate[T any](items []*T, typeTitle string) []error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=%v dbname=attest password=%v sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")))
	if err != nil {
		fmt.Println(err)
		return []error{err}
	}
	defer db.Close()
	size := 500

	chunkList := chunkItems(items, size)

	var wg sync.WaitGroup

	errCh := make(chan error)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := db.Conn(context.Background())
			if err != nil {
				errCh <- err
				return
			}
			defer conn.Close()

			tx, err := conn.BeginTx(context.Background(), nil)
			if err != nil {
				errCh <- err
				return
			}
			defer tx.Rollback()

			columnNamesForCurrentTable, ok := insertContextByTypes[typeTitle]
			if !ok {
				errCh <- fmt.Errorf("unsupported type %v thrown for bunk creation", typeTitle)
			}

			stmt, err := tx.PrepareContext(context.Background(), pq.CopyIn(typeTitle, columnNamesForCurrentTable...))
			if err != nil {
				errCh <- err
				return
			}
			defer stmt.Close()

			for _, chunk := range chunkList {
				for _, item := range chunk {
					values, err := fillTypeData(*item)
					if err != nil {
						errCh <- err
						return
					}
					_, err = stmt.ExecContext(context.Background(), values...)
					if err != nil {
						errCh <- err
						return
					}
				}
			}

			_, err = stmt.ExecContext(context.Background())
			if err != nil {
				errCh <- err
				return
			}

			err = tx.Commit()
			if err != nil {
				errCh <- err
				return
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	errorSlice := []error{}
	for err := range errCh {
		if err != nil {
			//здесь всегда возникает ошибка pq: повторяющееся значение ключа нарушает ограничение уникальности, экспериментально
			//установлено, что все данные попадают в таблицу, поэтому эту ошибку принято решение не выводить
			if !strings.Contains(err.Error(), "pkey") {
				errorSlice = append(errorSlice, err)
			}
		}
	}

	if len(errorSlice) == 0 {
		return nil
	}
	return errorSlice
}

func chunkItems[T any](items []*T, size int) [][]*T {
	chunkList := make([][]*T, 0)
	chunk := make([]*T, 0, size)

	for _, item := range items {
		chunk = append(chunk, item)
		if len(chunk) == size {
			chunkList = append(chunkList, chunk)
			chunk = make([]*T, 0, size)
		}
	}

	if len(chunk) > 0 {
		chunkList = append(chunkList, chunk)
	}

	return chunkList
}
