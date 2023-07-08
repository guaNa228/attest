package parsing

func parsedFacultiesChannelRead(fs *[]*FacultyParsed, facultiesChannel chan *FacultyParsed) {
	for faculty := range facultiesChannel {
		*fs = append(*fs, faculty)
	}
}

func parsedCoursesChannelRead(f *FacultyParsed, coursesChannel chan *CourseParsed) {
	for course := range coursesChannel {
		f.Courses = append(f.Courses, course)
	}
}

func parsedGroupsChannelRead(c *CourseParsed, groupsChannel chan *GroupParsed) {
	for group := range groupsChannel {
		c.Groups = append(c.Groups, group)
	}
}

// func readTeachersMailsData(parsedTeachersEmailsSlice *[]*ParsedTeachersEmails, channel *chan *ParsedTeachersEmails) {
// 	for teacherEmailRead := range *channel {
// 		*parsedTeachersEmailsSlice = append(*parsedTeachersEmailsSlice, teacherEmailRead)
// 	}
// }
