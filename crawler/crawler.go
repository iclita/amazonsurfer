package crawler

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"strings"

	"github.com/PuerkitoBio/goquery"
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

const (
	minSleep = 15
	maxSleep = 30
)

var wg sync.WaitGroup

var client = &http.Client{
	Timeout: 10 * time.Second,
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
	for i, sub := range cat.subs {
		sid := strconv.Itoa(int(sub.id))
		link := path.Join(base, sub.slug, "zgbs", cat.slug, sid)
		u, err := url.Parse(link)
		if err != nil {
			log.Fatal(err)
		}
		u.Scheme = "https"
		link = u.String()
		links[i] = link
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
		for i, cl := range clinks {
			links[i] = cl
		}
	}
	return links
}

// sleep simply puts the program to sleep for a random number of seconds
// between min and max
func sleep(min int, max int) {
	delay := min + rand.Intn(max-min)
	time.Sleep(time.Duration(delay) * time.Second)
}

// formatLink removes unncessary data from link
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

// scrape extracts all product links from a certain category
// When it finds suitable products it sends them through the prods channel
// and the main goroutine sends them in the frontend
func (crw *Crawler) scrape(link string, prods chan<- Product) {
	defer wg.Done()

	page := 1

	for {
		pg := strconv.Itoa(page)
		q := url.Values{}
		q.Set("_encoding", "UTF8")
		q.Set("pg", pg)
		q.Set("ajax", "1")
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
		if res.StatusCode == http.StatusNotFound {
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

	for _, link := range links {
		go crw.scrape(link, prods)
		sleep(minSleep, maxSleep)
	}

	wg.Wait()
	close(prods)
}
