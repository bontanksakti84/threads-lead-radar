package threads

import (
	"github.com/bontanksakti84/threads-lead-radar/internal/leads"
)

func ToLeadPost(post Post) leads.Post {
	return leads.Post{
		ExternalID: post.ID,
		Username:   post.Username,
		URL:        post.URL,
		Content:    post.Content,
	}
}
