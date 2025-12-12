package gox

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	KB   = 1 << 10
	MB   = 1 << 20
	GB   = 1 << 30
	TB   = 1 << 40
	Mbps = 10e6 / 8
	Gbps = 10e9 / 8
)

type ObjectID = primitive.ObjectID

func NewOID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func NewOIDHex() string {
	return NewOID().Hex()
}

func AsOID(hex string) (ObjectID, error) {
	return primitive.ObjectIDFromHex(hex)
}
