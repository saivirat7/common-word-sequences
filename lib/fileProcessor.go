package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

var endOfChunkElements [][]string

var phrases map[string]int

func ProcessFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	defer file.Close()

	fileStats, _ := file.Stat()
	totalSize := fileStats.Size()

	fileReader := bufio.NewReader(file)
	if totalSize < int64(ByteSize) {
		totalSize = int64(ByteSize)
	}
	noOfChunks := int(totalSize) / ByteSize
	endOfChunkElements = make([][]string, 10000)
	phrases = make(map[string]int)

	//assigning memory to n+1 blocks
	for i := range endOfChunkElements {
		endOfChunkElements[i] = make([]string, 4)
	}

	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}

	for chunk := 0; ; chunk++ {
		bufChunk := make([]byte, ByteSize)
		end, err := fileReader.Read(bufChunk)
		bufChunk = bufChunk[:end]
		// fmt.Println("Size of chunk:", len(bufChunk))

		if end == 0 {
			if err != nil {
				fmt.Println(err)
				break
			}
			if err == io.EOF {
				break
			}
			return
		}

		//process till end of line for the chunk
		lineTillEnd, err := fileReader.ReadBytes('\n')
		if err != io.EOF {
			bufChunk = append(bufChunk, lineTillEnd...)
		}

		wg.Add(1)
		go func(bufChunk []byte, chunk, noOfChunks int, wg *sync.WaitGroup, mtx *sync.Mutex) {
			defer wg.Done()

			//process text and form the map of phrases
			processText(string(bufChunk), chunk, noOfChunks, wg, mtx)
			//process end of chunks for new line cases
		}(bufChunk, chunk, noOfChunks, &wg, &mtx)

	}
	wg.Wait()

	processEndOfChunkElements()
	//calculating the frequencies and sorting by rank
	pairList := rankByWordCount(phrases)

	//print phrases
	printPhrases(pairList)
}
