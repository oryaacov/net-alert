package server

import (
	"fmt"
	"net-alert/pkg/db"
	"net-alert/pkg/dm"
	"net-alert/pkg/logging"
	"net-alert/pkg/sniffer"
	"net/http"

	"github.com/gin-gonic/gin"
)

//IsAlive returns "yes" if the server is up and running, part of DT standart
func (s *Server) IsAlive() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte(""))
	}

}

//IsAlive returns "yes" if the server is up and running, part of DT standart
func (s *Server) GetOwnerInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		if owner, err := db.GetOwner(s.DB); err != nil {
			logging.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, owner)
		}
	}
}

//UpdateOwner update the master info
func (s *Server) UpdateOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		var owner *dm.Owner
		if !readBody(c, &owner) {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		if err := ((owner).Update(s.DB)); err != nil {
			logging.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "text/html", []byte(""))
	}
}

//CreateOrUpdateProfile create or update mac profile
func (s *Server) CreateOrUpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		profiles := make([]dm.Profile, 0)
		if !readBody(c, &profiles) {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		for _, profile := range profiles {
			if err := ((&profile).CreateOrUpdate(s.DB)); err != nil {
				logging.LogError(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
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

//GetNetworkInfo return the machine's network cards, Service set and gateway information
func (s *Server) Sniff() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		res := make([]byte, 0)
		onMonitorMode, err := sniffer.IsOnMonitorMode(s.Config.Sniffer.DeviceMonName)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to determain if device is on monitor mode"))
			return
		}
		if !onMonitorMode {
			ch, err := sniffer.GetCurrentChannel(s.Config.Sniffer.DeviceName)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get current channel"))
				return
			}
			res, err = sniffer.StartMonitorMode(s.Config.Sniffer.DeviceName, ch, s.Config.Sniffer.DeviceMonName)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to start monitor mode"))
				return
			}
			logging.LogInfo(string(res))
		}
		if !sniffer.IsSniffing {
			res = append(res, []byte("started!")...)
			fmt.Println("init sniffer & opening pcap...")
			s.InitSniffer()
			go s.Sniffer.Analyze(s.DB)
			c.Data(http.StatusOK, "application/text", res)
		} else {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("net alert is already sniffing"))
		}
	}
}

//GetNetworkInfo return the machine's network cards, Service set and gateway information
func (s *Server) GetNetworkInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		if info, err := sniffer.GetNetworkInfo(s.Config.Sniffer.DeviceName); err != nil {
			logging.LogError(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, info)
		}
	}
}

// //GetGatewayInfo return the default mac address information
// func (s *Server) GetGatewayInfo() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if ip, mac, err := sniffer.GetDefualtGetwayMacAddresss(); err != nil {
// 			logging.LogError(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 		} else {
// 			c.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf("%s,%s", ip, mac)))
// 		}
// 	}
// }

// //GetServiceSetInfo returns the current network ssid and bssid
// func (s *Server) GetServiceSetInfo() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if ssid, bssid, err := sniffer.GetCurrentSSIDAndBSSID(); err != nil {
// 			logging.LogError(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 		} else {
// 			c.Data(http.StatusOK, "text/html", []byte(fmt.Sprintf("%s,%s", ssid, bssid)))
// 		}
// 	}
// }
