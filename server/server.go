package main

import (
	"database/sql"
	"log"
	grpcapi "mailinglist/grpc-api"
	jsonapi "mailinglist/json-api"
	"mailinglist/mdb"
	"sync"

	arg "github.com/alexflint/go-arg"
)

var args struct {
	DbPath   string `arg:"env:MAILINGLIST_DB"`
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
	BindGrpc string `arg:"env:MAILINGLIST_BIND_GRPC"`
}

func main() {
	arg.MustParse(&args)
	if args.DbPath == "" {
		args.DbPath = "list.db"
	}
	if args.BindJson == "" {
		args.BindJson = ":8080"
	}
	if args.BindGrpc == "" {
		args.BindGrpc = ":8081"
	}
	log.Printf("using database '%v'\n", args.DbPath)
	db, err := sql.Open("sqlite3", args.DbPath)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	mdb.TryCreate(db)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("starting json api server...\n")
		jsonapi.Serve(db, args.BindJson)
	}()
	wg.Add(1)
	go func() {
		log.Printf("starting grpc api server...\n")
		grpcapi.Serve(db, args.BindGrpc)
	}()
	wg.Wait()
}
