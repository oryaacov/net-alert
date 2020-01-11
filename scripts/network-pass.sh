SSID=`nmcli -f SSID,BSSID,ACTIVE dev wifi list | grep yes| cut -d ' ' -f1`
sudo cat /etc/NetworkManager/system-connections/'Auto '$SSID | grep psk= | cut -d '=' -f2
