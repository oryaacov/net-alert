package db

import (
	"net-alert/pkg/config"
	"net-alert/pkg/dm"
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
	createTables(db)
	return db
}

//CloseDBConnection closes the connection to the DB
func CloseDBConnection(db *gorm.DB) {
	db.Close()
}

//GetAllProfiles query with a simple get all query all of the profiles
func GetAllProfiles(db *gorm.DB) ([]dm.Profile, error) {
	profiles := []dm.Profile{}
	if result := db.Find(&profiles); result.Error != nil {
		return nil, result.Error
	}
	return profiles, nil
}

//GetOwner return the first (and should be only) DB row
func GetOwner(db *gorm.DB) (*dm.Owner, error) {
	owner := &dm.Owner{}
	if result := db.First(owner); result.Error != nil {
		return nil, result.Error
	}
	return owner, nil
}

func createTables(db *gorm.DB) {
	profile := &dm.Profile{}
	site := &dm.Site{}
	ipv4 := &dm.IPV4Record{}
	owner := &dm.Owner{}
	if !db.HasTable(owner) {
		db.CreateTable(owner)
	}
	if !db.HasTable(site) {
		db.CreateTable(site)
	}
	if !db.HasTable(ipv4) {
		db.CreateTable(ipv4)
	}
	if !db.HasTable(profile) {
		db.CreateTable(profile)
	}
	db.AutoMigrate(ipv4, site, profile)
}
