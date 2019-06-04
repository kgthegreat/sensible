package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

func classifyTweets(timelineTweets []anaconda.Tweet, keywordStore Keyword) map[string][]anaconda.Tweet {
	// TODO change the below to use Page model
	classifiedTweets := make(map[string][]anaconda.Tweet)
	var techTweets []anaconda.Tweet
	var politicsTweets []anaconda.Tweet
	var travelTweets []anaconda.Tweet
	var sportsTweets []anaconda.Tweet
	var businessTweets []anaconda.Tweet
	var otherTweets []anaconda.Tweet

	//	totalKeywordStore := mergeKeywords(keywordStore, populateKeywordStore())

	for _, tweet := range timelineTweets {
		if itIs(keywordStore.TechKeywords, tweet) {
			techTweets = append(techTweets, tweet)
		} else if itIs(keywordStore.PoliticsKeywords, tweet) {
			politicsTweets = append(politicsTweets, tweet)
		} else if itIs(keywordStore.TravelKeywords, tweet) {
			travelTweets = append(travelTweets, tweet)
		} else if itIs(keywordStore.SportsKeywords, tweet) {
			sportsTweets = append(sportsTweets, tweet)
		} else if itIs(keywordStore.BusinessKeywords, tweet) {
			businessTweets = append(businessTweets, tweet)
		} else {
			otherTweets = append(otherTweets, tweet)
		}

	}
	classifiedTweets["tech"] = techTweets
	classifiedTweets["politics"] = politicsTweets
	classifiedTweets["travel"] = travelTweets
	classifiedTweets["sports"] = sportsTweets
	classifiedTweets["business"] = businessTweets
	classifiedTweets["other"] = otherTweets
	return classifiedTweets
}

func mergeKeywords(keyword1 Keyword, keyword2 Keyword) Keyword {
	keyword1.PoliticsKeywords = append(keyword1.PoliticsKeywords, keyword2.PoliticsKeywords...)
	keyword1.TechKeywords = append(keyword1.TechKeywords, keyword2.TechKeywords...)
	keyword1.TravelKeywords = append(keyword1.TravelKeywords, keyword2.TravelKeywords...)
	keyword1.SportsKeywords = append(keyword1.SportsKeywords, keyword2.SportsKeywords...)
	keyword1.BusinessKeywords = append(keyword1.BusinessKeywords, keyword2.BusinessKeywords...)
	return keyword1

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
	// if no file then create from a template
	// probably dont create it here. its a waste. create when necessary. return from here and put a check in classify
	/*	if _, err := os.Stat("filename"); os.IsNotExist(err) {
		return nil
		// path/to/whatever does not exist
		//		os.Link(templateKeywordFilename, filename)
	}*/
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
