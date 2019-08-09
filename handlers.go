package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
		rootCategories := populateCategories(rootKeywordFilename)
		log.Print("Root categories : ", rootCategories)
		adminCategories := populateCategories(adminKeywordFile)
		log.Print("Admin categories : ", adminCategories)
		seedCategories := mergeKeywords(rootCategories, adminCategories)
		log.Print("Merged root and admin keyword store : ", seedCategories)

		userKeywordFilename := s.Values[userKeywordPresent].(string)
		userCategories := populateCategories(userKeywordFilename)
		if userKeywordFilename != templateKeywordFilename {
			log.Print("Merging seed and user keywords>>>>>>>>>>>>>>>")
			seedCategories = mergeKeywords(seedCategories, userCategories)
		}

		for categoryIndex, category := range seedCategories {
			log.Print(categoryIndex)
			log.Print(category.Show)
		}
		classifiedTweets := classifyTweets(timelineTweets, seedCategories)

		log.Print(len(classifiedTweets["others"]))

		p1 := &Page1{Tweets: classifiedTweets, Categories: populateCategories(userKeywordFilename)}

		if err := s.Save(r, w); err != nil {
			http.Error(w, "Error saving session, "+err.Error(), 500)
			return
		}

		renderTemplate(w, "index", p1)

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

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		log.Print(string(body))

		var keywordFile string

		log.Print("Keyword filename from cookie: ", s.Values[userKeywordPresent])

		// what happens if a person does not allow cookie? then twitter sign in wont work as well
		keywordFile = s.Values[userKeywordPresent].(string)
		categories := populateCategories(keywordFile)
		var keywordToAdd KeywordToAdd
		err = json.Unmarshal(body, &keywordToAdd)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories[keywordToAdd.Category].Keywords = append(categories[keywordToAdd.Category].Keywords, keywordToAdd.Phrase)

		log.Print("keywordstore has been appended: ", categories)
		b, err := json.Marshal(categories)

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

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func retweetHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Entering retweet handler")
	if r.Method == "POST" {

		s := getSession(r, sessionName)
		api1 := getAuthenticatedTwitterApi(s)

		twitterId := recieveTwitterId(w, r)
		rt, err := api1.Retweet(twitterId, false)

		if err != nil {
			log.Print(err)
			log.Print(rt)
		}

	}

}

func favHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Entering favourite handler")
	if r.Method == "POST" {

		s := getSession(r, sessionName)
		api1 := getAuthenticatedTwitterApi(s)

		twitterId := recieveTwitterId(w, r)
		rt, err := api1.Favorite(twitterId)

		if err != nil {
			log.Print(err)
			log.Print(rt)
			// send back not success
		}

	}

}

func saveCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("entering saveCategories handler")
	if r.Method == "POST" {
		log.Print("%+v\n", r.Form)
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//		fmt.Fprintf(w, "Post from website! r.PostForm = %v\n", r.PostForm)

		//	m := make(map[string]*Category)
		log.Print(">>>>>>>>> PostForm", r.Form)
		log.Print(r.Form["product"])
		categories := populateCategories("keyword_kgthegreat.json")
		for categoryIndex, keywords := range r.PostForm {
			log.Print(categoryIndex)
			// due to how checkbox and html form post works
			if len(keywords) == 2 {
				log.Print("This is the value of show: ", keywords[0])
				log.Print("This is the value of keywords: ", keywords[1])
				categories[categoryIndex].Show, _ = strconv.ParseBool(keywords[0]) //convert to bool
				//before saving, weed out empty strings. try to use map reduce if available
				splitFn := func(c rune) bool {
					return c == ','
				}
				//				keywords := strings.Split(keywords[1], ",") //convert to array
				keywords := strings.FieldsFunc(keywords[1], splitFn) //convert to array
				/*				for _, kwd := range keywords {
											if kwd == "" {
												//delete from array
								`				//delete(keywords, kwd)
											}
										}*/
				log.Print("Printing keywords after weeding >>>>>>>>>>>", keywords)
				categories[categoryIndex].Keywords = keywords
			} else {
				log.Print("This is the value of show: ", "false")
				log.Print("This is the value of keywords: ", keywords[0])
				categories[categoryIndex].Show = false
				categories[categoryIndex].Keywords = strings.Split(keywords[0], ",")
			}
			b, _ := json.Marshal(categories)

			filename := "keyword_kgthegreat.json" //" + s.Values[screenName].([]string)[0] + dotJson
			ioutil.WriteFile(filename, b, 0600)

			//			m[categoryIndex].Keywords = keywords
		}
		log.Print(">>>>Printing ", categories)
		//		log.Print(r.Form)
		//a := r.FormValue("tech")
		//log.Print(a)
		//		log.Print("It's a post")

	}
}

func manageHandler(w http.ResponseWriter, r *http.Request) {
	s := getSession(r, sessionName)
	userKeywordFilename := s.Values[userKeywordPresent].(string)
	p1 := &Page1{Categories: populateCategories(userKeywordFilename)}

	renderTemplate(w, "manage", p1)

}
