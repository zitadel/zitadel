package admin

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zitadel/logging"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	action_grpc "github.com/zitadel/zitadel/internal/api/grpc/action"
	"github.com/zitadel/zitadel/internal/api/grpc/authn"
	"github.com/zitadel/zitadel/internal/api/grpc/management"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	management_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/policy"
	v1_pb "github.com/zitadel/zitadel/pkg/grpc/v1"
)

type importResponse struct {
	ret   *admin_pb.ImportDataResponse
	count *counts
	err   error
}
type counts struct {
	humanUserCount          int
	humanUserLen            int
	machineUserCount        int
	machineUserLen          int
	userMetadataCount       int
	userMetadataLen         int
	userLinksCount          int
	userLinksLen            int
	projectCount            int
	projectLen              int
	oidcAppCount            int
	oidcAppLen              int
	apiAppCount             int
	apiAppLen               int
	actionCount             int
	actionLen               int
	projectRolesCount       int
	projectRolesLen         int
	projectGrantCount       int
	projectGrantLen         int
	userGrantCount          int
	userGrantLen            int
	projectMembersCount     int
	projectMembersLen       int
	orgMemberCount          int
	orgMemberLen            int
	projectGrantMemberCount int
	projectGrantMemberLen   int
	appKeysCount            int
	machineKeysCount        int
}

func (c *counts) getProgress() string {
	return "progress:" +
		"human_users " + strconv.Itoa(c.humanUserCount) + "/" + strconv.Itoa(c.humanUserLen) + ", " +
		"machine_users " + strconv.Itoa(c.machineUserCount) + "/" + strconv.Itoa(c.machineUserLen) + ", " +
		"user_metadata " + strconv.Itoa(c.userMetadataCount) + "/" + strconv.Itoa(c.userMetadataLen) + ", " +
		"user_links " + strconv.Itoa(c.userLinksCount) + "/" + strconv.Itoa(c.userLinksLen) + ", " +
		"projects " + strconv.Itoa(c.projectCount) + "/" + strconv.Itoa(c.projectLen) + ", " +
		"oidc_apps " + strconv.Itoa(c.oidcAppCount) + "/" + strconv.Itoa(c.oidcAppLen) + ", " +
		"api_apps " + strconv.Itoa(c.apiAppCount) + "/" + strconv.Itoa(c.apiAppLen) + ", " +
		"actions " + strconv.Itoa(c.actionCount) + "/" + strconv.Itoa(c.actionLen) + ", " +
		"project_roles " + strconv.Itoa(c.projectRolesCount) + "/" + strconv.Itoa(c.projectRolesLen) + ", " +
		"project_grant " + strconv.Itoa(c.projectGrantCount) + "/" + strconv.Itoa(c.projectGrantLen) + ", " +
		"user_grants " + strconv.Itoa(c.userGrantCount) + "/" + strconv.Itoa(c.userGrantLen) + ", " +
		"project_members " + strconv.Itoa(c.projectMembersCount) + "/" + strconv.Itoa(c.projectMembersLen) + ", " +
		"org_members " + strconv.Itoa(c.orgMemberCount) + "/" + strconv.Itoa(c.orgMemberLen) + ", " +
		"project_grant_members " + strconv.Itoa(c.projectGrantMemberCount) + "/" + strconv.Itoa(c.projectGrantMemberLen)
}

func (s *Server) ImportData(ctx context.Context, req *admin_pb.ImportDataRequest) (_ *admin_pb.ImportDataResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if req.GetDataOrgs() != nil || req.GetDataOrgsv1() != nil {
		timeoutDuration, err := time.ParseDuration(req.Timeout)
		if err != nil {
			return nil, err
		}
		ch := make(chan importResponse, 1)
		ctxTimeout, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		go func() {
			orgs := make([]*admin_pb.DataOrg, 0)
			if req.GetDataOrgsv1() != nil {
				dataOrgs, err := s.dataOrgsV1ToDataOrgs(ctx, req.GetDataOrgsv1())
				if err != nil {
					ch <- importResponse{ret: nil, err: err}
					return
				}
				orgs = dataOrgs.GetOrgs()
			} else {
				orgs = req.GetDataOrgs().GetOrgs()
			}

			ret, count, err := s.importData(ctx, orgs)
			ch <- importResponse{ret: ret, count: count, err: err}
		}()

		select {
		case <-ctxTimeout.Done():
			logging.Errorf("Import to response timeout: %v", ctxTimeout.Err())
			return nil, ctxTimeout.Err()
		case result := <-ch:
			logging.OnError(result.err).Errorf("error while importing: %v", result.err)
			if result.count != nil {
				logging.Infof("Import done: %s", result.count.getProgress())
			}
			return result.ret, result.err
		}
	} else {
		v1Transformation := false
		var gcsInput *admin_pb.ImportDataRequest_GCSInput
		var s3Input *admin_pb.ImportDataRequest_S3Input
		var localInput *admin_pb.ImportDataRequest_LocalInput
		if req.GetDataOrgsGcs() != nil {
			gcsInput = req.GetDataOrgsGcs()
		}
		if req.GetDataOrgsv1Gcs() != nil {
			gcsInput = req.GetDataOrgsv1Gcs()
			v1Transformation = true
		}
		if req.GetDataOrgsS3() != nil {
			s3Input = req.GetDataOrgsS3()
		}
		if req.GetDataOrgsv1S3() != nil {
			s3Input = req.GetDataOrgsv1S3()
			v1Transformation = true
		}
		if req.GetDataOrgsLocal() != nil {
			localInput = req.GetDataOrgsLocal()
		}
		if req.GetDataOrgsv1Local() != nil {
			localInput = req.GetDataOrgsv1Local()
			v1Transformation = true
		}

		timeoutDuration, err := time.ParseDuration(req.Timeout)
		if err != nil {
			return nil, err
		}
		dctx := authz.Detach(ctx)
		go func() {
			ch := make(chan importResponse, 1)
			ctxTimeout, cancel := context.WithTimeout(dctx, timeoutDuration)
			defer cancel()
			go func() {
				dataOrgs, err := s.transportDataFromFile(ctxTimeout, v1Transformation, gcsInput, s3Input, localInput)
				if err != nil {
					ch <- importResponse{nil, nil, err}
					return
				}
				resp, count, err := s.importData(ctxTimeout, dataOrgs)
				if err != nil {
					ch <- importResponse{nil, count, err}
					return
				}
				ch <- importResponse{resp, count, nil}
			}()

			select {
			case <-ctxTimeout.Done():
				logging.Errorf("Export to response timeout: %v", ctxTimeout.Err())
				return
			case result := <-ch:
				logging.OnError(result.err).Errorf("error while importing: %v", err)
				if result.count != nil {
					logging.Infof("Import done: %s", result.count.getProgress())
				}
			}
		}()
	}
	return &admin_pb.ImportDataResponse{}, nil
}

