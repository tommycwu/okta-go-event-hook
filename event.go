package main

import (
	"bytes"
	"encoding/json"
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

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return getHandler(req)
	case "POST":
		return postHandler(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func getHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var oneTimeChallenge = request.Headers["x-okta-verification-challenge"]
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"verification": oneTimeChallenge,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}

	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func postHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer

	createUser()

	body, err := json.Marshal(map[string]interface{}{
		"message": "UserCreated",
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}

	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(router)
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

	requestBody := strings.NewReader(`{"profile": {"firstName": "` + fName + `","lastName": "` +
		lName + `","email": "` + uName + `@mailinator.com","login": "` + uName + `@mailinator.com"}}`)

	url := "https://[your.okta.org]/api/v1/users?activate=false"

	req, err := http.NewRequest("POST", url, requestBody) //bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "SSWS [api_key]")
	log.Println("req.Header.Set(Authorization)")

	client := &http.Client{Timeout: time.Second * 10}
	log.Println("client.Timeout.Set")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	defer resp.Body.Close()
}
