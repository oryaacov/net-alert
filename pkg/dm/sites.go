package dm

//Site is a known address that can be assigned to a profile
type Site struct {
	IP        string
	Domain    string
	Nickname  string `primary_key"`
	ProfileID string `gorm:"foreignkey:Mac"`
}
