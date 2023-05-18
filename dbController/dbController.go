package dbController

import (
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

var DB Database

type Database struct {
	Db *leveldb.DB
}

var DbPath string = "\\data\\ValidatorNode"

func (database *Database) OpenDb() bool {
	cwd, _ := os.Getwd()
	db, err := leveldb.OpenFile(cwd+DbPath, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	database.Db = db
	return true
}

func (database *Database) CloseDb() bool {
	err := database.Db.Close()
	if err != nil {
		return false
	}
	return true
}

func (database *Database) SetValue(key string, value []byte) bool {
	err := database.Db.Put([]byte(key), value, nil)
	if err != nil {
		return false
	}
	return true
}

func (database *Database) GetValue(key string) ([]byte, bool) {
	value, err := database.Db.Get([]byte(key), nil)
	if err != nil {
		return value, false
	}
	return value, true
}

func (database *Database) RemoveValue(key string) bool {
	err := database.Db.Delete([]byte(key), nil)
	if err != nil {
		return false
	}
	return true
}
