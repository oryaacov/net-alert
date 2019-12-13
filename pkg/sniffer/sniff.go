package sniffer

import (
	"fmt"
	"log"
	"math"
	"net-alert/pkg/dm"
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
	Handler           *pcap.Handle
	Device            string
	Promiscuous       bool
	SnapshotLen       int32
	NewCaptureRequest bool
	PcapFolder        string
	Timeout           time.Duration
}

type State struct {
	HandshakeFound           bool
	HandshakeAddresses       []string
	HandshakedProfile        *dm.Profile
	CurrentPcapFilePath      string
	CurrentPcapFile          *os.File
	CurrentPcapFolder        string
	DecryptedPcapFilePath    string
	DecryptedPcapFile        *os.File
	ExternalRequestToNewFile bool
}

var (
	pcapFiles               map[string]time.Time = make(map[string]time.Time)
	macRegexExp             *regexp.Regexp
	packetWriter            *pcapgo.Writer
	lastHandshakePacketTime time.Time
	hsCounter               int8
	packetCounter           int64
	state                   State
)

func init() {
	state = State{}
	macRegexExp, _ = regexp.Compile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
}

func generatePcapFile(folder string) *os.File {
	currentTime := time.Now()
	f, err := os.Create(fmt.Sprintf("%s/na%d.pcap", folder, currentTime.Unix()))
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func createNewPacketWriter() *pcapgo.Writer {
	pcapw := pcapgo.NewWriter(generatePcapFile((state.CurrentPcapFolder)))
	if err := pcapw.WriteFileHeader(1600, layers.LinkTypeIEEE80211Radio); err != nil {
		log.Fatalf("WriteFileHeader: %v", err)
	}
	return pcapw
}

func getPacketWriter(pcapFolder string, externalRequest bool) *pcapgo.Writer {
	if packetWriter == nil {
		packetWriter = createNewPacketWriter()
		return packetWriter
	}
	if packetCounter%1000000 == 0 || externalRequest {
		packetCounter = 0
		if err := state.CurrentPcapFile.Close(); err != nil {
			fmt.Println(err)
		}
		packetWriter = createNewPacketWriter()
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

func containThreeZeros(bytes []byte) bool {
	return len(bytes) >= 3 && bytes[0] == 0 && bytes[1] == 0 && bytes[2] == 0
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
func isPartOfHandShake(packet gopacket.Packet) bool {
	initPackets := packet.GetInitialLayers()
	return len(initPackets) == 6 && initPackets[5] != nil && initPackets[5].LayerType().String() == "SNAP" &&
		containThreeZeros(initPackets[5].LayerContents())
}

//Analyze reciving the raw pcap packets and reading their information
func (sniffer *Sniffer) Analyze(db *gorm.DB) string {
	var hs bool
	state.CurrentPcapFolder = sniffer.PcapFolder
	packetSource := gopacket.NewPacketSource(sniffer.Handler, sniffer.Handler.LinkType())
	for packet := range packetSource.Packets() {
		writePacketToFile(packet, sniffer.NewCaptureRequest, sniffer.PcapFolder)
		if hs, state.HandshakeAddresses = handshakeCaptured(packet); hs {

		}
		// record := &dm.IPV4Record{Src: packet.NetworkLayer().NetworkFlow().Src().String(), Dst: packet.NetworkLayer().NetworkFlow().Dst().String(), TS: time.Now()}
		// db.Create(record)
	}
	return ""
}
