package parsing

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// const value of classes, that we don't wan't to add
var bannedClasses []string = []string{"Элективная физическая культура и спорт", "Иностранный язык: Русский язык как иностранный", "Модуль саморазвития (SoftSkills)"}
var classTypesToParse []string = []string{"Практика"}

type FacultyParsed struct {
	Name    string
	Courses []*CourseParsed
}

type CourseParsed struct {
	Number int16
	Groups []*GroupParsed
}

type GroupParsed struct {
	FullCode string
	Classes  []*Class
}

type Class struct {
	Name     string
	Teachers []*ParsedTeacher
}

type ParsedTeacher struct {
	Name string
	Id   int32
}

func StartParsing(logChan *chan string, errorChan *chan error, dateToParse string) *[]*FacultyParsed {

	var facultiesWg sync.WaitGroup
	facultyChannel := make(chan *FacultyParsed)
	facultiesToReturn := []*FacultyParsed{}
	go parsedFacultiesChannelRead(&facultiesToReturn, facultyChannel)

	facultiesWg.Add(1)
	go ParseFaculty("https://ruz.spbstu.ru/faculty/95/groups", logChan, errorChan, &facultiesWg, facultyChannel, dateToParse)

	facultiesWg.Wait()

	facultyChannelWG := sync.WaitGroup{}
	facultyChannelWG.Add(1)
	go func() {
		defer facultyChannelWG.Done()
		close(facultyChannel)
	}()

	facultyChannelWG.Wait()

	time.Sleep(time.Second * 2)
	logParsingResult(&facultiesToReturn, logChan)
	*logChan <- fmt.Sprintln("Parsing finished")

	return &facultiesToReturn
}

func ParseFaculty(url string, parsingLogsChannel *chan string, parsingErrorsChannel *chan error, facultiesWg *sync.WaitGroup, dataChannel chan *FacultyParsed, dateToParse string) {
	// Make an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		*parsingErrorsChannel <- err
		return
	}
	defer response.Body.Close()
	defer facultiesWg.Done()
	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		*parsingErrorsChannel <- err
		return
	}

	appElement := doc.Find(".app").First()

	_, exists := appElement.Attr("data-reactid")
	if !exists {
		*parsingErrorsChannel <- errors.New("unable to start parsing, framework has changed")
		return
	}

	facultyItem := doc.Find(".breadcrumb-item.active").Text()
	*parsingLogsChannel <- fmt.Sprintf("Starting to parse faculty %s", facultyItem)

	facultyToReturn := FacultyParsed{
		Name: facultyItem,
	}
	coursesWg := sync.WaitGroup{}
	coursesChannel := make(chan *CourseParsed)

	go parsedCoursesChannelRead(&facultyToReturn, coursesChannel)

	doc.Find(".faculty__level").Each(func(index int, element *goquery.Selection) {
		course, err := strconv.Atoi(element.Find(".page__h3 span").First().Text())
		if err != nil {
			*parsingErrorsChannel <- err
			return
		}
		coursesWg.Add(1)
		go ParseCourse(int16(course), element.Find(".groups-list"), parsingLogsChannel, parsingErrorsChannel, &coursesWg, coursesChannel, dateToParse)
	})

	coursesWg.Wait()

	parsedCoursesChannelWG := sync.WaitGroup{}
	parsedCoursesChannelWG.Add(1)
	go func() {
		defer parsedCoursesChannelWG.Done()
		close(coursesChannel)
	}()

	parsedCoursesChannelWG.Wait()

	dataChannel <- &facultyToReturn
}

