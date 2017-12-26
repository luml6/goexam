package tests

import (
	. "notepad-api/model"
	"testing"

	"github.com/satori/go.uuid"
)

func pubulicattachInitMongo(t *testing.T) error {
	return MSessionInit("localhost:27017", "cloudnoteuser", "cloudnoteuser", 100)
}

func newInitPublicAttach() *PublicAttach {
	attPath := uuid.NewV4().String()
	att := NewPublicAttach(attPath)
	att.Status = STATUS_NORMAL
	att.Type = "test"
	att.Key = "test"

	return att
}

func TestPublicAttachAdd(t *testing.T) {
	pubulicattachInitMongo(t)
	att := newInitPublicAttach()

	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestPublicAttachGet(t *testing.T) {
	pubulicattachInitMongo(t)
	att := newInitPublicAttach()

	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att1 := NewPublicAttach(att.ID)
	if _, err := att1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if att1.ID != att.ID || att1.Key != att.Key {
		t.Fail()
	}
}

func TestPublicAttachUpdate(t *testing.T) {
	pubulicattachInitMongo(t)

	att := newInitPublicAttach()
	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att1 := NewPublicAttach(att.ID)
	if _, err := att1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att.Key = "test1"

	if err := att.Update([]string{"Key"}); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := att1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if att.Key != att1.Key {
		t.Fail()
	}

}

func TestPublicAttachRemove(t *testing.T) {
	pubulicattachInitMongo(t)

	att := newInitPublicAttach()
	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att1 := NewPublicAttach(att.ID)
	if _, err := att1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if err := att.Remove(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if attRet, err := att1.Get(); err == nil {
		if attRet != nil {
			t.Fail()
		}
	}
}
