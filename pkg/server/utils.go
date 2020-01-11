package server

import (
	"encoding/json"
	"net-alert/pkg/logging"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func readBody(c *gin.Context, req interface{}) bool {
	var err error
	var body []byte
	if body, err = c.GetRawData(); err != nil {
		logging.LogError("failed to parse body", err)
		return false
	}
	if err = json.Unmarshal(body, &req); err != nil {
		logging.LogError("failed to unmarshel body", err)
		return false
	}
	return true
}

func openBrowser(url string) {
	if err := exec.Command("/usr/bin/xdg-open", url).Start(); err != nil {
		logging.LogError(err)
	} else {
		logging.LogInfo(url)
	}
}
