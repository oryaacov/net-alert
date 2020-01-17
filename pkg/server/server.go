package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/gopacket/pcap"
	cors "github.com/itsjamie/gin-cors"
	"github.com/jinzhu/gorm"
	"log"
	"net-alert/pkg/config"
	"net-alert/pkg/db"
	"net-alert/pkg/sniffer"
	"net-alert/pkg/utils"
	"net/http"
	"time"
)

//Server represnt the project core objects at one place
type Server struct {
	Config  *config.Configuration
	DB      *gorm.DB
	Router  *gin.Engine
	Sniffer *sniffer.Sniffer
}

//Start runs the net-alert web server application
func (s *Server) Start(path string) {
	fmt.Println("starting net-alert...\nreading configuraion file...")
	s.Config = config.ReadConfigutionFromFile(path)
	fmt.Println("done!\ninit twilio SMS client...")
	utils.InitTwilio(s.Config.SMS.AccountSid, s.Config.SMS.AuthToken, s.Config.SMS.Number)
	fmt.Println("done!\ninit SMTP client...")
	utils.InitSMTP(s.Config.SMTP.SMTPServer, s.Config.SMTP.Port, s.Config.SMTP.EmailAddress, s.Config.SMTP.EmailPassword)
	fmt.Println("done!\ninit db connection...")
	s.DB = db.InitDB(s.Config)
	defer s.DB.Close()
	fmt.Println("done!\ninit http-server...")
	go openBrowser(fmt.Sprintf("http://%s:%d/site", s.Config.WebServer.URL, s.Config.WebServer.Port))
	s.InitGin()
}

//Exit capture the interupt signal and close monitor mode
func (s *Server) Exit() {
	sniffer.Exit(s.Config.Sniffer.DeviceMonName)
}

//InitSniffer configuring the sniffer and opening the pcap
func (s *Server) InitSniffer() {
	var err error
	s.Sniffer = &sniffer.Sniffer{Device: s.Config.Sniffer.DeviceMonName, Promiscuous: s.Config.Sniffer.Promiscuous, Timeout: time.Duration(s.Config.Sniffer.Timeout) * time.Second, SnapshotLen: s.Config.Sniffer.SnapshotLen, PcapFolder: s.Config.Sniffer.PcapsFolder}
	s.Sniffer.Handler, err = pcap.OpenLive(s.Sniffer.Device, s.Sniffer.SnapshotLen, s.Sniffer.Promiscuous, s.Sniffer.Timeout)
	if err != nil {
		log.Fatal(err)
	}
}

//InitGin init the gin based webserver
func (s *Server) InitGin() {
	router := gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         s.Config.WebServer.AllowedMethods,
		RequestHeaders:  s.Config.WebServer.AllowedHeaders,
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	router.Use(gin.Recovery()) //skip logger setup, since we already have one
	router.Use(VerifyHeader())

	router.StaticFS(fmt.Sprintf("/%s", s.Config.WebServer.SiteURL), http.Dir(s.Config.WebServer.StaticFilesLocation))

	router.GET("/api/alive", s.IsAlive())
	router.GET("/api/master", s.GetOwnerInfo())
	router.POST("/api/master", s.UpdateOwner())
	router.GET("/api/network", s.GetNetworkInfo())
	router.GET("/api/start", s.Sniff())
	router.GET("/api/profiles", s.GetAllProfiles())
	router.POST("/api/profiles", s.CreateOrUpdateProfile())
	router.Run(fmt.Sprintf("%s:%d", s.Config.WebServer.URL, s.Config.WebServer.Port))
}
