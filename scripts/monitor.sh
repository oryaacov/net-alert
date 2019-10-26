sudo airmon-ng check kill
sudo airmon-ng start wlx002127fd59f3
sudo ip link set wlan0mon up    
sudo systemctl start NetworkManager
