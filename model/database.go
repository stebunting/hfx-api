package model

import (
	"context"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Model struct {
	Db *pg.DB
}

func (s *Model) Connect() {
	address := os.Getenv("DATABASE_URL")

	options, err := pg.ParseURL(address)
	if err != nil {
		panic(err)
	}

	db := pg.Connect(options)

	ctx := context.Background()
	err = db.Ping(ctx)
	if err != nil {
		panic(err)
	}

	s.Db = db
}

func (s *Model) Disconnect() {
	s.Db.Close()
}

func (s *Model) CreateSchema() {
	models := []interface{}{
		(*Currency)(nil),
		(*Exchange)(nil),
	}

	for _, model := range models {
		err := s.Db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			log.Fatal(err)
		}
	}
}
