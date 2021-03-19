package main
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
)
import "net/http"
import "github.com/gorilla/mux"



type event struct {
	ID 			string `json:"ID"`
	Title 		string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID: 	"1",
		Title: 	"Introduction to GoLang",
		Description: "This is my Web Service using GoLang",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request)  {
	//fmt.Println(w, "Welcome home!")
	log.Println("Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("/createEvent")

	var newEvent event

	requestBody, err  := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(requestBody, &newEvent)

	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(events)

}

func getOneEvent(w http.ResponseWriter, r *http.Request){
	log.Println("/getOneEvent")
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID{
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request){
	log.Println("/getAllEvents")
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request){
	log.Println("/updateEvent")
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events{
		if singleEvent.ID == eventID {
			singleEvent.Title 		= updatedEvent.Title
			singleEvent.Description = updatedEvent.Description

			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request){
	log.Println("/deleteEvent")

	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}

}


func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", 					homeLink)
	router.HandleFunc("/createEvent",		createEvent).Methods("POST")
	router.HandleFunc("/events", 			getAllEvents).Methods("GET")
	router.HandleFunc("/event/{id}",		getOneEvent).Methods("GET")
	router.HandleFunc("/updateEvent",		updateEvent).Methods("PATCH")
	router.HandleFunc("/deleteEvent/{id}",	deleteEvent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))

}

