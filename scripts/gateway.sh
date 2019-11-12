IP=`ip route | grep default | grep -m 1 -oE "\b([0-9]{1,3}\.){3}[0-9]{1,3}\b"` 
MAC=`arp $IP | grep  -m 1 -oE '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}'`
echo $IP,$MAC
