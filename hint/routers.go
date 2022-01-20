package hint

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func HintRegister(router *gin.RouterGroup) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("maxfiveletters", maxFiveLetters)
	}
	router.POST("/", HintFind)
}

var maxFiveLetters validator.Func = func(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) <= 5
}

type WordHelp struct {
	Include  string `json:"include" binding:"required,maxfiveletters"`
	Exclude  string `json:"exclude" binding:"required"`
	Template string `json:"template" binding:"required,len=5"`
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

func HintFind(c *gin.Context) {
	var word_help WordHelp
	if err := c.ShouldBindJSON(&word_help); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	words := c.MustGet("word_list").([]string)
	var matching_words []string
	for _, word := range words {
		if IsMatch(&word_help, word) == true {
			matching_words = append(matching_words, word)
		}
	}

	c.JSON(http.StatusOK, gin.H{"matches": matching_words})
}
