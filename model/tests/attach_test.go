package tests

import (
	. "notepad-api/model"
	"testing"

	"github.com/satori/go.uuid"
)

func attachInitMongo(t *testing.T) error {
	return MSessionInit("localhost:27017", "cloudnoteuser", "cloudnoteuser", 100)
}

func newInitAttach() *Attachment {
	attPath := uuid.NewV4().String()
	notePath := uuid.NewV4().String()
	att := NewAttachment(attPath)
	att.OK = "www.baidu.com/sdfas/qfqae4ef.jpg"
	att.OT = ATT_STORE_TYPE_ZEUSIS
	att.NotePath = notePath
	att.Type = "jpg"

	return att
}

func TestAttachAdd(t *testing.T) {
	attachInitMongo(t)
	att := newInitAttach()

	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestAttachGet(t *testing.T) {
	attachInitMongo(t)
	att := newInitAttach()

	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att1 := NewAttachment(att.ID)
	if _, err := att1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if att1.ID != att.ID || att1.OK != att.OK {
		t.Fail()
	}
}

func TestAttachUpdate(t *testing.T) {
	attachInitMongo(t)

	att := newInitAttach()
	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att1 := NewAttachment(att.ID)
	if _, err := att1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att.OK = "www.google.com/sdfas/qfqasf2fasdae4ef.mp3"
	att.OT = ATT_STORE_TYPE_ZEUSIS
	att.Type = "mp3"

	if err := att.Update([]string{"OK"}); err != nil {
		t.Log(err)
		t.Fail()
	}

	if _, err := att1.Get(); err != nil {
		t.Log(err)
		t.Fail()
	}

	if att.OK != att1.OK {
		t.Fail()
	}

}

func TestAttachRemove(t *testing.T) {
	attachInitMongo(t)

	att := newInitAttach()
	if err := att.Add(); err != nil {
		t.Log(err)
		t.Fail()
	}

	att1 := NewAttachment(att.ID)
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
