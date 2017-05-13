package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Product is a representation of an Amazon product
// This contains basic properties needed to represent it
type Product struct {
	Name    string  `json:"name"`
	Link    string  `json:"link"`
	Price   float64 `json:"price"`
	BSR     uint    `json:"bsr"`
	Reviews uint    `json:"reviews"`
	Length  float64 `json:"length"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Weight  float64 `json:"weight"`
}

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

// findName gets the product name from the parsed document
func findName(doc *goquery.Document) string {
	name := doc.Find("#productTitle").Text()
	return name
}

// findPrice gets the product price from the parsed document
func findPrice(doc *goquery.Document) float64 {
	var price float64
	var strPrice string
	// First look for the sale (discounted) price
	strPrice = doc.Find("#priceblock_saleprice").Text()
	if strPrice == "" {
		// If not discounted price is present we take the normal price
		strPrice = doc.Find("#priceblock_ourprice").Text()
	}
	// If no price was found return price 0
	if strPrice == "" {
		return price
	}
	// If the string does not start with $ return price 0
	// This is because all prices on Amazon start with $
	if !strings.HasPrefix(strPrice, "$") {
		return price
	}
	// Check if this a price range ($10.00 - $15.99)
	if strings.Contains(strPrice, "-") {
		ps := strings.Split(strPrice, "-")
		if len(ps) != 2 {
			log.Println("Error parsing price range", strPrice)
			return price
		}
		// Remove $ from the begining of the price and trim space
		lowStrPrice := strings.TrimSpace(ps[0][1:])
		highStrPrice := strings.TrimSpace(ps[1][1:])

		lowPrice, err := strconv.ParseFloat(lowStrPrice, 64)
		if err != nil {
			log.Println("Error parsing low price")
			return price
		}
		highPrice, err := strconv.ParseFloat(highStrPrice, 64)
		if err != nil {
			log.Println("Error parsing high price")
			return price
		}
		price = (lowPrice + highPrice) / 2
	} else {
		// No price range, just s single normal price
		// Remove $ from the begining of the price and trim space
		strPrice = strings.TrimSpace(strPrice[1:])
		// Replace , with empty string to be a valid number
		strPrice = strings.Replace(strPrice, ",", "", -1)
		numPrice, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			log.Println("Error parsing price", strPrice)
			return price
		}
		price = numPrice
	}

	return price
}

// getProduct fetches the product found at the given link
// It attaches all the necessary data to the product type
func getProduct(link string, client *http.Client) (Product, error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)

	if err != nil {
		log.Fatal(err)
	}
	// Send the request
	res, err := client.Do(req)
	if err != nil {
		return Product{}, fmt.Errorf("Request error at url %s", link)
	}
	// Return error if no product was found
	if res.StatusCode != http.StatusOK {
		return Product{}, fmt.Errorf("Product not found at url %s", link)
	}
	// Parse the DOM
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Product{}, fmt.Errorf("Could not parse document at url %s", link)
	}
	defer res.Body.Close()

	// Find product attributes
	name := findName(doc)
	price := findPrice(doc)

	prod := Product{
		Name:  name,
		Link:  link,
		Price: price,
	}

	return prod, nil
}
