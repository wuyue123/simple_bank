package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"pxsemic.com/simplebank/util"
)

var testQueries *Queries
var sqlDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can't load config app.env. err:", err)
	}
	sqlDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}
	testQueries = New(sqlDB)
	os.Exit(m.Run())
}
