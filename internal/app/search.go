package app

import (
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

func (data *WordsStore) search(query string) []models.WordData {
	data.mu.RLock()
	defer data.mu.RUnlock()

	results := make([]models.WordData, 0, len(data.words))
	for _, word := range data.words {
		// if strings.Contains(word.Word, query) {
		if fuzzy.Match(query, word.Word) {
			results = append(results, word)
		}
	}

	return results
}
