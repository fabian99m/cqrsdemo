package main

import (
	awsCfg "cqrsdemo/config/aws"
	propsCfg "cqrsdemo/config/props"
	restCfg "cqrsdemo/config/rest"

	"cqrsdemo/entrypoint"
	restH "cqrsdemo/handler/rest"

	"log/slog"
	"os"

	"sync"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	appConfig := propsCfg.ReadConfig()
	evtProps := appConfig.EventsProps

	sqsOperations := awsCfg.NewSqsClient()
	snsOperations := awsCfg.NewSnsClient()
	s3Operations := awsCfg.NewS3Client()

	wg := sync.WaitGroup{}
	wg.Add(1)

	sqsListener := entrypoint.NewSqsListener(sqsOperations)

	messageHandler := entrypoint.NewMessagesHandler(snsOperations, appConfig)
	commandRestHanlder := restH.NewCommandRestHandler(snsOperations, messageHandler.GetCmds(), evtProps.Get())
	fileRestHandler := restH.NewFileHandler(s3Operations, appConfig.Aws.GetBucketProps())
	groupHandler := restH.NewGruopHandler(commandRestHanlder, fileRestHandler)

	go sqsListener.Listen(messageHandler, evtProps.QueueName)
	go entrypoint.RestServer(restCfg.NewBaseRestServer(groupHandler))

	wg.Wait()
}