func (s *Server) transportDataFromFile(ctx context.Context, v1Transformation bool, gcsInput *admin_pb.ImportDataRequest_GCSInput, s3Input *admin_pb.ImportDataRequest_S3Input, localInput *admin_pb.ImportDataRequest_LocalInput) (_ []*admin_pb.DataOrg, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	dataOrgs := make([]*admin_pb.DataOrg, 0)
	data := make([]byte, 0)
	if gcsInput != nil {
		gcsData, err := getFileFromGCS(ctx, gcsInput)
		if err != nil {
			return nil, err
		}
		data = gcsData
	}
	if s3Input != nil {
		s3Data, err := getFileFromS3(ctx, s3Input)
		if err != nil {
			return nil, err
		}
		data = s3Data
	}
	if localInput != nil {
		localData, err := os.ReadFile(localInput.Path)
		if err != nil {
			return nil, err
		}
		data = localData
	}

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	if v1Transformation {
		dataImportV1 := new(v1_pb.ImportDataOrg)
		if err := jsonpb.Unmarshal(data, dataImportV1); err != nil {
			return nil, err
		}

		dataImport, err := s.dataOrgsV1ToDataOrgs(ctx, dataImportV1)
		if err != nil {
			return nil, err
		}
		dataOrgs = dataImport.Orgs
	} else {
		dataImport := new(admin_pb.ImportDataOrg)
		if err := jsonpb.Unmarshal(data, dataImport); err != nil {
			return nil, err
		}
		dataOrgs = dataImport.Orgs
	}

	return dataOrgs, nil
}

func getFileFromS3(ctx context.Context, input *admin_pb.ImportDataRequest_S3Input) (_ []byte, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	minioClient, err := minio.New(input.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(input.AccessKeyId, input.SecretAccessKey, ""),
		Secure: input.Ssl,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(ctx, input.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("bucket not existing: %v", err)
	}

	object, err := minioClient.GetObject(ctx, input.Bucket, input.Path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer object.Close()
	return io.ReadAll(object)
}

func getFileFromGCS(ctx context.Context, input *admin_pb.ImportDataRequest_GCSInput) (_ []byte, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	saJson, err := base64.StdEncoding.DecodeString(input.ServiceaccountJson)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(saJson))
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(input.Bucket)
	reader, err := bucket.Object(input.Path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func importOrg1(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, ctxData authz.CtxData, org *admin_pb.DataOrg, success *admin_pb.ImportDataSuccess, count *counts, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessInitCode crypto.Generator) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	_, err = s.command.AddOrgWithID(ctx, org.GetOrg().GetName(), ctxData.UserID, ctxData.ResourceOwner, org.GetOrgId(), []string{})
	if err != nil {
		*errors = append(*errors, &admin_pb.ImportDataError{Type: "org", Id: org.GetOrgId(), Message: err.Error()})
		if _, err := s.query.OrgByID(ctx, true, org.OrgId); err != nil {
			// TODO: Only nil if err != not found
			return nil
		}
	}
	successOrg := &admin_pb.ImportDataSuccessOrg{
		OrgId:               org.GetOrgId(),
		ProjectIds:          []string{},
		OidcAppIds:          []string{},
		ApiAppIds:           []string{},
		HumanUserIds:        []string{},
		MachineUserIds:      []string{},
		ActionIds:           []string{},
		ProjectGrants:       []*admin_pb.ImportDataSuccessProjectGrant{},
		UserGrants:          []*admin_pb.ImportDataSuccessUserGrant{},
		OrgMembers:          []string{},
		ProjectMembers:      []*admin_pb.ImportDataSuccessProjectMember{},
		ProjectGrantMembers: []*admin_pb.ImportDataSuccessProjectGrantMember{},
	}
	logging.Debugf("successful org: %s", successOrg.OrgId)
	success.Orgs = append(success.Orgs, successOrg)

	domainPolicy := org.GetDomainPolicy()
	if org.DomainPolicy != nil {
		_, err := s.command.AddOrgDomainPolicy(ctx, org.GetOrgId(), domainPolicy.UserLoginMustBeDomain, domainPolicy.ValidateOrgDomains, domainPolicy.SmtpSenderAddressMatchesInstanceDomain)
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "domain_policy", Id: org.GetOrgId(), Message: err.Error()})
		}
	}
	return importResources(ctx, s, errors, successOrg, org, count, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessInitCode)
}

