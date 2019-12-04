package sniffer

import (
	"log"
	"net-alert/pkg/dm"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
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
	f, err := os.Create("lo.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pcapw := pcapgo.NewWriter(f)
	if err := pcapw.WriteFileHeader(1600, layers.LinkTypeIEEE80211Radio); err != nil {
		log.Fatalf("WriteFileHeader: %v", err)
	}

	packetSource := gopacket.NewPacketSource(sniffer.Handler, sniffer.Handler.LinkType())
	for packet := range packetSource.Packets() {
		if packet != nil && packet.NetworkLayer() != nil {
			if err := pcapw.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
				log.Fatalf("pcap.WritePacket(): %v", err)
			}
			record := &dm.IPV4Record{Src: packet.NetworkLayer().NetworkFlow().Src().String(), Dst: packet.NetworkLayer().NetworkFlow().Dst().String(), TS: time.Now()}
			db.Create(record)
		}
	}
	return ""
}
