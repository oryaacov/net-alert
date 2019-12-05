package dm

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Owner the machine owner (only one master per machine)
type Owner struct {
	Mac            string `gorm:"primary_key"`
	Nickname       string
	IP             string
	Email          string
	Phone          string
	GetEmailAlerts bool
	GetSMSAlerts   bool
	LastLoginTime  time.Time
}

//Profile a someone on the network with a known mac address
type Profile struct {
	Mac        string `gorm:"primary_key"`
	NickName   string
	CreateDate time.Time
	Sites      []Site
}

//CreateOrUpdate insert in case of a new user or update an existing one
func (p *Profile) CreateOrUpdate(db *gorm.DB) error {
	tx := db.Begin()
	if err := tx.Where(Profile{Mac: p.Mac}).Assign(Profile{NickName: p.NickName, Mac: p.Mac, CreateDate: p.CreateDate}).FirstOrCreate(p).Error; err != nil {
		tx.Rollback()
		return err
	}
	if p.Sites != nil && len(p.Sites) > 0 {
		for i := 0; i < len(p.Sites); i++ {
			if err := tx.Where(Site{Nickname: p.Sites[i].Nickname}).Assign(Site{IP: p.Sites[i].IP, Domain: p.Sites[i].Domain, Nickname: p.Sites[i].Nickname, ProfileID: p.Mac}).FirstOrCreate(&p.Sites[i]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
