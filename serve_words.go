package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/anna-osipova/go-wordle/db"
	"github.com/anna-osipova/go-wordle/errorcheck"
	"github.com/anna-osipova/go-wordle/game"
	"github.com/anna-osipova/go-wordle/hint"
	"github.com/anna-osipova/go-wordle/letters"
	"github.com/anna-osipova/go-wordle/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	godotenv.Load()

	dbInstance := db.Init()
	game.AutoMigrate()
	service.AutoMigrate()
	dbConn, err := dbInstance.DB()
	defer dbConn.Close()
	if err != nil {
		panic(fmt.Sprintf("Error connecting to DB: %s", err))
	}

	r := gin.Default()
	r.LoadHTMLFiles("client/build/index.html")
	r.Static("/static", "./client/build/static")
	r.Static("/assets", "./client/build/assets")
	r.StaticFile("/manifest.json", "./client/build/manifest.json")
	r.StaticFile("/robots.txt", "./client/build/robots.txt")

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	r.Use(cors.New(config))

	r.Use(WordsMiddleware)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1 := r.Group("/api")
	v1.Use(WordsMiddleware)

	v1.GET("/words", func(c *gin.Context) {
		words := c.MustGet("word_list").([]string)
		wordsResponse := WordsResponse{
			Count: len(words),
			Words: words,
		}
		c.JSON(200, wordsResponse)
	})

	hintGroup := v1.Group("/hint")
	hint.HintRegister(hintGroup)

	lettersGroup := v1.Group("/letters")
	letters.LettersRegister(lettersGroup)

	gameGroup := v1.Group("/game")
	game.GameRegister(gameGroup)

	if os.Getenv("USE_LOCALHOST") == "true" {
		r.Run("localhost:" + os.Getenv("PORT"))
	} else {
		r.Run()
	}
}
