package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ariel17/jobberwocky/internal/adapters/clients"
	http2 "github.com/ariel17/jobberwocky/internal/adapters/http"
	jobHandler "github.com/ariel17/jobberwocky/internal/adapters/http/job"
	subscriptionHandler "github.com/ariel17/jobberwocky/internal/adapters/http/subscription"
	jobRepository "github.com/ariel17/jobberwocky/internal/adapters/repositories/job"
	"github.com/ariel17/jobberwocky/internal/adapters/repositories/keyword"
	"github.com/ariel17/jobberwocky/internal/adapters/repositories/subscription"
	"github.com/ariel17/jobberwocky/internal/core/services/job"
	"github.com/ariel17/jobberwocky/internal/core/services/notification"
	subscriptionService "github.com/ariel17/jobberwocky/internal/core/services/subscription"
	"github.com/ariel17/jobberwocky/resources/configs"
)

// @title           Jobberwocky API
// @version         0.1.0
// @description     A job posting and searching API.

// @contact.name   Ariel Gerardo Ríos
// @contact.url    http://ariel17.com.ar/
// @contact.email  arielgerardorios@gmail.com

// @license.name  MIT
// @license.url   https://github.com/ariel17/jobberwocky/blob/master/LICENSE.md

// @BasePath  /
func main() {
	db, err := gorm.Open(sqlite.Open(configs.GetDatabaseName()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	kr := keyword.NewKeywordRepository(db)
	kr.SyncSchemas()

	sr := subscription.NewSubscriptionRepository(db)
	sr.SyncSchemas()

	ec := clients.MockEmailProviderClient{}

	ns := notification.NewNotificationService(configs.GetNotificationWorkers(), sr, &ec, configs.GetEmailTemplate())
	ns.StartWorkers()
	defer ns.StopWorkers()

	httpClient := http.Client{}
	external := clients.NewJobberwockyExternalJobClient(&httpClient, configs.GetJobberwockyURL())

	jr := jobRepository.NewJobRepository(db)
	jr.SyncSchemas()

	js := job.NewJobService(jr, ns, external)

	jh := jobHandler.NewJobHTTPHandler(js)

	ss := subscriptionService.NewSubscriptionService(sr)
	sh := subscriptionHandler.NewSubscriptionHTTPHandler(ss)
	swh := http2.NewSwaggerHandler()

	router := gin.Default()
	jh.ConfigureRoutes(router)
	sh.ConfigureRoutes(router)
	swh.ConfigureRoutes(router)

	router.Run()
}