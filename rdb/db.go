package rdb

import (
	"database/sql"
	"log"
	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DbConnector () error {
	cfg := mysql.Config{
        User:   "root",
        Passwd: "",
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "recordings",
		AllowNativePasswords: true,
    }

	dbcfgstr := cfg.FormatDSN()
	var err error
	DB, err = sql.Open("mysql", dbcfgstr)
	if err != nil {
        log.Fatal(err)
		return err
    }

	pingErr := DB.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
		return pingErr
    }
    
	return nil
}