package app

import (
	"slices"
	"strings"

	"github.com/cubedhuang/lipu-lili/internal/models"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func defaultSort(a, b models.WordData) int {
	if u := models.UsageCategories[b.UsageCategory] - models.UsageCategories[a.UsageCategory]; u != 0 {
		return u
	}

	return strings.Compare(a.Word, b.Word)
}

type wordScore struct {
	word  *models.WordData
	score int
}

func (data *WordsStore) search(query string) []models.WordData {
	data.mu.RLock()
	defer data.mu.RUnlock()

	if query == "" {
		return data.words
	}

	scores := make([]wordScore, 0, len(data.words))

	for _, word := range data.words {
		score := fuzzy.RankMatchNormalizedFold(query, word.Word)

		if score != -1 {
			score = max(16-score, 0) * 100
		} else {
			score = 0
		}

		if d := fuzzy.RankMatchNormalizedFold(query, word.Translations["en"].Definition); d != -1 {
			score += max(16-d, 0) * 10
		}

		kuScore := 0
		for ku := range word.KuData {
			if k := fuzzy.RankMatchNormalizedFold(query, ku); k != -1 {
				kuScore = max(kuScore, 10-k, 0)
			}
		}
		kuScore *= 10

		score += kuScore

		if score > 0 {
			scores = append(scores, wordScore{word: &word, score: score})
		}
	}

	slices.SortFunc(scores, func(a, b wordScore) int {
		if a.score == b.score {
			return defaultSort(*a.word, *b.word)
		}

		return b.score - a.score
	})

	results := make([]models.WordData, 0, len(scores))
	for _, score := range scores {
		results = append(results, *score.word)
	}

	return results
}
