package lib

import (
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Pair struct {
	Key   string
	Value int
}

type WordRanks []Pair

func (p WordRanks) Len() int           { return len(p) }
func (p WordRanks) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p WordRanks) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func processText(text string, chunk int, lastChunk int, wg *sync.WaitGroup, mtx *sync.Mutex) {
	//Replacing all next lines with space
	text = strings.ReplaceAll(text, "\n", " ")

	//removing all symbols from text
	text = RemovePunctuations(text)

	//convert all text to lower case
	text = strings.ToLower(text)

	//splitting text to words
	allWords := strings.Split(text, " ")

	words := []string{}
	//ignoring extra spaces and next line space
	for i := range allWords {
		if allWords[i] != "" {
			words = append(words, allWords[i])
		}
	}
	allWords = nil

	//logic to store end of the chunk elements
	mtx.Lock()
	if chunk == 0 {
		endOfChunkElements[chunk][0] = words[len(words)-2]
		endOfChunkElements[chunk][1] = words[len(words)-1]
	} else if chunk != lastChunk {
		endOfChunkElements[chunk-1][2] = words[0]
		endOfChunkElements[chunk-1][3] = words[1]
		endOfChunkElements[chunk][0] = words[len(words)-2]
		endOfChunkElements[chunk][1] = words[len(words)-1]
	}
	mtx.Unlock()

	//form the phrases map
	for i := 0; i < len(words)-WordLength+1; i++ {
		currPhrase := GetPhraseFromWords(words[i : i+WordLength])
		if currPhrase != "" {
			insertMap(currPhrase, wg, mtx)
		}
	}
}

func insertMap(phrase string, wg *sync.WaitGroup, mtx *sync.Mutex) {
	mtx.Lock()
	if count, exists := phrases[phrase]; exists {
		phrases[phrase] = count + 1
	} else {
		phrases[phrase] = 1
	}
	mtx.Unlock()
}

func insertMap2(phrase string) {
	if count, exists := phrases[phrase]; exists {
		phrases[phrase] = count + 1
	} else {
		phrases[phrase] = 1
	}
}

func processEndOfChunkElements() {

	for _, element := range endOfChunkElements {
		phrase1 := GetPhraseFromWords(element[:3])
		phrase2 := GetPhraseFromWords(element[1:])
		if phrase1 != "" {
			insertMap2(phrase1)
		}
		if phrase2 != "" {
			insertMap2(phrase2)
		}
	}
}

func GetPhraseFromWords(words []string) string {
	finalWord := ""
	for i, word := range words {
		if word != "" {
			if i == 0 {
				finalWord += word
			} else {
				finalWord += " " + word
			}
		} else {
			return ""
		}
	}
	return finalWord
}

func RemovePunctuations(str string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func rankByWordCount(wordFrequencies map[string]int) WordRanks {
	wordRankCounts := make(WordRanks, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		wordRankCounts[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(wordRankCounts))
	return wordRankCounts
}
