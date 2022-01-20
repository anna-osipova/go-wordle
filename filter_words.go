package main

import (
	"bufio"
	"log"
	"os"

	"github.com/anna-osipova/go-wordle/errorcheck"
)

func filter_words() {
	file, err := os.Open("./words.txt")
	errorcheck.Check(err)
	defer file.Close()

	write_file, err := os.Create("./words_5.txt")
	errorcheck.Check(err)
	defer write_file.Close()
	writer := bufio.NewWriter(write_file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == 5 {
			_, err := writer.WriteString(word + "\n")
			errorcheck.Check(err)
		}
	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
