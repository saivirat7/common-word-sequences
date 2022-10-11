package main

import (
	"fmt"
	"os"

	"github.com/saivirat7/three-word-sequences/lib"
)

func main() {
	if os.Args == nil || len(os.Args) == 1 {
		fmt.Println("Provided stdIn")
		lib.ProcessStdIn()
	} else if os.Args != nil || len(os.Args) > 0 {
		fmt.Println("File provided in arguments")
		for i, fileName := range os.Args {
			if i != 0 {
				lib.ProcessFile(fileName)
			}
		}
	} else {
		fmt.Println("Wrong input!!")
	}
}
