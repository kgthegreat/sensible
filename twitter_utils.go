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

func classifyTweets(timelineTweets []anaconda.Tweet, keywordStore map[string][]string) map[string][]anaconda.Tweet {
	log.Print("This is the keyword store: ", keywordStore)
	classifiedTweets := make(map[string][]anaconda.Tweet)

	for _, tweet := range timelineTweets {

		flag := false
		for _, category := range categories {

			if itIs(keywordStore[category], tweet) {

				flag = true
				classifiedTweets[category] = append(classifiedTweets[category], tweet)
			}

		}
		if !flag {
			classifiedTweets["others"] = append(classifiedTweets["others"], tweet)
		}

	}

	return classifiedTweets
}

func mergeKeywords(categories1 []Category, categories2 []Category) []Category {

	for _, category := range categories1 {
		category.Keywords = append(category.Keywords, category.Keywords)
		keyword1[category] = append(keyword1[category], keyword2[category]...)
	}
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

func populateKeywordStore(filename string) map[string][]string {
	log.Print(" populating keyword from this file : ", filename)
	var categories []Category
	//	var store1 map[string][]string
	keyword_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error", err)
	}

	err = json.Unmarshal(keyword_bytes, &categories)
	if err != nil {
		log.Print("Error reading keyword file: ", err)
	}
	return categories
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
