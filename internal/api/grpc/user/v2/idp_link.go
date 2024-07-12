package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddIDPLink(ctx context.Context, req *user.AddIDPLinkRequest) (_ *user.AddIDPLinkResponse, err error) {
	details, err := s.command.AddUserIDPLink(ctx, req.UserId, "", &command.AddLink{
		IDPID:         req.GetIdpLink().GetIdpId(),
		DisplayName:   req.GetIdpLink().GetUserName(),
		IDPExternalID: req.GetIdpLink().GetUserId(),
	})
	if err != nil {
		return nil, err
	}
	return &user.AddIDPLinkResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListIDPLinks(ctx context.Context, req *user.ListIDPLinksRequest) (_ *user.ListIDPLinksResponse, err error) {
	queries, err := ListLinkedIDPsRequestToQuery(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.IDPUserLinks(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	res.RemoveNoPermission(ctx, s.checkPermission)
	return &user.ListIDPLinksResponse{
		Result:  IDPLinksToPb(res.Links),
		Details: object.ToListDetails(res.SearchResponse),
	}, nil
}

func ListLinkedIDPsRequestToQuery(req *user.ListIDPLinksRequest) (*query.IDPUserLinksSearchQuery, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	userQuery, err := query.NewIDPUserLinksUserIDSearchQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	return &query.IDPUserLinksSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: []query.SearchQuery{userQuery},
	}, nil
}

func IDPLinksToPb(res []*query.IDPUserLink) []*user.IDPLink {
	links := make([]*user.IDPLink, len(res))
	for i, link := range res {
		links[i] = IDPLinkToPb(link)
	}
	return links
}

func IDPLinkToPb(link *query.IDPUserLink) *user.IDPLink {
	return &user.IDPLink{
		IdpId:    link.IDPID,
		UserId:   link.ProvidedUserID,
		UserName: link.ProvidedUsername,
	}
}

func (s *Server) RemoveIDPLink(ctx context.Context, req *user.RemoveIDPLinkRequest) (*user.RemoveIDPLinkResponse, error) {
	objectDetails, err := s.command.RemoveUserIDPLink(ctx, RemoveIDPLinkRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &user.RemoveIDPLinkResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}

func RemoveIDPLinkRequestToDomain(ctx context.Context, req *user.RemoveIDPLinkRequest) *domain.UserIDPLink {
	return &domain.UserIDPLink{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.UserId,
		},
		IDPConfigID:    req.IdpId,
		ExternalUserID: req.LinkedUserId,
	}
}
