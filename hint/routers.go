package hint

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func HintRegister(router *gin.RouterGroup) {
	router.POST("/", HintFind)
}

type WordHelp struct {
	Include string `json:"include" binding:"required"`
	Exclude string `json:"exclude" binding:"required"`
	Template string `json:"template" binding:"required"`
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
	c.BindJSON(&word_help)

	words := c.MustGet("word_list").([]string)
	var matching_words []string
	for _, word := range words {
		if IsMatch(&word_help, word) == true {
			matching_words = append(matching_words, word)
		}
	}

	c.JSON(200, gin.H{ "matches": matching_words })
}
