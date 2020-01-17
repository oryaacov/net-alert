package sniffer

import (
	"fmt"
	"log"
	"math"
	"net-alert/pkg/db"
	"net-alert/pkg/dm"
	"net-alert/pkg/logging"
	"net-alert/pkg/utils"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
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
	db                       *gorm.DB
	HandshakeFound           bool
	HandshakeAddresses       []string
	HandshakedProfile        *dm.Profile
	CurrentPcapFilePath      string
	CurrentPcapFile          *os.File
	CurrentPcapFolder        string
	DecryptedPcapFilePath    string
	DecryptedPcapFile        *os.File
	ExternalRequestToNewFile bool
	ProfileIPv4              string
	ProfileIpv6              string
	BSSID                    string
	SSID                     string
	PSK                      string
	DecryptArgs              []string
	decMux                   sync.Mutex
	Captures                 map[string]*dm.IPV4Record
	master                   *dm.Owner
}

var (
	pcapFiles               map[string]time.Time = make(map[string]time.Time)
	macRegexExp             *regexp.Regexp
	ipv4RegexExp            *regexp.Regexp
	ipv6RegexExp            *regexp.Regexp
	packetWriter            *pcapgo.Writer
	lastHandshakePacketTime time.Time
	hsCounter               int8
	packetCounter           int64
	state                   State
	IsSniffing              bool
)

func init() {
	state = State{}
	macRegexExp, _ = regexp.Compile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	ipv4RegexExp, _ = regexp.Compile(`\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b`)
	ipv6RegexExp, _ = regexp.Compile(`(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`)
	state.Captures = make(map[string]*dm.IPV4Record, 0)
}

func generatePcapFile(folder string) *os.File {
	currentTime := time.Now()
	state.CurrentPcapFilePath = fmt.Sprintf("%s/na%d.pcap", folder, currentTime.Unix())
	f, err := os.Create(state.CurrentPcapFilePath)
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

func IsOnMonitorMode(deviceMonName string) (bool, error) {
	out, err := exec.Command("/sbin/iwconfig").Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(out), deviceMonName), nil
}

