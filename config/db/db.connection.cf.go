package config

import (
	props "github.com/fabian99m/cqrsdemo/config/props"
	
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDbConnection(appConfig *props.AppConfig) *gorm.DB {
	conf := appConfig.BdConnnection

	query := &url.Values{
		"sslmode":     []string{conf.SslMode},
		"search_path": []string{conf.Schema},
	}

	dsn := url.URL{
		User:     url.UserPassword(conf.User, conf.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:     conf.DbName,
		RawQuery: query.Encode(),
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  true,
			},
		),
	})

	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}

	return db
}
