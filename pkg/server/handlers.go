package server

import (
	"net-alert/pkg/dm"
	"net-alert/pkg/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

//IsAlive returns "yes" if the server is up and running, part of DT standart
func (s *Server) IsAlive() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte(""))
	}
}

//CreateOrUpdateProfile create or update mac profile
func (s *Server) CreateOrUpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		profile := dm.Profile{}
		if !readBody(c, &profile) {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		if err := ((&profile).CreateOrUpdate(s.DB)); err != nil {
			logging.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Data(http.StatusOK, "text/html", []byte(""))
	}
}
