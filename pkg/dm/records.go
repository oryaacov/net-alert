package dm

import "time"

//IPV4Record contains a src dst ipv4 record at the DB
type IPV4Record struct {
	ID  uint `gorm:"PRIMARY_KEY`
	Src string
	Dst string
	TS  time.Time
}
