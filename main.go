package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type response struct {
	CallHomeResponse int    `json:"callHomeResponse"`
	Hostname         string `json:"hostname"`
	IPAddress        string `json:"ip"`
	LatestVersion    string `json:"latestVersion"`
	MacAddress       string `json:"MAC"`
	Version          string `json:"version"`
}

type githubRelease struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string      `json:"tarball_url"`
	ZipballURL string      `json:"zipball_url"`
	Body       interface{} `json:"body"`
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

func getMacAddress() string {
	ifas, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, ifa := range ifas {
		n := ifa.Name
		a := ifa.HardwareAddr.String()
		isEth := strings.HasPrefix(n, "en") || strings.HasPrefix(n, "eth")
		if a != "" && isEth {
			return a
		}
	}
	return ""
}

func getJson(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	logger.Println("Fetching latest release")

	release := &githubRelease{}
	err := getJson("http://api.github.com/repos/josegonzalez/dc2/releases/latest", release)
	latestVersion := "unknown"
	callHomeResponse := 404
	if err == nil {
		latestVersion = strings.TrimPrefix(release.TagName, "v")
		callHomeResponse = 200
	}

	logger.Println("Server is starting")

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		resp := response{
			Hostname:         getHostname(),
			IPAddress:        getIPAddress(),
			MacAddress:       getMacAddress(),
			Version:          Version,
			LatestVersion:    latestVersion,
			CallHomeResponse: callHomeResponse,
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(resp)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8765"
	}

	listenAddr := fmt.Sprintf(":%s", port)
	logger.Println("Server is ready to handle requests at", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}
