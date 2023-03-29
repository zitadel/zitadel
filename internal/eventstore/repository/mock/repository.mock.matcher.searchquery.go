package mock

import (
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type filterQueryMatcher repository.SearchQuery

func (f *filterQueryMatcher) String() string {
	var filterLists []string
	for _, filterSlice := range f.Filters {
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
	if len(f.Filters) != len(other.Filters) {
		return false
	}
	for filterSliceIdx, filterSlice := range f.Filters {
		if len(filterSlice) != len(other.Filters[filterSliceIdx]) {
			return false
		}
		for filterIdx, filter := range f.Filters[filterSliceIdx] {
			if !(*filterMatcher)(filter).Matches(other.Filters[filterSliceIdx][filterIdx]) {
				return false
			}
		}
	}
	return true
}

func (f *filterQueryMatcher) Got(got interface{}) string {
	return (*filterQueryMatcher)(got.(*repository.SearchQuery)).String()
}
