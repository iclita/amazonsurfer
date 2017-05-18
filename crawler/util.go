package crawler

import (
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// formatLink removes unncessary data from product link
// This is done to make sure unique links are retained
// and we do not have duplicate urls
func formatLink(link string) string {
	s := strings.Split(link, "/")
	// Remove first part which is "" and last parts with ref= and other param
	s = s[1 : len(s)-2]
	// Rebuild the link
	link = strings.Join(s, "/")
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = base
	return u.String()
}

// sleep simply puts the program to sleep for a random number of seconds
// between min and max when looking for products
func sleep(min int, max int) {
	// Seed the random source to get truly random numbers
	rand.Seed(time.Now().UTC().UnixNano())
	// Calculate random delay between requests
	delay := min + rand.Intn(max-min)
	time.Sleep(time.Duration(delay) * time.Second)
}
