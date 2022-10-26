package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func formatDsn(conf *Cfg) (string, string) {
	var suffix string

	if conf.Database.Host == "localhost" {
		suffix = "/" + conf.Database.Dbname
	} else {
		suffix = fmt.Sprintf("tcp(%s:3306)/%s",
			conf.Database.Host,
			conf.Database.Dbname)
	}

	return "mysql", fmt.Sprintf("%s:%s@%s",
		conf.Database.User,
		conf.Database.Password, suffix)
}

func Connect(conf *Cfg) (*sql.DB, error) {
	conn, err := sql.Open(formatDsn(conf))
	if err != nil {
		log.Fatalf("Can't connect to DB `%s` at `%s` as `%s`",
			conf.Database.Dbname, conf.Database.Host, conf.Database.User)
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("Can't connect to DB `%s` at `%s` as `%s`",
			conf.Database.Dbname, conf.Database.Host, conf.Database.User)
	}

	return conn, nil
}
