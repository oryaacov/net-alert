sudo airmon-ng check kill
sudo airmon-ng start $1 $2
sudo ip link set $3 up    
sudo systemctl start NetworkManager
