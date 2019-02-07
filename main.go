package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
)

type response struct {
	Version      string            `json:"version"`
	Hostname     string            `json:"hostname"`
	IPAddress    string            `json:"ip_address"`
	MacAddresses map[string]string `json:"mac_addresses"`
}

// Version for sshd-config
var Version string

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

func getMacAddresses() map[string]string {
	as := make(map[string]string)
	ifas, err := net.Interfaces()
	if err != nil {
		return as
	}

	for _, ifa := range ifas {
		n := ifa.Name
		a := ifa.HardwareAddr.String()
		if a != "" {
			as[n] = a
		}
	}
	return as
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		resp := response{
			Hostname:     getHostname(),
			IPAddress:    getIPAddress(),
			MacAddresses: getMacAddresses(),
			Version:      Version,
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(resp)
	})

	http.ListenAndServe(":8765", nil)
}
