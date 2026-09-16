package threads

import "github.com/bontanksakti84/threads-lead-radar/internal/leads"

func DetectCandidate(
	post Post,
	keywords []leads.Keyword,
) *Candidate {
	matches := MatchKeywords(post.Content, keywords)

	if len(matches) == 0 {
		return nil
	}

	return &Candidate{
		Post:    post,
		Matches: matches,
	}
}
