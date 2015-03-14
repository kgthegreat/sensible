package main

import (
	"testing"
	"github.com/ChimeraCoder/anaconda"
	_ "fmt"
)

var politicsTweet = anaconda.Tweet {
		Text: "Narendra modi going to Sri Lanka",
	}

var techTweet1 = anaconda.Tweet {
		Text: "Looking for android developer",
	}
var techTweet2 = anaconda.Tweet {
		Text: "The newest thing in golang",
	}

var nonTechTweet = anaconda.Tweet {
		Text: "This is not related",
	}
var techTweetMixed = anaconda.Tweet {
		Text: "Looking for android developer golang",
	}

var nonTechTweet2 = anaconda.Tweet {
		Text: "The newest thing in happy",
	}

func TestCanary(t *testing.T) {
}

func TestCanClassifyTwoDifferentTweets(t *testing.T) {
	timelineTweets := []anaconda.Tweet{techTweet1, politicsTweet}
	expected := map[string][]anaconda.Tweet{
		"tech": []anaconda.Tweet{techTweet1},
		"politics": []anaconda.Tweet{politicsTweet},
	}

	actual := classifyTweets(timelineTweets)
 	if len(actual["tech"]) != len(expected["tech"]) {
		t.Errorf("Tech did not match. Actual length : %v, Expected length : %v",len(actual["tech"]),len(expected["tech"]))
	}
 	if len(actual["politics"]) != len(expected["politics"]) {
		t.Errorf("Politics did not match. Actual length : %v, Expected length : %v",len(actual["politics"]),len(expected["politics"]))
	}

	if actual["tech"] != nil && actual["tech"][0].Text != expected["tech"][0].Text {
		t.Errorf("Did not match. Actual tweet : %v, Expected tweet : %v",actual["tech"][0].Text,expected["tech"][0].Text)
	}
	if actual["politics"] != nil && actual["politics"][0].Text != expected["politics"][0].Text {
		t.Errorf("Did not match. Actual tweet : %v, Expected tweet : %v",actual["politics"][0].Text,expected["politics"][0].Text)
	}
}

func TestCanClassifyOneTechTweet(t *testing.T) {
	timelineTweets := []anaconda.Tweet{techTweet1}
	expected := map[string][]anaconda.Tweet{
		"tech": []anaconda.Tweet{techTweet1},
	}

	actual := classifyTweets(timelineTweets)
	if len(actual["tech"]) != len(expected["tech"]) {
		t.Errorf("Did not match. Actual length : %v, Expected length : %v",len(actual["tech"]),len(expected["tech"]))
	}
	if actual["tech"] != nil && actual["tech"][0].Text != expected["tech"][0].Text {
		t.Errorf("Did not match. Actual tweet : %v, Expected tweet : %v",actual["tech"][0].Text,expected["tech"][0].Text)
	}
}


func TestATechTweetIsIdentifiedAsTech(t *testing.T) {
	if !itIs("tech", techTweet1) {
		t.Errorf("Did not classify")
	}
}

func TestAnotherTechTweetIsIdentifiedAsTech(t *testing.T) {
	if !itIs("tech", techTweet2) {
		t.Errorf("Did not classify")
	}
}

func TestContextBasedIdentificationOfTweet(t *testing.T) {
	if !itIs("tech", techTweet1) {
		t.Errorf("Did not classify")
	}
	
	if !itIs("politics", politicsTweet) {
		t.Errorf("Could not classify politics tweet")
	}
}
/*
func TestCanClassifyOneTweet(t *testing.T) {
	techTweet := anaconda.Tweet {
		Text: "Looking for android developer golang",
	}
	politicsTweet := anaconda.Tweet {
		Text: "The latest thing in India is BJP government",
	}
	
	timelineTweets := []anaconda.Tweet{techTweet, politicsTweet}
	expected := map[string][]anaconda.Tweet{
		"tech": []anaconda.Tweet{techTweet1},
		"politics": []anaconda.Tweet{politicsTweet},
	}

	actual := getTechTweets(timelineTweets)
	if len(actual) != len(expected) {
		t.Errorf("Did not match")
	}


	
}*/
