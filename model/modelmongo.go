package model

import (
	"gopkg.in/mgo.v2"
)

var (
	mgoSession *mgo.Session
)

func MSessionInit(url string, user string, pwd string, maxPoolSize int) (err error) {
	if mgoSession != nil {
		return nil
	}

	mgoSession, err = mgo.Dial(url)
	if err != nil {
		panic(err) //直接终止程序运行
	}

	if maxPoolSize > 0 {
		mgoSession.SetPoolLimit(maxPoolSize)
	}

	//	mgoSession.Login(&mgo.Credential{Username: user, Password: pwd})

	return nil
}

func MSessionGet() *mgo.Session {
	if mgoSession == nil {
		return nil
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

func MSessionClose(session *mgo.Session) {
	session.Close()
}