func importLabelPolicy(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.LabelPolicy == nil {
		return nil
	}
	_, err = s.command.AddLabelPolicy(ctx, org.GetOrgId(), management.AddLabelPolicyToDomain(org.GetLabelPolicy()))
	if err != nil {
		*errors = append(*errors, &admin_pb.ImportDataError{Type: "label_policy", Id: org.GetOrgId(), Message: err.Error()})
		if isCtxTimeout(ctx) {
			return err
		}
	} else {
		_, err = s.command.ActivateLabelPolicy(ctx, org.GetOrgId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "label_policy", Id: org.GetOrgId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
		}
	}
	return nil
}

func importLockoutPolicy(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.LockoutPolicy == nil {
		return
	}
	_, err := s.command.AddLockoutPolicy(ctx, org.GetOrgId(), management.AddLockoutPolicyToDomain(org.GetLockoutPolicy()))
	if err != nil {
		*errors = append(*errors, &admin_pb.ImportDataError{Type: "lockout_policy", Id: org.GetOrgId(), Message: err.Error()})
	}
}

func importOidcIdps(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.OidcIdps == nil {
		return nil
	}
	for _, idp := range org.OidcIdps {
		logging.Debugf("import oidcidp: %s", idp.IdpId)
		_, err := s.command.ImportIDPConfig(ctx, management.AddOIDCIDPRequestToDomain(idp.Idp), idp.IdpId, org.GetOrgId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "oidc_idp", Id: idp.IdpId, Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		logging.Debugf("successful oidcidp: %s", idp.GetIdpId())
		successOrg.OidcIpds = append(successOrg.OidcIpds, idp.GetIdpId())
	}
	return nil
}

func importJwtIdps(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.JwtIdps == nil {
		return nil
	}
	for _, idp := range org.JwtIdps {
		logging.Debugf("import jwtidp: %s", idp.IdpId)
		_, err := s.command.ImportIDPConfig(ctx, management.AddJWTIDPRequestToDomain(idp.Idp), idp.IdpId, org.GetOrgId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "jwt_idp", Id: idp.IdpId, Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		logging.Debugf("successful jwtidp: %s", idp.GetIdpId())
		successOrg.JwtIdps = append(successOrg.JwtIdps, idp.GetIdpId())
	}
	return nil
}

func importLoginPolicy(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.LoginPolicy == nil {
		return
	}
	_, err := s.command.AddLoginPolicy(ctx, org.GetOrgId(), management.AddLoginPolicyToCommand(org.GetLoginPolicy()))
	if err != nil {
		*errors = append(*errors, &admin_pb.ImportDataError{Type: "login_policy", Id: org.GetOrgId(), Message: err.Error()})
	}
}

func importPwComlexityPolicy(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.PasswordComplexityPolicy == nil {
		return
	}
	_, err := s.command.AddPasswordComplexityPolicy(ctx, org.GetOrgId(), management.AddPasswordComplexityPolicyToDomain(org.GetPasswordComplexityPolicy()))
	if err != nil {
		*errors = append(*errors, &admin_pb.ImportDataError{Type: "password_complexity_policy", Id: org.GetOrgId(), Message: err.Error()})
	}
}

func importPrivacyPolicy(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.PrivacyPolicy == nil {
		return
	}
	_, err := s.command.AddPrivacyPolicy(ctx, org.GetOrgId(), management.AddPrivacyPolicyToDomain(org.GetPrivacyPolicy()))
	if err != nil {
		*errors = append(*errors, &admin_pb.ImportDataError{Type: "privacy_policy", Id: org.GetOrgId(), Message: err.Error()})
	}
}

func importHumanUsers(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessInitCode crypto.Generator) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.HumanUsers == nil {
		return nil
	}
	for _, user := range org.GetHumanUsers() {
		logging.Debugf("import user: %s", user.GetUserId())
		human, passwordless, links := management.ImportHumanUserRequestToDomain(user.User)
		human.AggregateID = user.UserId
		_, _, err := s.command.ImportHuman(ctx, org.GetOrgId(), human, passwordless, links, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessInitCode)
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "human_user", Id: user.GetUserId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
		} else {
			count.humanUserCount += 1
			logging.Debugf("successful user %d: %s", count.humanUserCount, user.GetUserId())
			successOrg.HumanUserIds = append(successOrg.HumanUserIds, user.GetUserId())
		}

		if user.User.OtpCode != "" {
			logging.Debugf("import user otp: %s", user.GetUserId())
			if err := s.command.ImportHumanTOTP(ctx, user.UserId, "", org.GetOrgId(), user.User.OtpCode); err != nil {
				*errors = append(*errors, &admin_pb.ImportDataError{Type: "human_user_otp", Id: user.GetUserId(), Message: err.Error()})
				if isCtxTimeout(ctx) {
					return err
				}
			} else {
				logging.Debugf("successful user otp: %s", user.GetUserId())
			}
		}
	}
	return nil
}

func importMachineUsers(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.MachineUsers == nil {
		return nil
	}
	for _, user := range org.GetMachineUsers() {
		logging.Debugf("import user: %s", user.GetUserId())
		_, err := s.command.AddMachine(ctx, management.AddMachineUserRequestToCommand(user.GetUser(), org.GetOrgId()))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "machine_user", Id: user.GetUserId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.machineUserCount += 1
		logging.Debugf("successful user %d: %s", count.machineUserCount, user.GetUserId())
		successOrg.MachineUserIds = append(successOrg.MachineUserIds, user.GetUserId())
	}
	return nil
}

