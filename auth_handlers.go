package main

import (
	"log"
	"net/http"

	"github.com/garyburd/go-oauth/oauth"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("entering signinhandler")
	callback := "http://" + r.Host + "/callback"
	tempCred, err := signinOAuthClient.RequestTemporaryCredentials(nil, callback, nil)

	//	authURL, tempCred, err := api.AuthorizationURL(callback)
	if err != nil {
		http.Error(w, "Error getting temp cred, "+err.Error(), 500)
		return
	}
	//	ctx := context.Context
	//	context.WithValue(ctx, key interface{}, val interface{})
	s := getSession(r, sessionName)
	s.Values[tempCredKey] = tempCred
	if err := s.Save(r, w); err != nil {
		http.Error(w, "Error saving sessions, "+err.Error(), 500)
		return
	}

	http.Redirect(w, r, signinOAuthClient.AuthorizationURL(tempCred, nil), 302)
}

func twitterCallbackHandler(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, "twit")
	if err != nil {
		log.Print("We have an error getting the session cookie: ", err)

	}

	//	s := getSession(r, sessionName)
	tempCred, _ := s.Values[tempCredKey].(oauth.Credentials)

	//	t, ok1 := tempCred1.(*oauth.Credentials)
	//	tempCred := t.(oauth.Credentials)

	//	log.Print("Printing resifual from type assertion: ", ok1)
	//	log.Print("Trying to get tempCred1 in callback: ", tempCred1)
	log.Print("Trying to get tempCred in callback: ", tempCred)
	log.Print("Trying to print Token ", tempCred.Token)

	token := r.FormValue("oauth_token")
	verifier := r.FormValue("oauth_verifier")

	if tempCred.Token != token {
		http.Error(w, "Unknown oauth_token", 500)
	}

	//	tokenCred, _, err := api.GetCredentials(tempCred, verifier)
	tokenCred, urlValues, err := oauthClient.RequestToken(nil, &tempCred, verifier)

	log.Print("We have the screen name!: ", urlValues["screen_name"])
	if err != nil {
		http.Error(w, "Error getting request token, "+err.Error(), 500)
		return
	}

	delete(s.Values, tempCredKey)
	s.Values[tokenCredKey] = tokenCred
	s.Values[screenName] = urlValues["screen_name"]
	if err := s.Save(r, w); err != nil {
		http.Error(w, "Error saving sessions, "+err.Error(), 500)
		return
	}

	log.Print("tokenCred is :", tokenCred)
	http.Redirect(w, r, "/", 302)

}

// serveLogout clears the authentication cookie.
func twitterLogoutHandler(w http.ResponseWriter, r *http.Request) {

	s, err := store.Get(r, "twit")
	if err != nil {
		log.Print("We have an error getting the session cookie: ", err)

	}

	delete(s.Values, tokenCredKey)
	if err := s.Save(r, w); err != nil {
		http.Error(w, "Error saving session , "+err.Error(), 500)
		return
	}
	log.Print("Logged out!!")
	http.Redirect(w, r, "/", 302)
}
