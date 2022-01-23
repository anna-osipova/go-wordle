package game

import (
	"fmt"
	"net/http"

	"github.com/anna-osipova/go-wordle/logic"
	"github.com/anna-osipova/go-wordle/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type newGamePayload struct {
	Word string `json:"word" binding:"required"`
}

func AuthorizeJWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		token, err := jwt.ParseWithClaims(tokenString, &service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtService.GetSecretKey()), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		if token.Valid {
			claims := token.Claims.(*service.CustomClaims)
			c.Set("secret_word", claims.Word)
			c.Set("attempts", int(claims.Attempts))
		} else {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func GameRegister(router *gin.RouterGroup) {
	var jwtService service.JWTService = service.JWTAuthService()
	router.POST("/new", GameStart(jwtService))
	router.Use(AuthorizeJWT(jwtService))
	router.POST("/guess/:word", GameGuess(jwtService))
}

type GameGuessResponse struct {
	Token   string         `json:"token"`
	Letters []logic.Letter `json:"letters"`
}

func GameGuess(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		word := c.MustGet("secret_word").(string)
		attempts := c.MustGet("attempts").(int)
		words := c.MustGet("word_list").([]string)
		if len(word) != 5 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong",
			})
			return
		}

		wordGuess := c.Param("word")
		if len(wordGuess) != 5 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Guess must be 5 letters long",
			})
			return
		}

		if logic.CheckWordExists(words, wordGuess) == false {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Word does not exist",
			})
			return
		}

		if attempts >= 6 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Out of tries",
			})
			return
		}

		letters := logic.MakeGuess(wordGuess, word)
		token := jwtService.GenerateToken(word, attempts+1)

		gameGuessResponse := GameGuessResponse{
			Token:   token,
			Letters: letters,
		}
		c.JSON(200, gameGuessResponse)
	}
}

func GameStart(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newGame newGamePayload
		if err := c.ShouldBindJSON(&newGame); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token := jwtService.GenerateToken(newGame.Word, 0)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
