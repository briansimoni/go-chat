
package main

import (
	"github.com/emicklei/go-restful"
	//"io"
	"io/ioutil"
	//"encoding/json"
	"net/http"
	"os"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"strings"
	"encoding/json"
)


var chat[] string

// the key:value will be uuid:token
var sessions map[string] string


// curl http://localhost:4000/test -H "Content-Type: application/json" -X POST -d '{"username":"xyz","password":"xyz"}'
func main() {

	sessions = make(map[string]string)

	port := "4000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)

	d := new(restful.WebService)
	ws.Route(d.GET("/").To(home))
	ws.Route(d.GET("/login").To(githubLogin))
	ws.Route(d.POST("/send").To(sendMessages))
	ws.Route(d.GET("/getMessages").To(getMessages))
	fmt.Println("Starting http server on", port)
	http.ListenAndServe(":"+port, nil)
}


func sendMessages(req *restful.Request, resp *restful.Response) {
	cookie, err := req.Request.Cookie("session-id")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("This is the value of the cookie!:", cookie.Value)
	token := sessions[cookie.Value] // access token

	apiURL := "https://api.github.com/user?" + string(token)
	userData, err := http.Get(apiURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	var chatUser map[string]interface{}

	j, err := ioutil.ReadAll(userData.Body)
	err = json.Unmarshal(j, &chatUser)
	if err != nil {
		fmt.Println(err.Error())
	}

	username := chatUser["login"].(string)


	data, _ := ioutil.ReadAll(req.Request.Body)
	m := username + ": " + string(data)
	//m = template.HTMLEscapeString(m) // This would prevent people from script injecting
	chat = append(chat, m)
	fmt.Println(chat)
	fmt.Println(len(chat))
}

func getMessages(req *restful.Request, resp *restful.Response) {
	resp.WriteAsJson(chat)
}

func hello(req *restful.Request, resp *restful.Response) {
	if verifyLogin(req, resp) {
		cookie, err := req.Request.Cookie("session-id")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("This is the value of the cookie!:", cookie.Value)
		fmt.Println("this is the value of the token!:", sessions[cookie.Value]) // access token
		body, _ := ioutil.ReadFile("home.html")
		fmt.Fprint(resp, string(body))
	}
}



func githubLogin(req *restful.Request, resp *restful.Response) {
	cookie, err := req.Request.Cookie("session-id")
	id := cookie.Value;
	if err != nil {
		fmt.Println("COOKIE NOT SET", err.Error())
	}

	auth := req.Request.URL.Query().Get("code")
	clientId := "fda60467458c62443d52"
	clientSecret := "256b07b11e40e50f8ea34c3af5b6c9dae678a490"
	// https://github.com/login/oauth/access_token?client_id=fda60467458c62443d52&client_secret=256b07b11e40e50f8ea34c3af5b6c9dae678a490&code=53f778a136cbb86b39ce
	url := "https://github.com/login/oauth/access_token?client_id=" + clientId +"&client_secret=" + clientSecret + "&code=" + auth
	r, err := http.Post(url, "multipart/form-data", strings.NewReader(""))
	if err != nil {
		fmt.Println(err.Error())
	}

	token, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	sessions[id] = string(token)

	apiURL := "https://api.github.com/user?" + string(token)
	userData, err := http.Get(apiURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	var chatUser map[string]interface{}

	j, err := ioutil.ReadAll(userData.Body)
	err = json.Unmarshal(j, &chatUser)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(sessions)

	//userLogin := chatUser["login"].(string)

	fmt.Println("This is the value of the session in the login function .... " , sessions[id])
	if sessions[id] != "" {
		http.Redirect(resp.ResponseWriter, req.Request, "/hello", 302)
	}
	//fmt.Fprint(resp.ResponseWriter, "This is your github username " + userLogin)
}

func home(req *restful.Request, resp *restful.Response) {
	cookie, err := req.Request.Cookie("session-id")
	fmt.Println(cookie, err)
	if err != nil {
		id, _ := uuid.NewV4()
		cookie = &http.Cookie {
			Name: "session-id",
			Value: id.String(),
		}
		http.SetCookie(resp.ResponseWriter, cookie)
	}


	body, _ := ioutil.ReadFile("login.html")
	fmt.Fprint(resp, string(body))
}


// if the user is not logged in, redirect to the home page
func verifyLogin(req *restful.Request, resp *restful.Response ) bool {
	cookie, err := req.Request.Cookie("session-id")
	if cookie.Value != "" {
		_, exists := sessions[cookie.Value]
		if !exists {
			http.Redirect(resp.ResponseWriter, req.Request, "/", 302)
			return false
		}
		return true
	} else if err != nil {
		fmt.Println(err.Error())
		http.Redirect(resp.ResponseWriter, req.Request, "/", 302)
		return false
	} else {
		http.Redirect(resp.ResponseWriter, req.Request, "/", 302)
		return false
	}
}