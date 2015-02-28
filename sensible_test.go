package main

import (
	"testing"
	"github.com/ChimeraCoder/anaconda"
)

func TestCanary(t *testing.T) {
}
func TestTechTweetsShouldIncludeAndroidTweets(t *testing.T) {
	techTweet1 := anaconda.Tweet {
		Text: "Looking for android developer",
	}
	nonTechTweet := anaconda.Tweet {
		Text: "This is not related",
	}
	expected := []anaconda.Tweet{techTweet1}
	timelineTweets := []anaconda.Tweet{techTweet1, nonTechTweet}
	actual := getTechTweets(timelineTweets)
	if len(actual) != len(expected) {
		t.Errorf("Wanted: %v, Got: %v", expected, actual)
	}


}
	
func TestTechTweetsShouldIncludeAllRelevantTweets(t *testing.T) {
	techTweet1 := anaconda.Tweet {
		Text: "Looking for android developer",
	}
	techTweet2 := anaconda.Tweet {
		Text: "The newest thing in golang",
	}

	nonTechTweet := anaconda.Tweet {
		Text: "This is not related",
	}
	timelineTweets := []anaconda.Tweet{techTweet1, techTweet2, nonTechTweet}
	expected := []anaconda.Tweet{techTweet1, techTweet2}

	actual := getTechTweets(timelineTweets)
	if len(actual) != len(expected) {
		t.Errorf("Wanted: %v, Got: %v", expected, actual)
	}


}

func TestTechTweetsShouldIncludeARelevantTweetOnlyOnce(t *testing.T) {
	techTweet1 := anaconda.Tweet {
		Text: "Looking for android developer golang",
	}
	techTweet2 := anaconda.Tweet {
		Text: "The newest thing in golang",
	}

	nonTechTweet := anaconda.Tweet {
		Text: "This is not related",
	}
	timelineTweets := []anaconda.Tweet{techTweet1, techTweet2, nonTechTweet}
	expected := []anaconda.Tweet{techTweet1, techTweet2}

	actual := getTechTweets(timelineTweets)
	if len(actual) != len(expected) {
		t.Errorf("Wanted: %v, Got: %v", expected, actual)
	}


}
