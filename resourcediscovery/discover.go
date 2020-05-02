package resourcediscovery

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindFiles - Finds all the files that subscribe to the pattern in the target directory
func FindFiles(targetDir string, pattern ...string) ([]string, error) {

	var matches []string
	for _, v := range pattern {
		matches, err := filepath.Glob(targetDir + v)

		if err != nil {
			fmt.Println(err)
		}

		if len(matches) != 0 {
			fmt.Println("Found : ", matches)
		}
	}
	return matches, nil
}

func main() {

	if len(os.Args) <= 2 {
		fmt.Printf("USAGE : %s <target_directory> <target_filename or part of filename> \n", os.Args[0])
		os.Exit(0)
	}

	targetDirectory := os.Args[1] // get the target directory
	fileName := os.Args[2:]       // to handle wildcard such as filename*.go

	FindFiles(targetDirectory, fileName...)

}
