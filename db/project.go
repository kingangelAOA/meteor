package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	. "meteor/common"
	. "meteor/models"
	"strings"
)

const (
	projectCollection   = "project"
	OriginRefField      = "$ref"
	ProcessRefField     = "#ref#"
	SwaggerVersionField = "swaggers.version"
)

//GetAllProjects get all projects
func GetAllProjects(qp QueryProject) (*ProjectResponse, error) {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	count := make(chan int64, 1)
	var filter interface{}
	if qp.Name == "" {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"name", qp.Name}}
	}
	go func() {
		if c, err := p.CountDocuments(ctx, filter); err == nil {
			count <- c
		} else {
			count <- -1
		}
	}()
	if c, err := p.Find(ctx, filter, qp.GetOptions()); err != nil {
		return nil, err
	} else {
		var projects []Project
		if err = c.All(ctx, &projects); err != nil {
			return nil, err
		}
		total := <-count
		if total == -1 {
			return nil, ProjectTotalError
		}
		return &ProjectResponse{Projects: projects, Total: total}, nil
	}
}

func createProject(project Project) error {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	if _, err := p.InsertOne(ctx, NewProject{Name: project.Name, Comments: project.Comments}); err != nil {
		return err
	}
	return nil
}

func updateProject(project Project) error {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	update := bson.D{{"$set", bson.M{"name": project.Name, "comments": project.Comments}}}
	if _, err := p.UpdateOne(ctx, bson.M{"_id": project.ID}, update); err != nil {
		return err
	}
	return nil
}

//CreateOrUpdateProject create or update project
func CreateOrUpdateProject(project Project) error {
	if project.ID.IsZero() {
		return createProject(project)
	} else {
		return updateProject(project)
	}
}

func swaggerExsit(id primitive.ObjectID, version string) (bool, error) {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	count, err := p.CountDocuments(ctx, bson.M{"_id": id, "swaggers": bson.M{"$elemMatch": bson.M{"version": version}}})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func updateSwagger(id primitive.ObjectID, stb SwaggerFromDB) error {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	if _, err := p.UpdateOne(ctx, bson.M{"_id": id, "swaggers.version": stb.Version},
		bson.M{"$set": bson.M{"swaggers.$.swagger": stb.Swagger}}); err != nil {
		return err
	}
	return nil
}

func pushSwagger(id primitive.ObjectID, stb SwaggerFromDB) error {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	if r, err := p.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"swaggers": stb}}); err != nil {
		return err
	} else {
		fmt.Println(r)
	}
	return nil
}

//CreateOrUpdateSwagger create or update swagger
func CreateOrUpdateSwagger(cusr CreateUpdateSwaggerRequest) error {
	swagger, err := cusr.GetSwagger()
	if err != nil {
		return err
	}
	id := cusr.ID
	stb := SwaggerFromDB{
		Version: swagger.Info.Version,
		Swagger: processSwaggerJson(cusr.Swagger),
	}
	if r, e := swaggerExsit(id, stb.Version); e != nil {
		return err
	} else {
		if r {
			if err := updateSwagger(id, stb); err != nil {
				return err
			}
		} else {
			if err := pushSwagger(id, stb); err != nil {
				return err
			}
		}

	}
	return nil
}

func GetSwaggerVersions(id primitive.ObjectID) ([]interface{}, error) {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	return p.Distinct(ctx, SwaggerVersionField, bson.M{"_id": id})
}

func GetSwagger(id primitive.ObjectID, version string) (string, error) {
	ctx, cancel, p := getCollection(projectCollection)
	defer cancel()
	matchProject := bson.D{{"$match", bson.D{{"_id", id}}}}
	unwind := bson.D{{"$unwind", "$swaggers"}}
	matchVersion := bson.D{{"$match", bson.D{{"swaggers.version", version}}}}
	project := bson.D{{"$project", getSwaggerFieldsCondition()}}
	c, err := p.Aggregate(ctx, mongo.Pipeline{matchProject, unwind, matchVersion, project})
	if err != nil {
		return "", err
	}
	var sfds []SwaggerFromDB
	if err := c.All(ctx, &sfds); err != nil {
		return "", err
	}
	if sfds == nil {
		return "", nil
	}
	return reductionSwaggerJson(sfds[0].Swagger), nil
}

func getSwaggerFieldsCondition() bson.D {
	swaggerFields := []string{"version", "swagger"}
	var r bson.D
	for _, field := range swaggerFields {
		r = append(r, bson.E{Key: field, Value: fmt.Sprintf("$swaggers.%s", field)})
	}
	return r
}

func processSwaggerJson(sj string) string {
	return strings.Replace(sj, OriginRefField, ProcessRefField, -1)
}

func reductionSwaggerJson(sj string) string {
	return strings.Replace(sj, ProcessRefField, OriginRefField, -1)
}
