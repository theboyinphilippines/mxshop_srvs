package main

import (
	"fmt"
	"mxshop_srvs/pick_cards/prep"
)

// 挑选卡牌，在（1，100000）中生成不重复的10万个数字，每个数字都是5位数，不足的补0
func main() {

	numbers := prep.InitNum()

	// Function to get the next batch of 3 numbers
	getNextBatch := func(numbers []string, batchSize int) ([]string, []string) {
		if len(numbers) < batchSize {
			return numbers, nil
		}
		batch := numbers[:batchSize]
		numbers = numbers[batchSize:]
		return batch, numbers
	}

	// Loop until all numbers are picked
	var batches [][]string
	for len(numbers) > 0 {
		var batch []string
		batch, numbers = getNextBatch(numbers, 5)
		batches = append(batches, batch)
		fmt.Printf("Picked batch: %v\n", batch)
	}

	// Print the total number of batches
	fmt.Printf("Total batches: %d\n", len(batches))
}
