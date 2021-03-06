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
	_ssid      = ""
	_bssid     = ""
	_psk       = ""
	_channel   = ""
)

//GetNetworkInfo gather all of the netword information using the .sh scripts
func GetNetworkInfo(device string) (*dm.NetworkInfo, error) {
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
	if result.NetworkPass, err = getLinuxNetworkPassword(); err != nil {
		return nil, err
	}
	if result.Channel, err = GetCurrentChannel(device); err != nil {
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

func GetCurrentChannel(device string) (string, error) {
	if len(_channel) > 0 {
		return _channel, nil
	}
	out, err := exec.Command("/sbin/iwlist", "channel").Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	output := string(out)
	i := strings.Index(string(output), device)
	if i > -1 {
		output = output[i:]
		if strings.Contains(output, "(Channel") {
			_channel = output[strings.Index(output, "(Channel")+8 : strings.Index(output, ")")]
			return _channel, nil
		}

	}
	return "", nil
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
	if len(_bssid) > 0 && len(_ssid) > 0 {
		return _ssid, _bssid, nil
	}
	out, err := exec.Command("/bin/bash", "./scripts/bssid.sh").Output()
	if err != nil {
		panic(err)
	}
	res := strings.Split(string(out), " ")
	if res == nil || len(res) == 0 {
		return "", "", errors.New("failed to obtain gatway ip & mac address")
	}
	for _, s := range res {
		if len(s) > 0 {
			if len(_ssid) == 0 {
				_ssid = s
			} else if len(_bssid) == 0 {
				_bssid = s
			} else {
				break
			}
		}
	}
	return _ssid, _bssid, nil
}

func getLinuxNetworkPassword() (string, error) {
	if len(_psk) > 0 {
		return _psk, nil
	}
	out, err := exec.Command("/bin/bash", "./scripts/network-pass.sh").Output()
	if err != nil {
		return "", err
	}
	_psk = strings.TrimSpace(string(out))
	return _psk, nil
}

func getLinuxGatewayIP() (string, string, error) {
	out, err := exec.Command("/bin/bash", "./scripts/gateway.sh").Output()
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
