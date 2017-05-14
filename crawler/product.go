package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
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
	name = strings.TrimSpace(name)
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
		log.Println("Error parsing price", strPrice)
		return price
	}
	// If the string does not start with $ return price 0
	// This is because all prices on Amazon start with $
	if !strings.HasPrefix(strPrice, "$") {
		log.Println("Error parsing price", strPrice)
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

// findReviews gets the product number of reviews from the parsed document
func findReviews(doc *goquery.Document) uint {
	var reviews uint
	strReviews := doc.Find("#acrCustomerReviewText").Text()
	// If reviews text does not contain 'customer review' then it is something else
	// This also acts for 'customer reviews'
	if !strings.Contains(strReviews, "customer review") {
		log.Println("Error parsing reviews", strReviews)
		return reviews
	}
	// If so, carry on with extracting the number of reviews
	// We will have something like '150 customer reviews'
	rs := strings.Split(strReviews, " ")
	// Get the price from this text
	strReviews = strings.TrimSpace(rs[0])
	// Parse the price number
	numReviews, err := strconv.ParseUint(strReviews, 10, 64)
	if err != nil {
		log.Println("Error parsing reviews", strReviews)
		return reviews
	}

	return uint(numReviews)
}

// findDimensions gets the product dimensions from the parsed document
// It searches into the HTML container for a certain pattern and returns all dimensions
func findDimensions(container string) (float64, float64, float64) {
	// Compute the regex to find the dimensions pattern in the container
	re := regexp.MustCompile("[0-9]+\\.?[0-9]*\\s+x\\s+[0-9]+\\.?[0-9]*\\s+x\\s+[0-9]+\\.?[0-9]*\\s+inches")
	// We return something like '12.3 x 14 x 23 inches'
	strDim := re.FindString(container)
	if strDim == "" {
		log.Println("Error parsing dimensions", strDim)
		return 0, 0, 0
	}
	ds := strings.Split(strDim, "x")
	if len(ds) != 3 {
		log.Println("Error parsing dimensions", strDim)
		return 0, 0, 0
	}
	// Extract all 3 dimensions as strings first
	strLength := strings.TrimSpace(ds[0])
	strWidth := strings.TrimSpace(ds[1])
	strHeight := strings.TrimSpace(strings.Replace(ds[2], "inches", "", -1))

	numLength, err := strconv.ParseFloat(strLength, 64)
	if err != nil {
		log.Println("Error parsing length", strLength)
		return 0, 0, 0
	}

	numWidth, err := strconv.ParseFloat(strWidth, 64)
	if err != nil {
		log.Println("Error parsing width", strWidth)
		return 0, 0, 0
	}

	numHeight, err := strconv.ParseFloat(strHeight, 64)
	if err != nil {
		log.Println("Error parsing height", strHeight)
		return 0, 0, 0
	}

	return numLength, numWidth, numHeight
}

// findWeight gets the product weight from the parsed document
// It searches into the HTML container for a certain pattern and returns the weight
func findWeight(container string) float64 {
	var strWeight string
	re := regexp.MustCompile("[0-9]+\\.?[0-9]*\\s+(ounces|pounds)")
	// We return something like '[23.45 ounces|pounds, 24 ounces|pounds]'
	// The first is the item weight and the second is the shipping weight
	// We are interested in the shipping weight
	sliceWeight := re.FindAllString(container, 2)
	// Get the last item which is the shipping weight if something was found
	if len(sliceWeight) > 0 {
		strWeight = sliceWeight[len(sliceWeight)-1]
	}

	if strWeight == "" {
		log.Println("Error parsing weight", strWeight)
		return 0
	}
	// Split the found string
	ws := strings.Split(strWeight, " ")

	strWeight = strings.TrimSpace(ws[0])

	numWeight, err := strconv.ParseFloat(strWeight, 64)
	if err != nil {
		log.Println("Error parsing weight", strWeight)
		return 0
	}

	return numWeight
}

func findBSR(container string) uint {
	re := regexp.MustCompile("#[0-9]+\\.?[0-9]*\\s+in\\s+.+\\s+")
	// We return something like '#45 in Kitchen (See Top 100 Kitchen)'
	strBSR := re.FindString(container)

	if strBSR == "" {
		log.Println("Error parsing BSR", strBSR)
		return 0
	}
	// Split the string
	bs := strings.Split(strBSR, " ")
	// Extract the first element which is the BSR
	strBSR = strings.TrimSpace(bs[0])
	if !strings.HasPrefix(strBSR, "#") {
		log.Println("Error parsing BSR", strBSR)
		return 0
	}
	// Remove the first character which is #
	strBSR = strBSR[1:]
	// Parse the price number
	numBSR, err := strconv.ParseUint(strBSR, 10, 64)
	if err != nil {
		log.Println("Error parsing BSR", strBSR)
		return 0
	}

	return uint(numBSR)
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
	reviews := findReviews(doc)

	// Get the container from the HTML document
	container := doc.Find("#dp-container").Text()
	// Replace all , with empty space to easily find every number
	container = strings.Replace(container, ",", "", -1)
	// Fetch all 3 dimensions
	length, width, height := findDimensions(container)
	// Fetch product shipping weight
	weight := findWeight(container)
	// Fetch BSR
	bsr := findBSR(container)

	prod := Product{
		Name:    name,
		Link:    link,
		Price:   price,
		BSR:     bsr,
		Reviews: reviews,
		Length:  length,
		Width:   width,
		Height:  height,
		Weight:  weight,
	}

	return prod, nil
}
