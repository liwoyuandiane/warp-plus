package wiresocks

import (
	"crypto/rand"
	"fmt"
	"github.com/bepass-org/ipscanner"
	"github.com/go-ini/ini"
	"log"
	"net"
	"time"
)

func canConnectIPv6(remoteAddr string) bool {
	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.Dial("tcp6", remoteAddr)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func RunScan() (result []string) {
	cfg, err := ini.Load("./primary/wgcf-profile.ini")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Reading the private key from the 'Interface' section
	privateKey := cfg.Section("Interface").Key("PrivateKey").String()

	// Reading the public key from the 'Peer' section
	publicKey := cfg.Section("Peer").Key("PublicKey").String()

	// new scanner
	scanner := ipscanner.NewScanner(
		ipscanner.WithWarpPing(),
		ipscanner.WithWarpPrivateKey(privateKey),
		ipscanner.WithWarpPeerPublicKey(publicKey),
		ipscanner.WithUseIPv6(canConnectIPv6("[2001:4860:4860::8888]:80")),
		ipscanner.WithUseIPv4(true),
		ipscanner.WithMaxDesirableRTT(500),
		ipscanner.WithCidrList([]string{
			"162.159.192.0/24",
			"162.159.193.0/24",
			"162.159.195.0/24",
			"188.114.96.0/24",
			"188.114.97.0/24",
			"188.114.98.0/24",
			"188.114.99.0/24",
			"2606:4700:d0::/48",
			"2606:4700:d1::/48",
		}),
	)
	scanner.Run()
	var ipList []net.IP
	for {
		ipList = scanner.GetAvailableIPS()
		if len(ipList) > 1 {
			scanner.Stop()
			break
		}
		time.Sleep(1 * time.Second)
	}
	for i := 0; i < 2; i++ {
		result = append(result, ipToAddress(ipList[i]))
	}
	return
}

func ipToAddress(ip net.IP) string {
	ports := []int{500, 854, 859, 864, 878, 880, 890, 891, 894, 903, 908, 928, 934, 939, 942,
		943, 945, 946, 955, 968, 987, 988, 1002, 1010, 1014, 1018, 1070, 1074, 1180, 1387, 1701,
		1843, 2371, 2408, 2506, 3138, 3476, 3581, 3854, 4177, 4198, 4233, 4500, 5279,
		5956, 7103, 7152, 7156, 7281, 7559, 8319, 8742, 8854, 8886}

	// Pick a random port number
	b := make([]byte, 8)
	n, err := rand.Read(b)
	if n != 8 {
		panic(n)
	} else if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s:%d", ip.String(), ports[int(b[0])%len(ports)])
}
