package auth

import (
	"context"
	"github.com/caos/zitadel/internal/user/model"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyUser(ctx context.Context, _ *empty.Empty) (*auth.UserView, error) {
	user, err := s.repo.MyUser(ctx)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) GetMyUserProfile(ctx context.Context, _ *empty.Empty) (*auth.UserProfileView, error) {
	profile, err := s.repo.MyProfile(ctx)
	if err != nil {
		return nil, err
	}
	return profileViewFromModel(profile), nil
}

func (s *Server) GetMyUserEmail(ctx context.Context, _ *empty.Empty) (*auth.UserEmailView, error) {
	email, err := s.repo.MyEmail(ctx)
	if err != nil {
		return nil, err
	}
	return emailViewFromModel(email), nil
}

func (s *Server) GetMyUserPhone(ctx context.Context, _ *empty.Empty) (*auth.UserPhoneView, error) {
	phone, err := s.repo.MyPhone(ctx)
	if err != nil {
		return nil, err
	}
	return phoneViewFromModel(phone), nil
}

func (s *Server) RemoveMyUserPhone(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.RemoveMyPhone(ctx)
	return &empty.Empty{}, err
}

func (s *Server) GetMyUserAddress(ctx context.Context, _ *empty.Empty) (*auth.UserAddressView, error) {
	address, err := s.repo.MyAddress(ctx)
	if err != nil {
		return nil, err
	}
	return addressViewFromModel(address), nil
}

func (s *Server) GetMyMfas(ctx context.Context, _ *empty.Empty) (*auth.MultiFactors, error) {
	mfas, err := s.repo.MyUserMfas(ctx)
	if err != nil {
		return nil, err
	}
	return &auth.MultiFactors{Mfas: mfasFromModel(mfas)}, nil
}

func (s *Server) UpdateMyUserProfile(ctx context.Context, request *auth.UpdateUserProfileRequest) (*auth.UserProfile, error) {
	profile, err := s.repo.ChangeMyProfile(ctx, updateProfileToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return profileFromModel(profile), nil
}

func (s *Server) ChangeMyUserName(ctx context.Context, request *auth.ChangeUserNameRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.repo.ChangeMyUsername(ctx, request.UserName)
}

func (s *Server) ChangeMyUserEmail(ctx context.Context, request *auth.UpdateUserEmailRequest) (*auth.UserEmail, error) {
	email, err := s.repo.ChangeMyEmail(ctx, updateEmailToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return emailFromModel(email), nil
}

func (s *Server) VerifyMyUserEmail(ctx context.Context, request *auth.VerifyMyUserEmailRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyEmail(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyEmailVerificationMail(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyEmailVerificationMail(ctx)
	return &empty.Empty{}, err
}

func (s *Server) ChangeMyUserPhone(ctx context.Context, request *auth.UpdateUserPhoneRequest) (*auth.UserPhone, error) {
	phone, err := s.repo.ChangeMyPhone(ctx, updatePhoneToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return phoneFromModel(phone), nil
}

func (s *Server) VerifyMyUserPhone(ctx context.Context, request *auth.VerifyUserPhoneRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyPhone(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyPhoneVerificationCode(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyPhoneVerificationCode(ctx)
	return &empty.Empty{}, err
}

func (s *Server) UpdateMyUserAddress(ctx context.Context, request *auth.UpdateUserAddressRequest) (*auth.UserAddress, error) {
	address, err := s.repo.ChangeMyAddress(ctx, updateAddressToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) ChangeMyPassword(ctx context.Context, request *auth.PasswordChange) (*empty.Empty, error) {
	err := s.repo.ChangeMyPassword(ctx, request.OldPassword, request.NewPassword)
	return &empty.Empty{}, err
}

func (s *Server) SearchMyExternalIDPs(ctx context.Context, request *auth.ExternalIDPSearchRequest) (*auth.ExternalIDPSearchResponse, error) {
	externalIDP, err := s.repo.SearchMyExternalIDPs(ctx, externalIDPSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return externalIDPSearchResponseFromModel(externalIDP), nil
}

func (s *Server) RemoveMyExternalIDP(ctx context.Context, request *auth.ExternalIDPRemoveRequest) (*empty.Empty, error) {
	err := s.repo.RemoveMyExternalIDP(ctx, externalIDPRemoveToModel(ctx, request))
	return &empty.Empty{}, err
}

func (s *Server) GetMyPasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*auth.PasswordComplexityPolicy, error) {
	policy, err := s.repo.GetMyPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyFromModel(policy), nil
}

func (s *Server) AddMfaOTP(ctx context.Context, _ *empty.Empty) (_ *auth.MfaOtpResponse, err error) {
	otp, err := s.repo.AddMyMfaOTP(ctx)
	if err != nil {
		return nil, err
	}
	return otpFromModel(otp), nil
}

func (s *Server) VerifyMfaOTP(ctx context.Context, request *auth.VerifyMfaOtp) (*empty.Empty, error) {
	err := s.repo.VerifyMyMfaOTPSetup(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) RemoveMfaOTP(ctx context.Context, _ *empty.Empty) (_ *empty.Empty, err error) {
	s.repo.RemoveMyMfaOTP(ctx)
	return &empty.Empty{}, err
}

func (s *Server) AddMfaU2F(ctx context.Context, _ *empty.Empty) (_ *auth.MfaU2FResponse, err error) {
	u2f, err := s.repo.AddMyMfaU2F(ctx)
	return verifyMfaU2FFromModel(u2f), err
}

func (s *Server) VerifyMfaU2F(ctx context.Context, request *auth.VerifyMfaU2F) (*empty.Empty, error) {
	err := s.repo.VerifyMyMfaU2FSetup(ctx, request.PublicKeyCredential)
	return &empty.Empty{}, err
}

func verifyMfaU2FFromModel(u2f *model.WebAuthNToken) *auth.MfaU2FResponse {
	return &auth.MfaU2FResponse{
		Id:        u2f.WebAuthNTokenID,
		PublicKey: u2f.PublicKey,
		State:     mfaStateFromModel(u2f.State),
	}
}

//
//func publicKeyFromModel(response protocol.PublicKeyCredentialCreationOptions) *auth.U2FPublicKey {
//	return &auth.U2FPublicKey{
//		Challenge:              response.Challenge,
//		Rp:                     rpFromModel(response.RelyingParty),
//		User:                   userFromModel(response.User),
//		PubKeyCredParams:       publicKeyCredParamsFromModel(response.Parameters),
//		AuthenticatorSelection: authenticatorSelectionFromModel(response.AuthenticatorSelection),
//		Timeout:                int32(response.Timeout),
//		ExcludeCredentials:     excludeCredentialsFromModel(response.CredentialExcludeList),
//		Extensions:             extensionsFromModel(response.Extensions),
//		Attestation:            attestionFromModel(response.Attestation),
//	}
//}
//
//func attestionFromModel(attestation protocol.ConveyancePreference) auth.ConveyancePreference {
//	switch attestation {
//	case protocol.PreferNoAttestation:
//		return auth.ConveyancePreference_ConveyancePreferenceNoAttestation
//	case protocol.PreferDirectAttestation:
//		return auth.ConveyancePreference_ConveyancePreferenceDirectAttestation
//	case protocol.PreferIndirectAttestation:
//		return auth.ConveyancePreference_ConveyancePreferenceIndirectAttestation
//	default:
//		return auth.ConveyancePreference_ConveyancePreferenceNoAttestation
//	}
//}
//
//func extensionsFromModel(extensions protocol.AuthenticationExtensions) map[string]string {
//	if extensions == nil {
//		return nil
//	}
//	exts := make(map[string]string)
//	for key, value := range extensions {
//		exts[key] = value.(string)
//	}
//	return exts
//}
//
//func excludeCredentialsFromModel(list []protocol.CredentialDescriptor) []*auth.CredentialDescriptor {
//	if list == nil {
//		return nil
//	}
//	creds := make([]*auth.CredentialDescriptor, len(list))
//	for i, desc := range list {
//		creds[i] = &auth.CredentialDescriptor{
//			Type:       credentialTypeFromModel(desc.Type),
//			Id:         desc.CredentialID,
//			Transports: transportsFromModel(desc.Transport),
//		}
//	}
//	return creds
//}
//
//func transportsFromModel(transports []protocol.AuthenticatorTransport) []auth.AuthenticatorTransport {
//	if transports == nil {
//		return nil
//	}
//	trans := make([]auth.AuthenticatorTransport, len(transports))
//	for i, t := range transports {
//		trans[i] = transportFromModel(t)
//	}
//	return trans
//}
//
//func transportFromModel(trans protocol.AuthenticatorTransport) auth.AuthenticatorTransport {
//	switch trans {
//	case protocol.USB:
//		return auth.AuthenticatorTransport_AuthenticatorTransportUSB
//	case protocol.NFC:
//		return auth.AuthenticatorTransport_AuthenticatorTransportNFC
//	case protocol.BLE:
//		return auth.AuthenticatorTransport_AuthenticatorTransportBLE
//	case protocol.Internal:
//		return auth.AuthenticatorTransport_AuthenticatorTransportInternal
//	default:
//		return auth.AuthenticatorTransport_AuthenticatorTransportUnspecified
//	}
//}
//
//func authenticatorSelectionFromModel(selection protocol.AuthenticatorSelection) *auth.AuthenticatorSelection {
//	return &auth.AuthenticatorSelection{
//		AuthenticatorAttachment: authenticatorAttachementFromModel(selection.AuthenticatorAttachment),
//		RequireResidentKey:      *selection.RequireResidentKey,
//		UserVerification:        userVerificationFromModel(selection.UserVerification),
//	}
//}
//
//func userVerificationFromModel(verification protocol.UserVerificationRequirement) auth.UserVerificationRequirement {
//	switch verification {
//	case protocol.VerificationDiscouraged:
//		return auth.UserVerificationRequirement_UserVerificationRequirementDiscouraged
//	case protocol.VerificationPreferred:
//		return auth.UserVerificationRequirement_UserVerificationRequirementPreferred
//	case protocol.VerificationRequired:
//		return auth.UserVerificationRequirement_UserVerificationRequirementRequired
//	default:
//		return auth.UserVerificationRequirement_UserVerificationRequirementPreferred
//	}
//}
//
//func authenticatorAttachementFromModel(attachment protocol.AuthenticatorAttachment) auth.AuthenticatorAttachment {
//	switch attachment {
//	case protocol.Platform:
//		return auth.AuthenticatorAttachment_AuthenticatorAttachmentPlatform
//	case protocol.CrossPlatform:
//		return auth.AuthenticatorAttachment_AuthenticatorAttachmentCrossPlatform
//	default:
//		return auth.AuthenticatorAttachment_AuthenticatorAttachmentUnspecified
//	}
//}
//
//func publicKeyCredParamsFromModel(parameters []protocol.CredentialParameter) []*auth.CredentialParameter {
//	if parameters == nil {
//		return nil
//	}
//	creds := make([]*auth.CredentialParameter, len(parameters))
//	for i, param := range parameters {
//		creds[i] = &auth.CredentialParameter{
//			Type:      credentialTypeFromModel(param.Type),
//			Algorithm: int32(param.Algorithm),
//		}
//	}
//	return creds
//}
//
//func credentialTypeFromModel(credentialType protocol.CredentialType) auth.CredentialType {
//	switch credentialType {
//	case protocol.PublicKeyCredentialType:
//		return auth.CredentialType_CredentialTypePublicKey
//	default:
//		return auth.CredentialType_CredentialTypePublicKey
//	}
//}
//
//func userFromModel(user protocol.UserEntity) *auth.UserEntity {
//	return &auth.UserEntity{
//		Name:        user.Name,
//		Icon:        user.Icon,
//		DisplayName: user.DisplayName,
//		Id:          user.ID,
//	}
//}
//
//func rpFromModel(party protocol.RelyingPartyEntity) *auth.RelyingParty {
//	return &auth.RelyingParty{
//		Name: party.Name,
//		Icon: party.Icon,
//		Id:   party.ID,
//	}
//}

//
//func attestionToModel(response *auth.Response) protocol.AuthenticatorAttestationResponse {
//	return protocol.AuthenticatorAttestationResponse{
//		AuthenticatorResponse: protocol.AuthenticatorResponse{
//			ClientDataJSON: response.ClientData_JSON,
//		},
//		AttestationObject: response.AttestionObject,
//	}
//}
//
//func credentialTypeToModel(credentialType auth.CredentialType) protocol.CredentialType {
//	switch credentialType {
//	case auth.CredentialType_CredentialTypePublicKey:
//		return protocol.PublicKeyCredentialType
//	default:
//		return protocol.PublicKeyCredentialType
//	}
//}

func (s *Server) GetMyUserChanges(ctx context.Context, request *auth.ChangesRequest) (*auth.Changes, error) {
	changes, err := s.repo.MyUserChanges(ctx, request.SequenceOffset, request.Limit, request.Asc)
	if err != nil {
		return nil, err
	}
	return userChangesToResponse(changes, request.GetSequenceOffset(), request.GetLimit()), nil
}
