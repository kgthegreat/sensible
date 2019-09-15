package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, sessionName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tokenCred, ok := s.Values[tokenCredKey].(oauth.Credentials)

	if (tokenCred.Token != "" && ok) || mode == "dev" {
		log.Print("Screen Name from cookie", s.Values[screenName])

		userFilename := keywordPrefix + s.Values[screenName].([]string)[0] + dotJson

		if _, err := os.Stat(userFilename); err == nil {
			s.Values[userKeywordPresent] = userFilename

		}

		if s.Values[userKeywordPresent] == nil && s.Values[screenName].([]string)[0] != adminUsername {
			log.Print("Copying template file")
			copyFile(templateKeywordFilename, userFilename)
			s.Values[userKeywordPresent] = userFilename
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
		log.Print(">>>>>>>^^^^^^^^ ", userKeywordFilename)
		userCategories := populateCategories(userKeywordFilename)
		seedCategories = mergeKeywords(seedCategories, userCategories)

		classifiedTweets := classifyTweets(timelineTweets, seedCategories)

		p1 := &Page1{Tweets: classifiedTweets, Categories: populateCategories(userKeywordFilename)}

		if err := s.Save(r, w); err != nil {
			http.Error(w, "Error saving session, "+err.Error(), 500)
			return
		}

		renderTemplate(w, "index", p1)

	} else {
		log.Print("This user is not logged in")
		e := &Page1{}
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
		filename := s.Values[userKeywordPresent].(string)
		ioutil.WriteFile(filename, b, 0600)

		if e := s.Save(r, w); e != nil {
			http.Error(w, "Error saving session, "+e.Error(), 500)
			return
		}

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
	s := getSession(r, sessionName)
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		//		fmt.Fprintf(w, "Post from website! r.PostForm = %v\n", r.PostForm)

		//	m := make(map[string]*Category)
		log.Print(">>>>>>>>> PostForm", r.Form)
		// TODO remove hardcoding
		filename := "keyword_" + s.Values[screenName].([]string)[0] + dotJson
		categories := populateCategories(filename)
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
				log.Print("Printing keywords after weeding >>>>>>>>>>>", keywords)
				categories[categoryIndex].Keywords = keywords
			} else {
				log.Print("This is the value of show: ", "false")
				log.Print("This is the value of keywords: ", keywords[0])
				categories[categoryIndex].Show = false
				categories[categoryIndex].Keywords = strings.Split(keywords[0], ",")
			}
			//			m[categoryIndex].Keywords = keywords
		}
		b, _ := json.Marshal(categories)

		ioutil.WriteFile(filename, b, 0600)
		http.Redirect(w, r, "/", http.StatusSeeOther)

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

func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	s := getSession(r, sessionName)
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		filename := s.Values[userKeywordPresent].(string)
		categories := populateCategories(filename)

		newCategoryName := r.PostForm["add-category-name"][0]
		log.Print(">>>>>>>>> PostFormCategories ", r.PostForm["add-category-keywords"])
		newCategoryKeywords := r.PostForm["add-category-keywords"]
		log.Print(">>>>>>>>> PostFormCategories ", newCategoryKeywords)
		log.Print(">>>>>>>>> PostFormCategories ", len(newCategoryKeywords))
		if len(newCategoryKeywords) == 0 || (len(newCategoryKeywords) == 1 && (newCategoryKeywords[0] == "" || newCategoryKeywords[0] == " ")) {
			newCategoryKeywords = nil
		}
		categories[newCategoryName] = &Category{Name: newCategoryName, Show: true, Keywords: newCategoryKeywords}
		log.Print(">>>>>>>>> PostFormCategories ", categories[newCategoryName].Keywords)

		b, _ := json.Marshal(categories)

		ioutil.WriteFile(filename, b, 0600)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
