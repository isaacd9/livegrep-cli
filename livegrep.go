package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Livegrep struct {
	URL      string
	UseHTTPS bool
	Client   *http.Client
}

type Query struct {
	Term     string
	FoldCase bool
	Regex    bool
}

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

func NewLivegrep(url string) Livegrep {
	return Livegrep{
		URL:      url,
		UseHTTPS: true,
		Client:   &http.Client{},
	}
}

func (lg *Livegrep) NewQuery(q string) Query {
	return Query{
		Term:     q,
		FoldCase: false,
		Regex:    true,
	}
}

func (lg *Livegrep) Query(q Query) (QueryResponse, error) {
	var protocol string
	if lg.UseHTTPS {
		protocol = "https"
	} else {
		protocol = "http"
	}
	uri := fmt.Sprintf(
		"%s://%s/api/v1/search/linux?q=%s&fold_case=%t&regex=%t",
		protocol,
		lg.URL,
		q.Term,
		q.FoldCase,
		q.Regex,
	)

	resp, err := lg.Client.Get(uri)
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
