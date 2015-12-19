package instance_test

import (
	"github.com/pivotal-pez/cfmgo"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
	"gopkg.in/mgo.v2"
)

type fakeCol struct {
	cfmgo.Collection
	SpyID      interface{}
	SpyUpdate  interface{}
	FakeResult []InstanceModel
}

func (s *fakeCol) Wake() {

}

func (s *fakeCol) UpsertID(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	s.SpyID = id
	s.SpyUpdate = update
	return
}

func (s *fakeCol) Find(params cfmgo.Params, result interface{}) (count int, err error) {
	for i := range s.FakeResult {
		(*(result.(*[]InstanceModel)))[i] = s.FakeResult[i]
	}
	return
}
