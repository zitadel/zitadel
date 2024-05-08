package eventstore

type Aggregate struct {
	ID       string
	Type     string
	Instance string
	Owner    string
}

func (agg *Aggregate) Equals(aggregate *Aggregate) bool {
	if aggregate.ID != "" && aggregate.ID != agg.ID {
		return false
	}
	if aggregate.Type != "" && aggregate.Type != agg.Type {
		return false
	}
	if aggregate.Instance != "" && aggregate.Instance != agg.Instance {
		return false
	}
	if aggregate.Owner != "" && aggregate.Owner != agg.Owner {
		return false
	}
	return true
}
