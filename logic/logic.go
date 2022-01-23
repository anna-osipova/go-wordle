package logic

import (
	"strings"
)

type Letter struct {
	Letter string `json:"letter"`
	Color  string `json:"color"`
}

func CheckWordExists(words []string, word string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

func CountExistingLetters(letters []Letter, letter string) int {
	count := 0
	for _, n := range letters {
		if n.Letter == letter {
			count++
		}
	}
	return count
}

func CountExactMatches(word string, guess_word string, letter string) int {
	count := 0
	for i, n := range word {
		l := string(n)
		if letter == l && l == string(guess_word[i]) {
			count++
		}
	}
	return count
}

func MakeGuess(wordGuess string, word string) []Letter {
	letters := make([]Letter, 0)
	for i, r := range wordGuess {
		letter := string(r)
		index := strings.Index(word, letter)
		color := "grey"
		// Do green first, then rest
		if letter == string(word[i]) {
			// Target word has the same letter in the same position
			color = "green"
		} else if index > -1 &&
			// Target word has more of the same letter than what has already been found
			strings.Count(word, letter) > CountExistingLetters(letters, letter) &&
			// Guess word has more of letter than there are exat matches
			strings.Count(wordGuess, letter) > CountExactMatches(word, wordGuess, letter) {
			color = "yellow"
		} else {
			color = "grey"
		}
		letters = append(letters, Letter{
			Color:  color,
			Letter: letter,
		})
	}
	return letters
}
