package twitter_web_api

import (
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"strings"

	"github.com/golibs/twitter-web-api/objects"
	"github.com/golibs/twitter-web-api/urls"
)

type Client struct {
	httpClient *http.Client
	jar        *cookiejar.Jar
}

func NewTwitterClient() (c *Client) {
	j, _ := cookiejar.New(nil)
	c = &Client{httpClient: &http.Client{Jar: j}}
	c.jar = j
	c.setCookies()
	return
}

func (c *Client) nojsRouter(address string) (respStr string, err error) {
	var resp *http.Response
	var req *http.Request

	req, err = http.NewRequest("POST", fmt.Sprintf(urls.NoJSRouter, address), nil)
	if err != nil {
		return
	}
	req.Header.Set("referer", fmt.Sprintf(urls.User, address))
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36")
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respStr = string(b)
	return
}

func (c *Client) setCookies() (err error) {
	// var resp *http.Response
	var req *http.Request

	req, err = http.NewRequest("GET", urls.MainMobile, nil)
	if err != nil {
		return
	}

	// req.Header.Set("referer", urls.MainMobile)
	req.Header.Set("user-agent", " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36")
	_, err = c.httpClient.Do(req)
	if err != nil {
		return
	}
	return
}

func (c *Client) GetUserByUsername(username string) (u *objects.User, err error) {
	var respStr string
	respStr, err = c.nojsRouter(username)
	if err != nil {
		return
	}

	profileIndex := strings.Index(respStr, `<table class="profile-details">`)
	if profileIndex == -1 {
		err = errors.New("not_found")
		return
	}

	profile := respStr[profileIndex:]
	profile = profile[:strings.Index(respStr, "</table>")]

	photoR := regexp.MustCompile(`img alt.+src="(.+)"`)
	photo := photoR.FindAllStringSubmatch(profile, -1)
	u = &objects.User{}
	if len(photo) > 0 {
		u.PhotoSrc = photo[0][1]
	}

	nameR := regexp.MustCompile(`<div class="fullname">(.*)[\s\S]+?</div>`)
	name := nameR.FindAllStringSubmatch(profile, -1)
	if len(name) > 0 {
		u.FullName = name[0][1]
	}

	usernameR := regexp.MustCompile(`<span class="screen-name">(.*)</span>`)
	uName := usernameR.FindAllStringSubmatch(profile, -1)
	if len(name) > 0 {
		u.Username = uName[0][1]
	}

	bioR := regexp.MustCompile(`<div class="dir-ltr" dir="ltr">\s(.*)\s+</div>`)
	bio := bioR.FindAllStringSubmatch(profile, -1)
	if len(bio) > 0 {
		u.Biography = bio[0][1]
		u.Biography = html.UnescapeString(strings.TrimSpace(u.Biography))
	}

	profileStatsIndex := strings.Index(respStr, `<table class="profile-stats">`)
	if profileStatsIndex != -1 {
		profileStats := respStr[profileStatsIndex:]
		profileStats = profileStats[:strings.Index(profileStats, "</table>")]

		statsR := regexp.MustCompile(`<div class="statnum">([\d,]+)</div>`)
		stats := statsR.FindAllStringSubmatch(profileStats, -1)

		if len(stats) == 3 {
			u.TwitsCount, _ = strconv.Atoi(strings.ReplaceAll(stats[0][1], ",", ""))
			u.FollowingsCount, _ = strconv.Atoi(strings.ReplaceAll(stats[1][1], ",", ""))
			u.FollowersCount, _ = strconv.Atoi(strings.ReplaceAll(stats[2][1], ",", ""))

		}
	}

	twitsHtmlR := regexp.MustCompile(`<table class="tweet[\s\S]+?</table>`)
	twitsHtml := twitsHtmlR.FindAllString(respStr, -1)
	var twits []objects.Twit
	if len(twitsHtml) > 0 {
		htmls := ""
		for _, twitHtml := range twitsHtml {
			twit := objects.Twit{}

			twitIdR := regexp.MustCompile(`<table class="tweet[\s\S]+?href="?/.+status/(\d+)`)
			twitId := twitIdR.FindAllStringSubmatch(twitHtml, -1)
			if len(twitId) == 1 {
				twit.Id, _ = strconv.Atoi(twitId[0][1])
				twit.Link = fmt.Sprintf(urls.Status, username, twit.Id)
			}

			timestampR := regexp.MustCompile(`<td class="timestamp">[\s\S]+?">(.+)</a>`)
			timestamp := timestampR.FindAllStringSubmatch(twitHtml, -1)
			if len(timestamp) == 1 {
				twit.Time = timestamp[0][1]
			}

			if strings.Index(twitHtml, `<tr class="tweet-content`) != -1 {
				twit.IsReTweet = true

				retweetUserInfo := twitHtml[strings.Index(twitHtml, `<td class="user-info">`):]
				retweetUserInfo = retweetUserInfo[:strings.Index(retweetUserInfo, "</td>")]

				fullnameR := regexp.MustCompile(`<strong class="fullname">(.+)</strong>`)
				fullname := fullnameR.FindAllStringSubmatch(retweetUserInfo, -1)
				ru := objects.RetweetUser{}
				if len(fullname) > 0 {
					ru.FullName = fullname[0][1]
				}

				rUsernameR := regexp.MustCompile(`<div class="username">[\s\S]+?<span.+</span>(.+)[\s\S]+</div>`)
				ruUsername := rUsernameR.FindAllStringSubmatch(retweetUserInfo, -1)
				if len(ruUsername) > 0 {
					ru.Username = ruUsername[0][1]
				}

				profilePicR := regexp.MustCompile(`<td class="avatar"[\s\S]+?<img[\s\S]+?src="(.+)"`)
				profilePic := profilePicR.FindAllStringSubmatch(twitHtml, -1)
				if len(profilePic) > 0 {
					ru.PhotoSrc = profilePic[0][1]
				}
				twit.OriginalTweetUser = ru
			}

			if replyIndex := strings.Index(twitHtml, `<div class="tweet-reply-context username">`); replyIndex != -1 {
				twit.IsReply = true

				reply := twitHtml[replyIndex:]
				reply = twitHtml[:strings.Index(reply, "</div>")]

				replyToUsersR := regexp.MustCompile(`<a href=".+@(.+)</a>`)
				replyToUsers := replyToUsersR.FindAllStringSubmatch(twitHtml, -1)
				if len(replyToUsers) > 0 {
					for _, user := range replyToUsers {
						twit.RepliedUsername = append(twit.RepliedUsername, user[1])
					}
				}
			}

			twitTextR := regexp.MustCompile(`<div class="tweet-text"[\s\S]+?<div class="(dir-rtl|dir-ltr)" dir="(rtl|ltr)">(.+)[\s\S]?</div>`)
			twitText := twitTextR.FindAllStringSubmatch(twitHtml, -1)
			if len(twitText) == 1 {
				text := html.UnescapeString(twitText[0][len(twitText[0])-1])

				hashTagsR := regexp.MustCompile(`<a.+?>(.+?)</a>`)
				hashtags := hashTagsR.FindAllStringSubmatch(text, -1)
				for _, hashtag := range hashtags {
					text = strings.ReplaceAll(text, hashtag[0], hashtag[1])
					if !strings.Contains(hashtag[0], "hashtag") {
						continue
					}
					twit.Hashtag = append(twit.Hashtag, strings.ReplaceAll(hashtag[1], "#", ""))
				}
				twit.Text = text
			}

			htmls += twitHtml
			twits = append(twits, twit)
		}
		// ioutil.WriteFile("twits.html", []byte(htmls), os.ModePerm)
	}
	u.Twits = twits
	// pretty.Println(u)
	// ioutil.WriteFile("username.html", []byte(respStr), os.ModePerm)
	return
}
