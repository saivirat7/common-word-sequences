package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

func ProcessStdIn() {
	bufChunk, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading stdin:", err)
		return
	}
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	wg.Add(1)

	endOfChunkElements = make([][]string, 10000)
	phrases = make(map[string]int)
	//assigning memory to n+1 blocks
	for i := range endOfChunkElements {
		endOfChunkElements[i] = make([]string, 4)
	}

	go func(bufChunk []byte, chunk, noOfChunks int, wg *sync.WaitGroup, mtx *sync.Mutex) {
		defer wg.Done()
		// fmt.Println("String:", string(bufChunk))

		//process text and form the map of phrases
		processText(string(bufChunk), chunk, noOfChunks, wg, mtx)
		//process end of chunks for new line cases
	}(bufChunk, 1, 1, &wg, &mtx)

	wg.Wait()

	processEndOfChunkElements()
	//calculating the frequencies and sorting by rank
	pairList := rankByWordCount(phrases)

	//print phrases
	printPhrases(pairList)
}
