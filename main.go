// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	neturl "net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var url, cftkn, email, cfsite, Version string
var dryrun, version bool
var interval int

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&cftkn, "cftkn", "", "Cloudflare API token")
	flag.StringVar(&cfsite, "cfsite", "", "The name of the site to purge in Cloudflare")
	flag.StringVar(&email, "email", "", "Cloudflare account email")
	flag.StringVar(&url, "url", "", "The url to watch for changes")
	flag.BoolVar(&dryrun, "dryrun", false, "Simulates a purging without hitting Cloudflare.")
	flag.BoolVar(&version, "version", false, "Prints version")
	flag.IntVar(&interval, "interval", 15, "The time in seconds to check for changes.")
}

var lastChecksum string

func main() {
	flag.Parse()

	if version {
		log.Println(Version)
		return
	}

	if url == "" || cftkn == "" || email == "" || cfsite == "" {
		flag.Usage()
		return
	}

	c := time.Tick(time.Duration(interval) * time.Second)
	log.Printf("[INFO] Waiting %d seconds...", interval)

	for _ = range c {
		log.Printf("[INFO] Checking for changes in %s", url)
		if res, ok := check(url); ok {
			log.Printf("[INFO] Cloudflare response: %+v", res)
		}
		log.Printf("[INFO] Waiting %d seconds...", interval)
	}
}

type CFResponse struct {
	Result  string `json:"result"`
	Message string `json:"msg"`
}

func purge(url string) (*CFResponse, bool) {
	u, err := neturl.Parse(url)
	if err != nil {
		log.Printf("[ERROR] %#v", err)
		return nil, false
	}

	// Full purge by default unless the URL has a path
	params := neturl.Values{
		"a":     {"fpurge_ts"},
		"tkn":   {cftkn},
		"email": {email},
		"z":     {cfsite},
		"v":     {"1"},
	}

	// Purges only the resource pointed by the url
	if u.Path != "" {
		params.Set("a", "zone_file_purge")
		params.Add("url", url)
	}

	if dryrun {
		log.Printf("[INFO] DryRun successful")
		return &CFResponse{Result: "success"}, true
	}

	resp, err := http.PostForm("https://www.cloudflare.com/api_json.html", params)
	if err != nil {
		log.Printf("[ERROR] %#v", err)
		return nil, false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] %#v", err)
		return nil, false
	}

	cfresp := new(CFResponse)
	err = json.Unmarshal(body, cfresp)
	if err != nil {
		log.Printf("[ERROR] %#v", err)
		log.Printf("[DEBUG] Response: %s", string(body[:]))
		return nil, false
	}

	return cfresp, true
}

func check(url string) (*CFResponse, bool) {
	// We need to get the body of the document because getting the raw content
	// of the response will always be different due to Cloudflare code injections.
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Printf("[ERROR] %#v\n", err)
		return nil, false
	}

	hash := sha512.New()
	hash.Write([]byte(doc.Find("body").Text()))
	md := hash.Sum(nil)
	checksum := hex.EncodeToString(md)
	if lastChecksum == checksum {
		return nil, false
	}

	log.Printf("[INFO] Checksum: %s", checksum)
	lastChecksum = checksum
	log.Printf("[INFO] %s changed, purging Cloudflare cache...", url)
	return purge(url)
}
