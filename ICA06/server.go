package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-martini/martini"

	"net/http"
	"net/url"
)

//Google Translator url
var baseURL = "https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=%s&client=tw-ob"

func main() {

	m := martini.Classic()

	m.Get("/speech/:text", func(params martini.Params, w http.ResponseWriter, r *http.Request) {
		text := params["text"]         //Get the text from the request params
		speech, _ := Speak(text, "no") //Get audio from text
		// Changing сontent-type, otherwise browser does not understand that this is an audio recording
		w.Header().Set("Content-Type", "audio/mpeg")

		speech.WriteTo(w) //In response we send an audio file
	})
	m.RunOnAddr(":8001")
	m.Run()
}

type Speech struct {
	bytes.Buffer
}

// Function for get audio from Google translator
func Speak(text, language string) (*Speech, error) {
	req := fmt.Sprintf(baseURL, url.QueryEscape(text), url.QueryEscape(language)) // Formatting URL with text and language
	res, err := http.Get(req)                                                     // Make GET request
	if err != nil {                                                               //Check for errors
		return nil, err
	}

	speech := &Speech{}                                          // It will be returned as response
	if _, err := io.Copy(&speech.Buffer, res.Body); err != nil { //Read response body and copy it to buffer,also сheck for errors
		return nil, err
	}

	return speech, nil
}
