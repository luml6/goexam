package tests

import (
	. "notepad-api/model"
	"testing"

	"github.com/satori/go.uuid"
)

func noteInitMongo(t *testing.T) error {
	return MSessionInit("10.0.12.104:27017", "cloudnoteuser", "cloudnoteuser", 100)
}

func addNoteToList(uid, catePath string) error {
	note := NewNote(uuid.NewV4().String())
	note.UID = uid
	note.CategoryPath = catePath

	return note.Add()
}

func TestNotelistGet(t *testing.T) {
	noteInitMongo(t)
	UID := uuid.NewV4().String()
	catePath := uuid.NewV4().String()

	var Num int = 3

	for i := 0; i < Num; i++ {
		addNoteToList(UID, catePath)
	}

	nl := NewNoteList(catePath, UID)

	noteList, err := nl.GetList()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if len(noteList) != Num {
		t.Fail()
	}
}

func TestNotelistDelete(t *testing.T) {
	UID := uuid.NewV4().String()
	catePath := uuid.NewV4().String()

	nl := NewNoteList(catePath, UID)
	if err := nl.Delete(); err != nil {
		t.Log(err)
		t.Fail()
	}
}
