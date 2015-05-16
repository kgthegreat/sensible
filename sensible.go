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

const (
	consumerKey       = "6CfGgaWC42SlfSHJ0AusSH64j"
	consumerSecret    = "WGGQwcaDZ6mMyHB9WoZgTyJFczHzy865GPiCL5EeukFoPanmvH"
	accessToken       = "14741206-HBQKBY4WHW93EWVv1bgLhBlLrrRqMwgdTpEoDreiI"
	accessTokenSecret = "NgvGtxH5UaGiX3hbcXKbjuB5LbVS4WXWyg960XC4w0Ksn"
)

var api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
var keywordMap = make(map[string][]string)
var mode = ""

type Page struct {
	Title  string
	TechTweets []anaconda.Tweet
	PoliticsTweets []anaconda.Tweet
	TravelTweets []anaconda.Tweet
	OtherTweets []anaconda.Tweet
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var timelineTweets = getTimelineTweets()
	classifiedTweets := classifyTweets(timelineTweets)
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

func classifyTweets(timelineTweets []anaconda.Tweet) map[string][]anaconda.Tweet {
	classifiedTweets := make(map[string][]anaconda.Tweet)
	var techTweets []anaconda.Tweet
	var politicsTweets []anaconda.Tweet
	var travelTweets []anaconda.Tweet
	var otherTweets []anaconda.Tweet
	for _, tweet := range timelineTweets {
		if itIs("tech", tweet) {
			techTweets = append(techTweets, tweet)
		} else if itIs("politics", tweet) {
			politicsTweets = append(politicsTweets, tweet)
		} else if itIs("travel", tweet) {
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

func itIs(context string, tweet anaconda.Tweet) bool {
	populateKeywordMap()
	for _, keyword := range keywordMap[context] {
		if strings.Contains(strings.ToLower(tweet.Text), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func populateKeywordMap() {
	keywordMap["tech"] = []string{"golang", "ruby", "devs", "developers", "android", "ios", "programming", "code", "java", "coders", "developer", "fullstack", "full stack", "product", "hack", "hacker", "bug", "technology", "software", "mvc"}
	keywordMap["politics"] = []string{"modi", "congress", "bjp", "rahul gandhi", "manmohan singh", "narendra modi"}
	keywordMap["travel"] = []string{"travel","#lp"}

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

//	renderTemplate(w, "index", p)
}

func main() {
	wordPtr := flag.String("mode", "", "which mode to run")
	flag.Parse()

	fmt.Println("word:", *wordPtr)

	if *wordPtr == "dev" {
		mode = "dev"
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
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
