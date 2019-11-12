package sniffer

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var (
	_getwayIP  = ""
	_getwayMAC = ""
)

//GetDefualtGetwayMacAddresss return the def
func GetDefualtGetwayMacAddresss() (string, string, error) {
	if len(_getwayIP) > 0 && len(_getwayMAC) > 0 {
		return _getwayIP, _getwayMAC, nil
	}
	switch runtime.GOOS {
	case "linux":
		ip, mac, err := getLinuxGatewayIP()
		if err != nil {
			return "", "", err
		}
		return ip, mac, nil
	default:
		return "", "", fmt.Errorf("unhandeled OS %s", runtime.GOOS)
	}

}

//GetLinuxNetworkCardsNameAndMac return a list on networks cads with mac address
func GetLinuxNetworkCardsNameAndMac() ([]string, error) {
	var err error
	var out []byte
	var res string
	if out, err = exec.Command("/sbin/ip", "link").Output(); err != nil {
		return nil, err
	}
	strs := strings.Split(string(out), "\n")
	if strs == nil || len(strs) == 0 {
		return nil, errors.New("no result")
	}
	results := make([]string, 0)
	for i, s := range strs {
		if i%2 == 0 {
			if splittedString := strings.Split(s, ":"); len(splittedString) < 2 {
				continue
			} else {
				res = splittedString[1] + "-"
			}
		} else {
			res = res + strings.Split(strings.TrimSpace(s), " ")[1]
			if strings.Contains(s, "ether") {
				results = append(results, res)
			}
		}
	}
	return results, nil
}

func getLinuxGatewayIP() (string, string, error) {
	out, err := exec.Command("/bin/bash", "/home/brain/Projects/src/net-alert/scripts/gateway.sh").Output()
	if err != nil {
		panic(err)
	}
	res := strings.Split(string(out), ",")
	if res == nil || len(res) == 0 {
		return "", "", errors.New("failed to obtain gatway ip & mac address")
	}
	_getwayIP = res[0]
	_getwayMAC = res[1]
	return res[0], res[1], nil
}
