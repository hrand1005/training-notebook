package handler

import (
	//"fmt"
	"log"
	"net/http"
	//"strconv"

	"github.com/hrand1005/training-notebook/data"
)

type Set struct {
	logger *log.Logger
}

func NewSet(l *log.Logger) *Set {
	return &Set{l}
}

func (s *Set) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.ReadAll(rw, r)
	case http.MethodPost:
		s.Create(rw, r)
	case http.MethodPut:
		s.Update(rw, r)
	/*case http.MethodDelete:
	  s.Delete(rw, r)*/
	default:
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
	}
	return
}

func (s *Set) ReadAll(rw http.ResponseWriter, r *http.Request) {
	if err := data.ToJSON(data.Sets(), rw); err != nil {
		s.logger.Printf("could not serialize sets to json: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
	}

	return
}

func (s *Set) Create(rw http.ResponseWriter, r *http.Request) {
	var newSet data.Set

	if err := data.FromJSON(&newSet, r.Body); err != nil {
		s.logger.Printf("could not bind json to set: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	// assigns ID to newSet
	data.AddSet(&newSet)

	rw.WriteHeader(http.StatusCreated)
	data.ToJSON(newSet, rw)
	return
}

func (s *Set) Update(rw http.ResponseWriter, r *http.Request) {
	var newSet data.Set

	if err := data.FromJSON(&newSet, r.Body); err != nil {
		s.logger.Printf("could not bind json to set: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := data.UpdateSet(&newSet); err != nil {
		s.logger.Printf("could not update set: %v", err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	data.ToJSON(newSet, rw)
	return
}

// REQUIRE ID PARAM
/*
func (s *Set) Read(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

    s, err := data.SetByID(id)
    if err != nil {
        log.Printf("could not read set: %v", err)
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
        return
    }

    c.IndentedJSON(http.StatusOK, s)
    return
}

func (s *Set) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))

    if err := data.DeleteSet(id); err != nil {
        log.Printf("could not delete set: %v", err)
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
        return
    }

    c.IndentedJSON(http.StatusNoContent, gin.H{})
    return
}
*/
