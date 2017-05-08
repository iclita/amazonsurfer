package crawler

import (
	"errors"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

// Crawler scrappes Amazon website for products
type Crawler struct {
	opts options
}

// Product is a representation of an Amazon product
// This contains basic properties
// New properties might be added over time
type Product struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// category represent a category/subcategory on Amazon
type category struct {
	id   uint64
	name string
	slug string
	subs []category
}

// options holds parameters necessary to filter products
type options struct {
	categories []category
	minPrice   float64
	maxPrice   float64
	minBSR     uint32
	maxBSR     uint32
	minReviews uint32
	maxReviews uint32
	maxVolume  float64
	maxWeight  float64
}

// MapOptions extracts the request data and maps the input to Crawler options
// This way the crawler knows which options to use when filtering products
func (crw *Crawler) MapOptions(r *http.Request) error {

	minPrice, err := strconv.ParseFloat(r.FormValue("min-price"), 64)
	if err != nil {
		return err
	}

	maxPrice, err := strconv.ParseFloat(r.FormValue("max-price"), 64)
	if err != nil {
		return err
	}

	minBSR, err := strconv.ParseUint(r.FormValue("min-bsr"), 10, 32)
	if err != nil {
		return err
	}

	maxBSR, err := strconv.ParseUint(r.FormValue("max-bsr"), 10, 32)
	if err != nil {
		return err
	}

	minReviews, err := strconv.ParseUint(r.FormValue("min-reviews"), 10, 32)
	if err != nil {
		return err
	}

	maxReviews, err := strconv.ParseUint(r.FormValue("max-reviews"), 10, 32)
	if err != nil {
		return err
	}

	length, err := strconv.ParseFloat(r.FormValue("length"), 64)
	if err != nil {
		return err
	}

	width, err := strconv.ParseFloat(r.FormValue("width"), 64)
	if err != nil {
		return err
	}

	height, err := strconv.ParseFloat(r.FormValue("height"), 64)
	if err != nil {
		return err
	}

	maxVolume := length * width * height

	maxWeight, err := strconv.ParseFloat(r.FormValue("max-weight"), 64)
	if err != nil {
		return err
	}

	cats, err := filterCategories(r.Form["categories"])

	if err != nil {
		return err
	}

	crw.opts.categories = cats

	crw.opts.minPrice = minPrice
	crw.opts.maxPrice = maxPrice
	crw.opts.minBSR = uint32(minBSR)
	crw.opts.maxBSR = uint32(maxBSR)
	crw.opts.minReviews = uint32(minReviews)
	crw.opts.maxReviews = uint32(maxReviews)
	crw.opts.maxVolume = maxVolume
	crw.opts.maxWeight = maxWeight

	return nil
}

// getLinks fetches all the links that need to be scrapped
// All these links belong to a certain category
func (cat *category) getLinks() ([]string, error) {
	links := make([]string, len(cat.subs))
	if cat.subs == nil {
		return nil, errors.New("No subcategories found")
	}
	for _, sub := range cat.subs {
		sid := strconv.Itoa(int(sub.id))
		link := path.Join(base, sub.slug, "zgbs", cat.slug, sid)
		links = append(links, link)
	}
	return links, nil
}

// getLinks delegates work to the function getLinks mentioned above
// The crawler must a have a list of all links waiting to be scrapped
// Every category comes with its links and are all accumulated here
func (crw *Crawler) getLinks() []string {
	var length uint16
	for _, cat := range crw.opts.categories {
		length += uint16(len(cat.subs))
	}
	links := make([]string, length)
	for _, cat := range crw.opts.categories {
		clinks, err := cat.getLinks()
		if err != nil {
			log.Println(err)
			continue
		}
		links = append(links, clinks...)
	}
	return links
}

// Run searches for products and sends them on the channel to be received in the main goroutine
// Now we only send some test products
func (crw *Crawler) Run(prods chan Product) {
	p1 := Product{
		Name: "Cool baby",
		Link: "https://skilldetector.com",
	}
	p2 := Product{
		Name: "Smart Ass",
		Link: "https://www.facebook.com",
	}
	p3 := Product{
		Name: "Top Gun",
		Link: "https://www.golang.org",
	}

	p := []Product{p1, p2, p3}

	for _, v := range p {
		time.Sleep(time.Second)
		prods <- v
		time.Sleep(time.Second)
	}
	close(prods)
}
