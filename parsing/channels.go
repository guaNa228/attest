package parsing

import "sync"

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

func readChannelData[T any](sliceToAdd *[]*T, channel *chan *T) {
	for dataPiece := range *channel {
		*sliceToAdd = append(*sliceToAdd, dataPiece)
	}
}

func ReadChannelDataWithWG[T any](sliceToAdd *[]*T, channel *chan *T, readingWg *sync.WaitGroup) {
	for dataPiece := range *channel {
		go func(pieceOfData *T) {
			defer readingWg.Done()
			*sliceToAdd = append(*sliceToAdd, pieceOfData)
		}(dataPiece)
	}
}
