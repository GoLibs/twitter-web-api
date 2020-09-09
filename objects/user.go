package objects

type User struct {
	/*Id        int64*/
	FullName        string
	Username        string
	Biography       string
	PhotoSrc        string
	TwitsCount      int
	FollowersCount  int
	FollowingsCount int
	Twits           []Twit
}
