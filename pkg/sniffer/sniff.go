package sniffer

import (
	"net-alert/pkg/dm"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/jinzhu/gorm"
)

//Sniffer contains the pcal handle and his configuration
type Sniffer struct {
	Handler     *pcap.Handle
	Device      string
	Promiscuous bool
	SnapshotLen int32
	Timeout     time.Duration
}

//Analyze reciving the raw pcap packets and reading their information
func (sniffer *Sniffer) Analyze(db *gorm.DB) string {
	packetSource := gopacket.NewPacketSource(sniffer.Handler, sniffer.Handler.LinkType())
	for packet := range packetSource.Packets() {
		if packet != nil && packet.NetworkLayer() != nil {
			record := &dm.IPV4Record{Src: packet.NetworkLayer().NetworkFlow().Src().String(), Dst: packet.NetworkLayer().NetworkFlow().Dst().String(), TS: time.Now()}
			db.Create(record)
		}
	}
	return ""
}
