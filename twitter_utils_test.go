package main

import (
	_ "fmt"
	"testing"

	"github.com/ChimeraCoder/anaconda"
)

const (
	politics = "politics"
	sports   = "sports"
)

var politicsCategory = &Category{Name: politics, Show: true, Keywords: []string{"modi"}}
var businessCategory = &Category{Name: "business", Show: true, Keywords: []string{"toyota"}}
var techCategory = &Category{Name: "ruby", Show: true, Keywords: []string{"ruby"}}
var sportsCategory = &Category{Name: sports, Show: true, Keywords: []string{"tennis"}}

var politicsTweet = anaconda.Tweet{
	FullText: "Narendra modi going to Sri Lanka",
}

var sportsTweet = anaconda.Tweet{
	FullText: "Roger Federer wins tennis",
}

var techTweet1 = anaconda.Tweet{
	FullText: "Looking for android developer",
}
var techTweet2 = anaconda.Tweet{
	FullText: "The newest thing in golang",
}

var nonTechTweet = anaconda.Tweet{
	FullText: "This is not related",
}
var techTweetMixed = anaconda.Tweet{
	FullText: "Looking for android developer golang",
}

var nonTechTweet2 = anaconda.Tweet{
	FullText: "The newest thing in happy",
}

var techPoliticsTweet = anaconda.Tweet{
	FullText: "Android Modi Travel",
}
var genericTweet = anaconda.Tweet{
	FullText: "What is happening",
}

var categoriesMap = make(map[string]*Category)

var timelineTweetsWithOnlyPolitics = []anaconda.Tweet{politicsTweet}

var timelineTweetsWithPoliticsAndSportsTweets = []anaconda.Tweet{politicsTweet, sportsTweet}

func TestCanary1(t *testing.T) {
}

func TestClassifiedTweetsHasAllCategories(t *testing.T) {
	categoriesMap[politics] = politicsCategory
	categoriesMap[sports] = sportsCategory
	classifiedTweets := classifyTweets(timelineTweetsWithOnlyPolitics, categoriesMap)
	//	fmt.Print(classifiedTweets)
	if classifiedTweets[sports] == nil {
		t.Errorf("Expected empty tweet struct when no sports tweet. Received nil")
	}

}

func TestClassifiedTweetsHasAllCategoriesAndIFTweetPresentThenCategoryNotEmpty(t *testing.T) {
	categoriesMap[politics] = politicsCategory
	categoriesMap[sports] = sportsCategory
	classifiedTweets := classifyTweets(timelineTweetsWithOnlyPolitics, categoriesMap)
	//	fmt.Print(classifiedTweets)
	if len(classifiedTweets[politics]) != 1 {
		t.Errorf("Expected politics not to be empty, expected: %v, actuals: %v", 1, len(classifiedTweets[politics]))
	}

}

func TestClassifiedTweetsHasTweetsUnderCorrectCategory(t *testing.T) {
	categoriesMap[politics] = politicsCategory
	categoriesMap[sports] = sportsCategory

	classifiedTweets := classifyTweets(timelineTweetsWithPoliticsAndSportsTweets, categoriesMap)
	//	fmt.Print(classifiedTweets)
	if len(classifiedTweets[politics]) != 1 {
		t.Errorf("Expected politics not to be empty, expected: %v, actuals: %v", 1, len(classifiedTweets[politics]))
	}
	if classifiedTweets[politics][0].FullText != politicsTweet.FullText {
		t.Errorf("Expected politics tweet under politics category, expected: %v, actuals: %v", politicsTweet.FullText, classifiedTweets[politics][0].FullText)
	}
	if len(classifiedTweets[sports]) != 1 {
		t.Errorf("Expected sports not to be empty, expected: %v, actuals: %v", 1, len(classifiedTweets[sports]))
	}
	if classifiedTweets[sports][0].FullText != sportsTweet.FullText {
		t.Errorf("Expected sports tweet under sports category, expected: %v, actuals: %v", sportsTweet.FullText, classifiedTweets[sports][0].FullText)
	}

}
