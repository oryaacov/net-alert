
    ip link | grep ether | grep -oE '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}'
