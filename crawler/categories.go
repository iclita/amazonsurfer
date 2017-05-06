package crawler

import "strconv"

var appliances = Category{
	ID:            1,
	Name:          "Appliances",
	Link:          "https://www.amazon.com/Best-Sellers-Appliances/zgbs/appliances",
	SubCategories: nil,
}

var appsGames = Category{
	ID:            2,
	Name:          "Apps & Games",
	Link:          "https://www.amazon.com/Best-Sellers-Appstore-Android/zgbs/mobile-apps",
	SubCategories: nil,
}

var artsCraftsSewing = Category{
	ID:            3,
	Name:          "Arts, Crafts & Sewing",
	Link:          "https://www.amazon.com/Best-Sellers-Arts-Crafts-Sewing/zgbs/arts-crafts",
	SubCategories: nil,
}

var automotive = Category{
	ID:            4,
	Name:          "Automotive",
	Link:          "https://www.amazon.com/Best-Sellers-Automotive/zgbs/automotive",
	SubCategories: nil,
}

var baby = Category{
	ID:            5,
	Name:          "Baby",
	Link:          "https://www.amazon.com/Best-Sellers-Baby/zgbs/baby-products",
	SubCategories: nil,
}

var beautyPersonalCare = Category{
	ID:            6,
	Name:          "Beauty & Personal Care",
	Link:          "https://www.amazon.com/Best-Sellers-Beauty/zgbs/beauty",
	SubCategories: nil,
}

var books = Category{
	ID:            7,
	Name:          "Books",
	Link:          "https://www.amazon.com/best-sellers-books-Amazon/zgbs/books",
	SubCategories: nil,
}

var cDsVinyl = Category{
	ID:            8,
	Name:          "CDs & Vinyl",
	Link:          "https://www.amazon.com/best-sellers-music-albums/zgbs/music",
	SubCategories: nil,
}

var cameraPhoto = Category{
	ID:            9,
	Name:          "Camera & Photo",
	Link:          "https://www.amazon.com/best-sellers-camera-photo/zgbs/photo",
	SubCategories: nil,
}

var cellPhonesAccessories = Category{
	ID:            10,
	Name:          "Cell Phones & Accessories",
	Link:          "https://www.amazon.com/Best-Sellers-Cell-Phones-Accessories/zgbs/wireless",
	SubCategories: nil,
}

var clothingShoesJewelry = Category{
	ID:            11,
	Name:          "Clothing, Shoes & Jewelry",
	Link:          "https://www.amazon.com/Best-Sellers/zgbs/fashion",
	SubCategories: nil,
}

var collectibleCoins = Category{
	ID:            12,
	Name:          "Collectible Coins",
	Link:          "https://www.amazon.com/Best-Sellers-Collectible-Coins/zgbs/coins",
	SubCategories: nil,
}

var computersAccessories = Category{
	ID:            13,
	Name:          "Computers & Accessories",
	Link:          "https://www.amazon.com/Best-Sellers-Computers-Accessories/zgbs/pc",
	SubCategories: nil,
}

var digitalMusic = Category{
	ID:            14,
	Name:          "Digital Music",
	Link:          "https://www.amazon.com/Best-Sellers-MP3-Downloads/zgbs/dmusic",
	SubCategories: nil,
}

var electronics = Category{
	ID:            15,
	Name:          "Electronics",
	Link:          "https://www.amazon.com/Best-Sellers-Electronics/zgbs/electronics",
	SubCategories: nil,
}

var entertainmentCollectibles = Category{
	ID:            16,
	Name:          "Entertainment Collectibles",
	Link:          "https://www.amazon.com/Best-Sellers-Entertainment-Collectibles/zgbs/entertainment-collectibles",
	SubCategories: nil,
}

var giftCards = Category{
	ID:            17,
	Name:          "Gift Cards",
	Link:          "https://www.amazon.com/Best-Sellers-Gift-Cards/zgbs/gift-cards",
	SubCategories: nil,
}

var groceryGourmetFood = Category{
	ID:            18,
	Name:          "Grocery & Gourmet Food",
	Link:          "https://www.amazon.com/Best-Sellers-Grocery-Gourmet-Food/zgbs/grocery",
	SubCategories: nil,
}

var healthHousehold = Category{
	ID:            19,
	Name:          "Health & Household",
	Link:          "https://www.amazon.com/Best-Sellers-Health-Personal-Care/zgbs/hpc",
	SubCategories: nil,
}

