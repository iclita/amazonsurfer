package crawler

import (
	"math/rand"
	"time"
)

// sleep simply puts the program to sleep for a random number of seconds
// between min and max when looking for products
func sleep(min int, max int) {
	delay := min + rand.Intn(max-min)
	time.Sleep(time.Duration(delay) * time.Second)
}
