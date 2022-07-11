package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Event struct {
	Id          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Start       string     `json:"start"`
	End         string     `json:"end"`
	Schedules   []Schedule `json:"schedules"`
	Attendees   []Attendee `json:"attendees"`
	IsOpen      bool       `json:"isOpen"`
}

type Schedule struct {
	End   string `json:"end"`
	Start string `json:"start"`
	Title string `json:"title"`
	Note  string `json:"note"`
}

type Attendee struct {
	Id      string `json:"id"`
	EventId string `json:"eventId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type Events []Event

var events = Events{
	{
		Id:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to event to eventually try it out.",
	},
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/events", Index).Methods("GET")
	router.HandleFunc("/events", Store).Methods("POST")
	router.HandleFunc("/events/{id}", Show).Methods("GET")
	router.HandleFunc("/events/{id}", Update).Methods("PATCH")
	router.HandleFunc("/events/{id}", Delete).Methods("DELETE")

	router.HandleFunc("/events/{event_id}/attendees", StoreAttendee).Methods("POST")
	router.HandleFunc("/events/{event_id}/attendees/{id}", removeAttendee).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		return
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	Id := mux.Vars(r)["id"]

	var event Event

	for _, e := range events {
		if e.Id == Id {
			event = e
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(event)
	if err != nil {
		return
	}

}

func Store(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("unprocessable entity")
	}
	var event Event
	err = json.Unmarshal(reqBody, &event)
	if err != nil {
		fmt.Println("Error", err)
	}
	events = append(events, event)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func Update(w http.ResponseWriter, r *http.Request) {
	Id := mux.Vars(r)["id"]
	reqBody, err := ioutil.ReadAll(r.Body)
	var updatedEvent Event
	if err != nil {
		return
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, e := range events {
		if e.Id == Id {
			updatedEvent.Id = Id
			events = append(events[:i], updatedEvent)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEvent)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	Id := mux.Vars(r)["id"]
	for i, e := range events {
		if e.Id == Id {
			events = append(events[:i], events[i+1:]...)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

}

func StoreAttendee(w http.ResponseWriter, r *http.Request) {
	EventId := mux.Vars(r)["event_id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var attendee Attendee

	json.Unmarshal(reqBody, &attendee)

	for i, e := range events {
		if e.Id == EventId {
			attendee.EventId = EventId
			e.Attendees = append(events[i].Attendees, attendee)
			events = append(events[:i], e)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attendee)
}

func removeAttendee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	EventId := params["event_id"]
	AttendeeId := params["id"]

	var event Event
	for _, e := range events {
		if e.Id == EventId {
			event = e
		}
	}

	for i, a := range event.Attendees {
		if a.Id == AttendeeId {
			if len(event.Attendees) <= 1 {
				event.Attendees = nil
			} else {
				event.Attendees = append(event.Attendees[:i], event.Attendees[i+i:]...)
			}
			events = append(events[:i], event)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(event)

}
