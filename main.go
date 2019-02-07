package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
)

type Response struct {
	Hostname     string   `json:"hostname"`
	IPAddress    string   `json:"ip_address"`
	MacAddresses []string `json:"mac_addresses"`
}

func getHostname() (name string) {
	name, _ = os.Hostname()
	return name
}

func getIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getMacAddresses() (as []string) {
	ifas, err := net.Interfaces()
	if err != nil {
		return as
	}

	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Hostname: getHostname(), IPAddress: getIPAddress(), MacAddresses: getMacAddresses()})
	})

	http.ListenAndServe(":8765", nil)
}
