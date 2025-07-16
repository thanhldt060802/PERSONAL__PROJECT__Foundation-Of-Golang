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

type IBunSqlClientConn interface {
	GetDB() *bun.DB
}

type BunSqlConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type BunSqlClientConn struct {
	BunSqlConfig
	DB *bun.DB
}

func NewBunSqlClient(config BunSqlConfig) IBunSqlClientConn {
	client := &BunSqlClientConn{}
	client.BunSqlConfig = config
	if err := client.Connect(); err != nil {
		log.Fatalf("ping to postgres failed: %v", err.Error())
	}
	return client
}

func (c *BunSqlClientConn) Connect() error {
	postgresConn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%v:%v", c.Host, c.Port)),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithUser(c.Username),
		pgdriver.WithPassword(c.Password),
		pgdriver.WithDatabase(c.Database),
		pgdriver.WithInsecure(true),
	)
	postgresDB := sql.OpenDB(postgresConn)

	db := bun.NewDB(postgresDB, pgdialect.New(), bun.WithDiscardUnknownColumns())
	if err := db.Ping(); err != nil {
		return err
	}
	c.DB = db

	return nil
}

func (c *BunSqlClientConn) GetDB() *bun.DB {
	return c.DB
}
