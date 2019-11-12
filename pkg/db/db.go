package db

import (
	"net-alert/pkg/config"
	"net-alert/pkg/dm"
	"net-alert/pkg/logging"
	"net-alert/pkg/sniffer"

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
	createTables(db)
	return db
}

//CloseDBConnection closes the connection to the DB
func CloseDBConnection(db *gorm.DB) {
	db.Close()
}

func createTables(db *gorm.DB) {
	profile := &dm.Profile{}
	ipv4 := &sniffer.IPV4Record{Src: "127.0.0.1", Dst: "31.13.22.44"}
	if !db.HasTable(ipv4) {
		db.CreateTable(ipv4)
	}
	if !db.HasTable(profile) {
		db.CreateTable(profile)
	}
}
