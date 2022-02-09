package handler

import (
	//"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hrand1005/training-notebook/data"
)

func CreateSet(rw http.ResponseWriter, r *http.Request) {
	var newSet data.Set

	if err := data.FromJSON(&newSet, r.Body); err != nil {
		http.Error(rw, "could not bind json to set", http.StatusBadRequest)
		return
	}
	// assigns ID to newSet
	data.AddSet(&newSet)

	rw.WriteHeader(http.StatusCreated)
	data.ToJSON(newSet, rw)
	return
}

func ReadSet(rw http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "bad request id", http.StatusBadRequest)
		return
	}

	s, err := data.SetByID(id)
	if err != nil {
		http.Error(rw, "could not find set", http.StatusNotFound)
		return
	}

	data.ToJSON(s, rw)
	return
}

func ReadSets(rw http.ResponseWriter, r *http.Request) {
	if err := data.ToJSON(data.Sets(), rw); err != nil {
		http.Error(rw, "could not serialize sets to json", http.StatusBadRequest)
	}

	return
}

func UpdateSet(rw http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "bad request id", http.StatusBadRequest)
		return
	}

	var newSet data.Set

	if err := data.FromJSON(&newSet, r.Body); err != nil {
		http.Error(rw, "could not bind json to set", http.StatusBadRequest)
		return
	}

	if err := data.UpdateSet(id, &newSet); err != nil {
		http.Error(rw, "could not find set", http.StatusNotFound)
		return
	}

	data.ToJSON(newSet, rw)
	return
}

func DeleteSet(rw http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "bad request id", http.StatusBadRequest)
		return
	}

	if err := data.DeleteSet(id); err != nil {
		http.Error(rw, "could not find set", http.StatusNotFound)
		return
	}

	return
}
