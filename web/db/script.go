package db

import (
	"errors"
	"web/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ScriptTable = "script"

type Script struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type string             `json:"type,omitempty" bson:"type,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Code string             `json:"code,omitempty" bson:"code,omitempty"`
}

func CreateScript(bs models.ScriptVO) (primitive.ObjectID, error) {
	ctx, cancel, c := getCollection(ScriptTable)
	defer cancel()
	r, err := c.InsertOne(ctx, bs.GetCreateBson())
	if err != nil {
		return primitive.NilObjectID, errors.New("insert script error")
	}
	if id, err := getIDFromInsertOneResult(*r); err != nil {
		return primitive.NilObjectID, err
	} else {
		return id, nil
	}
}

func UpdateScript(bs models.ScriptVO) error {
	ctx, cancel, c := getCollection(ScriptTable)
	defer cancel()
	objectID, err := ConverObjectID(bs.ID)
	if err != nil {
		return err
	}
	_, err = c.UpdateByID(ctx, objectID, bs.GetUpdateBson())
	if err != nil {
		return errors.New("insert script error")
	}
	return nil
}

func GetScriptByID(id primitive.ObjectID) (*Script, error) {
	ctx, cancel, c := getCollection(ScriptTable)
	defer cancel()
	sr := c.FindOne(ctx, bson.M{"_id": id})
	var s Script
	err := sr.Decode(&s)
	if err != nil {
		return nil, errors.New("decode script error")
	}
	return &s, nil
}

func GetScripts() ([]Script, error) {
	ctx, cancel, c := getCollection(ScriptTable)
	defer cancel()
	cursor, err := c.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("find scripts error")
	}
	var ss []Script
	if err := cursor.All(ctx, &ss); err != nil {
		return nil, errors.New("decode scripts error")
	}
	return ss, nil
}
