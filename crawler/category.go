package crawler

import (
	"errors"
	"log"
	"net/url"
	"path"
	"strconv"
)

// category represent a category/subcategory on Amazon
type category struct {
	id   uint64
	name string
	slug string
	subs []category
}

// Base URL for Amazon
const base = "www.amazon.com"

// Below we define all main categories mentioned on Amazon
var appliances = category{
	id:   1,
	name: "Appliances",
	slug: "appliances",
	subs: []category{
		category{
			3741261,
			"Cooktops",
			"Best-Sellers-Appliances-Cooktops",
			nil,
		},
	},
}

var appsGames = category{
	id:   2,
	name: "Apps & Games",
	slug: "mobile-apps",
	subs: nil,
}

var artsCraftsSewing = category{
	id:   3,
	name: "Arts, Crafts & Sewing",
	slug: "arts-crafts",
	subs: nil,
}

var automotive = category{
	id:   4,
	name: "Automotive",
	slug: "automotive",
	subs: nil,
}

var baby = category{
	id:   5,
	name: "Baby",
	slug: "baby-products",
	subs: nil,
}

var beautyPersonalCare = category{
	id:   6,
	name: "Beauty & Personal Care",
	slug: "beauty",
	subs: nil,
}

var books = category{
	id:   7,
	name: "Books",
	slug: "books",
	subs: nil,
}

var cDsVinyl = category{
	id:   8,
	name: "CDs & Vinyl",
	slug: "music",
	subs: nil,
}

var cameraPhoto = category{
	id:   9,
	name: "Camera & Photo",
	slug: "photo",
	subs: nil,
}

var cellPhonesAccessories = category{
	id:   10,
	name: "Cell Phones & Accessories",
	slug: "wireless",
	subs: nil,
}

var clothingShoesJewelry = category{
	id:   11,
	name: "Clothing, Shoes & Jewelry",
	slug: "fashion",
	subs: nil,
}

var collectibleCoins = category{
	id:   12,
	name: "Collectible Coins",
	slug: "coins",
	subs: nil,
}

var computersAccessories = category{
	id:   13,
	name: "Computers & Accessories",
	slug: "pc",
	subs: nil,
}

var digitalMusic = category{
	id:   14,
	name: "Digital Music",
	slug: "dmusic",
	subs: nil,
}

var electronics = category{
	id:   15,
	name: "Electronics",
	slug: "electronics",
	subs: nil,
}

var entertainmentCollectibles = category{
	id:   16,
	name: "Entertainment Collectibles",
	slug: "entertainment-collectibles",
	subs: nil,
}

var giftCards = category{
	id:   17,
	name: "Gift Cards",
	slug: "gift-cards",
	subs: nil,
}

var groceryGourmetFood = category{
	id:   18,
	name: "Grocery & Gourmet Food",
	slug: "grocery",
	subs: nil,
}

var healthHousehold = category{
	id:   19,
	name: "Health & Household",
	slug: "hpc",
	subs: nil,
}

var homeKitchen = category{
	id:   20,
	name: "Home & Kitchen",
	slug: "home-garden",
	subs: nil,
}

var industrialScientific = category{
	id:   21,
	name: "Industrial & Scientific",
	slug: "industrial",
	subs: nil,
}

var kindleStore = category{
	id:   22,
	name: "Kindle Store",
	slug: "digital-text",
	subs: nil,
}

var kitchenDining = category{
	id:   23,
	name: "Kitchen & Dining",
	slug: "kitchen",
	subs: nil,
}

var magazineSubscriptions = category{
	id:   24,
	name: "Magazine Subscriptions",
	slug: "magazines",
	subs: nil,
}

var moviesTV = category{
	id:   25,
	name: "Movies & TV",
	slug: "movies-tv",
	subs: nil,
}

var musicalInstruments = category{
	id:   26,
	name: "Musical Instruments",
	slug: "musical-instruments",
	subs: nil,
}

var officeProducts = category{
	id:   27,
	name: "Office Products",
	slug: "office-products",
	subs: nil,
}

var patioLawnGarden = category{
	id:   28,
	name: "Patio, Lawn & Garden",
	slug: "lawn-garden",
	subs: nil,
}

var petSupplies = category{
	id:   29,
	name: "Pet Supplies",
	slug: "pet-supplies",
	subs: nil,
}

var primePantry = category{
	id:   30,
	name: "Prime Pantry",
	slug: "pantry",
	subs: nil,
}

var software = category{
	id:   31,
	name: "Software",
	slug: "software",
	subs: nil,
}

var sportsOutdoors = category{
	id:   32,
	name: "Sports & Outdoors",
	slug: "sporting-goods",
	subs: nil,
}

var sportsCollectibles = category{
	id:   33,
	name: "Sports Collectibles",
	slug: "sports-collectibles",
	subs: nil,
}

var toolsHomeImprovement = category{
	id:   34,
	name: "Tools & Home Improvement",
	slug: "hi",
	subs: nil,
}

var toysGames = category{
	id:   35,
	name: "Toys & Games",
	slug: "toys-and-games",
	subs: nil,
}

var videoGames = category{
	id:   36,
	name: "Video Games",
	slug: "videogames",
	subs: nil,
}

// categories holds all categories defined above in a collection
var categories = [...]category{
	appliances,
	appsGames,
	artsCraftsSewing,
	automotive,
	baby,
	beautyPersonalCare,
	books,
	cDsVinyl,
	cameraPhoto,
	cellPhonesAccessories,
	clothingShoesJewelry,
	collectibleCoins,
	computersAccessories,
	digitalMusic,
	electronics,
	entertainmentCollectibles,
	giftCards,
	groceryGourmetFood,
	healthHousehold,
	homeKitchen,
	industrialScientific,
	kindleStore,
	kitchenDining,
	magazineSubscriptions,
	moviesTV,
	musicalInstruments,
	officeProducts,
	patioLawnGarden,
	petSupplies,
	primePantry,
	software,
	sportsOutdoors,
	sportsCollectibles,
	toolsHomeImprovement,
	toysGames,
	videoGames,
}

// getLinks fetches all the links that need to be scrapped
// All these links belong to a certain category
func (cat *category) getLinks() ([]string, error) {
	links := make([]string, len(cat.subs))
	// Every category must have at least 1 subcateogry
	if cat.subs == nil {
		return nil, errors.New("No subcategories found")
	}
	// Get links for all subcategories belonging to this category
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

// filterCategories holds the categories that the user chose on search request
// This way the application knows which categories to scrappe
func filterCategories(catIDs []string) ([]category, error) {
	if len(catIDs) == 0 {
		return nil, nil
	}
	cats := make([]category, len(catIDs))
	for _, id := range catIDs {
		id, err := strconv.ParseUint(id, 10, 8)
		if err != nil {
			return nil, err
		}
		catID := uint8(id)
		for i, c := range categories {
			if catID == uint8(c.id) {
				cats[i] = c
				break
			}
		}
	}

	return cats, nil
}

// GetCategories fetches a list of all main categories in a map
// After we load this map in template to be rendered as a HTML select
func GetCategories() map[uint8]string {
	m := make(map[uint8]string)
	for _, v := range categories {
		m[uint8(v.id)] = v.name
	}

	return m
}
