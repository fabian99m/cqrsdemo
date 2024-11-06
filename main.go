package main

import (
	"github.com/fabian99m/cqrsdemo/config/auto"
	awsCfg "github.com/fabian99m/cqrsdemo/config/aws"
//	propsCfg "github.com/fabian99m/cqrsdemo/config/props"
	restCfg "github.com/fabian99m/cqrsdemo/config/rest"
	"github.com/go-chi/chi/v5"

	"github.com/fabian99m/cqrsdemo/entrypoint"
//	restH "github.com/fabian99m/cqrsdemo/handler/rest"

	//"log/slog"
	//"os"

	"sync"
)

func main() {
	router := chi.NewRouter()

	s3Client := awsCfg.NewS3Client()
	restCfg.NewFileRouter(router, auto.FileHandler(s3Client))

	wg := sync.WaitGroup{}
	wg.Add(1)

	go entrypoint.RestServer(router)

	wg.Wait()
}

func main2() {
	/* logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	appConfig := propsCfg.ReadConfig[propsCfg.AppConfig]("app.yml")
	evtProps := appConfig.EventsProps

	sqsOperations := awsCfg.NewSqsClient()
	snsOperations := awsCfg.NewSnsClient()
	s3Client := awsCfg.NewS3Client()

	wg := sync.WaitGroup{}
	wg.Add(1)

	sqsListener := entrypoint.NewSqsListener(sqsOperations)

	messageHandler := entrypoint.NewMessagesHandler(snsOperations, appConfig)
	commandRestHanlder := restH.NewCommandRestHandler(snsOperations, messageHandler.GetCmds(), evtProps.Get())
	fileRestHandler := auto.FileHandler(s3Client)
	groupHandler := restH.NewGruopHandler(commandRestHanlder, fileRestHandler)

	go sqsListener.Listen(messageHandler, evtProps.QueueName)
	go entrypoint.RestServer(restCfg.NewBaseRestServer(groupHandler))

	wg.Wait() */
}
