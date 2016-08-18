
package main

import (
	"github.com/emicklei/go-restful"
	//"io"
	"io/ioutil"
	//"encoding/json"
	"net/http"
	"os"
	"fmt"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello
type message struct {
	username string
	password string
}

var m message

var chat[] string


// curl http://localhost:4000/test -H "Content-Type: application/json" -X POST -d '{"username":"xyz","password":"xyz"}'
func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)

	d := new(restful.WebService)
	ws.Route(d.POST("/test").To(testfunc))
	ws.Route(d.POST("/send").To(testfunc))
	ws.Route(d.GET("/getMessages").To(getMessages))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

var dat map[string]interface{}

func testfunc(req *restful.Request, resp *restful.Response) {
	data, _ := ioutil.ReadAll(req.Request.Body)
	m := string(data)
	chat = append(chat, m)
	fmt.Println(chat)
	fmt.Println(len(chat))
}

func getMessages(req *restful.Request, resp *restful.Response){
	resp.WriteAsJson(chat)
}

func hello(req *restful.Request, resp *restful.Response) {
	body, _ := ioutil.ReadFile("home.html")
	fmt.Fprint(resp, string(body))
}

