package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/parsing"
)

func (apiCfg *apiConfig) parsingResultToDBInstances(parsingResult *[]*parsing.FacultyParsed) (map[string]interface{}, error) {
	uniqueStreams := []*db.Stream{}
	uniqueGroups := []*db.Group{}
	uniqueTeachers := []*db.User{}
	uniqueClasses := []*db.Class{}
	workloads := []*db.Workload{}
	//Для обработки других программ обучения(магистратура и т.д.) требуется изменить парсинг
	bachelorsProgam, err := apiCfg.DB.GetProgramsIdByName(context.Background(), "Бакалавриат")
	if err != nil {
		return map[string]interface{}{}, errors.New("bachelor program is not supported")
	}

	for _, faculty := range *parsingResult {
		for _, course := range (*faculty).Courses {
			for _, group := range (*course).Groups {
				currentStreamId := addOrGetStream(&uniqueStreams, strings.Split(group.FullCode, "/")[0], bachelorsProgam)
				currentGroupId := addGroup(&uniqueGroups, strings.Split(group.FullCode, "/")[1], currentStreamId, course.Number)
				for _, class := range (*group).Classes {
					currentClassId := addOrGetClass(&uniqueClasses, class.Name)
					for _, teacher := range (*class).Teachers {
						currentTeacherId := apiCfg.addOrGetTeacher(&uniqueTeachers, teacher.Name, teacher.Id)
						newWorkload := createWorkload(currentGroupId, currentClassId, currentTeacherId)
						workloads = append(workloads, &newWorkload)
					}
				}
			}
		}
	}

	makeUserLoginsUnique(&uniqueTeachers)

	return map[string]interface{}{
		"users":     &uniqueTeachers,
		"classes":   &uniqueClasses,
		"streams":   &uniqueStreams,
		"groups":    &uniqueGroups,
		"workloads": &workloads,
	}, nil
}

func (apiCfg *apiConfig) createTeacherInstance(name string, teacherId int32) (db.User, error) {
	teacherIdToPaste := sql.NullInt32{}
	if teacherId != -1 {
		teacherIdToPaste = sql.NullInt32{
			Valid: true,
			Int32: teacherId,
		}
	}

	uniqueCredentials, err := apiCfg.credentialsByName(name)
	if err != nil {
		return db.User{}, fmt.Errorf("error generating user: %s", err)
	}

	return db.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		TeacherID: teacherIdToPaste,
		Login:     uniqueCredentials.login,
		Password:  uniqueCredentials.password,
		Role:      "teacher",
		GroupID:   uuid.NullUUID{},
	}, nil
}

var StreamByCode = map[string]string{
	"3530903": "Прикладная информатика",
	"3530203": "Математическое обеспечение и администрирование информационных систем",
	"3530201": "Математика и компьютерные науки",
	"3530202": "Фундаментальная информатика и информационные технологии",
	"3530901": "Информатика и вычислительная техника",
	"3530902": "Информационные системы и технологии",
	"3530904": "Программная инженерия",
	"3532701": "Стандартизация и метрология",
	"3532702": "Управление качеством",
	"3532703": "Системный анализ и управление",
	"3532704": "Управление в технических системах",
	"3532705": "Инноватика",
}

func createStream(code string, programId uuid.UUID) db.Stream {

	return db.Stream{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      StreamByCode[code],
		Code:      code,
		Program:   programId,
	}
}

func createClass(name string) db.Class {
	return db.Class{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
}

func createGroup(subcode string, course int16, streamId uuid.UUID) db.Group {
	return db.Group{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Subcode:   subcode,
		Course:    course,
		Stream:    streamId,
	}
}

func createWorkload(groupId uuid.UUID, classId uuid.UUID, teacherId uuid.UUID) db.Workload {
	return db.Workload{
		ID:      uuid.New(),
		GroupID: groupId,
		Class:   classId,
		Teacher: teacherId,
	}
}

//utils

func makeUserLoginsUnique(uniqueUsers *[]*db.User) {
	sameLoginsCount := make(map[string]*[]*db.User)

	for _, user := range *uniqueUsers {
		usersWithDuplicatedLogin, ok := sameLoginsCount[user.Login]
		if ok {
			*usersWithDuplicatedLogin = append(*usersWithDuplicatedLogin, user)
		} else {
			users := []*db.User{user}
			sameLoginsCount[user.Login] = &users
		}
	}

	for login, usersWithTheLogin := range sameLoginsCount {
		if len(*usersWithTheLogin) > 1 {
			for index, user := range *usersWithTheLogin {
				if index > 0 {
					user.Login = fmt.Sprintf("%s%v", login, index+1)
				}
			}
		}
	}

	for login, user := range sameLoginsCount {
		fmt.Println(login, len(*user))
	}
}

func addOrGetStream(streams *[]*db.Stream, code string, programId uuid.UUID) uuid.UUID {
	for _, stream := range *streams {
		if stream.Code == code && stream.Program == programId {
			return stream.ID
		}
	}

	newStream := createStream(code, programId)
	*streams = append(*streams, &newStream)
	return newStream.ID
}

func addGroup(groups *[]*db.Group, subcode string, streamId uuid.UUID, course int16) uuid.UUID {
	newGroup := createGroup(subcode, course, streamId)
	*groups = append(*groups, &newGroup)
	return newGroup.ID
}

func addOrGetClass(classes *[]*db.Class, name string) uuid.UUID {
	for _, class := range *classes {
		if class.Name == name {
			return class.ID
		}
	}

	newClass := createClass(name)
	*classes = append(*classes, &newClass)
	return newClass.ID
}

func (apiCfg *apiConfig) addOrGetTeacher(teachers *[]*db.User, name string, teacherId int32) uuid.UUID {
	teacherIdToPaste := sql.NullInt32{
		Valid: false,
	}
	if teacherId != -1 {
		teacherIdToPaste = sql.NullInt32{
			Valid: true,
			Int32: teacherId,
		}
	}

	foundTeacherId, err := apiCfg.DB.GetTeacherIDByNameAndTeacherId(context.Background(), db.GetTeacherIDByNameAndTeacherIdParams{
		Name:      name,
		TeacherID: teacherIdToPaste,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			for _, teacher := range *teachers {
				if teacher.Name == name && teacher.TeacherID == teacherIdToPaste {
					return teacher.ID
				}
			}

			newTeacher, err := apiCfg.createTeacherInstance(name, teacherId)
			if err != nil {
				log.Fatal("broken db")
				return uuid.UUID{}
			}
			*teachers = append(*teachers, &newTeacher)
			return newTeacher.ID
		} else {
			log.Fatal("broken db")
			return uuid.UUID{}
		}
	}

	return foundTeacherId
}
