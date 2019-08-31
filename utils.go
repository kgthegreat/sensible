package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/gorilla/sessions"
)

func getSession(r *http.Request, name string) *sessions.Session {
	s, err := store.Get(r, name)
	if err != nil {
		log.Print("We have an error getting the session cookie: ", err)

	}
	return s

}

func recieveTwitterId(w http.ResponseWriter, r *http.Request) int64 {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}
	log.Print(string(body))
	var tweet *anaconda.Tweet
	err = json.Unmarshal(body, &tweet)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Print(tweet.IdStr)
	id, er := strconv.ParseInt(tweet.IdStr, 10, 64)

	if er != nil {
		log.Print(er)
	}

	return id

}

func getAuthenticatedTwitterApi(s *sessions.Session) *anaconda.TwitterApi {
	tokenCred, ok := s.Values[tokenCredKey].(oauth.Credentials)
	if !ok {
		log.Print("This user is not logged in")
	}

	token1 := getTokens()

	return anaconda.NewTwitterApiWithCredentials(tokenCred.Token, tokenCred.Secret, token1.ConsumerKey, token1.ConsumerSecret)

}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	templates := template.Must(template.ParseGlob("templates/*"))
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func copyFile(src string, dest string) {
	from, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}
