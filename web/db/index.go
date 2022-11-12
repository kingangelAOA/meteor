package db

import (
	"common"
	"common/bindmodels"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"web/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DataBase = "meteor"

var Mongo *MongoDB

var once sync.Once

type MongoDB struct {
	client      *mongo.Client
	timeout     int64
	collections []configs.Collection
}

// InitMongo *mongo.Client
func InitMongo(m configs.Mongo) {
	once.Do(func() {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.GetMongoUrl()).SetSocketTimeout(time.Duration(m.Timeout)*time.Millisecond))
		if err != nil {
			panic(err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			panic(err)
		}
		Mongo = &MongoDB{
			client:      client,
			timeout:     m.Timeout,
			collections: m.Collections,
		}
		if err := Mongo.InitCollections(); err != nil {
			panic(err)
		}
	})
}

func (m *MongoDB) getConnectCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), time.Duration(m.timeout)*time.Millisecond)
}

func (m *MongoDB) InitCollections() error {
	ctx, cancel := Mongo.getConnectCtx()
	defer cancel()
	db := m.client.Database(DataBase)
	nameOnly := true
	for _, c := range m.collections {
		names, err := db.ListCollectionNames(ctx, bson.M{}, &options.ListCollectionsOptions{
			NameOnly: &nameOnly,
		})
		if err != nil {
			panic(err)
		}
		flag := false
		if c.Opt.Timeseries.TimeField != "" {
			for _, name := range names {
				if c.Name == name {
					flag = true
					break
				}
			}
		}
		if flag {
			continue
		}
		if err := db.CreateCollection(ctx, c.Name, c.Opt.GetMongoCollectionOpt()); err != nil {
			return err
		}
		indexs := c.GetIndexModels()
		if len(indexs) > 0 {
			if _, err := db.Collection(c.Name).Indexes().CreateMany(ctx, indexs); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

func getCollection(c string) (context.Context, context.CancelFunc, *mongo.Collection) {
	ctx, cancel := Mongo.getConnectCtx()
	col := Mongo.client.Database(DataBase).Collection(c)
	return ctx, cancel, col
}

func getGridFS() (*gridfs.Bucket, error) {
	bucket, err := gridfs.NewBucket(
		Mongo.client.Database(DataBase),
	)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func commonCreate(doc interface{}, collection string) (primitive.ObjectID, error) {
	ctx, cancel, c := getCollection(collection)
	defer cancel()
	b, err := convertToBson(doc)
	b["createTime"] = time.Now().In(common.GetLocalZone())
	if err != nil {
		return primitive.NilObjectID, err
	}
	if ior, err := c.InsertOne(ctx, &b); err != nil {
		return primitive.NilObjectID, fmt.Errorf("create %s data error", collection)
	} else {
		return getIDFromInsertOneResult(*ior)
	}
}

func commontUpdate(doc interface{}, id primitive.ObjectID, collection string) error {
	if id == primitive.NilObjectID {
		return errors.New("must have id when updating node")
	}
	b, err := convertToBson(doc)
	b["updateTime"] = time.Now().In(common.GetLocalZone())
	if err != nil {
		return err
	}
	ctx, cancel, c := getCollection(collection)
	defer cancel()
	if _, err := c.UpdateByID(ctx, id, bson.M{"$set": b}); err != nil {
		return fmt.Errorf("update %s data error", collection)
	}
	return nil
}

func commonGet[T any](collection string, filterOption bson.M, options ...*options.FindOptions) ([]T, error) {
	ctx, cancel, c := getCollection(collection)
	defer cancel()
	cursor, err := c.Find(ctx, filterOption, options...)
	var tt T
	if err != nil {
		return nil, fmt.Errorf("find %s error: %s", common.GetType(tt), err.Error())
	}
	return decodeCursor[T](ctx, cursor)
}

func commonAggregate[T any](collection string, pipeline []bson.M, options ...*options.AggregateOptions) ([]T, error) {
	ctx, cancel, c := getCollection(collection)
	defer cancel()
	var tt T
	cursor, err := c.Aggregate(ctx, pipeline, options...)
	if err != nil {
		return nil, fmt.Errorf("aggregate %s error: %s", common.GetType(tt), err.Error())
	}
	return decodeCursor[T](ctx, cursor)
}

func commontGetByID[T any](id primitive.ObjectID, filterOption bson.M, collection string, opts ...*options.FindOneOptions) (T, error) {
	ctx, cancel, c := getCollection(collection)
	defer cancel()
	var t T
	filter := bson.M{"_id": id}
	filter = common.MergeBsonM(filter, filterOption)
	sr := c.FindOne(ctx, filter, opts...)
	if err := sr.Decode(&t); err != nil {
		return t, fmt.Errorf("decode %s error in  CommontGetByID, msg: %s", common.GetType(t), err.Error())
	}
	return t, nil
}

func decodeCursor[T any](ctx context.Context, cursor *mongo.Cursor) ([]T, error) {
	var ts []T
	var tt T
	for cursor.Next(ctx) {
		var t T
		if err := cursor.Decode(&t); err != nil {
			return nil, fmt.Errorf("cursor decode %s error: %s", common.GetType(tt), err.Error())
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func decodeCursorSingleResult[T any](ctx context.Context, cursor *mongo.Cursor) (T, error) {
	var t T
	if err := cursor.Decode(&t); err != nil {
		var zero T
		return zero, fmt.Errorf("cursor decode %s error: %s", common.GetType(t), err.Error())
	}
	return t, nil
}

func getAggregateOpt() *options.AggregateOptions {
	return options.Aggregate().SetMaxTime(2 * time.Second)
}

type DBRef struct {
	Ref string             `bson:"$ref"`
	ID  primitive.ObjectID `bson:"$id"`
	DB  string             `bson:"$db"`
}

func newDBRef(ref string, id primitive.ObjectID) DBRef {
	return DBRef{
		Ref: ref,
		ID:  id,
		DB:  DataBase,
	}
}

func getIDFromInsertOneResult(ior mongo.InsertOneResult) (primitive.ObjectID, error) {
	return getIDFromInterface(ior.InsertedID)
}

func getIDFromInsertManyResult(imr *mongo.InsertManyResult) ([]primitive.ObjectID, error) {
	ps := make([]primitive.ObjectID, len(imr.InsertedIDs))
	for _, v := range imr.InsertedIDs {
		if oid, err := getIDFromInterface(v); err != nil {
			return ps, err
		} else {
			ps = append(ps, oid)
		}
	}
	return ps, nil
}

func getIDFromInterface(v interface{}) (primitive.ObjectID, error) {
	switch r := v.(type) {
	case string:
		objID, err := primitive.ObjectIDFromHex(r)
		if err != nil {
			return primitive.NilObjectID, errors.New("not objectid.ObjectID return")
		}
		return objID, nil
	case primitive.ObjectID:
		return r, nil
	default:
		return primitive.NilObjectID, errors.New("v type must be string or primitive.ObjectID")
	}
}

func GetIDFromInterface(v interface{}) (primitive.ObjectID, error) {
	return getIDFromInterface(v)
}

func getPagingOption(q bindmodels.Query) *options.FindOptions {
	option := getPaging(q.Limit, q.Offset)
	if q.SortField == "" {
		option.SetSort(bson.M{"createTime": -1})
	} else {
		option.SetSort(bson.M{q.SortField: q.SortType})
	}
	return option
}

func getFilterOption(q bindmodels.Query) (bson.M, error) {
	filterOption := bson.M{}
	if q.Search != "" {
		kv, err := parseSearch(q.Search)
		if err != nil {
			return nil, err
		}
		filterOption = common.MergeBsonM(filterOption, getRegexM(kv[0], kv[1]))
	}
	return filterOption, nil
}

func getPaging(limit, offset int64) *options.FindOptions {
	findOptions := &options.FindOptions{}
	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	return findOptions
}

func getRegexM(key, value string) bson.M {
	return bson.M{key: primitive.Regex{Pattern: value, Options: ""}}
}

func getRegexD(key, value string) bson.D {
	return bson.D{primitive.E{Key: key, Value: primitive.Regex{Pattern: value, Options: ""}}}
}

func parseSearch(s string) ([]string, error) {
	kv := strings.Split(s, ":")
	if len(kv) != 2 {
		return nil, errors.New("search params format error, must be key:value")
	}
	return kv, nil
}

func convertToBson(v interface{}) (doc bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func Trancaction(ctx context.Context, f func(sc mongo.SessionContext) error) error {
	return Mongo.client.UseSession(ctx, f)
}

func ConverObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, errors.New("not objectid.ObjectID return")
	}
	return objectID, nil
}
