package dbController

import "github.com/syndtr/goleveldb/leveldb"

type Database struct {
	Db *leveldb.DB
}

func (database *Database) OpenDb() bool {
	db, err := leveldb.OpenFile("data", nil)
	if err != nil {
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
