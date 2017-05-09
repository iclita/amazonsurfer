package crawler

import (
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// sleep simply puts the program to sleep for a random number of seconds
// between min and max when looking for products
func sleep(min int, max int) {
	delay := min + rand.Intn(max-min)
	time.Sleep(time.Duration(delay) * time.Second)
}

// formatLink removes unncessary data from product link
// This is done to make sure unique links are retained
// and we do not have duplicate urls
func formatLink(link string) string {
	s := strings.Split(link, "/")
	// Remove first part which is "" and last part with ref=
	s = s[1 : len(s)-1]
	link = strings.Join(s, "/")
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "https"
	u.Host = base
	return u.String()
}
