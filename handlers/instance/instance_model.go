package instance

import (
	"github.com/pivotal-pez/cfmgo"
	"gopkg.in/mgo.v2/bson"
)

//Save - saves the model to the given collection
func (s InstanceModel) Save(col cfmgo.Collection) {
	id := bson.NewObjectId().Hex()
	col.UpsertID(id, s)
}
