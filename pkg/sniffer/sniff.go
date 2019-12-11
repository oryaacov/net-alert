package sniffer

import (
	"fmt"
	"log"
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

var (
	pcapFiles map[string]time.Time = make(map[string]time.Time)
)

func generatePcapFile(folder string) *os.File {
	currentTime := time.Now()
	f, err := os.Create(fmt.Sprintf("%s/na%d.pcap", folder, currentTime.Unix()))
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func createNewPacketWriter(pcapFolder string) *pcapgo.Writer {
	pcapw := pcapgo.NewWriter(generatePcapFile((pcapFolder)))
	if err := pcapw.WriteFileHeader(1600, layers.LinkTypeIEEE80211Radio); err != nil {
		log.Fatalf("WriteFileHeader: %v", err)
	}
	return pcapw
}

//Analyze reciving the raw pcap packets and reading their information
func (sniffer *Sniffer) Analyze(db *gorm.DB, pcapFolder string) string {
	var pcapw *pcapgo.Writer
	packetCounter := 0
	packetSource := gopacket.NewPacketSource(sniffer.Handler, sniffer.Handler.LinkType())
	fmt.Println("creating new packetsouce")
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
		if packetCounter%100000 == 0 {
			fmt.Println("creating new file")
			pcapw = createNewPacketWriter(pcapFolder)
		}
		fmt.Println(packetCounter)
		packetCounter++
		if err := pcapw.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
			log.Fatalf("pcap.WritePacket(): %v", err)
		}
		// record := &dm.IPV4Record{Src: packet.NetworkLayer().NetworkFlow().Src().String(), Dst: packet.NetworkLayer().NetworkFlow().Dst().String(), TS: time.Now()}
		// db.Create(record)
	}
	return ""
}
