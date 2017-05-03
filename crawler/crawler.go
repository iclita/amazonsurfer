package crawler

import (
	"net/http"
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

type Category struct {
	ID            uint8      `json:"id"`
	Name          string     `json:"name"`
	Link          string     `json:"link"`
	SubCategories []Category `json:"subcategories"`
}

type options struct {
	minPrice   int
	maxPrice   int
	minBSR     int
	maxBSR     int
	minReviews int
	maxReviews int
	maxVolume  float64
	maxWeight  float64
}

func (crw *Crawler) MapOptions(r *http.Request) error {

	minPrice, err := strconv.Atoi(r.FormValue("min-price"))
	if err != nil {
		return err
	}

	maxPrice, err := strconv.Atoi(r.FormValue("max-price"))
	if err != nil {
		return err
	}

	minBSR, err := strconv.Atoi(r.FormValue("min-bsr"))
	if err != nil {
		return err
	}

	maxBSR, err := strconv.Atoi(r.FormValue("max-bsr"))
	if err != nil {
		return err
	}

	minReviews, err := strconv.Atoi(r.FormValue("min-reviews"))
	if err != nil {
		return err
	}

	maxReviews, err := strconv.Atoi(r.FormValue("max-reviews"))
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

	crw.opts.minPrice = minPrice
	crw.opts.maxPrice = maxPrice
	crw.opts.minBSR = minBSR
	crw.opts.maxBSR = maxBSR
	crw.opts.minReviews = minReviews
	crw.opts.maxReviews = maxReviews
	crw.opts.maxVolume = maxVolume
	crw.opts.maxWeight = maxWeight

	return nil
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
