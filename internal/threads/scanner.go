package threads

import (
	"context"
	"log"
	"time"

	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
	"github.com/bontanksakti84/threads-lead-radar/internal/scoring"
)

type Scanner struct {
	client      Client
	leadService *leads.Service
	keywordRepo *leads.KeywordRepository
	classifier  ai.Classifier
	scanRunRepo *leads.ScanRunRepository
}

type ScanStats struct {
	KeywordsScanned int
	PostsFound      int
	CandidatesFound int
	NewPosts        int
	ExistingPosts   int
	GroqCalls       int
}

func NewScanner(
	client Client,
	leadService *leads.Service,
	keywordRepo *leads.KeywordRepository,
	classifier ai.Classifier,
	scanRunRepo *leads.ScanRunRepository,
) *Scanner {
	return &Scanner{
		client:      client,
		leadService: leadService,
		keywordRepo: keywordRepo,
		classifier:  classifier,
		scanRunRepo: scanRunRepo,
	}
}

func (s *Scanner) Scan(ctx context.Context, keyword string) error {
	posts, err := s.client.SearchPosts(ctx, keyword)
	if err != nil {
		return err
	}

	keywords, err := s.keywordRepo.GetActiveKeywords(ctx)
	if err != nil {
		return err
	}

	log.Printf(
		"Found %d posts for keyword: %s",
		len(posts),
		keyword,
	)

	for _, post := range posts {
		candidate := DetectCandidate(post, keywords)

		if candidate == nil {
			log.Printf(
				"Skipping non-candidate post: %s",
				post.ID,
			)
			continue
		}

		classification, err := s.classifier.Classify(
			ctx,
			post.Content,
		)
		if err != nil {
			return err
		}

		classification = ai.NormalizeClassification(
			classification,
		)
		if err := ai.ValidateClassification(classification); err != nil {
			return err
		}

		matchedKeywords := make(
			[]string,
			0,
			len(candidate.Matches),
		)

		for _, match := range candidate.Matches {
			matchedKeywords = append(
				matchedKeywords,
				match.Keyword,
			)
		}

		leadPost := ToLeadPost(post)

		leadPost.Category = classification.Category
		leadPost.Intent = classification.Intent
		leadPost.LeadScore = scoring.Calculate(
			scoring.Input{
				Classification: classification,
				KeywordMatches: len(candidate.Matches),
			},
		)
		leadPost.HasBudget = classification.HasBudget
		leadPost.HasUrgency = classification.HasUrgency
		leadPost.NeedsDeveloper = classification.NeedsDeveloper
		leadPost.AISummary = classification.Summary
		leadPost.MatchedKeywords = matchedKeywords

		if err := s.leadService.CreatePost(
			ctx,
			leadPost,
		); err != nil {
			return err
		}

		log.Printf(
			"Saved candidate: %s | category=%s | intent=%s | needs_developer=%v | budget=%v | urgency=%v | keywords=%d | score=%d",
			post.ID,
			classification.Category,
			classification.Intent,
			classification.NeedsDeveloper,
			classification.HasBudget,
			classification.HasUrgency,
			len(candidate.Matches),
			leadPost.LeadScore,
		)
	}

	return nil
}

func (s *Scanner) ScanAll(ctx context.Context) error {
	startedAt := time.Now()

	keywords, err := s.keywordRepo.GetActiveKeywords(ctx)
	if err != nil {
		return err
	}

	stats := ScanStats{
		KeywordsScanned: len(keywords),
	}

	scanRunID, err := s.scanRunRepo.Start(
		ctx,
		stats.KeywordsScanned,
	)
	if err != nil {
		return err
	}

	processed := make(map[string]bool)

	failScan := func(scanErr error) error {
		if err := s.scanRunRepo.Fail(
			ctx,
			scanRunID,
			scanErr.Error(),
		); err != nil {
			log.Printf(
				"Failed to update scan run: %v",
				err,
			)
		}

		return scanErr
	}

	for _, keyword := range keywords {
		log.Printf(
			"Scanning keyword: %s",
			keyword.Keyword,
		)

		posts, err := s.client.SearchPosts(
			ctx,
			keyword.Keyword,
		)
		if err != nil {
			return failScan(err)
		}

		stats.PostsFound += len(posts)

		for _, post := range posts {
			if processed[post.ID] {
				log.Printf(
					"Skipping duplicate post in current scan: %s",
					post.ID,
				)
				continue
			}

			candidate := DetectCandidate(
				post,
				keywords,
			)

			if candidate == nil {
				log.Printf(
					"Skipping non-candidate post: %s",
					post.ID,
				)
				continue
			}

			processed[post.ID] = true
			stats.CandidatesFound++

			existingPost, err := s.leadService.GetPostByExternalID(
				ctx,
				post.ID,
			)
			if err != nil {
				return failScan(err)
			}

			if existingPost != nil {
				matchedKeywords := make(
					[]string,
					0,
					len(candidate.Matches),
				)

				for _, match := range candidate.Matches {
					matchedKeywords = append(
						matchedKeywords,
						match.Keyword,
					)
				}

				if err := s.leadService.MergeMatchedKeywords(
					ctx,
					post.ID,
					matchedKeywords,
				); err != nil {
					return failScan(err)
				}

				stats.ExistingPosts++

				log.Printf(
					"Skipping existing post: %s | Groq skipped | new_keywords=%d",
					post.ID,
					len(matchedKeywords),
				)

				continue
			}

			// New post: classify with Groq.
			stats.GroqCalls++

			classification, err := s.classifier.Classify(
				ctx,
				post.Content,
			)
			if err != nil {
				return failScan(err)
			}

			classification = ai.NormalizeClassification(
				classification,
			)

			if err := ai.ValidateClassification(
				classification,
			); err != nil {
				return failScan(err)
			}

			matchedKeywords := make(
				[]string,
				0,
				len(candidate.Matches),
			)

			for _, match := range candidate.Matches {
				matchedKeywords = append(
					matchedKeywords,
					match.Keyword,
				)
			}

			leadPost := ToLeadPost(post)

			leadPost.Category = classification.Category
			leadPost.Intent = classification.Intent

			leadPost.LeadScore = scoring.Calculate(
				scoring.Input{
					Classification: classification,
					KeywordMatches: len(candidate.Matches),
				},
			)

			leadPost.HasBudget = classification.HasBudget
			leadPost.HasUrgency = classification.HasUrgency
			leadPost.NeedsDeveloper = classification.NeedsDeveloper
			leadPost.AISummary = classification.Summary
			leadPost.MatchedKeywords = matchedKeywords

			if err := s.leadService.CreatePost(
				ctx,
				leadPost,
			); err != nil {
				return failScan(err)
			}

			stats.NewPosts++

			log.Printf(
				"Saved new candidate: %s | category=%s | intent=%s | score=%d",
				post.ID,
				classification.Category,
				classification.Intent,
				leadPost.LeadScore,
			)
		}
	}

	if err := s.scanRunRepo.Complete(
		ctx,
		scanRunID,
		stats.PostsFound,
		stats.NewPosts,
	); err != nil {
		return err
	}

	duration := time.Since(startedAt)

	log.Printf(
		"=== SCAN RESULT === keywords=%d posts=%d candidates=%d new=%d existing=%d groq=%d duration=%s",
		stats.KeywordsScanned,
		stats.PostsFound,
		stats.CandidatesFound,
		stats.NewPosts,
		stats.ExistingPosts,
		stats.GroqCalls,
		duration.Round(time.Millisecond),
	)

	return nil
}
