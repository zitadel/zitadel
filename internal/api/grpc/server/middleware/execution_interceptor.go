package middleware

/*
func ExecutionInterceptor(queries *query.Queries) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		typeSearchQuery, err := query.NewExecutionTypeSearchQuery(domain.ExecutionTypeRequest)
		if err != nil {
			return nil, err
		}
		typeSearchQuery, err := query.new(domain.ExecutionTypeRequest)
		if err != nil {
			return nil, err
		}

		searchQuery := &query.ExecutionSearchQueries{
			Queries: []query.SearchQuery{typeSearchQuery},
		}
		queries.SearchExecutions(ctx)

		return

	}
}
*/