func importUserMetadata(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.UserMetadata == nil {
		return nil
	}
	for _, userMetadata := range org.GetUserMetadata() {
		logging.Debugf("import usermetadata: %s", userMetadata.GetId()+"_"+userMetadata.GetKey())
		_, err := s.command.SetUserMetadata(ctx, &domain.Metadata{Key: userMetadata.GetKey(), Value: userMetadata.GetValue()}, userMetadata.GetId(), org.GetOrgId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "user_metadata", Id: userMetadata.GetId() + "_" + userMetadata.GetKey(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.userMetadataCount += 1
		logging.Debugf("successful usermetadata %d: %s", count.userMetadataCount, userMetadata.GetId()+"_"+userMetadata.GetKey())
		successOrg.UserMetadata = append(successOrg.UserMetadata, &admin_pb.ImportDataSuccessUserMetadata{UserId: userMetadata.GetId(), Key: userMetadata.GetKey()})
	}
	return nil
}

func importMachineKeys(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.MachineKeys == nil {
		return nil
	}
	for _, key := range org.GetMachineKeys() {
		logging.Debugf("import machine_user_key: %s", key.KeyId)
		_, err := s.command.AddUserMachineKey(ctx, &command.MachineKey{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   key.UserId,
				ResourceOwner: org.GetOrgId(),
			},
			KeyID:          key.KeyId,
			Type:           authn.KeyTypeToDomain(key.Type),
			ExpirationDate: key.ExpirationDate.AsTime(),
			PublicKey:      key.PublicKey,
		})
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "machine_user_key", Id: key.KeyId, Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.machineKeysCount += 1
		logging.Debugf("successful machine_user_key %d: %s", count.machineKeysCount, key.KeyId)
		successOrg.MachineKeys = append(successOrg.MachineKeys, key.KeyId)
	}
	return nil
}

func importUserLinks(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.UserLinks == nil {
		return nil
	}
	for _, userLinks := range org.GetUserLinks() {
		logging.Debugf("import userlink: %s", userLinks.GetUserId()+"_"+userLinks.GetIdpId()+"_"+userLinks.GetProvidedUserId()+"_"+userLinks.GetProvidedUserName())
		externalIDP := &command.AddLink{
			IDPID:         userLinks.IdpId,
			IDPExternalID: userLinks.ProvidedUserId,
			DisplayName:   userLinks.ProvidedUserName,
		}
		// TBD: why not command.BulkAddedUserIDPLinks?
		if _, err := s.command.AddUserIDPLink(ctx, userLinks.UserId, org.GetOrgId(), externalIDP); err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "user_link", Id: userLinks.UserId + "_" + userLinks.IdpId, Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.userLinksCount += 1
		logging.Debugf("successful userlink %d: %s", count.userLinksCount, userLinks.GetUserId()+"_"+userLinks.GetIdpId()+"_"+userLinks.GetProvidedUserId()+"_"+userLinks.GetProvidedUserName())
		successOrg.UserLinks = append(successOrg.UserLinks, &admin_pb.ImportDataSuccessUserLinks{UserId: userLinks.GetUserId(), IdpId: userLinks.GetIdpId(), ExternalUserId: userLinks.GetProvidedUserId(), DisplayName: userLinks.GetProvidedUserName()})
	}
	return nil

}

func importProjects(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.Projects == nil {
		return nil
	}
	for _, project := range org.GetProjects() {
		logging.Debugf("import project: %s", project.GetProjectId())
		_, err := s.command.AddProjectWithID(ctx, management.ProjectCreateToDomain(project.GetProject()), org.GetOrgId(), project.GetProjectId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "project", Id: project.GetProjectId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.projectCount += 1
		logging.Debugf("successful project %d: %s", count.projectCount, project.GetProjectId())
		successOrg.ProjectIds = append(successOrg.ProjectIds, project.GetProjectId())
	}
	return nil
}

func importOIDCApps(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.OidcApps == nil {
		return nil
	}
	for _, app := range org.GetOidcApps() {
		logging.Debugf("import oidcapplication: %s", app.GetAppId())
		oidcApp, err := management.AddOIDCAppRequestToDomain(app.App)
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "oidc_app", Id: app.GetAppId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		_, err = s.command.AddOIDCApplicationWithID(ctx, oidcApp, org.GetOrgId(), app.GetAppId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "oidc_app", Id: app.GetAppId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.oidcAppCount += 1
		logging.Debugf("successful oidcapplication %d: %s", count.oidcAppCount, app.GetAppId())
		successOrg.OidcAppIds = append(successOrg.OidcAppIds, app.GetAppId())
	}
	return nil
}

func importAPIApps(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.ApiApps == nil {
		return nil
	}
	for _, app := range org.GetApiApps() {
		logging.Debugf("import apiapplication: %s", app.GetAppId())
		_, err := s.command.AddAPIApplicationWithID(ctx, management.AddAPIAppRequestToDomain(app.GetApp()), org.GetOrgId(), app.GetAppId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "api_app", Id: app.GetAppId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.apiAppCount += 1
		logging.Debugf("successful apiapplication %d: %s", count.apiAppCount, app.GetAppId())
		successOrg.ApiAppIds = append(successOrg.ApiAppIds, app.GetAppId())
	}
	return nil
}

func importAppKeys(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.AppKeys == nil {
		return nil
	}
	for _, key := range org.GetAppKeys() {
		logging.Debugf("import app_key: %s", key.Id)
		_, err := s.command.AddApplicationKeyWithID(ctx, &domain.ApplicationKey{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   key.ProjectId,
				ResourceOwner: org.GetOrgId(),
			},
			ApplicationID:  key.AppId,
			ClientID:       key.ClientId,
			KeyID:          key.Id,
			Type:           authn.KeyTypeToDomain(key.Type),
			ExpirationDate: key.ExpirationDate.AsTime(),
			PublicKey:      key.PublicKey,
		}, org.GetOrgId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "app_key", Id: key.Id, Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.appKeysCount += 1
		logging.Debugf("successful app_key %d: %s", count.appKeysCount, key.Id)
		successOrg.AppKeys = append(successOrg.AppKeys, key.Id)
	}
	return nil
}

func importActions(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.Actions == nil {
		return nil
	}
	for _, action := range org.GetActions() {
		logging.Debugf("import action: %s", action.GetActionId())
		_, _, err := s.command.AddActionWithID(ctx, management.CreateActionRequestToDomain(action.GetAction()), org.GetOrgId(), action.GetActionId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "action", Id: action.GetActionId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.actionCount += 1
		logging.Debugf("successful action %d: %s", count.actionCount, action.GetActionId())
		successOrg.ActionIds = append(successOrg.ActionIds, action.ActionId)
	}
	return nil
}
func importProjectRoles(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.ProjectRoles == nil {
		return nil
	}
	for _, role := range org.GetProjectRoles() {
		logging.Debugf("import projectroles: %s", role.ProjectId+"_"+role.RoleKey)

		// TBD: why not command.BulkAddProjectRole?
		_, err := s.command.AddProjectRole(ctx, management.AddProjectRoleRequestToDomain(role), org.GetOrgId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "project_role", Id: role.ProjectId + "_" + role.RoleKey, Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.projectRolesCount += 1
		logging.Debugf("successful projectroles %d: %s", count.projectRolesCount, role.ProjectId+"_"+role.RoleKey)
		successOrg.ProjectRoles = append(successOrg.ProjectRoles, successOrg.ActionIds...)
		successOrg.ProjectRoles = append(successOrg.ProjectRoles, role.ProjectId+"_"+role.RoleKey)
	}
	return nil
}

func importResources(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg, count *counts, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessInitCode crypto.Generator) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := importOrgDomains(ctx, s, errors, successOrg, org); err != nil {
		return err
	}
	if err := importLabelPolicy(ctx, s, errors, org); err != nil {
		return err
	}
	importLockoutPolicy(ctx, s, errors, org)
	if err := importOidcIdps(ctx, s, errors, successOrg, org); err != nil {
		return err
	}
	if err := importJwtIdps(ctx, s, errors, successOrg, org); err != nil {
		return err
	}
	importLoginPolicy(ctx, s, errors, org)
	importPwComlexityPolicy(ctx, s, errors, org)
	importPrivacyPolicy(ctx, s, errors, org)
	importLoginTexts(ctx, s, errors, org)
	importInitMessageTexts(ctx, s, errors, org)
	importPWResetMessageTexts(ctx, s, errors, org)
	importVerifyEmailMessageTexts(ctx, s, errors, org)
	importVerifyPhoneMessageTexts(ctx, s, errors, org)
	importDomainClaimedMessageTexts(ctx, s, errors, org)
	importPasswordlessRegistrationMessageTexts(ctx, s, errors, org)
	importInviteUserMessageTexts(ctx, s, errors, org)
	if err := importHumanUsers(ctx, s, errors, successOrg, org, count, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessInitCode); err != nil {
		return err
	}
	if err := importMachineUsers(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importUserMetadata(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importMachineKeys(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importUserLinks(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importProjects(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importOIDCApps(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importAPIApps(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importAppKeys(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importActions(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	if err := importProjectRoles(ctx, s, errors, successOrg, org, count); err != nil {
		return err
	}
	return nil
}

func importOrgDomains(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.Domains == nil {
		return nil
	}
	for _, domainR := range org.Domains {
		orgDomain := &domain.OrgDomain{
			ObjectRoot: models.ObjectRoot{
				AggregateID: org.GetOrgId(),
			},
			Domain:   domainR.DomainName,
			Verified: domainR.IsVerified,
			Primary:  domainR.IsPrimary,
		}
		_, err := s.command.AddOrgDomain(ctx, org.GetOrgId(), domainR.DomainName, []string{})
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "domain", Id: org.GetOrgId() + "_" + domainR.DomainName, Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		logging.Debugf("successful domain: %s", domainR.DomainName)
		successOrg.Domains = append(successOrg.Domains, domainR.DomainName)

		if domainR.IsVerified {
			if _, err := s.command.VerifyOrgDomain(ctx, org.GetOrgId(), domainR.DomainName); err != nil {
				*errors = append(*errors, &admin_pb.ImportDataError{Type: "domain_isverified", Id: org.GetOrgId() + "_" + domainR.DomainName, Message: err.Error()})
			}
		}
		if domainR.IsPrimary {
			if _, err := s.command.SetPrimaryOrgDomain(ctx, orgDomain); err != nil {
				*errors = append(*errors, &admin_pb.ImportDataError{Type: "domain_isprimary", Id: org.GetOrgId() + "_" + domainR.DomainName, Message: err.Error()})
			}
		}
	}
	return nil
}

func importLoginTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.LoginTexts == nil {
		return
	}
	for _, text := range org.GetLoginTexts() {
		_, err := s.command.SetOrgLoginText(ctx, org.GetOrgId(), management.SetLoginCustomTextToDomain(text))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "login_texts", Id: org.GetOrgId() + "_" + text.Language, Message: err.Error()})
		}
	}
}

func importInitMessageTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.InitMessages == nil {
		return
	}
	for _, message := range org.GetInitMessages() {
		_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetInitCustomTextToDomain(message))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "init_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
		}
	}
}

func importPWResetMessageTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.PasswordResetMessages == nil {
		return
	}
	for _, message := range org.GetPasswordResetMessages() {
		_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetPasswordResetCustomTextToDomain(message))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "password_reset_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
		}
	}
}

func importVerifyEmailMessageTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.VerifyEmailMessages == nil {
		return
	}
	for _, message := range org.GetVerifyEmailMessages() {
		_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetVerifyEmailCustomTextToDomain(message))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "verify_email_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
		}
	}
}

func importVerifyPhoneMessageTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.VerifyPhoneMessages != nil {
		return
	}
	for _, message := range org.GetVerifyPhoneMessages() {
		_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetVerifyPhoneCustomTextToDomain(message))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "verify_phone_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
		}
	}
}

func importDomainClaimedMessageTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.DomainClaimedMessages == nil {
		return
	}
	for _, message := range org.GetDomainClaimedMessages() {
		_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetDomainClaimedCustomTextToDomain(message))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "domain_claimed_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
		}
	}
}

func importPasswordlessRegistrationMessageTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.PasswordlessRegistrationMessages == nil {
		return
	}
	for _, message := range org.GetPasswordlessRegistrationMessages() {
		_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetPasswordlessRegistrationCustomTextToDomain(message))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "passwordless_registration_message", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
		}
	}
}

func importInviteUserMessageTexts(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, org *admin_pb.DataOrg) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if org.PasswordlessRegistrationMessages == nil {
		return
	}
	for _, message := range org.GetInviteUserMessages() {
		_, err := s.command.SetOrgMessageText(ctx, authz.GetCtxData(ctx).OrgID, management.SetInviteUserCustomTextToDomain(message))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "invite_user_messages", Id: org.GetOrgId() + "_" + message.Language, Message: err.Error()})
		}
	}
}

func importOrg2(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, success *admin_pb.ImportDataSuccess, count *counts, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	successOrg := findOldOrg(success, org.OrgId)
	if successOrg == nil {
		return nil
	}
	if org.TriggerActions != nil {
		for _, triggerAction := range org.GetTriggerActions() {
			_, err := s.command.SetTriggerActions(ctx, action_grpc.FlowTypeToDomain(triggerAction.FlowType), action_grpc.TriggerTypeToDomain(triggerAction.TriggerType), triggerAction.ActionIds, org.GetOrgId())
			if err != nil {
				*errors = append(*errors, &admin_pb.ImportDataError{Type: "trigger_action", Id: triggerAction.FlowType + "_" + triggerAction.TriggerType, Message: err.Error()})
				continue
			}
			successOrg.TriggerActions = append(successOrg.TriggerActions, &management_pb.SetTriggerActionsRequest{FlowType: triggerAction.FlowType, TriggerType: triggerAction.TriggerType, ActionIds: triggerAction.GetActionIds()})
		}
	}
	if org.ProjectGrants != nil {
		for _, grant := range org.GetProjectGrants() {
			logging.Debugf("import projectgrant: %s", grant.GetGrantId()+"_"+grant.GetProjectGrant().GetProjectId()+"_"+grant.GetProjectGrant().GetGrantedOrgId())
			_, err := s.command.AddProjectGrantWithID(ctx, management.AddProjectGrantRequestToDomain(grant.GetProjectGrant()), grant.GetGrantId(), org.GetOrgId())
			if err != nil {
				*errors = append(*errors, &admin_pb.ImportDataError{Type: "project_grant", Id: org.GetOrgId() + "_" + grant.GetProjectGrant().GetProjectId() + "_" + grant.GetProjectGrant().GetGrantedOrgId(), Message: err.Error()})
				if isCtxTimeout(ctx) {
					return err
				}
				continue
			}
			count.projectGrantCount += 1
			logging.Debugf("successful projectgrant %d: %s", count.projectGrantCount, grant.GetGrantId()+"_"+grant.GetProjectGrant().GetProjectId()+"_"+grant.GetProjectGrant().GetGrantedOrgId())
			successOrg.ProjectGrants = append(successOrg.ProjectGrants, &admin_pb.ImportDataSuccessProjectGrant{GrantId: grant.GetGrantId(), ProjectId: grant.GetProjectGrant().GetProjectId(), OrgId: grant.GetProjectGrant().GetGrantedOrgId()})
		}
	}
	if org.UserGrants != nil {
		for _, grant := range org.GetUserGrants() {
			logging.Debugf("import usergrant: %s", grant.GetProjectId()+"_"+grant.GetUserId())
			_, err := s.command.AddUserGrant(ctx, management.AddUserGrantRequestToDomain(grant), org.GetOrgId())
			if err != nil {
				*errors = append(*errors, &admin_pb.ImportDataError{Type: "user_grant", Id: org.GetOrgId() + "_" + grant.GetProjectId() + "_" + grant.GetUserId(), Message: err.Error()})
				if isCtxTimeout(ctx) {
					return err
				}
				continue
			}
			count.userGrantCount += 1
			logging.Debugf("successful usergrant %d: %s", count.userGrantCount, grant.GetProjectId()+"_"+grant.GetUserId())
			successOrg.UserGrants = append(successOrg.UserGrants, &admin_pb.ImportDataSuccessUserGrant{ProjectId: grant.GetProjectId(), UserId: grant.GetUserId()})
		}
	}
	return nil
}

