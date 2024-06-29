package prep

import (
	"fmt"
	"math/rand"
	"time"
)

func InitNum() []string {
	// Initialize the slice with numbers from 1 to 100000
	numbers := make([]string, 100000)
	for i := 0; i < 100000; i++ {
		//numbers[i] = i + 1
		numbers[i] = fmt.Sprintf("%05d", i+1)
	}

	// Shuffle the slice
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(numbers), func(i, j int) { numbers[i], numbers[j] = numbers[j], numbers[i] })
	return numbers
}
