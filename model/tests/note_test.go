package tests

import (
	. "notepad-api/model"
	"testing"

	"github.com/satori/go.uuid"
)

func noteInitMongo(t *testing.T) error {
	return MSessionInit("localhost:27017", "cloudnoteuser", "cloudnoteuser", 100)
}

func newInitNote() *Note {
	notePath := uuid.NewV4().String()
	UID := uuid.NewV4().String()
	catePath := uuid.NewV4().String()

	note := NewNote(notePath)
	note.CategoryName = "资料"
	note.CategoryPath = catePath
	note.UID = UID
	note.Author = "gary.geng"
	note.Content = "content test"
	note.Summary = "summary"
	note.Title = "title"

	return note
}

func TestNoteAdd(t *testing.T) {
	noteInitMongo(t)

	note := newInitNote()

	if err := note.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestNoteGet(t *testing.T) {
	noteInitMongo(t)

	note := newInitNote()

	if err := note.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	note1 := NewNote(note.Path)

	if _, err := note1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestNoteGetByPathAndStatus(t *testing.T) {
	noteInitMongo(t)

	note := newInitNote()

	if err := note.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	note1 := NewNote(note.Path)
	if _, err := note1.GetByPathAndStatus(STATUS_NORMAL); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := note.Delete(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := note1.GetByPathAndStatus(STATUS_DELETED); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestNoteSetStatus(t *testing.T) {
	noteInitMongo(t)

	note := newInitNote()

	if err := note.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	note1 := NewNote(note.Path)

	if err := note.SetStatus(STATUS_DELETED); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := note1.GetByPathAndStatus(STATUS_DELETED); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := note.SetStatus(STATUS_NORMAL); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := note1.GetByPathAndStatus(STATUS_NORMAL); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestNoteDelete(t *testing.T) {
	noteInitMongo(t)

	note := newInitNote()

	if err := note.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	note1 := NewNote(note.Path)

	if err := note.Delete(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if noteRet, err := note1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		if noteRet != nil {
			t.Fail()
		}
	}

	if _, err := note1.GetByPathAndStatus(STATUS_DELETED); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestNoteRecover(t *testing.T) {
	noteInitMongo(t)

	note := newInitNote()

	if err := note.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	note1 := NewNote(note.Path)

	if err := note.Delete(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if noteRet, err := note1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	} else {
		if noteRet != nil {
			t.Fail()
		}
	}

	if err := note.Recover(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := note1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestNoteMove(t *testing.T) {
	noteInitMongo(t)

	note := newInitNote()

	if err := note.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	newCatePath := uuid.NewV4().String()
	newCateName := "电影"

	if err := note.Move(newCatePath, newCateName); err != nil {
		t.Log(err)
		t.Fail()
	}

	note1 := NewNote(note.Path)

	if _, err := note1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if note1.CategoryPath != newCatePath || note1.CategoryName != newCateName {
		t.Fail()
	}
}
