package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// mapping URL with short url
var urlDB = make(map[string]URL)

// func that short the original url
func generateShortURL(OriginalURL string) string {
	hasher := md5.New()               //md5 type of cheksum that help string to convert to particular hash
	hasher.Write([]byte(OriginalURL)) //It converts the originalURL string to a byte slice
	data := hasher.Sum(nil)           // Sum appends the current hash to data and returns the resulting slice
	hash := hex.EncodeToString(data)  //converts the data into string
	fmt.Println("Encode to string:", hash)
	fmt.Println("final string:", hash[:8]) //only 8 characters
	return hash[:8]
}

// after shorting create url
func createURL(originalURL string) string {
	shortUrl := generateShortURL(originalURL)
	id := shortUrl   //in database some id need to be stored any id for simplicity we use this
	urlDB[id] = URL{ //map containing key=id and value struct URL
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

// any string given and finding which url is ther in correspondance to it
func getURL(id string) (URL, error) {
	url, ok := urlDB[id] //finding data in map
	if !ok {             //if map is empty
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

// golang handle page by handle function
func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!") // write this on server
}

// handle short url and convert it into actual url
func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct { //data expect of url type and var bcz only one data is there
		URL string `json."url"`
	}
	//request to decode and then enter in data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest) //--.badreq... i.e 400 send krna
		return
	}
	shortUrl_ := createURL(data.URL)
	fmt.Fprintf(w, shortUrl_)
	response := struct {
		ShortURL string `json:"short_url`
	}{ShortURL: shortUrl_}
	w.Header().Set("Content_Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):] //redirect ke jitna b hoga le lega
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// fmt.Println("Starting URL shortener")
	// OriginalURL := "https://github.com/KambojKhushi"
	// generateShortURL(OriginalURL)

	//Register the handler function to handle all requests to the root URL("/")
	http.HandleFunc("/shorten", shortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)
	http.HandleFunc("/", RootPageURL) // keep this last

	//Start the HTTP server on port 8080 or 3000
	fmt.Println("Starting server on port 3000....")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error on starting server :", err)
	}

}