func importOrg3(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, success *admin_pb.ImportDataSuccess, count *counts, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	successOrg := findOldOrg(success, org.OrgId)
	if successOrg == nil {
		return nil
	}
	if err := importOrgMembers(ctx, s, errors, successOrg, count, org); err != nil {
		return err
	}
	if err := importProjectGrantMembers(ctx, s, errors, successOrg, count, org); err != nil {
		return err
	}
	return importProjectMembers(ctx, s, errors, successOrg, count, org)
}

func importOrgMembers(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, count *counts, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.OrgMembers == nil {
		return nil
	}
	for _, member := range org.GetOrgMembers() {
		logging.Debugf("import orgmember: %s", member.GetUserId())
		_, err := s.command.AddOrgMember(ctx, org.GetOrgId(), member.GetUserId(), member.GetRoles()...)
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "org_member", Id: org.GetOrgId() + "_" + member.GetUserId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.orgMemberCount += 1
		logging.Debugf("successful orgmember %d: %s", count.orgMemberCount, member.GetUserId())
		successOrg.OrgMembers = append(successOrg.OrgMembers, member.GetUserId())
	}
	return nil
}

func importProjectGrantMembers(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, count *counts, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.ProjectGrantMembers == nil {
		return nil
	}
	for _, member := range org.GetProjectGrantMembers() {
		logging.Debugf("import projectgrantmember: %s", member.GetProjectId()+"_"+member.GetGrantId()+"_"+member.GetUserId())
		_, err := s.command.AddProjectGrantMember(ctx, management.AddProjectGrantMemberRequestToDomain(member))
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "project_grant_member", Id: org.GetOrgId() + "_" + member.GetProjectId() + "_" + member.GetGrantId() + "_" + member.GetUserId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.projectGrantMemberCount += 1
		logging.Debugf("successful projectgrantmember %d: %s", count.projectGrantMemberCount, member.GetProjectId()+"_"+member.GetGrantId()+"_"+member.GetUserId())
		successOrg.ProjectGrantMembers = append(successOrg.ProjectGrantMembers, &admin_pb.ImportDataSuccessProjectGrantMember{ProjectId: member.GetProjectId(), GrantId: member.GetGrantId(), UserId: member.GetUserId()})
	}
	return nil
}

