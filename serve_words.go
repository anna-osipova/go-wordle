package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	C "./check_error"
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
			index := strings.Index(word, string(r))
			color := "Grey"
			if index == i {
				color = "Green"
			} else if index > 0 {
				color = "Yellow"
			} else {
				color = "Grey"
			}
			letters = append(letters, Letter{
				Color: color,
				Letter: string(r),
			})
		}
		c.JSON(200, letters)
	})
	r.Run()
}
