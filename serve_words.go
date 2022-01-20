package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/anna-osipova/go-wordle/errorcheck"
	"github.com/anna-osipova/go-wordle/game"
	"github.com/anna-osipova/go-wordle/hint"
	"github.com/anna-osipova/go-wordle/letters"
	"github.com/gin-gonic/gin"
)

type WordsResponse struct {
	Count int      `json:"count"`
	Words []string `json:"words"`
}

type WordResponse struct {
	Word string `json:"word"`
}

func WordsMiddleware(c *gin.Context) {
	file, err := os.Open("./words_5.txt")
	errorcheck.Check(err)
	defer file.Close()

	var words []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}
	log.Println("Finished reading:", len(words))
	c.Set("word_list", words)
}

func GetRandomWord(words []string) string {
	words_count := len(words)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return words[random.Intn(words_count)]
}

func main() {
	r := gin.Default()

	r.Use(WordsMiddleware)

	r.GET("/words", func(c *gin.Context) {
		words := c.MustGet("word_list").([]string)
		words_response := WordsResponse{
			Count: len(words),
			Words: words,
		}
		c.JSON(200, words_response)
	})

	r.GET("/words/random", func(c *gin.Context) {
		words := c.MustGet("word_list").([]string)

		word_response := WordResponse{
			Word: GetRandomWord(words),
		}
		c.JSON(200, word_response)
	})

	v1 := r.Group("/api")
	v1.Use(WordsMiddleware)

	hint_group := v1.Group("/hint")
	hint.HintRegister(hint_group)

	letters_group := v1.Group("/letters")
	letters.LettersRegister(letters_group)
	r.Run("localhost:8080")
}
