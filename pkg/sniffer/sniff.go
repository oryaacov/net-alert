package sniffer

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/jinzhu/gorm"
)

type IPV4Record struct {
	ID  int64 `gorm:"PRIMARY_KEY`
	Src string
	Dst string
	TS  time.Time
}
type Sniffer struct {
	Handler     *pcap.Handle
	Device      string
	Promiscuous bool
	SnapshotLen int32
	Timeout     time.Duration
}

//Start a
func (sniffer *Sniffer) Start(db *gorm.DB) string {
	packetSource := gopacket.NewPacketSource(sniffer.Handler, sniffer.Handler.LinkType())
	for packet := range packetSource.Packets() {
		if packet != nil && packet.NetworkLayer() != nil {
			record := &IPV4Record{Src: packet.NetworkLayer().NetworkFlow().Src().String(), Dst: packet.NetworkLayer().NetworkFlow().Dst().String(), TS: time.Now()}
			db.Create(record)
		}
	}
	return ""
}