func importProjectMembers(ctx context.Context, s *Server, errors *[]*admin_pb.ImportDataError, successOrg *admin_pb.ImportDataSuccessOrg, count *counts, org *admin_pb.DataOrg) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if org.ProjectMembers == nil {
		return nil
	}
	for _, member := range org.GetProjectMembers() {
		logging.Debugf("import orgmember: %s", member.GetProjectId()+"_"+member.GetUserId())
		_, err := s.command.AddProjectMember(ctx, management.AddProjectMemberRequestToDomain(member), org.GetOrgId())
		if err != nil {
			*errors = append(*errors, &admin_pb.ImportDataError{Type: "project_member", Id: org.GetOrgId() + "_" + member.GetProjectId() + "_" + member.GetUserId(), Message: err.Error()})
			if isCtxTimeout(ctx) {
				return err
			}
			continue
		}
		count.projectMembersCount += 1
		logging.Debugf("successful orgmember %d: %s", count.projectMembersCount, member.GetProjectId()+"_"+member.GetUserId())
		successOrg.ProjectMembers = append(successOrg.ProjectMembers, &admin_pb.ImportDataSuccessProjectMember{ProjectId: member.GetProjectId(), UserId: member.GetUserId()})
	}
	return nil
}

func findOldOrg(success *admin_pb.ImportDataSuccess, orgId string) *admin_pb.ImportDataSuccessOrg {
	for _, oldOrd := range success.Orgs {
		if orgId == oldOrd.OrgId {
			return oldOrd
		}
	}
	return nil
}

func (s *Server) importData(ctx context.Context, orgs []*admin_pb.DataOrg) (_ *admin_pb.ImportDataResponse, _ *counts, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	errors := make([]*admin_pb.ImportDataError, 0)
	success := &admin_pb.ImportDataSuccess{}
	count := &counts{}

	initCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeInitCode, s.userCodeAlg)
	if err != nil {
		return nil, nil, err
	}
	emailCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, nil, err
	}
	phoneCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, s.userCodeAlg)
	if err != nil {
		return nil, nil, err
	}
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, nil, err
	}

	ctxData := authz.GetCtxData(ctx)
	for _, org := range orgs {
		count.humanUserLen += len(org.GetHumanUsers())
		count.machineUserLen += len(org.GetMachineUsers())
		count.userMetadataLen += len(org.GetUserMetadata())
		count.userLinksLen += len(org.GetUserLinks())
		count.projectLen += len(org.GetProjects())
		count.oidcAppLen += len(org.GetOidcApps())
		count.apiAppLen += len(org.GetApiApps())
		count.actionLen += len(org.GetActions())
		count.projectRolesLen += len(org.GetProjectRoles())
		count.projectGrantLen += len(org.GetProjectGrants())
		count.userGrantLen += len(org.GetUserGrants())
		count.projectMembersLen += len(org.GetProjectMembers())
		count.orgMemberLen += len(org.GetOrgMembers())
		count.projectGrantMemberLen += len(org.GetProjectGrantMembers())
		count.machineKeysCount += len(org.GetMachineKeys())
		count.appKeysCount += len(org.GetAppKeys())
	}
	for _, org := range orgs {
		if err = importOrg1(ctx, s, &errors, ctxData, org, success, count, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessInitCode); err != nil {
			return &admin_pb.ImportDataResponse{Errors: errors, Success: success}, count, err
		}
	}
	for _, org := range orgs {
		if err = importOrg2(ctx, s, &errors, success, count, org); err != nil {
			return &admin_pb.ImportDataResponse{Errors: errors, Success: success}, count, err
		}
	}
	for _, org := range orgs {
		if err = importOrg3(ctx, s, &errors, success, count, org); err != nil {
			return &admin_pb.ImportDataResponse{Errors: errors, Success: success}, count, err
		}
	}
	return &admin_pb.ImportDataResponse{
		Errors:  errors,
		Success: success,
	}, count, nil
}

