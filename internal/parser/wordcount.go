package parser

import (
	"os"
	"strings"
)

// findFiles recursively finds files in a given directory
func findFiles(dirname string) chan string {
	output := make(chan string)

	go func() {
		findFilesRecursively(dirname, output)
		close(output)
	}()

	return output
}

func findFilesRecursively(dirname string, output chan string) {
	dir, _ := os.Open(dirname)
	dirnames, _ := dir.Readdirnames(-1)

	for i := 0; i < len(dirnames); i++ {
		fullpath := dirname + "/" + dirnames[i]
		file, _ := os.Stat(fullpath)

		if file.IsDir() {
			findFilesRecursively(fullpath, output)
		} else {
			output <- fullpath
		}
	}
}

// wordCount parses a file and splits the contents into "words" using various
// delimiters (whitespace, colons ":", semicolons ";",  pipes "|")
// If a string contains quotes (single or double), it counts the string inside the quotes, as well
// as the string with quotes.

// TODO: probably should incorporate the filename at least here somehow
func wordCount(filename string, output chan map[string]int) {
	results := make(map[string]int)

	for line := range line(filename) {
		words := strings.Fields(line)
		for _, v := range words {
			// Count the whole word
			results[v]++

			// If the word is in quotes or other containers, strip them and count that as well
			// TODO: make this cleaner, move to own function
			rv := []rune(v)
			if len(rv) > 1 {
				if string(rv[0]) == "\"" && string(rv[len(rv)-1]) == "\"" {
					results[string(rv[1:len(rv)-1])]++
				}

				if string(rv[0]) == "[" && string(rv[len(rv)-1]) == "]" {
					results[string(rv[1:len(rv)-1])]++
				}

				if string(rv[0]) == "<" && string(rv[len(rv)-1]) == ">" {
					results[string(rv[1:len(rv)-1])]++
				}

				if string(rv[0]) == "'" && string(rv[len(rv)-1]) == "'" {
					results[string(rv[1:len(rv)-1])]++
				}

				// If the word contains one of our splitters, split, and add to results
				splits := [...]string{":", ";", "|"}
				for _, s := range splits {
					if strings.Contains(v, s) {
						ss := strings.Split(v, s)
						for _, w := range ss {
							results[w]++
						}
					}
				}

			}
		}
	}
	output <- results
}

func reducer(input chan map[string]int, output chan map[string]int) {
	results := make(map[string]int)

	for newMatches := range input {
		for key, value := range newMatches {
			previousCount, exists := results[key]

			if !exists {
				results[key] = value
			} else {
				results[key] = previousCount + value
			}
		}
	}

	output <- results
}
