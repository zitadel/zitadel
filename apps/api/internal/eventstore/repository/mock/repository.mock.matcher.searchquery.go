package mock

import (
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type filterQueryMatcher repository.SearchQuery

func (f *filterQueryMatcher) String() string {
	var filterLists []string
	for _, filterSlice := range f.SubQueries {
		var str string
		for _, filter := range filterSlice {
			str += "," + (*filterMatcher)(filter).String()
		}
		filterLists = append(filterLists, fmt.Sprintf("[%s]", strings.TrimPrefix(str, ",")))

	}
	return fmt.Sprintf("Filters: %s", strings.Join(filterLists, " "))
}

func (f *filterQueryMatcher) Matches(x interface{}) bool {
	other := x.(*repository.SearchQuery)
	if len(f.SubQueries) != len(other.SubQueries) {
		return false
	}
	for filterSliceIdx, filterSlice := range f.SubQueries {
		if len(filterSlice) != len(other.SubQueries[filterSliceIdx]) {
			return false
		}
		for filterIdx, filter := range f.SubQueries[filterSliceIdx] {
			if !(*filterMatcher)(filter).Matches(other.SubQueries[filterSliceIdx][filterIdx]) {
				return false
			}
		}
	}
	return true
}

func (f *filterQueryMatcher) Got(got interface{}) string {
	return (*filterQueryMatcher)(got.(*repository.SearchQuery)).String()
}
