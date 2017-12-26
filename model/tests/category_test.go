package tests

import (
	. "notepad-api/model"
	"testing"

	"github.com/satori/go.uuid"
)

func cateInitMongo(t *testing.T) error {
	return MSessionInit("localhost:27017", "cloudnoteuser", "cloudnoteuser", 100)
}

func TestCateAdd(t *testing.T) {
	cateInitMongo(t)

	UID := uuid.NewV4().String()
	path := uuid.NewV4().String()

	cate := NewCategory(UID, path)
	cate.Name = "资料"

	if err := cate.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestCateGet(t *testing.T) {
	cateInitMongo(t)

	UID := uuid.NewV4().String()
	path := uuid.NewV4().String()

	cate := NewCategory(UID, path)
	cate.Name = "资料"

	if err := cate.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	cate2 := NewCategory(UID, path)

	if _, err := cate2.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if cate.Name != cate2.Name {
		t.Fail()
	}
}

func TestCateUpdate(t *testing.T) {
	cateInitMongo(t)

	UID := uuid.NewV4().String()
	path := uuid.NewV4().String()

	cate := NewCategory(UID, path)
	cate.Name = "资料"

	if err := cate.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	cate2 := NewCategory(UID, path)

	if _, err := cate2.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if cate.Name != cate2.Name {
		t.Fail()
	}

	cate.Name = "笔记"

	if err := cate.Update([]string{"Name"}); err != nil {
		t.Log(err)
		t.Fail()
	}

	cate3 := NewCategory(UID, path)

	if _, err := cate3.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if cate3.Name != cate.Name {
		t.Log("更新错误")
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	cateInitMongo(t)

	UID := uuid.NewV4().String()
	path := uuid.NewV4().String()

	cate := NewCategory(UID, path)
	cate.Name = "资料"

	if err := cate.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := cate.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if cate.Status != STATUS_NORMAL {
		t.Fail()
	}

	if err := cate.Delete(); err != nil {
		t.Log(err)
		t.Fail()
	}

	cate2 := NewCategory(UID, path)

	if _, err := cate2.GetByPathAndStatus(path, STATUS_DELETED); err != nil {
		t.Log(err)
		t.Fail()
	}

	if cate2.Status != STATUS_DELETED {
		t.Fail()
	}
}

func TestRecover(t *testing.T) {
	cateInitMongo(t)

	UID := uuid.NewV4().String()
	path := uuid.NewV4().String()

	cate := NewCategory(UID, path)
	cate.Name = "资料"

	if err := cate.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := cate.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if cate.Status != STATUS_NORMAL {
		t.Fail()
	}

	if err := cate.Delete(); err != nil {
		t.Log(err)
		t.Fail()
	}

	cate2 := NewCategory(UID, path)

	if _, err := cate2.GetByPathAndStatus(path, STATUS_DELETED); err != nil {
		t.Log(err)
		t.Fail()
	}

	if cate2.Status != STATUS_DELETED {
		t.Fail()
	}

	if err := cate2.Recover(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := cate2.Get(); err != nil {
		t.Log()
		t.Fail()
	}
}

func TestGetByNameAndUID(t *testing.T) {
	cateInitMongo(t)

	UID := uuid.NewV4().String()
	path := uuid.NewV4().String()

	cate := NewCategory(UID, path)
	cate.Name = "资料"

	if err := cate.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := cate.GetByNameAndUID(cate.Name, UID); err != nil {
		t.Log(err)
		t.Fail()
	}

}

func GetByPathAndStatus(t *testing.T) {
	cateInitMongo(t)

	UID := uuid.NewV4().String()
	path := uuid.NewV4().String()

	cate := NewCategory(UID, path)
	cate.Name = "资料"

	if err := cate.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := cate.GetByPathAndStatus(path, STATUS_NORMAL); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := cate.GetByPathAndStatus(path, STATUS_DELETED); err == nil {
		t.Fail()
	}
}
