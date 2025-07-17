package sqlclient

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var BunSqlClientConnInstance IBunSqlClientConn

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
		log.Fatalf("Ping to postgres failed: %v", err.Error())
	}
	return client
}

func (c *BunSqlClientConn) Connect() error {
	postgresConn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%v:%v", c.Host, c.Port)),
		pgdriver.WithDatabase(c.Database),
		pgdriver.WithUser(c.Username),
		pgdriver.WithPassword(c.Password),
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithInsecure(true),
	)

	db := bun.NewDB(sql.OpenDB(postgresConn), pgdialect.New(), bun.WithDiscardUnknownColumns())
	if err := db.Ping(); err != nil {
		return err
	}
	c.DB = db

	return nil
}

func (c *BunSqlClientConn) GetDB() *bun.DB {
	return c.DB
}
