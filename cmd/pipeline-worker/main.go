package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	clickhouserepository "sales-record-orchestration/internal/adapters/clickhouse"
	kafkaproducer "sales-record-orchestration/internal/adapters/kafka"
	postgresrepository "sales-record-orchestration/internal/adapters/postgres"
	"sales-record-orchestration/internal/adapters/temporal"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create Postgres Conn
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("PG_USER"),
		viper.GetString("PG_PASSWORD"),
		viper.GetString("PG_DB_NAME"),
		viper.GetString("PG_SSL_MODE"),
	)

	pg, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	pgRepository := postgresrepository.InitPostgresRepository(pg)

	// Create Kafka Conn
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	brokerUrl := fmt.Sprintf(
		"%s:%s",
		viper.GetString("KAFKA_HOST"),
		viper.GetString("KAFKA_PORT"),
	)
	brokersUrl := []string{brokerUrl}
	producer, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		fmt.Println(err)
	}
	kafkaProducer := kafkaproducer.InitKafkaProducer(producer)

	// Create Clickhouse Conn
	address := fmt.Sprintf(
		"%s:%s",
		viper.GetString("CH_HOST"),
		viper.GetString("CH_PORT"),
	)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{address},
		Auth: clickhouse.Auth{
			Database: viper.GetString("CH_DB_NAME"),
			Username: viper.GetString("CH_USER"),
			Password: viper.GetString("CH_PASSWORD"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		log.Fatal(err)
	}

	chRepository := clickhouserepository.InitClickhouseRepository(&conn)

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort: "localhost:17233",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, "order-worker", worker.Options{})

	temporalActivity := temporal.InitActivities(
		pgRepository,
		kafkaProducer,
		chRepository,
	)

	temporalWorkflow := temporal.InitWorkflowWorker(temporalActivity)

	w.RegisterWorkflow(temporalWorkflow.SalesETLWorkflow)
	w.RegisterActivity(temporalActivity.FetchSalesActivity)
	w.RegisterActivity(temporalActivity.PublishSalesActivity)
	w.RegisterActivity(temporalActivity.ProcessSalesActivity)

	if err := w.Run(worker.InterruptCh()); err != nil {
		panic(err)
	}
}
