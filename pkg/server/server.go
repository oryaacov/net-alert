package server

import (
	"fmt"
	"log"
	"net-alert/pkg/config"
	"net-alert/pkg/db"
	"net-alert/pkg/sniffer"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/gopacket/pcap"
	cors "github.com/itsjamie/gin-cors"
	"github.com/jinzhu/gorm"
)

type Server struct {
	Config  *config.Configuration
	DB      *gorm.DB
	Router  *gin.Engine
	Sniffer *sniffer.Sniffer
}

func (s *Server) Start(path string) {
	fmt.Println("starting net-alert...\nreading configuraion file...")
	s.Config = config.ReadConfigutionFromFile(path)
	fmt.Println("done!\ninit db connection...")
	s.DB = db.InitDB(s.Config)
	defer s.DB.Close()
	fmt.Println("done!\ninit sniffer & opening pcap...")
	s.InitSniffer()
	defer s.Sniffer.Handler.Close()
	fmt.Println("done!\ninit http-server...")
	go s.Sniffer.Start(s.DB)
	s.InitGin()
}

//InitSniffer configuring the sniffer and opening the pcap
func (s *Server) InitSniffer() {
	var err error
	s.Sniffer = &sniffer.Sniffer{Device: s.Config.Sniffer.Device, Promiscuous: s.Config.Sniffer.Promiscuous, Timeout: time.Duration(s.Config.Sniffer.Timeout) * time.Second, SnapshotLen: s.Config.Sniffer.SnapshotLen}
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

	router.GET("/isAlive", s.IsAlive())

	router.Run(fmt.Sprintf("%s:%d", s.Config.WebServer.URL, s.Config.WebServer.Port))
}
