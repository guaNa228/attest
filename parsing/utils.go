package parsing

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func addWeekToDate(dateStr string) string {
	// Parse the input date string into a time.Time value
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Handle parsing error
		return ""
	}

	// Add a week (7 days) to the date
	newDate := date.AddDate(0, 0, 7)

	// Format the new date as "yyyy-mm-dd" string
	newDateStr := newDate.Format("2006-01-02")

	return newDateStr
}

func Contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
}

func getNumberAfterSecondSlash(s string) (int32, error) {
	parts := strings.Split(s, "/")
	if len(parts) < 3 {
		return -1, fmt.Errorf("invalid string format")
	}

	number, err := strconv.Atoi(parts[2])
	if err != nil {
		return -1, err
	}

	return int32(number), nil
}

func PrintParsingResult(parsingResult []*FacultyParsed) {
	for _, faculty := range parsingResult {
		log.Println(faculty.Name)
		for _, course := range (*faculty).Courses {
			log.Println(" ", course.Number)
			for _, group := range (*course).Groups {
				log.Println("  ", group.FullCode)
				for _, class := range (*group).Classes {
					log.Println("   ", class.Name)
					for _, teacher := range (*class).Teachers {
						log.Println("    ", teacher.Name, "/", teacher.Id)
					}
				}
			}
		}
	}
}

func logParsingResult(parsingResult *[]*FacultyParsed, logChannel *chan string) {
	*logChannel <- fmt.Sprintf("Found %v faculties:", len(*parsingResult))
	for _, faculty := range *parsingResult {
		*logChannel <- fmt.Sprintf(" Found %v courses on faculty %s:", len(faculty.Courses), faculty.Name)
		for _, course := range (*faculty).Courses {
			*logChannel <- fmt.Sprintf("  On course %v found %v groups", course.Number, len(course.Groups))
			for _, group := range (*course).Groups {
				*logChannel <- fmt.Sprintf("   In group %s found %v classes:", group.FullCode, len(group.Classes))
				for _, class := range (*group).Classes {
					for _, teacher := range (*class).Teachers {
						*logChannel <- fmt.Sprintf("    Class %s is taught by %s(%v):", class.Name, teacher.Name, teacher.Id)
					}
				}
			}
		}
	}
}

func safeString(str string) string {
	str = reduceConsecutiveSpaces(str)
	return strings.Trim(str, " ")
}

func reduceConsecutiveSpaces(str string) string {
	var builder strings.Builder
	prevSpace := true

	for _, char := range str {
		if char == ' ' {
			if prevSpace {
				continue
			}
			prevSpace = true
		} else {
			prevSpace = false
		}

		builder.WriteRune(char)
	}

	return builder.String()
}
