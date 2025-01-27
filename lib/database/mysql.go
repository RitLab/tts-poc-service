package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
)

const (
	optionSingleStatement = "?parseTime=true&loc=UTC&multiStatements=false"
	optionMultiStatements = "?parseTime=true&loc=UTC&multiStatements=true"
)

func newDB(connect, option string) (*sql.DB, error) {
	dbms := "mysql"
	connect = strings.Join([]string{connect, option}, "")
	return sql.Open(dbms, connect)
}

func newConnect(host, database, user, password, port string) string {
	return strings.Join([]string{user, ":", password, "@", "tcp(", host, ":", port, ")/", database}, "")
}

func NewSqlHandler(log *baselogger.Logger, config *config.Cfg) *sql.DB {
	host := config.Database.Host
	database := config.Database.DbName
	user := config.Database.User
	password := config.Database.Password
	port := config.Database.Port

	connect := newConnect(host, database, user, password, port)
	db, err := newDB(connect, optionSingleStatement)
	if err != nil {
		log.Panic(err)
	}

	db.SetMaxIdleConns(config.Database.MaxIdle)
	db.SetMaxOpenConns(config.Database.MaxConn)
	db.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdletime))
	db.SetConnMaxLifetime(time.Duration(config.Database.MaxLifetime))

	return db
}
