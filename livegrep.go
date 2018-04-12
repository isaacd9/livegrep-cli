package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const livegrepPath = "/api/v1/search"

// Livegrep is a client object for connecting to a Livegrep instance
type Livegrep struct {
	Host     string
	UseHTTPS bool
	Client   *http.Client
}

// Query represents a single query against a Livegrep instance
type Query struct {
	Term     string
	FoldCase bool
	Regex    bool
}

// QueryResult is a single match in a result from Livegrep
type QueryResult struct {
	Tree          string   `json:"tree"`
	Version       string   `json:"version"`
	Path          string   `json:"path"`
	Lno           int      `json:"lno"`
	ContextBefore []string `json:"context_before"`
	Line          string   `json:"line"`
	ContextAfter  []string `json:"context_after"`
	Bounds        []int    `json:"bounds"`
}

// QueryResponse is the overall result from a query against a Livegrep instance
type QueryResponse struct {
	Re2Tme      int    `json:"re2_time"`
	GitTme      int    `json:"git_time"`
	SortTime    int    `json:"sort_time"`
	IndexTime   int    `json:"index_time"`
	AnalyzeTime int    `json:"analyze_time"`
	Why         string `json:"why"`

	Results     []QueryResult `json:"results"`
	FileResults []QueryResult `json:"file_results"`

	SearchType string `json:"search_type"`
}

// NewLivegrep returns a new Livegrep client
func NewLivegrep(host string) Livegrep {
	return Livegrep{
		Host:     host,
		UseHTTPS: true,
		Client:   &http.Client{},
	}
}

// NewQuery returns a new query for Livegrep
func (lg *Livegrep) NewQuery(q string) Query {
	return Query{
		Term:     q,
		FoldCase: false,
		Regex:    true,
	}
}

// Query runs the specified Query against the given Livegrep instance
func (lg *Livegrep) Query(q Query) (QueryResponse, error) {
	var protocol string
	if lg.UseHTTPS {
		protocol = "https"
	} else {
		protocol = "http"
	}

	query := fmt.Sprintf(
		"q=%s&fold_case=%t&regex=%t",
		q.Term,
		q.FoldCase,
		q.Regex,
	)

	uri := &url.URL{
		Scheme:   protocol,
		Host:     lg.Host,
		Path:     livegrepPath,
		RawQuery: query,
	}

	resp, err := lg.Client.Get(uri.String())
	if err != nil {
		return QueryResponse{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return QueryResponse{}, err
	}

	var results QueryResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return QueryResponse{}, err
	}

	return results, nil
}
