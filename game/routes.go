package game

import (
	"fmt"
	"net/http"

	"github.com/anna-osipova/go-wordle/db"
	"github.com/anna-osipova/go-wordle/logic"
	"github.com/anna-osipova/go-wordle/models"
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
			c.AbortWithStatus(http.StatusBadRequest)
		}
		if token.Valid {
			claims := token.Claims.(*service.CustomClaims)
			c.Set("session_id", claims.SessionId)
		} else {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func GameRegister(router *gin.RouterGroup, dbInstance db.Database) {
	var jwtService service.JWTService = service.JWTAuthService()
	router.POST("/new", GameNew(jwtService, dbInstance))

	router.Use(AuthorizeJWT(jwtService))
	router.POST("/start", GameStart)
	router.POST("/guess/:word", GameGuess(jwtService, dbInstance))
}

type GameGuessResponse struct {
	Letters []logic.Letter `json:"letters"`
}

func GameGuess(jwtService service.JWTService, dbInstance db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.MustGet("session_id").(string)
		session, err := dbInstance.GetSessionById(sessionId)
		if err != nil {
			c.AbortWithStatus(401)
		}

		words := c.MustGet("word_list").([]string)
		if len(session.Word) != 5 {
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

		if session.Attempts >= 6 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Out of tries",
			})
			return
		}

		letters := logic.MakeGuess(wordGuess, session.Word)
		dbInstance.UpdateSessionAttemptCount(sessionId, session.Attempts+1)

		gameGuessResponse := GameGuessResponse{
			Letters: letters,
		}
		c.JSON(200, gameGuessResponse)
	}
}

func GameNew(jwtService service.JWTService, dbInstance db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newGame newGamePayload
		if err := c.ShouldBindJSON(&newGame); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		session := &models.Session{Word: newGame.Word, Attempts: 0}
		if err := dbInstance.CreateSession(session); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		token := jwtService.GenerateToken(session)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func GameStart(c *gin.Context) {
	c.Status(http.StatusOK)
}
