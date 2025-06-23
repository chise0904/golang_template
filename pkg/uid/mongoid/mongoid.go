package mongoid

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type objectID struct {
}

func NewMongoIDGenerator() *objectID {
	return &objectID{}
}

func (u *objectID) GenUID() string {

	return primitive.NewObjectID().Hex()

}
