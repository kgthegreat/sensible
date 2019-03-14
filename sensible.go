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
	TechKeywords []string
	PoliticsKeywords []string
	TravelKeywords []string
}
var api *anaconda.TwitterApi
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
	keywordStore := populateKeywordStore()
	classifiedTweets := classifyTweets(timelineTweets, keywordStore)
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

func classifyTweets(timelineTweets []anaconda.Tweet, keywordStore Keyword) map[string][]anaconda.Tweet {
	classifiedTweets := make(map[string][]anaconda.Tweet)
	var techTweets []anaconda.Tweet
	var politicsTweets []anaconda.Tweet
	var travelTweets []anaconda.Tweet
	var otherTweets []anaconda.Tweet
	for _, tweet := range timelineTweets {
		if itIs(keywordStore.TechKeywords, tweet) {
			techTweets = append(techTweets, tweet)
		} else if itIs(keywordStore.PoliticsKeywords, tweet) {
			politicsTweets = append(politicsTweets, tweet)
		} else if itIs(keywordStore.TravelKeywords, tweet) {
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
		if strings.Contains(strings.ToLower(tweet.FullText), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func populateKeywordStore() Keyword {
	var keywordStore Keyword
	filename := "keyword.json"
	keyword_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error", err)
	}
	
	_ = json.Unmarshal(keyword_bytes, &keywordStore)
	return keywordStore
}


func getTimelineTweets() []anaconda.Tweet{
	v := url.Values{}
	v.Set("count", "200")
	v.Set("tweet_mode", "extended")
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

func signinHandler(w http.ResponseWriter, r *http.Request) {
	callback := "http://" + r.Host + "/callback"
	authURL, tempCred, err := anaconda.AuthorizationURL(callback)
	if err != nil {
		http.Error(w, "Error getting temp cred, "+err.Error(), 500)
		return 
	}
	s := session.Get(r)
	s[tempCredKey] = tempCred
	if err := session.Save(w, r, s); err != nil {
		http.Error(w, "Error saving session, " +err.Error(), 500)
		return
	}
	http.Redirect(w, r, authURL, 302)
}

func main() {
	wordPtr := flag.String("mode", "", "which mode to run")
	flag.Parse()

	fmt.Println("word:", *wordPtr)

	if *wordPtr == "dev" {
		mode = "dev"
	}
	token := getTokens()
	api = anaconda.NewTwitterApiWithCredentials(token.AccessToken, token.AccessTokenSecret, token.ConsumerKey, token.ConsumerSecret)
	anaconda.SetConsumerKey(token.ConsumerKey)
	anaconda.SetConsumerSecret(token.ConsumerSecret)
	cssHandler := http.FileServer(http.Dir("./css/"))
	jsHandler := http.FileServer(http.Dir("./js/"))
	imagesHandler := http.FileServer(http.Dir("./images/"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/dump", dumpHandler)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
