package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cronHandler "ara-server/internal/handler/cron"
	httpHandler "ara-server/internal/handler/http"
	mqHandler "ara-server/internal/handler/mq"
	"ara-server/internal/infrastructure"
	"ara-server/internal/infrastructure/configuration"
	"ara-server/internal/repository/db"
	"ara-server/internal/repository/mq"
	"ara-server/internal/usecase"
	"ara-server/util/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

const (
	maxConnectionRetryAttempts = 5
)

func main() {
	ctx := context.Background()

	// initialize config
	config, err := configuration.InitializeConfig()
	if err != nil {
		log.Fatal(ctx, nil, err, "init config got error")
	}

	// init logger
	log.NewLogger(config)

	// initialize db
	appConfig := config.GetConfig()
	dbInstance, err := initDB(ctx, appConfig.DB)
	if err != nil {
		log.Fatal(ctx, nil, err, "init DB got error")
	}

	mqttClient, err := initMQTT(ctx, appConfig.MQTT)
	if err != nil {
		log.Fatal(ctx, nil, err, "init MQTT got error")
	}

	router := initHTTPServer(config)

	// initialize layers
	infra := infrastructure.NewInfrastructure(config)
	repoDB := db.NewRepository(dbInstance, infra)
	repoMQ := mq.NewRepository(infra, mqttClient)
	usecase := usecase.NewUsecase(infra, repoDB, repoMQ)

	// initialize MQ handler
	mqHandler.InitHandler(usecase, mqttClient)

	// initialize HTTP handler
	handlerHTTP := httpHandler.NewHandler(usecase)
	handlerHTTP.RegisterHTTPHandler(router)

	// initialize cron handler
	cron := cron.New()
	cronHandler.InitHandler(usecase, cron)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// starting the server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(ctx, nil, err, "")
		}
	}()
	cron.Start()
	log.Info(ctx, nil, nil, "Service is running...")
	fmt.Println("Service is running...") // To notify docker

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(ctx, nil, nil, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully stop the server and its dependencies
	if err := server.Shutdown(ctx); err != nil {
		log.Error(ctx, nil, err, "Server Shutdown Error")
	}

	// Closing mqtt connection
	mqttClient.Disconnect(1000)

	// Stop cron
	cron.Stop()

	select {
	case <-ctx.Done():
		log.Info(ctx, nil, nil, "timeout of 5 seconds.")
	default:
		log.Info(ctx, nil, nil, "Server exiting")
	}
}

func initDB(ctx context.Context, config configuration.DBConfig) (*sqlx.DB, error) {
	var connectingError error
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s TimeZone=Asia/Jakarta sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.DBName,
	)

	for i := 0; i < maxConnectionRetryAttempts; i++ {
		log.Info(ctx, nil, nil, fmt.Sprintf("Connecting to DB (%d/%d)", i+1, maxConnectionRetryAttempts))
		db, err := sqlx.Connect("postgres", dbConnectionString)
		if err != nil {
			connectingError = err
			time.Sleep(1 * time.Second)
			continue
		}

		log.Info(ctx, nil, nil, "Connected to DB")
		return db, nil
	}

	return nil, connectingError
}

func initHTTPServer(config configuration.Config) *gin.Engine {
	if !config.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	ginConfig := cors.DefaultConfig()
	ginConfig.AllowAllOrigins = true
	ginConfig.AllowHeaders = []string{"Authorization", "Content-Type"}

	engine.Use(cors.New(ginConfig))

	return engine
}

func initMQTT(ctx context.Context, config configuration.MQTTConfig) (mqtt.Client, error) {
	// Now we establish the connection to the mqtt broker
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)

	opts.SetOrderMatters(false)       // Allow out of order messages (use this option unless in order delivery is essential)
	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick

	// Automate connection management (will keep trying to connect and will reconnect if network drops)
	opts.ConnectRetry = true
	opts.AutoReconnect = true

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	log.Info(ctx, nil, nil, "Connected to MQTT Broker")
	return client, nil
}
