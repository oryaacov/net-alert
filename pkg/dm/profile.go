package dm

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Profile a someone on the network with a known mac address
type Profile struct {
	Mac        string `gorm:"primary_key"`
	NickName   string
	CreateDate time.Time
}

//CreateOrUpdate insert in case of a new user or update an existing one
func (p *Profile) CreateOrUpdate(db *gorm.DB) error {
	return db.Where(Profile{Mac: p.Mac}).Assign(Profile{NickName: p.NickName, Mac: p.Mac, CreateDate: p.CreateDate}).FirstOrCreate(p).Error
}
