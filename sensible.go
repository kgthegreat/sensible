package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	_ "regexp"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/sessions"
)

var api *anaconda.TwitterApi
var mode string

var store = sessions.NewCookieStore([]byte("asdaskdhasdhgsajdgasdsadksakdhasidoajsdousdasf"))

//var store = sessions.NewCookieStore([]byte(storeGUID))

// Session state keys.
const (
	tempCredKey             = "tempCred"
	tokenCredKey            = "tokenCred"
	screenName              = "screenName"
	sessionName             = "twit"
	rootKeywordFilename     = "keyword.json"
	templateKeywordFilename = "keyword_template.json"
	keywordPrefix           = "keyword_"
	dotJson                 = ".json"
	userKeywordPresent      = "filePresent"
	adminKeywordFile        = "keyword_kgthegreat.json"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var signinOAuthClient oauth.Client

func main() {
	modePtr := flag.String("mode", "", "which mode to run")

	portPtr := flag.String("port", "8081", "Which port to run")
	flag.Parse()
	fmt.Println("word:", *modePtr)
	fmt.Println("port:", *portPtr)

	if *modePtr == "dev" {
		mode = "dev"
	}
	token := getTokens()
	//	api = anaconda.NewTwitterApiWithCredentials(token.AccessToken, token.AccessTokenSecret, token.ConsumerKey, token.ConsumerSecret)
	oauthClient.Credentials.Token = token.ConsumerKey
	oauthClient.Credentials.Secret = token.ConsumerSecret
	signinOAuthClient = oauthClient
	signinOAuthClient.ResourceOwnerAuthorizationURI = "https://api.twitter.com/oauth/authenticate"
	gob.Register(oauth.Credentials{})
	staticHandler := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/callback", twitterCallbackHandler)
	http.HandleFunc("/logout", twitterLogoutHandler)
	http.HandleFunc("/dump", dumpHandler)
	http.HandleFunc("/classify", classifyHandler)
	http.HandleFunc("/categorise", categoriseHandler)
	http.HandleFunc("/retweet", retweetHandler)
	http.HandleFunc("/fav", favHandler)
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
		// systemd run
		f := os.NewFile(3, "from systemd")
		l, err := net.FileListener(f)
		if err != nil {
			log.Fatal(err)
		}
		http.Serve(l, nil)
	} else {
		// manual run
		//		log.Fatal(http.ListenAndServe(":8080", nil))
		log.Fatal(http.ListenAndServe(":"+*portPtr, nil))
	}
}
