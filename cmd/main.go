package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"orders"
	"orders/pkg/client/repositoryNats"
	"orders/pkg/client/serviceNats"
	"orders/pkg/server/handler"
	"orders/pkg/server/repository"
	"orders/pkg/server/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgres(repository.ConfigDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.db_name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("error initializing db: %s", err.Error())
	}

	cache, err := repository.CreateCache(db)
	if err != nil {
		log.Fatalf("error initializing cache: %s", err.Error())
	}

	repos := repository.NewRepository(cache)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(orders.Server)
	go func() {

		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	reposNats := repositoryNats.NewRepository(db, cache)
	servicesNats := serviceNats.NewService(reposNats)

	err = orders.Connect(orders.ConfigNast{
		ClusterID:   viper.GetString("nats.cluster_id"),
		ClientID:    viper.GetString("nats.client_id"),
		Subject:     viper.GetString("nats.subject"),
		QGroup:      viper.GetString("nats.q_group"),
		DurableName: viper.GetString("nats.durable_name"),
	}, servicesNats.InitHandler())

	if err != nil {
		log.Fatalf("error creating connect with nats server: %s", err.Error())
	}
	log.Print("Application Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Application Shutting Down")

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occured on server shutting down: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		log.Printf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
