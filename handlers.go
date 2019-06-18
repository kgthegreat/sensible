package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"gopkg.in/jdkato/prose.v2"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, sessionName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tokenCred, ok := s.Values[tokenCredKey].(oauth.Credentials)

	if !ok {
		log.Print("This user is not logged in")
	}

	if tokenCred.Token != "" || mode == "dev" {
		if s.Values[userKeywordPresent] == nil {
			s.Values[userKeywordPresent] = templateKeywordFilename
		}

		log.Print("Printing tokenCred:", tokenCred)
		token1 := getTokens()

		api1 := anaconda.NewTwitterApiWithCredentials(tokenCred.Token, tokenCred.Secret, token1.ConsumerKey, token1.ConsumerSecret)

		timelineTweets := getTimelineTweets(api1)
		keywordStore := populateKeywordStore(rootKeywordFilename)

		adminKeywordStore := populateKeywordStore(adminKeywordFile)

		keywordStore = mergeKeywords(keywordStore, adminKeywordStore)

		userKeywordFilename := s.Values[userKeywordPresent].(string)

		if userKeywordFilename != templateKeywordFilename {
			keywordStore = mergeKeywords(keywordStore, populateKeywordStore(userKeywordFilename))

		}
		//		userKeywordStore := populateKeywordStore("keyword_" + s.Values[screenName].([]string)[0] + dotJson)
		classifiedTweets := classifyTweets(timelineTweets, keywordStore)
		p := &Page{Title: "Tech Tweets", TechTweets: classifiedTweets["tech"], BusinessTweets: classifiedTweets["business"], PoliticsTweets: classifiedTweets["politics"], SportsTweets: classifiedTweets["sports"], TravelTweets: classifiedTweets["travel"], OtherTweets: classifiedTweets["other"]}

		if err := s.Save(r, w); err != nil {
			http.Error(w, "Error saving session, "+err.Error(), 500)
			return
		}

		renderTemplate(w, "index", p)

	} else {
		e := &Page{}
		renderTemplate(w, "login", e)
	}

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

func categoriseHandler(w http.ResponseWriter, r *http.Request) {
	s := getSession(r, sessionName)
	log.Print("Printing test: ", s.Values["test"])
	log.Print("Keyword filename from cookie just after getting session: ", s.Values[screenName])
	log.Print("a new variable from cookie just after getting session: ", s.Values["some"])
	//	s := getSession(sessionName, r)
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		log.Print(string(body))
		//		results = append(results, string(body))

		//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//		w.WriteHeader(http.StatusOK)

		var keywordFile string

		log.Print("Keyword filename from cookie: ", s.Values[userKeywordPresent])

		// what happens if a person does not allow cookie? then twitter sign in wont work as well
		keywordFile = s.Values[userKeywordPresent].(string)
		/*
			if s.Values[userKeywordPresent] == templateKeywordFilename {
				keywordFile = templateKeywordFilename

				//			keywordFile = keywordPrefix + s.Values[screenName].([]string)[0] + dotJson
			} else {
				//			keywordFile =
				keywordFile = s.Values[userKeywordPresent].(string)
			}*/
		keywordStore := populateKeywordStore(keywordFile)

		//		fmt.Fprint(w, "POST done")

		//		b, error = Json.Unmarshal()
		var keywordToAdd KeywordToAdd
		err = json.Unmarshal(body, &keywordToAdd)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//reflection or metaprogramming in golang
		if keywordToAdd.Category == "politics" {
			keywordStore.PoliticsKeywords = append(keywordStore.PoliticsKeywords, keywordToAdd.Phrase)

		} else if keywordToAdd.Category == "travel" {
			keywordStore.TravelKeywords = append(keywordStore.TravelKeywords, keywordToAdd.Phrase)
		} else if keywordToAdd.Category == "tech" {
			keywordStore.TechKeywords = append(keywordStore.TechKeywords, keywordToAdd.Phrase)
		} else if keywordToAdd.Category == "sports" {
			keywordStore.SportsKeywords = append(keywordStore.SportsKeywords, keywordToAdd.Phrase)
		} else if keywordToAdd.Category == "business" {
			keywordStore.BusinessKeywords = append(keywordStore.BusinessKeywords, keywordToAdd.Phrase)
		}

		log.Print("keywordstore has been appended: ", keywordStore)
		b, err := json.Marshal(keywordStore)
		//			filename := "keyword.json"
		log.Print("what are we getting", s.Values[screenName].([]string)[0])
		filename := "keyword_" + s.Values[screenName].([]string)[0] + dotJson
		ioutil.WriteFile(filename, b, 0600)
		s.Values[userKeywordPresent] = filename
		s.Values["test"] = "Hi this is test 2"
		s.Values["some"] = "else"
		log.Print("Fetching from cookie before saving: ", s.Values[userKeywordPresent])

		if e := s.Save(r, w); e != nil {
			http.Error(w, "Error saving session, "+e.Error(), 500)
			return
		}
		log.Print("Fetching from cookie after saving: ", s.Values[userKeywordPresent])
		//		http.Redirect(w, r, "/", 302)
		//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
