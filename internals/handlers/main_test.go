package handlers_test

import (
	"context"
	"log"
	"os"
	"testing"
	"user-favorites-service/pkg/postgresql"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDb *postgresql.DB

func TestMain(m *testing.M) {
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"
	host := "0.0.0.0"

	postgresPort := nat.Port("5432/tcp")
	postgresContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17-alpine",
			ExposedPorts: []string{postgresPort.Port()},
			Env: map[string]string{
				"POSTGRES_USER":     dbUser,
				"POSTGRES_PASSWORD": dbPassword,
				"POSTGRES_DB":       dbName,
			},

			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort(postgresPort),
			),
		},
		Started: true,
	})

	if err != nil {
		log.Fatalf("generic container :%v", err)
	}
	mapedport, err := postgresContainer.MappedPort(context.Background(), postgresPort)
	if err != nil {
		log.Fatalf("portmap :%v", err)
	}
	db := postgresql.NewDB(postgresql.DbConfig{Dbname: dbName, Dbuser: dbUser, Dbpassword: dbPassword, Host: host, Port: mapedport.Port()})
	testDb = db
	byteSql, err := os.ReadFile("../../db.sql")
	if err != nil {
		log.Fatalf("reading sql file :%v", err)
	}
	err = testDb.GetConnection().Exec(string(byteSql)).Error
	if err != nil {
		log.Fatalf("get connection :%v", err)
	}
	os.Exit(m.Run())
}
