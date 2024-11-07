package auto

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fabian99m/cqrsdemo/adapter"
	config "github.com/fabian99m/cqrsdemo/config/props"
	restH "github.com/fabian99m/cqrsdemo/handler/rest"
	"github.com/fabian99m/cqrsdemo/util"
)

func FileHandler(client *s3.Client) *restH.FileHandler {
	app := config.ReadAppConfig()
	if err := util.ValidateStruct(app.Aws.S3); err != nil {
		log.Fatal("incomplete configuration: ", err)
	}

	s3Operations := adapter.NewS3Actions(client)
	return restH.NewFileHandler(s3Operations, app.Aws.GetBucketProps())
}
