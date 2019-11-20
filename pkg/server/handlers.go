package server

import (
	"fmt"
	"net-alert/pkg/db"
	"net-alert/pkg/dm"
	"net-alert/pkg/logging"
	"net-alert/pkg/sniffer"
	"net/http"
	"strings"

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

//GetAllProfiles return all of the known profiles
func (s *Server) GetAllProfiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		if profiles, err := db.GetAllProfiles(s.DB); err != nil {
			logging.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, profiles)
		}
	}
}

//GetNetworkCardsInfo return the machine's network cards information
func (s *Server) GetNetworkCardsInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		if devices, err := sniffer.GetLinuxNetworkCardsNameAndMac(); err != nil {
			logging.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.Data(http.StatusOK, "text/html", []byte(strings.Join(devices, ",")))
		}
	}
}

//GetGatewayInfo return the default mac address information
func (s *Server) GetGatewayInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ip, mac, err := sniffer.GetDefualtGetwayMacAddresss(); err != nil {
			logging.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf("%s,%s", ip, mac)))
		}
	}
}
