package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	_ "regexp"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/sessions"
)

type Token struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type Keyword struct {
	TechKeywords     []string
	PoliticsKeywords []string
	TravelKeywords   []string
}

var api *anaconda.TwitterApi
var mode string

var store = sessions.NewCookieStore([]byte("asdaskdhasdhgsajdgasdsadksakdhasidoajsdousdasf"))

//var store = sessions.NewCookieStore([]byte(storeGUID))

// Session state keys.
const (
	tempCredKey  = "tempCred"
	tokenCredKey = "tokenCred"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var signinOAuthClient oauth.Client

type Page struct {
	Title          string
	TechTweets     []anaconda.Tweet
	PoliticsTweets []anaconda.Tweet
	TravelTweets   []anaconda.Tweet
	OtherTweets    []anaconda.Tweet
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, "twit")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tokenCred, ok := s.Values[tokenCredKey].(oauth.Credentials)

	if !ok {
		log.Print("Cannot get tokenCred")
	}

	if tokenCred.Token != "" {
		log.Print("Printing tokenCred:", tokenCred)
		token1 := getTokens()

		api1 := anaconda.NewTwitterApiWithCredentials(tokenCred.Token, tokenCred.Secret, token1.ConsumerKey, token1.ConsumerSecret)

		timelineTweets := getTimelineTweets(api1)
		keywordStore := populateKeywordStore()
		classifiedTweets := classifyTweets(timelineTweets, keywordStore)
		p := &Page{Title: "Tech Tweets", TechTweets: classifiedTweets["tech"], PoliticsTweets: classifiedTweets["politics"], TravelTweets: classifiedTweets["travel"], OtherTweets: classifiedTweets["other"]}
		renderTemplate(w, "index", p)

	} else {
		e := &Page{}
		renderTemplate(w, "login", e)
	}

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
		if strings.Contains(strings.ToLower(tweet.FullText), strings.ToLower(" "+keyword+" ")) {
			//		if strings.ToLower(tweet.FullText) == strings.ToLower(keyword) {
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

func getTimelineTweets(ap *anaconda.TwitterApi) []anaconda.Tweet {
	v := url.Values{}
	v.Set("count", "200")
	v.Set("tweet_mode", "extended")
	if mode == "dev" {
		timelineTweets := getDummyTimeline()
		return timelineTweets
	} else {
		timelineTweets, err := ap.GetHomeTimeline(v)
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
	log.Print("entering signinhandler")
	callback := "http://" + r.Host + "/callback"
	tempCred, err := signinOAuthClient.RequestTemporaryCredentials(nil, callback, nil)

	//	authURL, tempCred, err := api.AuthorizationURL(callback)
	if err != nil {
		http.Error(w, "Error getting temp cred, "+err.Error(), 500)
		return
	}
	//	ctx := context.Context
	//	context.WithValue(ctx, key interface{}, val interface{})
	s, _ := store.Get(r, "twit")
	s.Values[tempCredKey] = tempCred

	if err := s.Save(r, w); err != nil {
		http.Error(w, "Error saving session, "+err.Error(), 500)
		return
	}

	http.Redirect(w, r, signinOAuthClient.AuthorizationURL(tempCred, nil), 302)
}

func twitterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, "twit")
	if err != nil {
		log.Print("We have an error getting the session cookie: ", err)

	}

	tempCred, _ := s.Values[tempCredKey].(oauth.Credentials)

	//	t, ok1 := tempCred1.(*oauth.Credentials)
	//	tempCred := t.(oauth.Credentials)

	//	log.Print("Printing resifual from type assertion: ", ok1)
	//	log.Print("Trying to get tempCred1 in callback: ", tempCred1)
	log.Print("Trying to get tempCred in callback: ", tempCred)
	log.Print("Trying to print Token ", tempCred.Token)

	token := r.FormValue("oauth_token")
	verifier := r.FormValue("oauth_verifier")

	if tempCred.Token != token {
		http.Error(w, "Unknown oauth_token", 500)
	}

	//	tokenCred, _, err := api.GetCredentials(tempCred, verifier)
	tokenCred, _, err := oauthClient.RequestToken(nil, &tempCred, verifier)

	if err != nil {
		http.Error(w, "Error getting request token, "+err.Error(), 500)
		return
	}

	delete(s.Values, tempCredKey)
	s.Values[tokenCredKey] = tokenCred
	if err := s.Save(r, w); err != nil {
		http.Error(w, "Error saving sessions, "+err.Error(), 500)
		return
	}

	log.Print("tokenCred is :", tokenCred)
	http.Redirect(w, r, "/", 302)

}

// serveLogout clears the authentication cookie.
func twitterLogoutHandler(w http.ResponseWriter, r *http.Request) {

	s, err := store.Get(r, "twit")
	if err != nil {
		log.Print("We have an error getting the session cookie: ", err)

	}

	delete(s.Values, tokenCredKey)
	if err := s.Save(r, w); err != nil {
		http.Error(w, "Error saving session , "+err.Error(), 500)
		return
	}
	log.Print("Logged out!!")
	http.Redirect(w, r, "/", 302)
}

func main() {
	wordPtr := flag.String("mode", "", "which mode to run")
	flag.Parse()

	fmt.Println("word:", *wordPtr)

	if *wordPtr == "dev" {
		mode = "dev"
	}
	token := getTokens()
	//	api = anaconda.NewTwitterApiWithCredentials(token.AccessToken, token.AccessTokenSecret, token.ConsumerKey, token.ConsumerSecret)
	oauthClient.Credentials.Token = token.ConsumerKey
	oauthClient.Credentials.Secret = token.ConsumerSecret
	signinOAuthClient = oauthClient
	signinOAuthClient.ResourceOwnerAuthorizationURI = "https://api.twitter.com/oauth/authenticate"
	gob.Register(oauth.Credentials{})
	//	anaconda.SetConsumerKey(token.ConsumerKey)
	//	anaconda.SetConsumerSecret(token.ConsumerSecret)
	cssHandler := http.FileServer(http.Dir("./css/"))
	jsHandler := http.FileServer(http.Dir("./js/"))
	imagesHandler := http.FileServer(http.Dir("./images/"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/callback", twitterCallbackHandler)
	http.HandleFunc("/logout", twitterLogoutHandler)
	http.HandleFunc("/dump", dumpHandler)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))
	http.Handle("/images/", http.StripPrefix("/images/", imagesHandler))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
