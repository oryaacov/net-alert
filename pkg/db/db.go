package db

import (
	"net-alert/pkg/config"
	"net-alert/pkg/logging"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//InitDB initilaize the DB connection using gorm
func InitDB(config *config.Configuration) *gorm.DB {
	var db *gorm.DB
	var err error
	db, err = gorm.Open(config.DB.Driver, config.DB.ConnectionString)
	if err != nil {
		logging.LogFatal(err)
	}
	createDBIfNotExists()
	return db
}

//CloseDBConnection closes the connection to the DB
func CloseDBConnection(db *gorm.DB) {
	db.Close()
}

func createDBIfNotExists() {

}
