/**
* @description:
* @author
* @date 2026-03-25 23:27:11
* @version 1.0
*
* Change Logs:
* Date           Author       Notes
*
 */

package main

import (
	"database/sql"
	"log"

	"pxsemic.com/simplebank/api"
	db "pxsemic.com/simplebank/db/sqlc"
	"pxsemic.com/simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config app.env. err:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot open db", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
