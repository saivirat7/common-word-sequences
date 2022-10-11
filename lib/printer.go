package lib

import "fmt"

func printPhrases(pairList WordRanks) {
	fmt.Println("3 word Phrases in the given file:")
	count := 0
	if len(pairList) < printCount {
		count = len(pairList)
	} else {
		count = printCount
	}
	for i := 0; i < count; i++ {
		fmt.Println("String:", pairList[i].Key, " Frequency:", pairList[i].Value)
	}
}
