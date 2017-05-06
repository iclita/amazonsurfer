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
