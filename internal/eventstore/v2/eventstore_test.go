package eventstore

import "testing"

func TestPushAggregates(t *testing.T) {
	type res struct{}
	type args struct{}
	tests := []struct {
		name string
		args args
		res  *SearchQueryFactory
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// factory := NewSearchQueryFactory(tt.args.aggregateTypes...)
			// for _, setter := range tt.args.setters {
			// 	factory = setter(factory)
			// }
			// if !reflect.DeepEqual(factory, tt.res) {
			// 	t.Errorf("NewSearchQueryFactory() = %v, want %v", factory, tt.res)
			// }
		})
	}
}
