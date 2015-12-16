package instance_test

import (
	"github.com/pivotal-pez/cfmgo"
	"gopkg.in/mgo.v2"
)

type fakeCol struct {
	cfmgo.Collection
	SpyID     interface{}
	SpyUpdate interface{}
}

func (s *fakeCol) UpsertID(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	s.SpyID = id
	s.SpyUpdate = update
	return
}
