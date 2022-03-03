package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/leejones/netrc"
)

type alfredCollection struct {
	Items []alfredItem `json:"items"`
}
type alfredItem struct {
	Arg      string `json:"arg"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Match    string `json:"match"`
	UID      string `json:"uid"`
}

type dashboard struct {
	// {
	// 	"id": 48,
	// 	"uid": "24Xy_QsZz",
	// 	"title": "(C1 2020) Selective Bulk Edit",
	// 	"uri": "db/c1-2020-selective-bulk-edit",
	// 	"url": "/d/24Xy_QsZz/c1-2020-selective-bulk-edit",
	// 	"slug": "",
	// 	"type": "dash-db",
	// 	"tags": [],
	// 	"isStarred": false
	// }
	UID   string `json:"uid"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

func main() {
	grafanaHost := os.Getenv("GRAFANA_HOST")
	apiURL, err := url.Parse(grafanaHost)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	apiURL.Path = path.Join(apiURL.Path, "api/search")

	grafanaUser := os.Getenv("GRAFANA_BASIC_AUTH_USER")
	grafanaPassword := os.Getenv("GRAFANA_BASIC_AUTH_PASSWORD")
	if grafanaUser == "" || grafanaPassword == "" {
		fmt.Fprintf(os.Stderr, "load credentials: ENV vars not set: GRAFANA_BASIC_AUTH_USER, GRAFANA_BASIC_AUTH_PASSWORD\n")
		basicAuth, err := netrc.Get(apiURL.Host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load credentials: unable to load from netrc: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "load credentials: found credentials in netrc\n")
			grafanaUser = basicAuth.Username
			grafanaPassword = basicAuth.Password
		}
	}

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	req.SetBasicAuth(grafanaUser, grafanaPassword)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR: HTTP Response:", resp.StatusCode)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	var dashboards []dashboard
	err = json.Unmarshal(body, &dashboards)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	var items []alfredItem
	for _, dashboard := range dashboards {
		targetURL, err := url.Parse(grafanaHost)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
		targetURL.Path = path.Join(targetURL.Path, dashboard.URL)
		match := strings.ReplaceAll(dashboard.Title, "(", "")
		match = strings.ReplaceAll(match, ")", "")
		match = strings.ReplaceAll(match, "/", "")
		item := alfredItem{
			Arg:      targetURL.String(),
			Match:    match,
			Subtitle: dashboard.Title,
			Title:    dashboard.Title,
			UID:      dashboard.UID,
		}
		items = append(items, item)
	}
	collection := alfredCollection{
		Items: items,
	}
	jsonData, _ := json.Marshal(collection)
	fmt.Println(string(jsonData))
}
