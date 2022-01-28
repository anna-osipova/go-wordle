package game

import (
	"fmt"
	"net/http"

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

	router.Use(AuthorizeJWT(jwtService))
	router.POST("/start", GameStart)
	router.POST("/guess/:word", GameGuess(jwtService))
}

type GameGuessResponse struct {
	Letters []Letter `json:"letters"`
}

func GameGuess(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
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

		if session.Attempts >= 6 {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{
				Message:   "Out of tries",
				ErrorCode: "NO_TRIES",
			})
			return
		}

		processingError := common.ErrorResponse{
			Message:   "Some issue",
			ErrorCode: "ERROR",
		}

		letters := MakeGuess(wordGuess, session.Word)
		// _, err = dbInstance.UpdateSessionAttemptCount(sessionId, session.Attempts+1)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, processingError)
			return
		}

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

		gameGuessResponse := GameGuessResponse{
			Letters: letters,
		}
		c.JSON(200, gameGuessResponse)
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
		session := &service.Session{Word: newGame.Word, Attempts: 0}
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

func GameStart(c *gin.Context) {
	c.Status(http.StatusOK)
}
