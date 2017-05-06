package crawler

import (
	"errors"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
)

type Crawler struct {
	opts  options
	Prods chan Product
}

type Product struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type category struct {
	id   uint64
	name string
	slug string
	subs []category
}

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
func (crw *Crawler) Run() {
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
		crw.Prods <- v
		time.Sleep(time.Second)
	}
	close(crw.Prods)
}
