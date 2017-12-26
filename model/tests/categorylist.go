package tests

import (
	. "notepad-api/model"
	"testing"

	"github.com/satori/go.uuid"
)

func addCategoryToList(uid string, State uint8) error {
	category := NewCategory(uid, uuid.NewV4().String())
	category.Status = State
	return category.Add()
}

func TestCategorylistGet(t *testing.T) {
	UID := uuid.NewV4().String()

	var Num int = 3

	for i := 0; i < Num; i++ {
		addCategoryToList(UID, STATUS_NORMAL)
	}

	nl := NewCategoryList(UID)

	cateList, err := nl.GetList()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if len(cateList) != Num {
		t.Fail()
	}
}

func TestCategorylistTrash(t *testing.T) {
	UID := uuid.NewV4().String()

	nl := NewCategoryList(UID)
	if err := nl.TrashAll(); err != nil {
		t.Log(err)
		t.Fail()
	}
}
