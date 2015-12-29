package instance

import (
	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/cfmgo/params"
	"github.com/xchapter7x/lo"
	"gopkg.in/mgo.v2/bson"
)

//GetTaskID - get a task id from a instanceid on a given collection
var GetTaskID = func(instanceID string, collection cfmgo.Collection) (taskID string, err error) {
	var instance InstanceModel

	if instance, err = getInstance(instanceID, collection); err == nil {
		taskID = instance.TaskGUID
	}
	return
}

//GetRequestID - get a request id from a instanceid on a given collection
var GetRequestID = func(instanceID string, collection cfmgo.Collection) (requestID string, err error) {
	var instance InstanceModel

	if instance, err = getInstance(instanceID, collection); err == nil {
		requestID = instance.RequestID
	}
	return
}

func getInstance(instanceID string, collection cfmgo.Collection) (instance InstanceModel, err error) {
	instance = InstanceModel{}
	var firstResultIndex = 0
	if instanceID == "" {
		err = ErrInvalidInstanceID

	} else {
		query := new(params.RequestParams)
		query.Q = bson.M{
			CollectionInstanceIDQueryField: instanceID,
		}
		var result = make([]InstanceModel, 1)

		if _, err = collection.Find(query, &result); err == nil {

			if len(result) > firstResultIndex {
				instance = result[firstResultIndex]

			} else {
				lo.G.Error("no records found in:", result)
				err = ErrNoRecordsInResult
			}
		}
	}
	return
}
