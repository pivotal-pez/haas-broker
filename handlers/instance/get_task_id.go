package instance

import (
	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/cfmgo/params"
	"gopkg.in/mgo.v2/bson"
)

//GetTaskID - get a task id from a instanceid on a given collection
var GetTaskID = func(instanceID string, collection cfmgo.Collection) (taskID string, err error) {

	if instanceID == "" {
		err = ErrInvalidInstanceID

	} else {
		query := new(params.RequestParams)
		query.Q = bson.M{
			CollectionInstanceIDQueryField: instanceID,
		}
		var result = make([]InstanceModel, 1)

		if _, err = collection.Find(query, &result); err == nil {
			taskID = result[0].TaskGUID
		}
	}
	return
}
