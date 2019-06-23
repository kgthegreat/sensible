package main

import "github.com/ChimeraCoder/anaconda"

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
	SportsKeywords   []string
	BusinessKeywords []string
}

type Page struct {
	Title          string
	TechTweets     []anaconda.Tweet
	PoliticsTweets []anaconda.Tweet
	TravelTweets   []anaconda.Tweet
	SportsTweets   []anaconda.Tweet
	BusinessTweets []anaconda.Tweet
	OtherTweets    []anaconda.Tweet
}

type TweetToClassify struct {
	Text         string
	Type         string
	SelectedTags []string
}

type KeywordToAdd struct {
	Phrase   string
	Category string
}