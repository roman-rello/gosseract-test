package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	//"github.com/otiai10/gosseract/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	path := os.Getenv("IMAGE_PATH")
	runWithEmbeddedOCR(path)
	runWithOCRServer(path)
}

func runWithEmbeddedOCR(path string) {
	//client := gosseract.NewClient()
	//defer client.Close()
	//client.SetImage(path)
	//text, _ := client.Text()
	//fmt.Println(text)
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func preparePayload(path string) string {
	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += toBase64(bytes)

	return base64Encoding
}

type OCRRequest struct {
	Base64 string `json:"base64"`
	Trim   string `json:"trim"`
}

func runWithOCRServer(path string) {
	url := "http://localhost:8080/base64"
	method := "POST"

	payload := OCRRequest{
		Base64: preparePayload(path),
		Trim:   "\n",
	}

	// Marshal it into JSON prior to requesting
	payloadJSON, err := json.Marshal(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
