package config

import (
	"github.com/fabian99m/cqrsdemo/model"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	EventsProps   Events        `yaml:"events"`
	Aws           Aws           `yaml:"aws"`
	BdConnnection BdConnnection `yaml:"bdConnnection"`
}

type BdConnnection struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslmode"`
	DbName   string `yaml:"dbname"`
	Schema   string `yaml:"schema"`
}

type Events struct {
	QueueName string `yaml:"queueName"`
	TopicArn  string `yaml:"topicArn"`
}

func (r *Events) Get() *model.EventProps {
	return &model.EventProps{
		QueueName: r.QueueName,
		TopicArn:  r.TopicArn,
	}
}

type Aws struct {
	S3 S3 `yaml:"s3"`
}

func (r *Aws) GetBucketProps() *model.BucketProps {
	return &model.BucketProps{
		Name:    r.S3.Bucket,
		MaxSize: r.S3.MaxSize,
	}
}

type S3 struct {
	Bucket  string `yaml:"bucket"`
	MaxSize int32  `yaml:"maxSize"`
}

func ReadConfig() *AppConfig {
	yamlFile, err := os.ReadFile("app.yml")
	if err != nil {
		panic(err)
	}

	var app AppConfig
	err = yaml.Unmarshal(yamlFile, &app)
	if err != nil {
		panic(err)
	}

	return &app
}
