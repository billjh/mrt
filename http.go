package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//// v1 navigate by stops
type navigateV1Request struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	All         bool   `json:"all"`
}

type navigateV1Response struct {
	Source            string   `json:"source"`
	Destination       string   `json:"destination"`
	StationsTravelled int      `json:"stations_travelled"`
	Route             []string `json:"route"`
	Instructions      []string `json:"instructions"`
}

func makeV1Response(paths []Path) []navigateV1Response {
	res := []navigateV1Response{}
	for _, path := range paths {
		l := len(path.Stops)
		res = append(res, navigateV1Response{
			Source:            path.Stops[0].(Station).name,
			Destination:       path.Stops[l-1].(Station).name,
			StationsTravelled: l - 1,
			Route:             makeRoute(path),
			Instructions:      makeInstructions(path),
		})
	}
	return res
}

func makeRoute(path Path) []string {
	r := []string{}
	for _, s := range path.Stops {
		r = append(r, s.(Station).id.String())
	}
	return r
}

func makeInstructions(path Path) []string {
	r := []string{}
	for i := 1; i < len(path.Stops); i++ {
		prev := path.Stops[i-1].(Station)
		next := path.Stops[i].(Station)
		if prev.id.line == next.id.line {
			r = append(r, fmt.Sprintf("Take %s line from %s to %s", prev.id.line, prev.name, next.name))
		} else {
			r = append(r, fmt.Sprintf("Change from %s line to %s line", prev.id.line, next.id.line))
		}
	}
	return r
}

func (n *Navigator) handleV1(w http.ResponseWriter, r *http.Request) {
	// decode body for request
	nr := navigateV1Request{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&nr); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// run navigator
	paths, err := n.NavigateByStops(nr.Source, nr.Destination, nr.All)
	if err != nil {
		switch err {
		case ErrorSourceNotFound, ErrorDestinationNotFound, ErrorSourceDestinationSame:
			respondError(w, http.StatusBadRequest, err.Error())
		case ErrorPathNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			// unexpected errors
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, makeV1Response(paths))
}

//// v2 navigate by time
type navigateV2Request struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Time        string `json:"time"`
	All         bool   `json:"all"`
}

type navigateV2Response struct {
	Source       string   `json:"source"`
	Destination  string   `json:"destination"`
	Minutes      int      `json:"minutes"`
	Route        []string `json:"route"`
	Instructions []string `json:"instructions"`
}

func makeV2Response(paths []Path) []navigateV2Response {
	res := []navigateV2Response{}
	for _, path := range paths {
		l := len(path.Stops)
		res = append(res, navigateV2Response{
			Source:       path.Stops[0].(Station).name,
			Destination:  path.Stops[l-1].(Station).name,
			Minutes:      int(path.Weight),
			Route:        makeRoute(path),
			Instructions: makeInstructions(path),
		})
	}
	return res
}

func (n *Navigator) handleV2(w http.ResponseWriter, r *http.Request) {
	// decode body for request
	nr := navigateV2Request{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&nr); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	t, err := time.Parse("2006-01-02T15:04", nr.Time)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// run navigator
	paths, err := n.NavigateByTime(nr.Source, nr.Destination, t, nr.All)
	if err != nil {
		switch err {
		case ErrorSourceNotFound, ErrorDestinationNotFound, ErrorSourceDestinationSame:
			respondError(w, http.StatusBadRequest, err.Error())
		case ErrorPathNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			// unexpected errors
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, makeV2Response(paths))
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
