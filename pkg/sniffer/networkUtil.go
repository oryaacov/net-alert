package sniffer

import (
	"errors"
	"fmt"
	"net-alert/pkg/dm"
	"os/exec"
	"runtime"
	"strings"
)

var (
	_getwayIP  = ""
	_getwayMAC = ""
)

//GetNetworkInfo gather all of the netword information using the .sh scripts
func GetNetworkInfo() (*dm.NetworkInfo, error) {
	var err error
	var result dm.NetworkInfo
	if result.SSID, result.BSSID, err = GetCurrentSSIDAndBSSID(); err != nil {
		return nil, err
	}
	if result.GatewayIP, result.GatewayMAC, err = GetDefualtGetwayMacAddresss(); err != nil {
		return nil, err
	}
	if result.NetworkCards, err = GetLinuxNetworkCardsNameAndMac(); err != nil {
		return nil, err
	}
	return &result, nil
}

//GetCurrentSSIDAndBSSID returns the current network ssid,bssid using nmcli
func GetCurrentSSIDAndBSSID() (string, string, error) {
	switch runtime.GOOS {
	case "linux":
		ssid, bssid, err := getLinuxBSSID()
		if err != nil {
			return "", "", err
		}
		return ssid, bssid, nil
	default:
		return "", "", fmt.Errorf("unhandeled OS %s", runtime.GOOS)
	}

}

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
func GetLinuxNetworkCardsNameAndMac() ([]dm.NetworkCardInfo, error) {
	var err error
	var out []byte
	if out, err = exec.Command("/sbin/ip", "link").Output(); err != nil {
		return nil, err
	}
	strs := strings.Split(string(out), "\n")
	if strs == nil || len(strs) == 0 {
		return nil, errors.New("no result")
	}
	results := make([]dm.NetworkCardInfo, 0)
	var res dm.NetworkCardInfo
	for i, s := range strs {
		if i%2 == 0 {
			if splittedString := strings.Split(s, ":"); len(splittedString) < 2 {
				continue
			} else {
				res.Name = splittedString[1]
			}
		} else {
			res.Mac = strings.Split(strings.TrimSpace(s), " ")[1]
			if strings.Contains(s, "ether") {
				results = append(results, res)
				res = dm.NetworkCardInfo{}
			}
		}
	}
	return results, nil
}

func getLinuxBSSID() (string, string, error) {
	out, err := exec.Command("/bin/bash", "/home/brain/Projects/src/net-alert/scripts/bssid.sh").Output()
	if err != nil {
		panic(err)
	}
	res := strings.Split(string(out), " ")
	if res == nil || len(res) == 0 {
		return "", "", errors.New("failed to obtain gatway ip & mac address")
	}
	var ssid, bssid string
	for _, s := range res {
		if len(s) > 0 {
			if len(ssid) == 0 {
				ssid = s
			} else if len(bssid) == 0 {
				bssid = s
			} else {
				break
			}
		}
	}
	return ssid, bssid, nil
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
	_getwayMAC = strings.Trim(res[1], " \n\t")
	return res[0], res[1], nil
}
