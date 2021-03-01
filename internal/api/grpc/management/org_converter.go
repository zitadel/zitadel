package management

//func ListOrgDomainsRequestToModel(req *mgmt_pb.ListOrgDomainsRequest) (*org_model.OrgDomainSearchRequest, error) {
//	queries, err := org_grpc.DomainQueriesToModel(req.Queries)
//	if err != nil {
//		return nil, err
//	}
//	return &org_model.OrgDomainSearchRequest{
//		Offset: req.MetaData.Offset,
//		Limit:  uint64(req.MetaData.Limit),
//		Asc:    req.MetaData.Asc,
//		//SortingColumn: //TODO: sorting
//		Queries: queries,
//	}, nil
//}
