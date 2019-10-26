package server

import (
	"fmt"
	"net-alert/pkg/config"
	"net-alert/pkg/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/jinzhu/gorm"
)

type Server struct {
	Config *config.Configuration
	DB     *gorm.DB
	Router *gin.Engine
}

func (s *Server) Start(path string) {
	fmt.Println("starting net-alert...\nreading configuraion file...")
	s.Config = config.ReadConfigutionFromFile(path)
	fmt.Println("done!\ninit db connection...")
	s.DB = db.InitDB(s.Config)
	fmt.Println("done!\nstarting http server...")

	s.InitGin()
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
