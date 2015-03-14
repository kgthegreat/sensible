package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"html/template"
	"log"
	"net/http"
	"net/url"
	_ "regexp"
	"strings"
)

const (
	consumerKey       = "6CfGgaWC42SlfSHJ0AusSH64j"
	consumerSecret    = "WGGQwcaDZ6mMyHB9WoZgTyJFczHzy865GPiCL5EeukFoPanmvH"
	accessToken       = "14741206-HBQKBY4WHW93EWVv1bgLhBlLrrRqMwgdTpEoDreiI"
	accessTokenSecret = "NgvGtxH5UaGiX3hbcXKbjuB5LbVS4WXWyg960XC4w0Ksn"
)

var api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
//var templates = template.Must(template.ParseFiles("index.html"))
var timelineTweets []anaconda.Tweet

type Page struct {
	Title  string
	TechTweets []anaconda.Tweet
	PoliticsTweets []anaconda.Tweet
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
//	classifiedTweets := getTechTweets
	classifiedTweets := classifyTweets(timelineTweets)
	p := &Page{Title: "Tech Tweets", TechTweets: classifiedTweets["tech"], PoliticsTweets: classifiedTweets["politics"]}
	t, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func classifyTweets(timelineTweets []anaconda.Tweet) map[string][]anaconda.Tweet {
	classifiedTweets := make(map[string][]anaconda.Tweet)
	var techTweets []anaconda.Tweet
	var politicsTweets []anaconda.Tweet
	for _, tweet := range timelineTweets {
		if itIs("tech", tweet) {
			techTweets = append(techTweets, tweet)
		}
		if itIs("politics", tweet) {
			politicsTweets = append(politicsTweets, tweet)
		}
	}
	classifiedTweets["tech"] = techTweets
	classifiedTweets["politics"] = politicsTweets
	return classifiedTweets
}

func itIs(context string, tweet anaconda.Tweet) bool {
	keywordMap := make(map[string][]string)
	keywordMap["tech"] = []string{"golang", "ruby", "devs", "developers", "android", "ios", "programming", "code", "java", "coders", "developer", "fullstack", "full stack", "product", "hack", "hacker", "bug", "technology", "software"}
	keywordMap["politics"] = []string{"modi", "congress", "bjp", "rahul gandhi", "manmohan singh", "narendra modi", "jashn"}
	for _, keyword := range keywordMap[context] {
		if strings.Contains(tweet.Text, keyword) {
			return true
		}
	}
	return false
}



func getTimelineTweets() []anaconda.Tweet{
	v := url.Values{}
	v.Set("count", "200")
	timelineTweets, err := api.GetHomeTimeline(v)
	if err != nil {
		fmt.Println(err)
	}
	return timelineTweets
}

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	timelineTweets = getTimelineTweets()
	http.HandleFunc("/", indexHandler)
//	http.HandleFunc("/politics", politicsHandler)
//	http.HandleFunc("/sports", sportsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
