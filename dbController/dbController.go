package dbController

import "github.com/syndtr/goleveldb/leveldb"

type DataBase struct {
	Db *leveldb.DB
}

func (database *DataBase) OpenDb() bool {
	db, err := leveldb.OpenFile("data", nil)
	if err != nil {
		return false
	}
	database.Db = db
	return true
}

func (database *DataBase) CloseDb() bool {
	err := database.Db.Close()
	if err != nil {
		return false
	}
	return true
}
