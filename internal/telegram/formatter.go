package telegram

import (
	"fmt"
	"html"
	"strings"

	"github.com/bontanksakti84/threads-lead-radar/internal/ai"
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
	"github.com/bontanksakti84/threads-lead-radar/internal/scoring"
)

func FormatLeadAlert(
	post leads.Post,
	classification ai.Classification,
	score int,
	matchedKeywords []string,
) string {
	tier := scoring.GetTier(score)

	var keywords string
	if len(matchedKeywords) == 0 {
		keywords = "-"
	} else {
		items := make([]string, 0, len(matchedKeywords))

		for _, keyword := range matchedKeywords {
			items = append(
				items,
				"• "+html.EscapeString(keyword),
			)
		}

		keywords = strings.Join(items, "\n")
	}

	signals := make([]string, 0)

	if classification.NeedsDeveloper {
		signals = append(signals, "✓ Needs Developer")
	}

	if classification.HasBudget {
		signals = append(signals, "✓ Has Budget")
	}

	if classification.HasUrgency {
		signals = append(signals, "✓ Urgent")
	}

	if len(signals) == 0 {
		signals = append(signals, "• No explicit signals")
	}

	return fmt.Sprintf(
		`🔥 <b>NEW THREADS LEAD</b>

<b>Score:</b> %d/100
<b>Tier:</b> %s

👤 @%s

📂 <b>Category</b>
%s

🎯 <b>Intent</b>
%s

🧠 <b>AI Summary</b>
%s

🔎 <b>Matched Keywords</b>
%s

⚡ <b>Signals</b>
%s

🔗 <b>Open Post</b>
%s`,
		score,
		html.EscapeString(strings.ToUpper(string(tier))),
		html.EscapeString(post.Username),
		html.EscapeString(post.Category),
		html.EscapeString(post.Intent),
		html.EscapeString(post.AISummary),
		keywords,
		strings.Join(signals, "\n"),
		html.EscapeString(post.URL),
	)
}
