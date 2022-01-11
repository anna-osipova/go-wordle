package main

import (
	"bufio"
	"log"
	"os"

	C "./check_error"
)

func filter_words() {
	file, err := os.Open("./words.txt")
	C.Check(err)
	defer file.Close()

	write_file, err := os.Create("./words_5.txt")
	C.Check(err)
	defer write_file.Close()
	writer := bufio.NewWriter(write_file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == 5 {
			_, err := writer.WriteString(word + "\n")
			C.Check(err)
		}
	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
			log.Fatal(err)
	}
}
