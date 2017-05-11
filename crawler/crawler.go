package crawler

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Crawler scrapes Amazon website for products
type Crawler struct {
	opts    options
	Timeout time.Duration
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
	maxLength  float64
	maxWidth   float64
	maxHeight  float64
	maxWeight  float64
}

const (
	minSleep = 15
	maxSleep = 30
)

var wg sync.WaitGroup

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

	maxLength, err := strconv.ParseFloat(r.FormValue("max-length"), 64)
	if err != nil {
		return err
	}

	maxWidth, err := strconv.ParseFloat(r.FormValue("max-width"), 64)
	if err != nil {
		return err
	}

	maxHeight, err := strconv.ParseFloat(r.FormValue("max-height"), 64)
	if err != nil {
		return err
	}

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
	crw.opts.maxLength = maxLength
	crw.opts.maxWidth = maxWidth
	crw.opts.maxHeight = maxHeight
	crw.opts.maxWeight = maxWeight

	return nil
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
		for i, cl := range clinks {
			links[i] = cl
		}
	}
	return links
}

// scrape extracts all product links from a certain category
// When it finds suitable products it sends them through the prods channel
// and the main goroutine sends them in the frontend
func (crw *Crawler) scrape(link string, prods chan<- Product, client *http.Client) {
	defer wg.Done()

	page := 1

	for {
		pg := strconv.Itoa(page)
		q := url.Values{}
		q.Set("_encoding", "UTF8")
		q.Set("ajax", "1")
		q.Set("pg", pg)
		plink := link + "?" + q.Encode()

		req, err := http.NewRequest(http.MethodGet, plink, nil)

		if err != nil {
			log.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		// Exit this goroutine when there are no more pages to scrape
		if res.StatusCode != http.StatusOK {
			return
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println(err)
		}
		res.Body.Close()
		prodLinks := make(map[string]bool)
		var products []Product
		// Find the product links
		doc.Find("div.zg_itemWrapper > div").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the url and name
			url, ok := s.Find("a").Attr("href")
			if !ok {
				log.Println(err)
			}
			url = formatLink(url)
			name := s.Find("a").Text()
			if !prodLinks[url] {
				prodLinks[url] = true
				p := Product{
					name,
					url,
				}
				products = append(products, p)
			}
		})

		for _, p := range products {
			prods <- p
		}

		page++
		sleep(minSleep, maxSleep)
	}

}

// Run searches for products and sends them on the channel to be received in the main goroutine
// It sends valid products in the frontend through the websocket connection
func (crw *Crawler) Run(prods chan Product) {

	rand.Seed(time.Now().UTC().UnixNano())
	links := crw.getLinks()
	wg.Add(len(links))

	// It is best not to use the default client which has no timeout
	// This way no request takes more then the the specified timeout
	// And the resources are not stuck
	httpClient := &http.Client{
		Timeout: crw.Timeout * time.Second,
	}

	for _, link := range links {
		go crw.scrape(link, prods, httpClient)
		sleep(minSleep, maxSleep)
	}

	wg.Wait()
	close(prods)
}
