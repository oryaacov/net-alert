package dm

//Site is a known address that can be assigned to a profile
type Site struct {
	IP        string `gorm:"primary_key"`
	Domain    string
	Nickname  string
	ProfileID string `gorm:"foreignkey:Mac"`
}
