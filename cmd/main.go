package main

import (
	"fmt"
	"github.com/cownetwork/indigo/internal/eventhandler"
	"github.com/cownetwork/indigo/internal/psql"
	"github.com/cownetwork/indigo/internal/rpc"
	"github.com/cownetwork/mooapis-go/cow/indigo/v1"
	"github.com/upper/db/v4/adapter/postgresql"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	log.Println("Hello World!")

	connUrl := &postgresql.ConnectionURL{
		Host:     getEnvOrDefault("INDIGO_SERVICE_POSTGRES_URL", "localhost:5432"),
		User:     getEnvOrDefault("INDIGO_SERVICE_POSTGRES_USER", "test"),
		Password: getEnvOrDefault("INDIGO_SERVICE_POSTGRES_PASSWORD", "password"),
		Database: getEnvOrDefault("INDIGO_SERVICE_POSTGRES_DB", "test_database"),
	}

	log.Printf("Connecting to PostgresSQL at %s ...", connUrl.Host)

	sess, err := postgresql.Open(connUrl)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer sess.Close()

	log.Println("Connected to PostgresSQL.")

	log.Printf("Connecting to cloudevents client ...")

	eventhandler.Initialize(
		getBrokersFromEnv(),
		getEnvOrDefault("INDIGO_SERVICE_KAFKA_TOPIC", "cow.global.indigo"),
		getEnvOrDefault("INDIGO_SERVICE_CLOUDEVENTS_SOURCE", "cow.global.indigo-service"),
	)

	log.Println("Connected to cloudevents client.")

	// setup grpc server
	address := fmt.Sprintf("%s:%s", getEnvOrDefault("INDIGO_SERVICE_HOST", ""), getEnvOrDefault("INDIGO_SERVICE_PORT", "6969"))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	s.RegisterService(&indigo.IndigoService_ServiceDesc, &rpc.IndigoServiceServer{
		Dao: &psql.DataAccessor{Session: sess},
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getBrokersFromEnv() []string {
	list := getEnvOrDefault("INDIGO_SERVICE_KAFKA_BROKERS", "127.0.0.1:9092")
	return strings.Split(list, ",")
}

func getEnvOrDefault(env string, def string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		return def
	}
	return value
}
