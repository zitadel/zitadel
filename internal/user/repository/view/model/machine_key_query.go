package model

//
//type MachineKeySearchRequest usr_model.MachineKeySearchRequest
//type MachineKeySearchQuery usr_model.MachineKeySearchQuery
//type MachineKeySearchKey usr_model.MachineKeySearchKey
//
//func (req MachineKeySearchRequest) GetLimit() uint64 {
//	return req.Limit
//}
//
//func (req MachineKeySearchRequest) GetOffset() uint64 {
//	return req.Offset
//}
//
//func (req MachineKeySearchRequest) GetSortingColumn() repository.ColumnKey {
//	if req.SortingColumn == usr_model.MachineKeyKeyUnspecified {
//		return nil
//	}
//	return MachineKeySearchKey(req.SortingColumn)
//}
//
//func (req MachineKeySearchRequest) GetAsc() bool {
//	return req.Asc
//}
//
//func (req MachineKeySearchRequest) GetQueries() []repository.SearchQuery {
//	result := make([]repository.SearchQuery, len(req.Queries))
//	for i, q := range req.Queries {
//		result[i] = MachineKeySearchQuery{Key: q.Key, Value: q.Value, Method: q.Method}
//	}
//	return result
//}
//
//func (req MachineKeySearchQuery) GetKey() repository.ColumnKey {
//	return MachineKeySearchKey(req.Key)
//}
//
//func (req MachineKeySearchQuery) GetMethod() global_model.SearchMethod {
//	return req.Method
//}
//
//func (req MachineKeySearchQuery) GetValue() interface{} {
//	return req.Value
//}
//
//func (key MachineKeySearchKey) ToColumnName() string {
//	switch usr_model.MachineKeySearchKey(key) {
//	case usr_model.MachineKeyKeyID:
//		return MachineKeyKeyID
//	case usr_model.MachineKeyObjectID:
//		return MachineKeyObjectID
//	case usr_model.MachineKeyObjectType:
//		return MachineKeyObjectType
//	default:
//		return ""
//	}
//}
