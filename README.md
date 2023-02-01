<h1>GO Facebook OAuth</h1>
<p>***This is unofficial package</p>
<p>This project is just to build a package to do the OAuth with Facebook from my pain to write a code by hand.</p>

<h2>Authorization Request</h2>
<img src="./screenshots/oauth_flow.png">

<h2>Callback</h2>
<img src="./screenshots/callback_flow.png">

<h2>Quickstart</h2>

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rayato159/go-facebook-oauth/src"
)

func main() {
	oauth := src.NewGoFacebookOauth(
		"15.0",
		"callback-url",
		"app-id",
		"client-secret",
	)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		url, err := oauth.GetCallbackUrl("test")

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error(),
			})
		}
		http.Redirect(w, r, url, http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")

		token, err := oauth.GetAccessToken(code)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error(),
			})
		}
		fmt.Println(token)

		profile, err := oauth.GetUserData(token.AccessToken)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error(),
			})
		}
		json.NewEncoder(w).Encode(profile)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		res, err := oauth.Logout("access_token")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error(),
			})
		}
		json.NewEncoder(w).Encode(res)
	})

	fmt.Println("Listening on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
```