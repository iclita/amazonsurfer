package crawler

// The min and max limits to sleep between requests
// The crawler will choose a random number of seconds to sleep
const (
	minSleep = 10
	maxSleep = 60
)

// Request headers to simulate a real browser
var headers = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Encoding": "gzip, deflate, sdch, br",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
}