func (s *Server) dataOrgsV1ToDataOrgs(ctx context.Context, dataOrgs *v1_pb.ImportDataOrg) (_ *admin_pb.ImportDataOrg, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	orgs := make([]*admin_pb.DataOrg, 0)
	for _, orgV1 := range dataOrgs.Orgs {
		triggerActions := make([]*management_pb.SetTriggerActionsRequest, 0)
		for _, action := range orgV1.GetTriggerActions() {
			triggerActions = append(triggerActions, &management_pb.SetTriggerActionsRequest{
				FlowType:    strconv.Itoa(int(action.GetFlowType().Number())),
				TriggerType: strconv.Itoa(int(action.GetTriggerType().Number())),
				ActionIds:   action.ActionIds,
			})
		}

		org := &admin_pb.DataOrg{
			OrgId:                            orgV1.GetOrgId(),
			Org:                              orgV1.GetOrg(),
			DomainPolicy:                     nil,
			LabelPolicy:                      orgV1.GetLabelPolicy(),
			LockoutPolicy:                    orgV1.GetLockoutPolicy(),
			LoginPolicy:                      orgV1.GetLoginPolicy(),
			PasswordComplexityPolicy:         orgV1.GetPasswordComplexityPolicy(),
			PrivacyPolicy:                    orgV1.GetPrivacyPolicy(),
			Projects:                         orgV1.GetProjects(),
			ProjectRoles:                     orgV1.GetProjectRoles(),
			ApiApps:                          orgV1.GetApiApps(),
			OidcApps:                         orgV1.GetOidcApps(),
			HumanUsers:                       orgV1.GetHumanUsers(),
			MachineUsers:                     orgV1.GetMachineUsers(),
			TriggerActions:                   triggerActions,
			Actions:                          orgV1.GetActions(),
			ProjectGrants:                    orgV1.GetProjectGrants(),
			UserGrants:                       orgV1.GetUserGrants(),
			OrgMembers:                       orgV1.GetOrgMembers(),
			ProjectMembers:                   orgV1.GetProjectMembers(),
			ProjectGrantMembers:              orgV1.GetProjectGrantMembers(),
			UserMetadata:                     orgV1.GetUserMetadata(),
			LoginTexts:                       orgV1.GetLoginTexts(),
			InitMessages:                     orgV1.GetInitMessages(),
			PasswordResetMessages:            orgV1.GetPasswordResetMessages(),
			VerifyEmailMessages:              orgV1.GetVerifyEmailMessages(),
			VerifyPhoneMessages:              orgV1.GetVerifyPhoneMessages(),
			DomainClaimedMessages:            orgV1.GetDomainClaimedMessages(),
			PasswordlessRegistrationMessages: orgV1.GetPasswordlessRegistrationMessages(),
			InviteUserMessages:               orgV1.GetInviteUserMessages(),
			OidcIdps:                         orgV1.GetOidcIdps(),
			JwtIdps:                          orgV1.GetJwtIdps(),
			UserLinks:                        orgV1.GetUserLinks(),
			Domains:                          orgV1.GetDomains(),
			AppKeys:                          orgV1.GetAppKeys(),
			MachineKeys:                      orgV1.GetMachineKeys(),
		}
		if orgV1.IamPolicy != nil {
			defaultDomainPolicy, err := s.query.DefaultDomainPolicy(ctx)
			if err != nil {
				return nil, err
			}

			org.DomainPolicy = &admin_pb.AddCustomDomainPolicyRequest{
				UserLoginMustBeDomain:                  orgV1.IamPolicy.UserLoginMustBeDomain,
				ValidateOrgDomains:                     defaultDomainPolicy.ValidateOrgDomains,
				SmtpSenderAddressMatchesInstanceDomain: defaultDomainPolicy.SMTPSenderAddressMatchesInstanceDomain,
			}
		}
		if org.LoginPolicy != nil {
			defaultLoginPolicy, err := s.query.DefaultLoginPolicy(ctx)
			if err != nil {
				return nil, err
			}
			org.LoginPolicy.ExternalLoginCheckLifetime = durationpb.New(time.Duration(defaultLoginPolicy.ExternalLoginCheckLifetime))
			org.LoginPolicy.MultiFactorCheckLifetime = durationpb.New(time.Duration(defaultLoginPolicy.MultiFactorCheckLifetime))
			org.LoginPolicy.SecondFactorCheckLifetime = durationpb.New(time.Duration(defaultLoginPolicy.SecondFactorCheckLifetime))
			org.LoginPolicy.PasswordCheckLifetime = durationpb.New(time.Duration(defaultLoginPolicy.PasswordCheckLifetime))
			org.LoginPolicy.MfaInitSkipLifetime = durationpb.New(time.Duration(defaultLoginPolicy.MFAInitSkipLifetime))

			if orgV1.SecondFactors != nil {
				org.LoginPolicy.SecondFactors = make([]policy.SecondFactorType, len(orgV1.SecondFactors))
				for i, factor := range orgV1.SecondFactors {
					org.LoginPolicy.SecondFactors[i] = factor.GetType()
				}
			}
			if orgV1.MultiFactors != nil {
				org.LoginPolicy.MultiFactors = make([]policy.MultiFactorType, len(orgV1.MultiFactors))
				for i, factor := range orgV1.MultiFactors {
					org.LoginPolicy.MultiFactors[i] = factor.GetType()
				}
			}
			if orgV1.Idps != nil {
				org.LoginPolicy.Idps = make([]*management_pb.AddCustomLoginPolicyRequest_IDP, len(orgV1.Idps))
				for i, idpR := range orgV1.Idps {
					org.LoginPolicy.Idps[i] = &management_pb.AddCustomLoginPolicyRequest_IDP{
						IdpId:     idpR.GetIdpId(),
						OwnerType: idpR.GetOwnerType(),
					}
				}
			}
		}
		orgs = append(orgs, org)
	}

	return &admin_pb.ImportDataOrg{
		Orgs: orgs,
	}, nil
}

func isCtxTimeout(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
