package game

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/anna-osipova/go-wordle/common"
	"github.com/anna-osipova/go-wordle/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type newGamePayload struct {
	Word string `json:"word" binding:"required"`
}

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := jwt.ParseWithClaims(tokenString, &service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtService.GetSecretKey()), nil
		})
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{
				ErrorCode: "INVALID_AUTH_TOKEN",
				Message:   "Invalid auth token",
			})
			return
		}
		if token.Valid {
			claims := token.Claims.(*service.CustomClaims)
			c.Set("session_id", claims.SessionId)
		} else {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorResponse{
				ErrorCode: "INVALID_AUTH_TOKEN",
				Message:   "Invalid auth token",
			})
			return
		}
	}
}

func GameRegister(router *gin.RouterGroup) {
	var jwtService service.JWTService = service.JWTAuthService()
	router.POST("/new", GameNew(jwtService))
	router.GET("/new/random", GameNewRandom(jwtService))

	router.Use(AuthorizeJWT(jwtService))
	router.GET("/status", GameStatus(jwtService))
	router.POST("/guess/:word", GameGuess(jwtService))
}

const MAX_ATTEMPTS = 6

type gameGuessResponse struct {
	Letters []Letter `json:"letters"`
	Word    string   `json:"word,omitempty"`
}

func GameGuess(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		processingError := common.ErrorResponse{
			Message:   "Some issue",
			ErrorCode: "ERROR",
		}

		sessionId := c.MustGet("session_id").(string)
		session, err := service.GetSessionById(sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorResponse{
				Message:   "Game does not exist",
				ErrorCode: "INVALID_SESSION",
			})
			return
		}

		words := c.MustGet("word_list").([]string)
		if len(session.Word) != 5 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorResponse{
				Message:   "Game does not exist",
				ErrorCode: "INVALID_SESSION",
			})
			return
		}

		wordGuess := c.Param("word")
		if len(wordGuess) != 5 || CheckWordExists(words, wordGuess) == false {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{
				Message:   "Word does not exist",
				ErrorCode: "INVALID_GUESS_WORD",
			})
			return
		}

		attempts, err := GetAttempts(session.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, processingError)
			return
		}
		if HasMadeSameAttempt(attempts, wordGuess) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{
				Message:   "Same word has already been attempted",
				ErrorCode: "DUPLICATE_ATTEMPT",
			})
			return
		}

		if len(attempts) >= MAX_ATTEMPTS {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{
				Message:   "Out of tries",
				ErrorCode: "NO_TRIES",
			})
			return
		}

		letters := MakeGuess(wordGuess, session.Word)

		attempt := Attempt{
			SessionId: session.ID,
			WordGuess: wordGuess,
			Letters:   letters,
		}
		err = CreateAttempt(&attempt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, processingError)
			return
		}

		guessResponse := gameGuessResponse{
			Letters: letters,
		}
		if len(attempts) == MAX_ATTEMPTS-1 {
			guessResponse.Word = session.Word
		}
		c.JSON(http.StatusOK, guessResponse)
	}
}

func GameNew(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newGame newGamePayload
		if err := c.ShouldBindJSON(&newGame); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{
				ErrorCode: "VALIDATION_ERROR",
				Message:   err.Error(),
			})
			return
		}

		words := c.MustGet("word_list").([]string)
		word := strings.ToLower(newGame.Word)
		if len(word) != 5 || CheckWordExists(words, word) == false {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{
				Message:   "Word does not exist",
				ErrorCode: "INVALID_GUESS_WORD",
			})
			return
		}
		session := &service.Session{Word: word}
		if err := service.CreateSession(session); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorResponse{
				ErrorCode: "DB_ERROR",
				Message:   err.Error(),
			})
			return
		}
		token := jwtService.GenerateToken(session)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func GameNewRandom(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		words := c.MustGet("word_list").([]string)
		wordsCount := len(words)

		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		word := words[random.Intn(wordsCount)]
		session := &service.Session{Word: word}
		if err := service.CreateSession(session); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrorResponse{
				ErrorCode: "DB_ERROR",
				Message:   err.Error(),
			})
			return
		}
		token := jwtService.GenerateToken(session)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

type gameStatusResponse struct {
	Attempts []Attempt `json:"attempts"`
	Word     string    `json:"word,omitempty"`
}

func GameStatus(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		processingError := common.ErrorResponse{
			Message:   "Some issue",
			ErrorCode: "ERROR",
		}

		sessionId := c.MustGet("session_id").(string)
		session, err := service.GetSessionById(sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ErrorResponse{
				Message:   "Game does not exist",
				ErrorCode: "INVALID_SESSION",
			})
			return
		}

		attempts, err := GetAttempts(session.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, processingError)
			return
		}

		response := gameStatusResponse{
			Attempts: attempts,
		}
		if len(attempts) == MAX_ATTEMPTS {
			response.Word = session.Word
		}

		c.JSON(http.StatusOK, response)
	}
}
