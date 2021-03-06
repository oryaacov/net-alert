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
	Sites      []Site `gorm:"foreignkey:profile_id"`
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
			if err := tx.Where(Site{IP: p.Sites[i].IP}).Assign(Site{IP: p.Sites[i].IP, Domain: p.Sites[i].Domain, Nickname: p.Sites[i].Nickname, ProfileID: p.Mac}).FirstOrCreate(&p.Sites[i]).Error; err != nil {
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

//Update update the owner data
func (o *Owner) Update(db *gorm.DB) error {
	var owner Owner
	tx := db.Begin()
	if err := tx.First(&owner).Error; err != nil {
		tx.Rollback()
		return err
	}
	owner.Nickname = o.Nickname
	owner.Email = o.Email
	owner.Phone = o.Phone
	owner.GetEmailAlerts = o.GetEmailAlerts
	owner.GetSMSAlerts = o.GetSMSAlerts
	if err := tx.Save(owner).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
