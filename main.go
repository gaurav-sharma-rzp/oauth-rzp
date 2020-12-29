package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	//"fmt"
	//"io/ioutil"
	//"log"
	//"net/http"
	//"net/url"
	//"os"
	//"strconv"
	"strings"
)



func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\n\t\"text\": \"Hello world\"\n}")
}

//
//func pageHandler(w http.ResponseWriter, r *http.Request) {
//	page := "<!DOCTYPE html>\n<html>\n<head>\n    <title>Page Title</title>\n</head>\n<body>\n\n<h1>This is a Heading</h1>\n<p>This is a paragraph.</p>\n<button id = \"myButton\">Click Here bro</button>\n\n<script type=\"text/javascript\">\n    document.getElementById(\"myButton\").onclick = function () {\n        location.href = \"https://auth.razorpay.com/authorize?client_id=FN2tdt8BCDmrwr&response_type=code&redirect_uri=https://oauth-test-rzp.herokuapp.com/public&scope=read_only&state=current_state\";\n    };\n</script>\n</body>\n</html>\n"
//
//	fmt.Fprintf(w, page)
//}

var code string
var accessToken string

func callback(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Callback")
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'code' is missing")
		//fmt.Fprintf(w, "Callback! code not present")
		//return
	} else {
		// Query()["key"] will return an array of items,
		// we only want the single item.
		key := keys[0]

		log.Println("Url Param 'code' is: " + string(key))
		//fmt.Fprintf(w, "Callback! "+ string(key))
		code = string(key)
	}

	page := "<!DOCTYPE html>\n<html>\n<head>\n    <title>OAuth Testing</title>\n</head>\n<body>\n\n<h1>OAuth Testing</h1>\n<p>Testing OAtuh to Razorpay Comunication</p>\n<button id = \"codeB\">Get Auth Code</button>\n<button id = \"tokenB\">Get Tokens</button>\n<button id = \"paymentsB\">Get Payments</button>\n<br>\n<textarea id=\"txtArea\" name=\"area\" rows=\"20\" cols=\"100\">\n" + code +
		"\n</textarea>\n<script type=\"text/javascript\">\n    document.getElementById(\"codeB\").onclick = function () {\n        location.href = \"https://auth.razorpay.com/authorize?client_id=FN2tdt8BCDmrwr&response_type=code&redirect_uri=http://localhost:3000/callback&scope=read_only&state=current_state\";\n    };\n\n    document.getElementById(\"tokenB\").onclick = function () {\n        var apiUrl = 'http://localhost:3000/token';\n        fetch(apiUrl).then(response => {\n            return response.json();\n        }).then(data => {\n            // Work with JSON data here\n            console.log(data);\n            document.getElementById(\"txtArea\").value = JSON.stringify(data);\n        }).catch(err => {\n            // Do something for an error here\n            console.log(err);\n        });\n    };\n\n    document.getElementById(\"paymentsB\").onclick = function () {\n        var apiUrl = 'http://localhost:3000/payments';\n        fetch(apiUrl).then(response => {\n            return response.json();\n        }).then(data => {\n            // Work with JSON data here\n            console.log(data);\n            document.getElementById(\"txtArea\").value = JSON.stringify(data);\n        }).catch(err => {\n            // Do something for an error here\n            console.log(err);\n        });\n    };\n\n\n\n</script>\n\n</body>\n</html>\n"

	fmt.Fprintf(w, page)

}

type Data struct {
	Access_token string `json:"access_token"`
}

var gdata string

func tokens(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tokens")
	//val := getPayments("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6IkZOT3V0dEhhOGRzU0EyIn0.eyJhdWQiOiJGTjJ0ZHQ4QkNEbXJ3ciIsImp0aSI6IkZOT3V0dEhhOGRzU0EyIiwiaWF0IjoxNTk2NzAwMjcxLCJuYmYiOjE1OTY3MDAyNzEsInN1YiI6IiIsImV4cCI6MTYwNDY0OTA3MSwidXNlcl9pZCI6IkFvVTZaMW5HQ0tMY3BDIiwibWVyY2hhbnRfaWQiOiJBb1U2WjdhRGdSbnpEOCIsInNjb3BlcyI6WyJyZWFkX29ubHkiXX0.OgxVcUICrzmPGFalPYpqrPFlsqdGGzOBkZdYh3M4rsfdRQLv9rPpgL3P22kToeL01AYCZZk5pfA6p94-nB_j1lEplOApPWRAnDLuS01IzgXKMTRrmI670iKvULsvUYBhNoiFtQpdS9Mm6KK6HkoKIXx3UvrTLMxF2-Hsgopatfo")
	fmt.Println("Code: ", code)
	val := getToken(code)
	fmt.Println(val)
	var data Data
	err := json.Unmarshal([]byte(val), &data)
	fmt.Println(err)
	fmt.Println(data)
	gdata = data.Access_token
	//fmt.Println(val["access_token"])
	fmt.Fprintf(w, val)
}

func payments(w http.ResponseWriter, r *http.Request) {
	fmt.Println("payments", gdata)
	val := getPayments(gdata)
	fmt.Fprintf(w, val)
}

func getPayments(token string) string {
	url := "https://api.razorpay.com/v1/payments"

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))
	return string([]byte(body))
}

func getToken(code string) string {
	fmt.Println("Start")
	u := "https://auth.razorpay.com/token"
	data := url.Values{}
	data.Set("client_id", "FN2tdt8BCDmrwr")
	data.Set("client_secret", "oSo1RIUii0L7B3lI0AIuabLv")
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "http://localhost:3000/callback")
	data.Set("code", code)
	data.Set("mode", "test")

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode())) // URL-encoded payload
	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)
	//fmt.Println(resp)

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))
	return string([]byte(body))
}



func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/payments", payments)
	http.HandleFunc("/token", tokens)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.ListenAndServe(":"+port, nil)

}

