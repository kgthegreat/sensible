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
	"gopkg.in/jdkato/prose.v2"
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
	tempCredKey        = "tempCred"
	tokenCredKey       = "tokenCred"
	screenName         = "screenName"
	sessionName        = "twit"
	rootKeywordFilenam = "keyword.json"
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

type TweetToClassify struct {
	Text         string
	Type         string
	SelectedTags []string
}

type KewordToAdd struct {
	Phrase   string
	Category string
}

func categoriseHandler(w http.ResponseWriter, r *http.Request) {
	s := getSession(sessionName, r)
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		log.Print(string(body))
		//		results = append(results, string(body))

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		//		fmt.Fprint(w, "POST done")
		keywordStore1 := populateKeywordStore(rootKeywordFilename)
		//		b, error = Json.Unmarshal()
		var keywordToAdd KewordToAdd
		err = json.Unmarshal(body, &keywordToAdd)
		if keywordToAdd.Category == "politics" {
			b, err := json.Marshal(keywordStore1)
			fmt.Println("json", b)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//			filename := "keyword.json"
			log.Print("what are we getting", s.Values[screenName].(string))
			filename := "keyword" + s.Values[screenName].(string) + ".json"
			ioutil.WriteFile(filename, b, 0600)

		}

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func classifyHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["tweet"]
	cat, ok := r.URL.Query()["type"]

	if !ok || len(keys[0]) < 1 || len(cat[0]) < 1 {
		log.Println("Url Param 'tweet' or 'type' is missing")
		return
	}
	//	tweetText := "Narendra Modi is astonishing. Virat Kohli is a good batsman. Madhya Pradesh polls are going to be exciting. Hum logon ko kuch nahi pata. (How), do we know this?"
	//	tweetText := "@jdkato, go to http://example.com thanks :)."
	doc, err := prose.NewDocument(keys[0])
	if err != nil {
		log.Fatal(err)
	}
	var selectedTags []string
	for _, ent := range doc.Tokens() {
		tag := ent.Tag
		text := ent.Text
		log.Print(text + " " + tag)
		if tag == "NNP" || tag == "NN" || tag == "JJ" {
			selectedTags = append(selectedTags, text+" "+tag)
		}
		// Go GPE
		// Google GPE
	}
	e := &TweetToClassify{Text: keys[0], Type: cat[0], SelectedTags: selectedTags}
	renderTemplate(w, "classify", e)
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

	if tokenCred.Token != "" || mode == "dev" {
		log.Print("Printing tokenCred:", tokenCred)
		token1 := getTokens()

		api1 := anaconda.NewTwitterApiWithCredentials(tokenCred.Token, tokenCred.Secret, token1.ConsumerKey, token1.ConsumerSecret)

		timelineTweets := getTimelineTweets(api1)
		keywordStore := populateKeywordStore(rootKeywordFilename)
		classifiedTweets := classifyTweets(timelineTweets, keywordStore)
		p := &Page{Title: "Tech Tweets", TechTweets: classifiedTweets["tech"], PoliticsTweets: classifiedTweets["politics"], TravelTweets: classifiedTweets["travel"], OtherTweets: classifiedTweets["other"]}
		renderTemplate(w, "index", p)

	} else {
		e := &Page{}
		renderTemplate(w, "login", e)
	}

}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	templates := template.Must(template.ParseGlob("templates/*"))
	//	t, err := template.ParseFiles(tmpl + ".html")
	/*	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	//	err = t.Execute(w, p)
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
			//		if strings.ToLower(tweet.FullText) == strings.ToLower(keyword) {
			return true
		}
	}
	return false
}

func populateKeywordStore(filename string) Keyword {
	var keywordStore Keyword
	//	filename := "keyword.json"
	keyword_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error", err)
	}

	err = json.Unmarshal(keyword_bytes, &keywordStore)
	if err != nil {
		log.Print("Error reading keyword file: ", err)
	}
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

func getSession(name string, r *http.Request) *sessions.Session {
	s, err := store.Get(r, "twit")
	if err != nil {
		log.Print("We have an error getting the session cookie: ", err)

	}
	return s

}

func twitterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	/*	s, err := store.Get(r, "twit")
		if err != nil {
			log.Print("We have an error getting the session cookie: ", err)


		}*/

	s := getSession("twit", r)
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
	tokenCred, urlValues, err := oauthClient.RequestToken(nil, &tempCred, verifier)

	log.Print("We have the screen name!: ", urlValues["screen_name"])
	if err != nil {
		http.Error(w, "Error getting request token, "+err.Error(), 500)
		return
	}

	delete(s.Values, tempCredKey)
	s.Values[tokenCredKey] = tokenCred
	s.Values[screenName] = urlValues["screen_name"]
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
	staticHandler := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/callback", twitterCallbackHandler)
	http.HandleFunc("/logout", twitterLogoutHandler)
	http.HandleFunc("/dump", dumpHandler)
	http.HandleFunc("/classify", classifyHandler)
	http.HandleFunc("/categorise", categoriseHandler)
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
