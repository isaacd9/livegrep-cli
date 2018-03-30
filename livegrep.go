package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Livegrep struct {
	URL string
}

type Query struct {
	Term     string
	FoldCase bool
	Regex    bool
	Context  bool
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

func (lg *Livegrep) NewQuery(q string) Query {
	return Query{
		Term:     q,
		FoldCase: false,
		Regex:    false,
		Context:  false,
	}
}

func (lg *Livegrep) Query(q Query) (QueryResponse, error) {
	uri := fmt.Sprintf(
		"https://%s/api/v1/search/linux?q=%s&fold_case=%t&regex=%t&context=%t",
		lg.URL,
		q.Term,
		q.FoldCase,
		q.Regex,
		q.Context,
	)

	resp, err := http.Get(uri)
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
