package sniffer

import (
	"fmt"
	"log"
	"math"
	"net-alert/pkg/db"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/jinzhu/gorm"
)

//Sniffer contains the pcal handle and his configuration
type Sniffer struct {
	Handler        *pcap.Handle
	Device         string
	Promiscuous    bool
	SnapshotLen    int32
	NewFileRequest bool
	pcapFolder     string
	db             *gorm.DB
	Timeout        time.Duration
}

func init() {
	macRegexExp, _ = regexp.Compile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)

}

var (
	pcapFiles               map[string]time.Time = make(map[string]time.Time)
	lastHandshakePacketTime time.Time
	hsCounter               int8
	macRegexExp             *regexp.Regexp
	packetWriter            *pcapgo.Writer
	packetCounter           int64
	pcapFile                *os.File
)

func generatePcapFile(folder string) *os.File {
	currentTime := time.Now()
	var err error
	filePath := fmt.Sprintf("%s/na%d.pcap", folder, currentTime.Unix())
	log.Printf("creating new pcap file %s", filePath)
	pcapFile, err = os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return pcapFile
}

func createNewPacketWriter(pcapFolder string) *pcapgo.Writer {
	pcapw := pcapgo.NewWriter(generatePcapFile((pcapFolder)))
	if err := pcapw.WriteFileHeader(1600, layers.LinkTypeIEEE80211Radio); err != nil {
		log.Fatalf("WriteFileHeader: %v", err)
	}
	return pcapw
}
func containThreeZeros(bytes []byte) bool {
	return len(bytes) >= 3 && bytes[0] == 0 && bytes[1] == 0 && bytes[2] == 0
}

func isPartOfHandShake(packet gopacket.Packet) bool {
	initPackets := packet.GetInitialLayers()
	return len(initPackets) == 6 && initPackets[5] != nil && initPackets[5].LayerType().String() == "SNAP" &&
		containThreeZeros(initPackets[5].LayerContents())
}

func detectHandshake(packet gopacket.Packet) bool {
	if isPartOfHandShake(packet) {
		if lastHandshakePacketTime.IsZero() || math.Abs(float64(packet.Metadata().Timestamp.Unix()-lastHandshakePacketTime.Unix())) > 1 {
			lastHandshakePacketTime = packet.Metadata().Timestamp
			fmt.Println(packet.Metadata().Timestamp)
			hsCounter = 1
		} else {
			hsCounter++
		}
	}
	return hsCounter == 4
}

func readPacketsFromFile(file string) error {
	if handle, err := pcap.OpenOffline(file); err != nil {
		return err
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			if isCaptured, addresses := handshakeCaptured(packet); isCaptured {
				fmt.Println(addresses)
			}
		}
	}
	return nil
}

func handshakeCaptured(packet gopacket.Packet) (bool, []string) {
	if detectHandshake(packet) {
		hsCounter = 0
		if len(packet.Layers()) > 2 {
			content := strings.Split(fmt.Sprint(packet.Layers()[1]), " ")
			address := make([]string, 0)
			for _, s := range content {
				if macRegexExp.Match([]byte(s)) {
					address = append(address, s)
				}
			}
			return true, address
		}
	}
	return false, nil
}

func getPacketWriter(pcapFolder string, externalRequest bool) *pcapgo.Writer {
	if packetWriter == nil {
		packetWriter = createNewPacketWriter(pcapFolder)
		return packetWriter
	}
	if packetCounter%10000 == 0 || externalRequest {
		packetCounter = 0
		if err := pcapFile.Close(); err != nil {
			fmt.Println(err)
		}
		packetWriter = createNewPacketWriter(pcapFolder)
	}
	return packetWriter
}

func writePacketToFile(packet gopacket.Packet, external bool, pcapFolder string) {
	pcapWriter := getPacketWriter(pcapFolder, external)
	packetCounter++
	if err := pcapWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
		log.Panicf("pcap.WritePacket(): %v", err)
	}
}

//Analyze reciving the raw pcap packets and reading their information
func (sniffer *Sniffer) Analyze(dbi *gorm.DB) string {
	packetSource := gopacket.NewPacketSource(sniffer.Handler, sniffer.Handler.LinkType())
	fmt.Println("creating new packetsouce")
	for packet := range packetSource.Packets() {
		writePacketToFile(packet, sniffer.NewFileRequest, sniffer.pcapFolder)
		if isCaptured, addresses := handshakeCaptured(packet); isCaptured {
			res, _ := db.GetProfileByMac(addresses[1], dbi)
			fmt.Println(addresses[1], res)
		}
		// fmt.Println(packet)
		// if packetCounter%100000 == 0 {
		// 	fmt.Println("creating new file")
		// 	pcapw = createNewPacketWriter(pcapFolder)
		// }
		// fmt.Println(packetCounter)
		// packetCounter++
		// if err := pcapw.WritePacket(packet.Metadata().CaptureInfo, packet.Data()); err != nil {
		// 	log.Fatalf("pcap.WritePacket(): %v", err)
		// }
		// // record := &dm.IPV4Record{Src: packet.NetworkLayer().NetworkFlow().Src().String(), Dst: packet.NetworkLayer().NetworkFlow().Dst().String(), TS: time.Now()}
		// db.Create(record)
	}
	return ""
}
