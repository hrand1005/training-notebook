package handler

import (
	//"fmt"
	"log"
	"net/http"
	"strconv"
    "regexp"

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
        rule := regexp.MustCompile(`/([0-9]+)`)
        capture := rule.FindAllStringSubmatch(r.URL.Path, -1)
        if len(capture) != 1 || len(capture[0]) != 2 {
            s.logger.Printf("Invalid URI: %v", r.URL.Path)
            http.Error(rw, "could not capture valid id from request URI", http.StatusBadRequest)
            return
        }

        id, err := strconv.Atoi(capture[0][1])
        if err != nil {
            s.logger.Println("Invalid URI, cannot convert to integer", capture[0][1])
            http.Error(rw, "could not capture valid id from request URI", http.StatusBadRequest)
            return
        }

		s.Update(id, rw, r)

	/*case http.MethodDelete:
	  s.Delete(rw, r)*/
	default:
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
	}
	return
}

func (s *Set) ReadAll(rw http.ResponseWriter, r *http.Request) {
    s.logger.Println("Reading All sets")

	if err := data.ToJSON(data.Sets(), rw); err != nil {
        http.Error(rw, "could not serialize sets to json", http.StatusBadRequest)
	}

	return
}

func (s *Set) Create(rw http.ResponseWriter, r *http.Request) {
    s.logger.Println("Creating set")
	var newSet data.Set

	if err := data.FromJSON(&newSet, r.Body); err != nil {
        http.Error(rw, "could not bind json to set", http.StatusBadRequest)
		return
	}
	// assigns ID to newSet
	data.AddSet(&newSet)

	rw.WriteHeader(http.StatusCreated)
	data.ToJSON(newSet, rw)
    s.logger.Printf("Created set: %#v\n", newSet)
	return
}

// This should be the handler for a URI that ends with /id
// ID fields, if not intended to be updated by the user, should not be present
// in the body of the request. 
func (s *Set) Update(id int, rw http.ResponseWriter, r *http.Request) {
    s.logger.Println("Updating set")
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
    s.logger.Printf("Updated set: %#v\n", newSet)
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
