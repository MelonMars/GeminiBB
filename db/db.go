package db

import (
	"sync/atomic"
)

type IDGenerator struct {
	counter uint64
}

func GenerateID(g *IDGenerator) uint64 {
	return atomic.AddUint64(&g.counter, 1)
}

type Comment struct {
	Content string
}

type Thread struct {
	ID       uint64
	Name     string
	Content  string
	Comments []Comment
}

type Section struct {
	ID       uint64
	Name     string
	Contents []Thread
}

type DB struct {
	Database []Section
}

type LoopError struct {
	message string
}

func (e LoopError) Error() string {
	return e.message
}

func Init() (DB, *IDGenerator) {
	return DB{[]Section{Section{123, "foo", []Thread{Thread{456, "thread1", "hello", []Comment{Comment{"FUCK"}}}}}}}, &IDGenerator{}
}

func NewComment(sectionID uint64, threadID uint64, db *DB, comment string) (error, DB) {
	var section *Section
	var thread *Thread
	for i := range db.Database {
		if db.Database[i].ID == sectionID {
			section = &db.Database[i]
			break
		}
	}
	if section == nil {
		return LoopError{"NO"}, *db
	}
	for i := range section.Contents {
		if section.Contents[i].ID == threadID {
			thread = &section.Contents[i]
			break
		}
	}
	if thread == nil {
		return LoopError{"NO"}, *db
	}
	thread.Comments = append(thread.Comments, Comment{comment})
	return nil, *db
}

func NewThread(sectionID uint64, db *DB, threadTitle string, threadContent string, Generator *IDGenerator) (error, DB) {
	var section *Section
	for i := range db.Database {
		if db.Database[i].ID == sectionID {
			section = &db.Database[i]
			break
		}
	}
	if section == nil {
		return LoopError{"NO"}, *db
	}
	section.Contents = append(section.Contents, Thread{GenerateID(Generator), threadTitle, threadContent, []Comment{}})
	return nil, *db
}