func StartMonitorMode(device, channel, deviceMonName string) ([]byte, error) {
	out, err := exec.Command("/bin/bash", "./scripts/monitor.sh", device, channel, deviceMonName).Output()
	if err != nil {
		return nil, err
	}
	return out, nil
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

func getProfileByHandshakedMac() {
	profile, err := db.GetProfileByMac(state.HandshakeAddresses[1], state.db)
	if err != nil {
		log.Panic(err, state.HandshakeAddresses[1])
	}
	if profile == nil {
		log.Panic("unknown profile handshake captured")
	}
	state.HandshakedProfile = profile
	log.Println(state.HandshakedProfile, "captured! sniff sniff!")
	for range time.Tick(time.Second * 5) {
		go decryptPackets()
	}
}

func decryptPackets() {
	args := strings.Split(fmt.Sprintf("-b %s -e %s -p %s %s", state.BSSID, state.SSID, state.PSK, state.CurrentPcapFilePath), " ")
	state.decMux.Lock()
	if _, err := exec.Command("/usr/bin/airdecap-ng", args...).Output(); err != nil {
		state.decMux.Unlock()
		log.Panic(err)
	} else {
		state.DecryptedPcapFilePath = fmt.Sprintf("%s-dec.pcap", state.CurrentPcapFilePath[:len(state.CurrentPcapFilePath)-5])
		readPacketsFromFile()
	}
}

func readPacketsFromFile() {
	if handle, err := pcap.OpenOffline(state.DecryptedPcapFilePath); err != nil {
		if handle != nil {
			handle.Close()
		}
		state.decMux.Unlock()
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		defer handle.Close()
		for packet := range packetSource.Packets() {
			if packet != nil && packet.NetworkLayer() != nil {
				src := packet.LinkLayer().LinkFlow().Src().String()
				dst := packet.LinkLayer().LinkFlow().Dst().String()
				if state.HandshakedProfile.Mac == dst || state.HandshakedProfile.Mac == src {
					if state.ProfileIPv4 == "" && state.ProfileIpv6 == "" {
						if state.HandshakedProfile.Mac == src {
							ip := packet.NetworkLayer().NetworkFlow().Src().String()
							if isIPV4(ip) {
								state.ProfileIPv4 = ip
							} else {
								state.ProfileIpv6 = ip
							}
						} else {
							ip := packet.NetworkLayer().NetworkFlow().Dst().String()
							if isIPV4(ip) {
								state.ProfileIPv4 = ip
							} else {
								state.ProfileIpv6 = ip
							}
						}
					}
					record := &dm.IPV4Record{Src: packet.NetworkLayer().NetworkFlow().Src().String(), Dst: packet.NetworkLayer().NetworkFlow().Dst().String(), TS: time.Now()}
					if !(record.Dst == state.ProfileIPv4 || strings.ToLower(record.Dst) == state.ProfileIpv6) {
						if _, exists := state.Captures[record.Dst]; !exists {
							handleNewRecord(record)
						}
					}
				} else {
					//	fmt.Println(src, dst, state.HandshakedProfile.Mac)
				}
			}
		}
		handle.Close()
		state.decMux.Unlock()
	}
}

func isIPV4(ip string) bool {
	return ipv4RegexExp.Match([]byte(ip))
}

func isIPV6(ip string) bool {
	return ipv6RegexExp.Match([]byte(ip))
}

func handleNewRecord(record *dm.IPV4Record) {
	state.Captures[record.Dst] = record
	state.db.Create(record)

	for _, site := range state.HandshakedProfile.Sites {
		siteIP := strings.ToLower(site.IP)
		if (siteIP == strings.ToLower(record.Dst) || siteIP == strings.ToLower(record.Src)) &&
			(siteIP != state.ProfileIPv4 && siteIP != strings.ToLower(state.ProfileIpv6)) {
			alertMaster(siteIP)
		}
	}
}

func alertMaster(siteIP string) {
	msg := fmt.Sprintf("ALERT! Profile %s Surfed to %s!", state.HandshakedProfile.NickName, siteIP)
	for i := 0; i < 5; i++ {
		fmt.Printf(msg)
	}
	if state.master.GetEmailAlerts {
		utils.SendEmail(state.master.Email, "NET ALERT!", msg)
	}
	if state.master.GetSMSAlerts {
		utils.SendSMS(state.master.Phone, msg)
	}
}

func Exit(device string) {
	mon, err := IsOnMonitorMode(device)
	if err != nil {
		logging.LogError(err)
		return
	}
	if mon {
		out, err := exec.Command("/bin/bash", "./scripts/close-monitor.sh", device).Output()
		if err != nil {
			logging.LogError(err)
		}
		logging.LogInfo(string(out))
	} else {
		logging.LogInfo("out of monitor mode")
	}
}

func handleHandshakeCapture() {
	if state.HandshakeAddresses != nil && len(state.HandshakeAddresses) == 3 {
		_, bssid, _ := GetCurrentSSIDAndBSSID()
		if strings.ToUpper(bssid) == strings.ToUpper(state.HandshakeAddresses[0]) {
			go getProfileByHandshakedMac()
		} else {
			log.Println("bssid:" + bssid)
			log.Println(state.HandshakeAddresses)
			log.Fatal("handshakes from another networks are not allowed.")
		}
	} else {
		log.Print("invalid handshake addresses pattern")
	}
}

//Analyze reciving the raw pcap packets and reading their information
func (sniffer *Sniffer) Analyze(dbi *gorm.DB) string {
	var hs bool
	var err error
	if state.SSID, state.BSSID, err = GetCurrentSSIDAndBSSID(); err != nil {
		log.Panic(err)
	}
	if state.PSK, err = getLinuxNetworkPassword(); err != nil {
		log.Panic(err)
	}
	state.db = dbi
	if state.master, err = db.GetOwner(state.db); err != nil {
		log.Panic(err)
	}
	IsSniffing = true
	state.CurrentPcapFolder = sniffer.PcapFolder
	packetSource := gopacket.NewPacketSource(sniffer.Handler, sniffer.Handler.LinkType())
	for packet := range packetSource.Packets() {
		writePacketToFile(packet, sniffer.NewCaptureRequest, sniffer.PcapFolder)
		if !state.HandshakeFound {
			if hs, state.HandshakeAddresses = handshakeCaptured(packet); hs {
				state.HandshakeFound = true
				handleHandshakeCapture()
			}
		}
	}
	return ""
}
