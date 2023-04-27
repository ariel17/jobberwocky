package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ariel17/jobberwocky/internal/adapters/clients"
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

	router := gin.Default()
	jh.ConfigureRoutes(router)
	sh.ConfigureRoutes(router)

	router.Run(configs.GetHTTPAddress())
}