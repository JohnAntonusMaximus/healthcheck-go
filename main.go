package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	sites := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://southwest.com",
		"http://amazon.com",
		"https://httpstat.us/500",
	}

	status := make(chan string)
	ticker := time.NewTicker(120 * time.Second)
	quit := make(chan bool)
	var buffer int

	go func() {
		for {
			select {
			case <-ticker.C:
				buffer = len(sites)
				checkLinks(sites, status, quit)
			case <-quit:
				ticker.Stop()
				os.Exit(1)
			}
		}
	}()

	for {
		fmt.Println(<-status)
	}

}

func checkLinks(sites []string, status chan string, quit chan bool) {
	for _, site := range sites {
		go func(s string) {
			resp, err := http.Get(s)
			if err != nil {
				fmt.Println("Error: ", err)
				quit <- true
			}
			if resp.StatusCode != 200 {
				status <- s + " - DOWN!"
			} else {
				status <- s + " - OK"
			}
		}(site)
	}
}
