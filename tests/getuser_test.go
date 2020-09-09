package examples

import (
	"testing"

	twitter_web_api "github.com/golibs/twitter-web-api"
)

func TestGetUser(t *testing.T) {
	c := twitter_web_api.NewTwitterClient()
	user, err := c.GetUserByUsername("jack")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(user)
}
