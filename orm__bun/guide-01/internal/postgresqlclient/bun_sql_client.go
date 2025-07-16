package postgresqlclient

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var PostgresClientConnInstance IPostgresClientConn

type IPostgresClientConn interface {
	GetDB() *bun.DB
}

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type PostgresClientConn struct {
	PostgresConfig
	DB *bun.DB
}

func NewPostgresClient(config PostgresConfig) IPostgresClientConn {
	client := &PostgresClientConn{}
	client.PostgresConfig = config
	if err := client.Connect(); err != nil {
		log.Fatalf("ping to postgres failed: %v", err.Error())
	}
	return client
}

func (c *PostgresClientConn) Connect() error {
	postgresConn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%v:%v", c.Host, c.Port)),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithUser(c.Username),
		pgdriver.WithPassword(c.Password),
		pgdriver.WithDatabase(c.Database),
		pgdriver.WithInsecure(true),
	)

	db := bun.NewDB(sql.OpenDB(postgresConn), pgdialect.New(), bun.WithDiscardUnknownColumns())
	if err := db.Ping(); err != nil {
		return err
	}
	c.DB = db

	return nil
}

func (c *PostgresClientConn) GetDB() *bun.DB {
	return c.DB
}
