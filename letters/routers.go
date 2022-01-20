package letters

import (
	"sort"

	"github.com/gin-gonic/gin"
)

func LettersRegister(router *gin.RouterGroup) {
	router.GET("/", LettersStats)
}

type LetterStat struct {
	Letter string `json:"letter" binding:"required"`
	Count  int    `json:"count" binding:"required"`
}

type LetterStatList []LetterStat

func (p LetterStatList) Len() int           { return len(p) }
func (p LetterStatList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p LetterStatList) Less(i, j int) bool { return p[i].Count > p[j].Count }

func LettersStats(c *gin.Context) {
	words := c.MustGet("word_list").([]string)
	m := make(map[string]int)
	var r rune = 97
	for r <= 122 {
		m[string(r)] = 0
		r++
	}
	for _, word := range words {
		for _, l := range word {
			m[string(l)]++
		}
	}
	list := make(LetterStatList, len(m))
	i := 0
	for k, v := range m {
		list[i] = LetterStat{k, v}
		i++
	}

	sort.Sort(list)

	c.JSON(200, gin.H{"result": list})
}