var homeKitchen = Category{
	ID:            20,
	Name:          "Home & Kitchen",
	Link:          "https://www.amazon.com/Best-Sellers-Home-Kitchen/zgbs/home-garden",
	SubCategories: nil,
}

var industrialScientific = Category{
	ID:            21,
	Name:          "Industrial & Scientific",
	Link:          "https://www.amazon.com/Best-Sellers-Industrial-Scientific/zgbs/industrial",
	SubCategories: nil,
}

var kindleStore = Category{
	ID:            22,
	Name:          "Kindle Store",
	Link:          "https://www.amazon.com/Best-Sellers-Kindle-Store/zgbs/digital-text",
	SubCategories: nil,
}

var kitchenDining = Category{
	ID:            23,
	Name:          "Kitchen & Dining",
	Link:          "https://www.amazon.com/Best-Sellers-Kitchen-Dining/zgbs/kitchen",
	SubCategories: nil,
}

var magazineSubscriptions = Category{
	ID:            24,
	Name:          "Magazine Subscriptions",
	Link:          "https://www.amazon.com/Best-Sellers-Magazines/zgbs/magazines",
	SubCategories: nil,
}

var moviesTV = Category{
	ID:            25,
	Name:          "Movies & TV",
	Link:          "https://www.amazon.com/best-sellers-movies-TV-DVD-Blu-ray/zgbs/movies-tv",
	SubCategories: nil,
}

var musicalInstruments = Category{
	ID:            26,
	Name:          "Musical Instruments",
	Link:          "https://www.amazon.com/Best-Sellers-Musical-Instruments/zgbs/musical-instruments",
	SubCategories: nil,
}

var officeProducts = Category{
	ID:            27,
	Name:          "Office Products",
	Link:          "https://www.amazon.com/Best-Sellers-Office-Products/zgbs/office-products",
	SubCategories: nil,
}

var patioLawnGarden = Category{
	ID:            28,
	Name:          "Patio, Lawn & Garden",
	Link:          "https://www.amazon.com/Best-Sellers-Patio-Lawn-Garden/zgbs/lawn-garden",
	SubCategories: nil,
}

var petSupplies = Category{
	ID:            29,
	Name:          "Pet Supplies",
	Link:          "https://www.amazon.com/Best-Sellers-Pet-Supplies/zgbs/pet-supplies",
	SubCategories: nil,
}

var primePantry = Category{
	ID:            30,
	Name:          "Prime Pantry",
	Link:          "https://www.amazon.com/Best-Sellers-Prime-Pantry/zgbs/pantry",
	SubCategories: nil,
}

var software = Category{
	ID:            31,
	Name:          "Software",
	Link:          "https://www.amazon.com/best-sellers-software/zgbs/software",
	SubCategories: nil,
}

var sportsOutdoors = Category{
	ID:            32,
	Name:          "Sports & Outdoors",
	Link:          "https://www.amazon.com/Best-Sellers-Sports-Outdoors/zgbs/sporting-goods",
	SubCategories: nil,
}

var sportsCollectibles = Category{
	ID:            33,
	Name:          "Sports Collectibles",
	Link:          "https://www.amazon.com/Best-Sellers-Sports-Collectibles/zgbs/sports-collectibles",
	SubCategories: nil,
}

var toolsHomeImprovement = Category{
	ID:            34,
	Name:          "Tools & Home Improvement",
	Link:          "https://www.amazon.com/Best-Sellers-Home-Improvement/zgbs/hi",
	SubCategories: nil,
}

var toysGames = Category{
	ID:            35,
	Name:          "Toys & Games",
	Link:          "https://www.amazon.com/Best-Sellers-Toys-Games/zgbs/toys-and-games",
	SubCategories: nil,
}

var videoGames = Category{
	ID:            36,
	Name:          "Video Games",
	Link:          "https://www.amazon.com/best-sellers-video-games/zgbs/videogames",
	SubCategories: nil,
}

var categories = [...]Category{
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

func GetCategories() map[uint8]string {
	m := make(map[uint8]string)
	for _, v := range categories {
		m[v.ID] = v.Name
	}

	return m
}

func filterCategories(catIDs []string) ([]Category, error) {
	if len(catIDs) == 0 {
		return nil, nil
	}
	cats := make([]Category, len(catIDs))
	for _, id := range catIDs {
		id, err := strconv.ParseUint(id, 10, 8)
		if err != nil {
			return nil, err
		}
		catID := uint8(id)
		for _, c := range categories {
			if catID == c.ID {
				cats = append(cats, c)
				break
			}
		}
	}

	return cats, nil
}
