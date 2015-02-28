package main
import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	_ "html/template"
	"github.com/ChimeraCoder/anaconda"
)
const (
	consumerKey = "6CfGgaWC42SlfSHJ0AusSH64j"
	consumerSecret = "WGGQwcaDZ6mMyHB9WoZgTyJFczHzy865GPiCL5EeukFoPanmvH"
	accessToken = "14741206-HBQKBY4WHW93EWVv1bgLhBlLrrRqMwgdTpEoDreiI" 
	accessTokenSecret = "NgvGtxH5UaGiX3hbcXKbjuB5LbVS4WXWyg960XC4w0Ksn"
)

	var api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	v := url.Values{}
	v.Set("count", "20")
	timelineTweets, err := api.GetHomeTimeline(v)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	techTweets := getTechTweets(timelineTweets)
	_ = techTweets
	
}

func getTechTweets(timelineTweets []anaconda.Tweet) []anaconda.Tweet {
	techKeywords := []string{"android", "golang"}
	var techTweets []anaconda.Tweet
	for _ , tweet := range timelineTweets {
		for _ , kwd := range techKeywords {
			if strings.Contains(tweet.Text, kwd) {
				techTweets = append(techTweets, tweet)
				break
			}
		}
	}
	return techTweets
}

 
func publishSearchResults(api *anaconda.TwitterApi, w http.ResponseWriter, v url.Values) {
	searchResult, err := api.GetSearch("gophercon", v)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	for _ , tweet := range searchResult.Statuses {
		fmt.Fprintln(w, tweet.Text)
	}

}

func getTimeline(api *anaconda.TwitterApi, w http.ResponseWriter, v url.Values) {
	timeline, err := api.GetHomeTimeline(v)
	techKeywords := []string{"golang", "ruby", "devs", "developers", "android", "ios", "app", "programming", "code", "java"}
	var techTweets []anaconda.Tweet
	var nonTechTweets []anaconda.Tweet
	if err != nil {
		fmt.Fprintln(w, err)
	}
	for _ , tweet := range timeline {
		for _ , kwd := range techKeywords {
			if strings.Contains(tweet.Text, kwd) {
				techTweets = append(techTweets, tweet)
				break
			} else {
				nonTechTweets = append(nonTechTweets, tweet)
				break
			}
		}
		
	}
	fmt.Fprintln(w, "-------------------Tech Tweets--------------------")
	for _ , t := range techTweets {
		fmt.Fprintln(w, t.Text)
	}
	fmt.Fprintln(w, "-------------------NonTech Tweets--------------------")
	for _ , nt := range nonTechTweets {
		fmt.Fprintln(w, nt.Text)
	}

}

func main() {
 	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
