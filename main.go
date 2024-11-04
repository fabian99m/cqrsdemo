package main

import (
	awsCfg "cqrsdemo/config/aws"
	propsCfg "cqrsdemo/config/props"
	restCfg "cqrsdemo/config/rest"

	"cqrsdemo/entrypoint"
	restH "cqrsdemo/handler/rest"
	//sqsH "cqrsdemo/handler/sqs"
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

	commandRestHanlder := restH.NewCommandRestHandler(snsOperations, entrypoint.CommandsHandler(), evtProps.Get())
	fileRestHandler := restH.NewFileHandler(s3Operations, appConfig.Aws.GetBucketProps())

	groupHandler := restH.NewGruopHandler(commandRestHanlder, fileRestHandler)

	//go sqsListener.Listen(handler.NewCarHandler(), "testcola")
	//go sqsListener.Listen(sqsH.NewUserHandler(), "usercola")
	go sqsListener.Listen(entrypoint.NewMessagesHandler(snsOperations, evtProps.Get()), evtProps.QueueName)

	go entrypoint.RestServer(restCfg.NewRestServer(groupHandler))

	wg.Wait()
}
