package tests

import (
	. "notepad-api/model"
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

func userInitMongo(t *testing.T) error {
	return MSessionInit("localhost:27017", "cloudnoteuser", "cloudnoteuser", 100)
}

func TestUserGet(t *testing.T) {
	userInitMongo(t)

	uid := uuid.NewV4().String()
	user := NewUser(uid)

	if err := user.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	newUser := NewUser(uid)

	if _, err := newUser.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestUserAdd(t *testing.T) {
	userInitMongo(t)

	uid := uuid.NewV4().String()
	user := NewUser(uid)

	if err := user.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestUserUpdate(t *testing.T) {
	userInitMongo(t)

	uid := uuid.NewV4().String()
	user := NewUser(uid)

	if err := user.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	t1 := user.UpdateTime
	t2 := t1.Add(-time.Hour * 24)

	user.UpdateTime = t2
	if err := user.Update([]string{"UpdateTime"}); err != nil {
		t.Log(err)
		t.Fail()
	}

	newUser := NewUser(uid)
	if _, err := newUser.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if newUser.UpdateTime == t1 {
		t.Fail()
	} else if !newUser.UpdateTime.Equal(t2) {
		t.Log(t2, newUser.UpdateTime)
		t.Fail()
	}
}

func TestUserUpdateAll(t *testing.T) {
	UID := uuid.NewV4().String()
	user := NewUser(UID)

	if err := user.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	the_time, _ := time.Parse(FOTMAT_TIME_STRING, FOTMAT_TIME_STRING)
	user.SyncTime = the_time
	user.UpdateTime = the_time

	if err := user.UpdateAll(); err != nil {
		t.Log(err)
		t.Fail()
	}

	user1 := NewUser(UID)
	if _, err := user1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if !user1.SyncTime.Equal(the_time) {
		t.Log(user1)
		t.Fail()
	}
}
