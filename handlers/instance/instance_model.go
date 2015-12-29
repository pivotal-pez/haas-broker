package instance

import (
	"github.com/pivotal-pez/cfmgo"
	"github.com/xchapter7x/lo"
	"gopkg.in/mgo.v2/bson"
)

//Save - saves the model to the given collection
func (s InstanceModel) Save(col cfmgo.Collection) {
	id := bson.NewObjectId().Hex()
	col.UpsertID(id, s)
}

//UpdateField - updates the existing record
func (s InstanceModel) UpdateField(col cfmgo.Collection, fieldname string, fieldvalue interface{}) {
	changeInfo, err := col.FindAndModify(
		bson.M{
			"instanceid": s.InstanceID,
		},
		bson.M{
			fieldname: fieldvalue,
		},
		nil,
	)

	if err != nil {
		lo.G.Error("InstanceModel update error: ", err)
	}
	lo.G.Debug("changinfo: ", changeInfo)
}
