package main

import (
	"time"
)

type Deleter struct {
	// delete tweets before this time
	Before time.Time

	// delete tweets before id
	BeforeId uint64

	// delete retweets
	Retweets bool

	// delete replies
	Replies bool
}

// Delete returns true if tweet should be deleted
func (d Deleter) Delete(t *Tweet) bool {
	if d.Before.After(t.Timestamp) || t.Id < d.BeforeId {
		if len(t.ReplyId) > 0 && !d.Replies {
			return false
		} else if len(t.RetweetId) > 0 && !d.Retweets {
			return false
		}

		return true
	}

	return false
}
