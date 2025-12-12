package jsonx

import (
	"github.com/fengzhi09/golibx/gox"

	"github.com/globalsign/mgo/bson"
)

var (
	Marshal             = gox.Marshal
	Unmarshal           = gox.Unmarshal
	UnsafeUnmarshal     = gox.UnsafeUnmarshal
	UnsafeMarshal       = gox.UnsafeMarshal
	UnsafeMarshalString = gox.UnsafeMarshalString
	MustMarshal         = gox.MustMarshal
	MustUnmarshal       = gox.MustUnmarshal
	IfElse              = gox.IfElse
	UnixMilli           = gox.UnixMilli
	ISODateTimeMs       = gox.ISODateTimeMs
	ISODateTime         = gox.ISODateTime
	AsLong              = gox.AsLong
	AsBool              = gox.AsBool
	AsDouble            = gox.AsDouble
	AsInt               = gox.AsInt
	AsStr               = gox.AsStr
	AsTime              = gox.AsTime
	AsULong             = gox.AsULong
	AsStrMap            = gox.AsStrMap
	AsArray             = gox.AsArray
	AsStrArr            = gox.AsStrArr
	AsMap               = gox.AsMap
	IsLong              = gox.IsLong
	IsInt               = gox.IsInt
	Dtoa                = gox.Dtoa
	Ftoa                = gox.Ftoa
	Ltoa                = gox.Ltoa
	Itoa                = gox.Itoa
	MapEq               = gox.MapEq
	AsOID               = gox.AsOID
	NewOID              = gox.NewOID
	NewOIDHex           = gox.NewOIDHex
)

type ObjectID = gox.ObjectID

func tryAsBsonId(val any) (string, bool) {
	switch val := val.(type) {
	case bson.ObjectId:
		return val.Hex(), true
	case ObjectID:
		return val.Hex(), true
	}
	return "", false
}
