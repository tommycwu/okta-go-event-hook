package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func handler(request context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello AWS Lambda and Netlify",
	}, nil
}

func main() {
	lambda.Start(handler)
	//http.HandleFunc("/userTransfer", handleRequests)

	//var err = http.ListenAndServe(":10000", nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	log.Println(time.Now().String() + " " + r.Method)

	if r.Method == "POST" {
		log.Println("POST - handleRequests()")
		createUser()
		w.Header().Set("Service", "Okta Event Hook")
		w.WriteHeader(200)
		fmt.Fprintf(w, "HTTP/1.1 200 OK")
		return
	}

	http.Error(w, r.Method+" is not supported.", http.StatusNotFound)
}

func createUser() {

	var fName = RandomString(4)
	var lName = RandomString(5)
	var uName = fName + "." + lName

	log.Println("uName - " + uName)
	requestBody := strings.NewReader(`{"profile": {"firstName": "` + fName + `","lastName": "` +
		lName + `","email": "` + uName + `@mailinator.com","login": "` + uName + `@mailinator.com"}}`)

	url := "https://dev-489843.okta.com/api/v1/users?activate=false"

	log.Println("POST - createUser()")
	req, err := http.NewRequest("POST", url, requestBody) //bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	log.Println("req.created")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "SSWS 00Yt1WYxQAXDJu_7Uj8ih0QRy_go01lCnKX93lp0su")
	log.Println("req.Header.Set(Authorization)")

	client := &http.Client{Timeout: time.Second * 10}
	log.Println("client.Timeout.Set")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	log.Println("response Status: " + resp.Status)
	defer resp.Body.Close()
	log.Println("resp.Body.Close()")
}
