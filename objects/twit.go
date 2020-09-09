package objects

type Twit struct {
	Id                int
	IsReTweet         bool
	OriginalTweetUser RetweetUser
	IsReply           bool
	RepliedUsername   []string
	Time              string
	Text              string
	Hashtag           []string
	Link              string
}
