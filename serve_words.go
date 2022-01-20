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
	wordsCount := len(words)

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return words[random.Intn(wordsCount)]
}

func main() {
	r := gin.Default()

	r.Use(WordsMiddleware)

	r.GET("/words", func(c *gin.Context) {
		words := c.MustGet("word_list").([]string)
		wordsResponse := WordsResponse{
			Count: len(words),
			Words: words,
		}
		c.JSON(200, wordsResponse)
	})

	r.GET("/words/random", func(c *gin.Context) {
		words := c.MustGet("word_list").([]string)

		wordResponse := WordResponse{
			Word: GetRandomWord(words),
		}
		c.JSON(200, wordResponse)
	})

	v1 := r.Group("/api")
	v1.Use(WordsMiddleware)

	hintGroup := v1.Group("/hint")
	hint.HintRegister(hintGroup)

	lettersGroup := v1.Group("/letters")
	letters.LettersRegister(lettersGroup)

	gameGroup := v1.Group("/game")
	game.GameRegister(gameGroup)
	r.Run("localhost:8080")
}
