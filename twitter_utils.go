package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
)

func classifyTweets(timelineTweets []anaconda.Tweet, categories map[string]*Category) map[string][]anaconda.Tweet {

	classifiedTweets := make(map[string][]anaconda.Tweet)
	//	var emptyTweets = []anaconda.Tweet{}
	for categoryIndex, _ := range categories {
		//		log.Println("Categorising this tweet ", tweet.Text)
		flag := false
		for i, tweet := range timelineTweets {

			//			classifiedTweets[categoryIndex] = emptyTweets
			//log.Print("for category, ", categoryIndex)
			//log.Print("The show attribute is ", categories[categoryIndex].Show)
			if len(categories[categoryIndex].Keywords) > 0 && itIs(categories[categoryIndex].Keywords, tweet) {
				log.Print("Are we ever here??")
				flag = true
				classifiedTweets[categoryIndex] = append(classifiedTweets[categoryIndex], tweet)
				if i < len(timelineTweets) {
					timelineTweets = timelineTweets[:i+copy(timelineTweets[i:], timelineTweets[i+1:])]
				}

				//				timelineTweets[i] = timelineTweets[len(timelineTweets)-1] // Replace it with the last o//ne. CAREFUL only works if you have enough elements.
				//				timelineTweets = timelineTweets[:len(timelineTweets)-1]
			}

		}
		if !flag {
			//			classifiedTweets["xOthers"] = append(classifiedTweets["xOthers"], tweet)
		}

	}
	classifiedTweets["xOthers"] = timelineTweets
	return classifiedTweets
}

func mergeKeywords(categories1 map[string]*Category, categories2 map[string]*Category) map[string]*Category {

	for category, _ := range categories2 {
		//		log.Print
		//		log.Print("Tryong to print somthhiong:  >>>>> ", categories1[category])
		if categories2[category].Show {
			if categories1[category] != nil {
				categories2[category].Keywords = append(categories2[category].Keywords, categories1[category].Keywords...)
			}
		} else {
			delete(categories2, category)
		}
	}
	return categories2

}

func itIs(keywords []string, tweet anaconda.Tweet) bool {
	for _, keyword := range keywords {
		//check for string empty before getting into if
		//	if strings.Contains(strings.ToLower(tweet.FullText), strings.ToLower(keyword)) {
		//TODO add ability to search with words starting with #, @ or ending with 's (plurals)
		var fullText string
		if tweet.RetweetedStatus != nil {
			fullText = tweet.RetweetedStatus.FullText
		} else {
			fullText = tweet.FullText
		}
		contains, _ := regexp.MatchString("(?i)\\b"+keyword+"\\b", fullText)
		if contains {
			//		if strings.ToLower(tweet.FullText) == strings.ToLower(keyword) {
			return true
		}
	}
	return false
}

func populateCategories(filename string) map[string]*Category {
	log.Print(" populating keyword from this file : ", filename)
	var categories map[string]*Category
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
