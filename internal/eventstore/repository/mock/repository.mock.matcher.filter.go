package mock

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

var _ gomock.Matcher = (*filterMatcher)(nil)
var _ gomock.GotFormatter = (*filterMatcher)(nil)

type filterMatcher repository.Filter

func (f *filterMatcher) String() string {
	jsonValue, err := json.Marshal(f.Value)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d %d (content=%+v,type=%T,json=%s)", f.Field, f.Operation, f.Value, f.Value, string(jsonValue))
}

func (f *filterMatcher) Matches(x interface{}) bool {
	other := x.(*repository.Filter)
	return f.Field == other.Field && f.Operation == other.Operation && reflect.DeepEqual(f.Value, other.Value)
}

func (f *filterMatcher) Got(got interface{}) string {
	return (*filterMatcher)(got.(*repository.Filter)).String()
}
