package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"pxsemic.com/simplebank/util"
)

var testStore Store

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("can't load config app.env. err:", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		panic(err)
	}
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
