package db

import (
	"bytes"
	"common"
	"common/bindmodels"
	"errors"
	"fmt"
	"os"
	"plugin/shared"
	"time"
	"web/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Plugin struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FileID     primitive.ObjectID `bson:"fileID,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Desc       string             `bson:"desc,omitempty"`
	Inputs     []Input            `bson:"inputs,omitempty"`
	Code       string             `bson:"code,omitempty"`
	Language   string             `bson:"language,omitempty"`
	CreateTime time.Time          `bson:"createTime,omitempty"`
	UpdateTime time.Time          `bson:"updateTime,omitempty"`
}

type Input struct {
	Name     string `bson:"name,omitempty"`
	Required bool   `bson:"required"`
	Desc     string `bson:"desc,omitempty"`
	Value    string `bson:"value,omitempty"`
}

type PluginFile struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	FileName string             `bson:"filename,omitempty"`
}

func GetPlugins(q bindmodels.Query) ([]Plugin, error) {
	filterOption, err := getFilterOption(q)
	if err != nil {
		return []Plugin{}, err
	}
	option := getPagingOption(q)
	return commonGet[Plugin](common.PluginCollection, filterOption, option)
}

func GetPluginList() ([]Plugin, error) {
	fo := &options.FindOptions{
		Projection: bson.M{"_id": 1, "name": 1, "language": 1, "fileID": 1},
	}
	return commonGet[Plugin](common.PluginCollection, bson.M{}, fo)
}

func CreatePlugin(p *Plugin) (primitive.ObjectID, error) {
	fID, err := CreatePluginFile(p)
	if err != nil {
		return primitive.NilObjectID, err
	}
	p.FileID = fID
	id, err := commonCreate(p, common.PluginCollection)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return id, nil
}

func GetPluginByID(id string) (Plugin, error) {
	oid, err := getIDFromInterface(id)
	if err != nil {
		return Plugin{}, err
	}
	return commontGetByID[Plugin](oid, bson.M{}, common.PluginCollection)
}

func UpdatePlugin(p *Plugin) error {
	if err := UpdatePluginFile(p); err != nil {
		return err
	}
	err := commontUpdate(p, p.ID, common.PluginCollection)
	if err != nil {
		return err
	}
	return nil
}

func DeletePlugin(id string) error {
	ctx, cancel, c := getCollection(common.PluginCollection)
	defer cancel()
	objectID, err := ConverObjectID(id)
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return errors.New("delete plugin error")
	}
	return deletePluginFile(objectID)
}

func GetPluginCount() (int64, error) {
	ctx, cancel, c := getCollection(common.PluginCollection)
	defer cancel()
	count, err := c.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, errors.New("get plugin count error")
	}
	return count, nil
}

func CreatePluginFile(p *Plugin) (primitive.ObjectID, error) {
	gb, err := getGridFS()
	if err != nil {
		return primitive.NilObjectID, err
	}
	uploadStream, err := gb.OpenUploadStream(p.Name)
	if err != nil {
		return primitive.NilObjectID, err
	}

	fID, err := getIDFromInterface(uploadStream.FileID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	pf := shared.NewPluginFactory(fID.Hex(), p.Code, p.Language, *configs.PluginRootPath)
	if err := pf.CheckFile(); err != nil {
		return primitive.NilObjectID, err
	}
	stream, err := pf.GetPluginStream()
	if err != nil {
		return primitive.NilObjectID, err
	}
	_, err = uploadStream.Write(stream)
	if err != nil {
		return primitive.NilObjectID, err
	}
	err = uploadStream.Close()
	if err != nil {
		return primitive.NilObjectID, err
	}
	return fID, nil
}

func UpdatePluginFile(p *Plugin) error {
	if err := deletePluginFile(p.FileID); err != nil {
		fmt.Println(err.Error())
	}
	gb, err := getGridFS()
	if err != nil {
		return err
	}
	uploadStream, err := gb.OpenUploadStreamWithID(p.FileID, p.ID.Hex())
	if err != nil {
		return err
	}
	pf := shared.NewPluginFactory(p.FileID.Hex(), p.Code, p.Language, *configs.PluginRootPath)
	stream, err := pf.GetPluginStream()
	if err != nil {
		return err
	}
	_, err = uploadStream.Write(stream)
	if err != nil {
		return err
	}
	err = uploadStream.Close()
	if err != nil {
		return err
	}
	return nil
}

func DownloadPluginFile(fileID, path string) error {
	gb, err := getGridFS()
	if err != nil {
		return err
	}
	ido, err := getIDFromInterface(fileID)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	_, err = gb.DownloadToStream(ido, &buf)
	if err != nil {
		return fmt.Errorf("download plugin(%s) file error, err: %s", path, err.Error())
	}
	if err := os.WriteFile(path, buf.Bytes(), 0755); err != nil {
		return err
	}
	return nil
}

func deletePluginFile(fileID primitive.ObjectID) error {
	gb, err := getGridFS()
	if err != nil {
		return err
	}
	err = gb.Delete(fileID)
	if err != nil {
		return errors.New("delete plugin file error")
	}
	return nil
}
