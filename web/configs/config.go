package configs

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var Conf *Config

func InitConfig() {
	Conf = getConfig()
}

type Config struct {
	Mongo Mongo `yaml:"mongo"`
	Log   Log   `yaml:"log"`
}

type FieldInfo struct {
	Type  string      `yaml:"type"`
	Name  string      `yaml:"name"`
	Extra []FieldInfo `yaml:"extra"`
}

type Log struct {
	FilePath string `yaml:"filePath"`
    FileName string `yaml:"fileName"`
}

type Mongo struct {
	User        string       `yaml:"user"`
	Password    string       `yaml:"password"`
	Host        string       `yaml:"host"`
	Port        string       `yaml:"port"`
	Timeout     int64        `yaml:"timeout"`
	DB          string       `yaml:"db"`
	Collections []Collection `yaml:"collections"`
}

type Collection struct {
	Name    string   `yaml:"name"`
	Indexes []string `yaml:"indexes"`
	Opt     Opt      `yaml:"opt"`
}

type Opt struct {
	Timeseries Timeseries `yaml:"timeseries"`
}

func (o *Opt) GetMongoCollectionOpt() *options.CreateCollectionOptions {
	return &options.CreateCollectionOptions{
		TimeSeriesOptions: &options.TimeSeriesOptions{
			TimeField:   o.Timeseries.TimeField,
			MetaField:   &o.Timeseries.MetaField,
			Granularity: &o.Timeseries.Granularity,
		},
	}
}

type Timeseries struct {
	TimeField   string `yaml:"timeField"`
	MetaField   string `yaml:"metaField"`
	Granularity string `yaml:"granularity"`
}

func (c *Collection) GetIndexModels() []mongo.IndexModel {
	var ms []mongo.IndexModel
	for _, index := range c.Indexes {
		ms = append(ms, mongo.IndexModel{Keys: bson.M{"name": index}})
	}
	return ms
}

func (m *Mongo) GetMongoUrl() string {
	if m.User == "" || m.Password == "" {
		return fmt.Sprintf("mongodb://%s:%s/%s?retryWrites=true", m.Host, m.Port, m.DB)
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?retryWrites=true", m.User, m.Password, m.Host, m.Port, m.DB)
}

func getConfig() *Config {
	yamlFile, err := os.ReadFile(*ConfigPath)
	if err != nil {
		panic(err)
	}
	var c Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	return &c
}
