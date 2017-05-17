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
	"github.com/gorilla/websocket"
)

// Crawler scrapes Amazon website for products
type Crawler struct {
	opts    options
	Done    chan struct{}
	conn    *websocket.Conn
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
	tolerance  float64
}

// The min and max limits to sleep between requests
// The crawler will choose a random number of seconds to sleep
const (
	minSleep = 10
	maxSleep = 60
)

// wg waits for all goroutines to finish
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

	tolerance, err := strconv.ParseFloat(r.FormValue("tolerance"), 64)
	if err != nil {
		return err
	}

	cats, err := filterCategories(r.Form["categories"])
	if err != nil {
		return err
	}
	// We save these options on the crawler
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
	crw.opts.tolerance = tolerance

	return nil
}

// getLinks delegates work to the function getLinks mentioned above
// The crawler must a have a list of all links waiting to be scrapped
// Every category comes with its links and are all accumulated here
func (crw *Crawler) getLinks() []string {
	// Calculate the total length of the links slice
	// This way it is very efficient because we make 1 allocation only
	var length uint16
	for _, cat := range crw.opts.categories {
		length += uint16(len(cat.subs))
	}
	links := make([]string, length)
	// We extract all links from every category and merge them in the final slice
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
	// Start from first page
	page := 1
	// Loop through all subcategory pages
	for {
		// Compute the link to be scraped for products
		pg := strconv.Itoa(page)
		q := url.Values{}
		q.Set("_encoding", "UTF8")
		q.Set("pg", pg)
		plink := link + "?" + q.Encode()

		req, err := http.NewRequest(http.MethodGet, plink, nil)

		if err != nil {
			log.Fatal(err)
		}
		// Send the request
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		// Exit this goroutine when there are no more pages to scrape
		if res.StatusCode != http.StatusOK {
			return
		}
		// Parse the DOM
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()
		// Hold the product links in a set like structure
		// This way we make sure that no duplicate links are inserted
		prodLinks := make(map[string]bool)
		// Find the product links
		sel := doc.Find("div.zg_itemWrapper > div")
		for i := range sel.Nodes {
			select {
			case <-crw.Done:
				return
			default:
				// Sleep between requests
				sleep(minSleep, maxSleep)
				// For each item found, get the url
				link, ok := sel.Eq(i).Find("a").Attr("href")
				if !ok {
					log.Println(err)
				} else {
					link = formatLink(link)
					if !prodLinks[link] {
						prodLinks[link] = true
						p, err := getProduct(link, client, crw.Done)
						if err != nil {
							log.Println(err)
						} else {
							// If product is valid send it
							if p.isValid(crw.opts) {
								prods <- p
							}
						}
					}
				}
			}
		}
		// Go to the next page
		page++
		// Sleep between requests
		sleep(minSleep, maxSleep)
	}
}

// Run searches for products and sends them on the channel to be received in the main goroutine
// It sends valid products in the frontend through the websocket connection
func (crw *Crawler) Run(conn *websocket.Conn, prods chan Product) {
	// Seed the random source to get truly random numbers
	rand.Seed(time.Now().UTC().UnixNano())
	// Hold a reference to the current connection
	crw.conn = conn
	// Reset Done channel to initial state so that calls to this channel block again
	crw.Done = nil
	// Get all the links that need to be scraped
	links := crw.getLinks()
	// Add all goroutines to the wait group
	wg.Add(len(links))
	// It is best not to use the default client which has no timeout
	// This way no request takes more then the the specified timeout
	// And the resources are not stuck
	httpClient := &http.Client{
		Timeout: crw.Timeout * time.Second,
	}
	// Scrape every subcateogry in its own goroutine
	for _, link := range links {
		go crw.scrape(link, prods, httpClient)
		sleep(minSleep, maxSleep)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// We're done. Close the channel
	close(prods)
}

// Stop signals all other goroutines to exit
// It then closes the current websockets connection
func (crw *Crawler) Stop() {
	// Make a new Done channel and close it to signal child goroutines to stop
	crw.Done = make(chan struct{})
	close(crw.Done)
	// Close the connection if it's still opened
	// The connection maybe already closed by the browser
	// So a check must be made to ensure the connection still exists
	if crw.conn != nil {
		crw.conn.Close()
	}
}
