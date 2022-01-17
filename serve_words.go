package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	C "github.com/anna-osipova/go-wordle/check_error"
	"github.com/gin-gonic/gin"
)

type WordsResponse struct {
	Count         int `json:"count"`
	Words         []string `json:"words"`
}

type WordResponse struct {
	Word          string `json:"word"`
}

type Letter struct {
	Letter        string `json:"letter`
	Color         string `json:"color`
}

type WordHelp struct {
	Include string `json:"include" binding:"required"`
	Exclude string `json:"exclude" binding:"required"`
	Template string `json:"template" binding:"required"`
}

func CountExistingLetters(letters []Letter, letter string) int {
	count := 0
	for _, n := range letters {
		if n.Letter == letter {
			count++
		}
	}
	return count
}

func CountExactMatches(word string, guess_word string, letter string) int {
	count := 0
	for i, n := range word {
		l := string(n)
		if  letter == l && l == string(guess_word[i]) {
			count++
		}
	}
	return count
}

func IsMatch(word_help *WordHelp, word string) bool {
	// word contains letters from "include"
	for _, n := range word_help.Include {
		if strings.Index(word, string(n)) < 0 {
			return false
		}
	}
	// word doesn't contain letters from "exclude"
	for _, n := range word_help.Exclude {
		if strings.Index(word, string(n)) > -1 {
			return false
		}
	}
	// word matches "template"
	if word_help.Template != "" {
		for i, n := range word_help.Template {
			letter := string(n)
			if letter != "_" && string(word[i]) != letter {
				return false
			}
		}
	}
	return true
}

func main() {
	file, err := os.Open("./words_5.txt")
	C.Check(err)
	defer file.Close()

	var words []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}
	log.Println("Finished reading")
	words_count := len(words)

	r := gin.Default()

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	word := words[random.Intn(words_count)]

	r.GET("/words", func(c *gin.Context) {
		words_response := WordsResponse{
			Count: words_count,
			Words: words,
		}
		c.JSON(200, words_response)
	})

	r.GET("/words/random", func(c *gin.Context) {
		word_response := WordResponse{
			Word: words[random.Intn(words_count)],
		}
		c.JSON(200, word_response)
	})

	r.GET("/words/:word", func(c *gin.Context) {
		word_guess := c.Param("word")
		letters := make([]Letter, 0)
		for i, r := range word_guess {
			letter := string(r)
			index := strings.Index(word, letter)
			color := "Grey"
			// Do green first, then rest
			if letter == string(word[i]) {
				// Target word has the same letter in the same position
				color = "Green"
			} else if index > -1 &&
				// Target word has more of the same letter than what has already been found
				strings.Count(word, letter) > CountExistingLetters(letters, letter) &&
				// Guess word has more of letter than there are exat matches
				strings.Count(word_guess, letter) > CountExactMatches(word, word_guess, letter) {
				color = "Yellow"
			} else {
				color = "Grey"
			}
			letters = append(letters, Letter{
				Color: color,
				Letter: letter,
			})
		}
		c.JSON(200, letters)
	})

	r.POST("/words/hint", func(c *gin.Context) {
		var word_help WordHelp
		c.BindJSON(&word_help)

		var matching_words []string
		for _, word := range words {

			if IsMatch(&word_help, word) == true {
				matching_words = append(matching_words, word)
			}
		}

		c.JSON(200, gin.H{ "matches": matching_words })
	})
	r.Run()
}
