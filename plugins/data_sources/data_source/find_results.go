package data_source

type FindResults struct {
	results map[string]*FindResult
}

func NewFindResults(results []*FindResult) *FindResults {
	fr := &FindResults{
		results: make(map[string]*FindResult),
	}

	for _, r := range results {
		fr.results[r.Id] = r
	}

	return fr
}

func (f FindResults) Get(id string) *FindResult {
	return f.results[id]
}

func (f FindResults) Len() int {
	return len(f.results)
}

func (f FindResults) All() []*FindResult {
	results := make([]*FindResult, 0)

	for _, r := range f.results {
		results = append(results, r)
	}

	return results
}
