package auto

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fabian99m/cqrsdemo/adapter"
	config "github.com/fabian99m/cqrsdemo/config/props"
	restH "github.com/fabian99m/cqrsdemo/handler/rest"
)

func FileHandler(client *s3.Client) *restH.FileHandler {
	app := config.ReadAppConfig()
	s3Operations := adapter.NewS3Actions(client)
	return restH.NewFileHandler(s3Operations, app.Aws.GetBucketProps())
}
