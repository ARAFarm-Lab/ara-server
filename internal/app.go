package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
)

const (
	maxConnectionRetryAttempts = 5
)

func main() {
	// initialize config
	config, err := configuration.InitializeConfig()
	if err != nil {
		log.Fatal(nil, err, "init config got error")
	}

	// init logger
	log.NewLogger(config)

	// initialize db
	appConfig := config.GetConfig()
	dbInstance, err := initDB(appConfig.DB)
	if err != nil {
		log.Fatal(nil, err, "init DB got error")
	}

	mqttClient, err := initMQTT(appConfig.MQTT)
	if err != nil {
		log.Fatal(nil, err, "init MQTT got error")
	}

	router := initHTTPServer(config)

	// initialize layers
	infra := infrastructure.NewInfrastructure(config)
	repoDB := db.NewRepository(dbInstance, infra)
	repoMQ := mq.NewRepository(infra, mqttClient)
	usecase := usecase.NewUsecase(infra, repoDB, repoMQ)

	handlerHTTP := httpHandler.NewHandler(usecase)
	mqHandler.InitHandler(usecase, mqttClient)

	handlerHTTP.RegisterHTTPHandler(router)

	server := http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(nil, err, "")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(nil, nil, "Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully stop the server and its dependencies
	if err := server.Shutdown(ctx); err != nil {
		log.Error(nil, err, "Server Shutdown Error")
	}

	// Closing mqtt connection
	mqttClient.Disconnect(1000)

	select {
	case <-ctx.Done():
		log.Info(nil, nil, "timeout of 5 seconds.")
	default:
		log.Info(nil, nil, "Server exiting")
	}
}

func initDB(config configuration.DBConfig) (*sqlx.DB, error) {
	var connectingError error
	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s TimeZone=Asia/Jakarta sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.DBName,
	)

	for i := 0; i < maxConnectionRetryAttempts; i++ {
		log.Info(nil, nil, fmt.Sprintf("Connecting to DB (%d/%d)", i+1, maxConnectionRetryAttempts))
		db, err := sqlx.Connect("postgres", dbConnectionString)
		if err != nil {
			connectingError = err
			time.Sleep(1 * time.Second)
			continue
		}

		log.Info(nil, nil, "Connected to DB")
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

func initMQTT(config configuration.MQTTConfig) (mqtt.Client, error) {
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

	log.Info(nil, nil, "Connected to MQTT Broker")
	return client, nil
}