func ParseCourse(courseNumber int16, element *goquery.Selection, parsingLogsChannel *chan string, parsingErrorsChannel *chan error, coursesWg *sync.WaitGroup, dataChannel chan *CourseParsed, dateToParse string) {

	*parsingLogsChannel <- fmt.Sprintf("Starting to parse course %v", courseNumber)

	courseToReturn := CourseParsed{
		Number: courseNumber,
	}

	defer coursesWg.Done()

	groupsWg := sync.WaitGroup{}
	groupsChannel := make(chan *GroupParsed)

	go parsedGroupsChannelRead(&courseToReturn, groupsChannel)
	element.Find(".groups-list__item").Each(func(index int, element *goquery.Selection) {
		linkToGroup, exists := element.Children().First().Attr("href")
		if !exists {
			*parsingErrorsChannel <- errors.New("didn't find group link, maybe template changed")
			return
		}
		groupsWg.Add(1)
		go ParseGroup(linkToGroup, parsingLogsChannel, parsingErrorsChannel, &groupsWg, groupsChannel, dateToParse)
	})

	groupsWg.Wait()

	parsedGroupsChannelWG := sync.WaitGroup{}
	parsedGroupsChannelWG.Add(1)
	go func() {
		defer parsedGroupsChannelWG.Done()
		close(groupsChannel)
	}()

	parsedGroupsChannelWG.Wait()

	dataChannel <- &courseToReturn
	*parsingLogsChannel <- fmt.Sprintf("finished parsing course %v", courseNumber)
}

func ParseGroup(link string, parsingLogsChannel *chan string, parsingErrorsChannel *chan error, groupsWg *sync.WaitGroup, dataChannel chan *GroupParsed, dateToParse string) {
	response, err := http.Get(fmt.Sprintf("https://ruz.spbstu.ru%s?date=%s", link, dateToParse))

	if err != nil {
		*parsingErrorsChannel <- err
		return
	}
	defer response.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		*parsingErrorsChannel <- err
		return
	}

	defer groupsWg.Done()
	groupCode := safeString(doc.Find(".breadcrumb-item.active").Children().Eq(1).Text())

	if groupCode == "" {
		*parsingErrorsChannel <- errors.New("didn't find group code, maybe template changed")
		return
	}

	*parsingLogsChannel <- fmt.Sprintf("Starting to parse group %v", groupCode)

	groupToReturn := GroupParsed{
		FullCode: groupCode,
	}

	resultingClassses := []*Class{}

	for i := 0; i < 4; i++ {
		parseGroupWeek(link, dateToParse, &resultingClassses, parsingErrorsChannel)
		dateToParse = addWeekToDate(dateToParse)
	}

	groupToReturn.Classes = resultingClassses

	dataChannel <- &groupToReturn
	*parsingLogsChannel <- fmt.Sprintf("finished parsing group %v", groupCode)
}

func parseGroupWeek(urlBase string, dateToParse string, classses *[]*Class, parsingErrorsChannel *chan error) {
	response, err := http.Get(fmt.Sprintf("https://ruz.spbstu.ru%s?date=%s", urlBase, dateToParse))
	if err != nil {
		*parsingErrorsChannel <- err
		return
	}
	defer response.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		*parsingErrorsChannel <- err
		return
	}

	doc.Find(".lesson").Each(func(index int, element *goquery.Selection) {
		lessonParams := element.Find(".lesson__params")
		classType := safeString(lessonParams.Find(".lesson__type").Text())
		if !Contains(classTypesToParse, classType) {
			return
		}

		className := safeString(element.Find(".lesson__subject").Children().Eq(2).Text())
		if Contains(bannedClasses, className) {
			return
		}
		if className == "" {
			*parsingErrorsChannel <- errors.New("broken timetable structure, class name not found")
			return
		}

		teacherParamsLink := lessonParams.Find(".lesson__teachers").Find(".lesson__link")
		teacherName := safeString(teacherParamsLink.Children().Eq(2).Text())
		if teacherName == "" {
			return
		}
		teacherIdLink, exists := teacherParamsLink.Attr("href")

		var teacherId int32 = -1
		if exists {
			teacherId, err = getNumberAfterSecondSlash(teacherIdLink)
		}
		for _, class := range *classses {
			if class.Name == className {
				for _, teacher := range class.Teachers {
					if teacher.Id == teacherId && teacher.Name == teacherName {
						return
					}
				}
			}
		}

		for _, class := range *classses {
			if class.Name == className {
				class.Teachers = append(class.Teachers, &ParsedTeacher{
					Id:   teacherId,
					Name: teacherName,
				})
				return
			}
		}

		*classses = append(*classses, &Class{
			Name: className,
			Teachers: []*ParsedTeacher{
				{
					Name: teacherName,
					Id:   teacherId,
				},
			},
		})

	})
}
