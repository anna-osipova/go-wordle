package main

import (
	"bufio"
	"log"
	"os"

	"github.com/anna-osipova/go-wordle/errorcheck"
)

func isPlural(list []string, item string) bool {
	substr := item[:len(item)-1]
	if substr+"s" != item {
		return false
	}
	for _, listItem := range list {
		if listItem == substr {
			return true
		}
	}
	return false
}

func filterWords() {
	file, err := os.Open("./words.txt")
	errorcheck.Check(err)
	defer file.Close()

	write_file, err := os.Create("./simple_words_5.txt")
	errorcheck.Check(err)
	defer write_file.Close()
	writer := bufio.NewWriter(write_file)

	scanner := bufio.NewScanner(file)

	var words []string

	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	var filteredWords []string

	for _, word := range words {
		if len(word) == 5 && !isPlural(words, word) {
			filteredWords = append(filteredWords, word)
		}
	}

	for _, word := range filteredWords {
		_, err := writer.WriteString(word + "\n")
		errorcheck.Check(err)
	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
