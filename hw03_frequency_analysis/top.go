package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type PopularWords struct {
	Word  string
	Count int
}

func Top10(s string) []string {
	if s == "" {
		return nil
	}

	re := regexp.MustCompile(`[\t\n\r]+`)
	strValue := re.ReplaceAllString(s, " ")

	strSlices := strings.Split(strValue, " ")

	words := CountWords(strSlices)
	sort.Slice(words, func(i, j int) bool {
		if words[i].Count != words[j].Count {
			return words[i].Count > words[j].Count
		}
		return words[i].Word < words[j].Word
	})

	var result []string
	maxWords := 10
	for i := 0; i < len(words) && i < maxWords; i++ {
		result = append(result, words[i].Word)
	}

	return result
}

func CountWords(words []string) []PopularWords {
	wordMap := make(map[string]int)
	for _, word := range words {
		wordMap[word]++
	}

	var result []PopularWords
	for word, count := range wordMap {
		if word != "" {
			result = append(result, PopularWords{
				Word:  word,
				Count: count,
			})
		}
	}

	return result
}
