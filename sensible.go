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
	"flag"
	"encoding/json"
	"io/ioutil"
)

type Token struct {
	ConsumerKey string
	ConsumerSecret string
	AccessToken string
	AccessTokenSecret string
}

type Keyword struct {
	keywordMap map[string][]string
}
var api *anaconda.TwitterApi
//var keywordMap = make(map[string][]string)
var mode string

type Page struct {
	Title  string
	TechTweets []anaconda.Tweet
	PoliticsTweets []anaconda.Tweet
	TravelTweets []anaconda.Tweet
	OtherTweets []anaconda.Tweet
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	timelineTweets := getTimelineTweets()
	keywordMap := populateKeywordMap()
	classifiedTweets := classifyTweets(timelineTweets, keywordMap)
	p := &Page{Title: "Tech Tweets", TechTweets: classifiedTweets["tech"], PoliticsTweets: classifiedTweets["politics"], TravelTweets: classifiedTweets["travel"], OtherTweets: classifiedTweets["other"]}
	renderTemplate(w, "index", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func classifyTweets(timelineTweets []anaconda.Tweet, keywordMap map[string][]string) map[string][]anaconda.Tweet {
	classifiedTweets := make(map[string][]anaconda.Tweet)
	var techTweets []anaconda.Tweet
	var politicsTweets []anaconda.Tweet
	var travelTweets []anaconda.Tweet
	var otherTweets []anaconda.Tweet
	for _, tweet := range timelineTweets {
		if itIs(keywordMap["tech"], tweet) {
			techTweets = append(techTweets, tweet)
		} else if itIs(keywordMap["politics"], tweet) {
			politicsTweets = append(politicsTweets, tweet)
		} else if itIs(keywordMap["travel"], tweet) {
			travelTweets = append(travelTweets, tweet)
		} else {
			otherTweets = append(otherTweets, tweet)
		}
		
	}
	classifiedTweets["tech"] = techTweets
	classifiedTweets["politics"] = politicsTweets
	classifiedTweets["travel"] = travelTweets
	classifiedTweets["other"] = otherTweets
	return classifiedTweets
}

func itIs(keywords []string, tweet anaconda.Tweet) bool {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(tweet.Text), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func populateKeywordMap() map[string][]string {
	var keywordMap = make(map[string][]string)

	keywordMap["tech"] = []string{"golang", "ruby", "devs", "developers", "android", "ios", "programming", "code", "java", "coders", "developer", "fullstack", "full stack", "product", "hack", "hacker", "bug", "technology", "software", "mvc"}
	keywordMap["politics"] = []string{"modi", "congress", "bjp", "rahul gandhi", "manmohan singh", "narendra modi"}
	keywordMap["travel"] = []string{"travel","#lp"}
	return keywordMap
}


func getTimelineTweets() []anaconda.Tweet{
	v := url.Values{}
	v.Set("count", "200")
	if mode == "dev" {
		timelineTweets := getDummyTimeline()
		return timelineTweets
	} else {
		timelineTweets, err := api.GetHomeTimeline(v)
		if err != nil {
			fmt.Println(err)
		}
		return timelineTweets
	}

}

func getDummyTimeline() []anaconda.Tweet {
	var dummyTimeline = []anaconda.Tweet{}
	filename := "timeline.json"
	timeline, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error", err)
	}
	_ = json.Unmarshal(timeline, &dummyTimeline)
	return dummyTimeline
}

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	v := url.Values{}
	v.Set("count", "200")
	timelineTweets, _ := api.GetHomeTimeline(v)
	fmt.Println("time", timelineTweets)
	b, err := json.Marshal(timelineTweets)
	fmt.Println("json", b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := "timeline.json"
	ioutil.WriteFile(filename, b, 0600)
}

func getTokens() Token {
	var token Token 
	filename := "token.json"
	token_bytes, err := ioutil.ReadFile(filename)
	err = json.Unmarshal(token_bytes, &token)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return token
}

func main() {
	wordPtr := flag.String("mode", "", "which mode to run")
	flag.Parse()

	fmt.Println("word:", *wordPtr)

	if *wordPtr == "dev" {
		mode = "dev"
	}
	//populateKeywordMap()
	token := getTokens()
	api = anaconda.NewTwitterApi(token.AccessToken, token.AccessTokenSecret)
	anaconda.SetConsumerKey(token.ConsumerKey)
	anaconda.SetConsumerSecret(token.ConsumerSecret)
	cssHandler := http.FileServer(http.Dir("./css/"))
	jsHandler := http.FileServer(http.Dir("./js/"))
	imagesHandler := http.FileServer(http.Dir("./images/"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dump", dumpHandler)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
