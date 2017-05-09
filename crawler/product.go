package crawler

import (
	"log"
	"net/url"
	"strings"
)

// Product is a representation of an Amazon product
// This contains basic properties
// New properties might be added over time
type Product struct {
	Name string `json:"name"`
	Link string `json:"link"`
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
