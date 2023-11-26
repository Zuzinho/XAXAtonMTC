package login

import (
	"fmt"
	"math/rand"
	"time"
)

func generateRandomCode() string {
	rand.Seed(time.Now().UnixNano())

	min := 10000
	max := 99999

	code := rand.Intn(max-min+1) + min

	return fmt.Sprintf("%05d", code)
}
