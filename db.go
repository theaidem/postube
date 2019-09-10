package main

import (
	"log"
	"sync"

	"github.com/recoilme/pudge"
)

type database struct {
	db *pudge.Db
}

var (
	dbOnce     sync.Once
	dbInstance database
)

func db() *database {
	dbOnce.Do(func() {
		log.Println("Init Database")
		db, err := pudge.Open("data", nil)
		if err != nil {
			log.Panic(err)
		}
		dbInstance = database{
			db: db,
		}
	})

	return &dbInstance
}

func (d *database) set(key string, timestamp int64) error {
	if err := d.db.Set(key, timestamp); err != nil {
		return err
	}
	return nil
}

func (d *database) has(key string) (bool, error) {
	exists, err := d.db.Has(key)
	if err != nil {
		return false, err
	}
	return exists, nil
}
