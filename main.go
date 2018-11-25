package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int             // Seconds
}

type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type eventRequest struct {
	EventType  string    `json:"eventType"`
	WebsiteUrl string    `json:"websiteUrl"`
	SessionID  string    `json:"sessionID"`
	Pasted     bool      `json:"pasted"`
	InputId    string    `json:"inputId"`
	Time       int       `json:"time"`
	ResizeFrom Dimension `json:"resizeFrom"`
	ResizeTo   Dimension `json:"resizeTo"`
}

var sessionEvents map[string]Data

func init() {
	sessionEvents = make(map[string]Data)
}

func main() {
	http.HandleFunc("/log", logHandler)

	http.ListenAndServe(":3030", nil)
}

func logHandler(w http.ResponseWriter, r *http.Request) {

	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Only allow POST Requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read body"))
		return
	}

	req := &eventRequest{}

	if err = json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to unmarshal JSON request"))
		return
	}

	session, sessionExists := sessionEvents[req.SessionID]

	if !sessionExists {
		session = Data{
			SessionId:          req.SessionID,
			WebsiteUrl:         req.WebsiteUrl,
			ResizeFrom:         Dimension{},
			ResizeTo:           Dimension{},
			CopyAndPaste:       make(map[string]bool),
			FormCompletionTime: 0,
		}
	}

	switch req.EventType {
	case "copyAndPaste":
		session.CopyAndPaste[req.InputId] = req.Pasted
	case "resize":
		session.ResizeFrom = req.ResizeFrom
		session.ResizeTo = req.ResizeTo
	case "elapsedTime":
		session.FormCompletionTime = req.Time
	}

	sessionEvents[req.SessionID] = session

	sessionReport, err := prettyPrint(session)
	if err != nil {
		fmt.Println(session)
	} else {
		fmt.Println(sessionReport)
	}

	w.WriteHeader(http.StatusOK)
}

func prettyPrint(data interface{}) (string, error) {
	s, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", err
	}
	return string(s), nil
}
