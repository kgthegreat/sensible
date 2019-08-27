package main

import (
	"github.com/ChimeraCoder/anaconda"
)

type Token struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
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

type Page1 struct {
	Tweets     map[string][]anaconda.Tweet
	Categories map[string]*Category
}

type Category struct {
	Name     string
	Show     bool
	Keywords []string
}
