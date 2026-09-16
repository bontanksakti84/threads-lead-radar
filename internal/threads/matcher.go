package threads

import (
	"strings"

	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
)

type KeywordMatch struct {
	Keyword  string
	Category string
}

func MatchKeywords(
	content string,
	keywords []leads.Keyword,
) []KeywordMatch {
	content = strings.ToLower(content)

	var matches []KeywordMatch

	for _, keyword := range keywords {
		target := strings.ToLower(
			strings.TrimSpace(keyword.Keyword),
		)

		if target == "" {
			continue
		}

		if strings.Contains(content, target) {
			matches = append(matches, KeywordMatch{
				Keyword:  keyword.Keyword,
				Category: keyword.Category,
			})
		}
	}

	return matches
}
