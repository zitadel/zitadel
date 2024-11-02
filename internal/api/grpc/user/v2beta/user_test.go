package user

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func Test_idpIntentToIDPIntentPb(t *testing.T) {
	decryption := func(err error) crypto.EncryptionAlgorithm {
		mCrypto := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
		mCrypto.EXPECT().Algorithm().Return("enc")
		mCrypto.EXPECT().DecryptionKeyIDs().Return([]string{"id"})
		mCrypto.EXPECT().DecryptString(gomock.Any(), gomock.Any()).DoAndReturn(
			func(code []byte, keyID string) (string, error) {
				if err != nil {
					return "", err
				}
				return string(code), nil
			})
		return mCrypto
	}

	type args struct {
		intent *command.IDPIntentWriteModel
		alg    crypto.EncryptionAlgorithm
	}
	type res struct {
		resp *user.RetrieveIdentityProviderIntentResponse
		err  error
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"decryption invalid key id error",
			args{
				intent: &command.IDPIntentWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:       "intentID",
						ProcessedSequence: 123,
						ResourceOwner:     "ro",
						InstanceID:        "instanceID",
						ChangeDate:        time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local),
					},
					IDPID:       "idpID",
					IDPUser:     []byte(`{"userID": "idpUserID", "username": "username"}`),
					IDPUserID:   "idpUserID",
					IDPUserName: "username",
					IDPAccessToken: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("accessToken"),
					},
					IDPIDToken:         "idToken",
					IDPEntryAttributes: map[string][]string{},
					UserID:             "userID",
					State:              domain.IDPIntentStateSucceeded,
				},
				alg: decryption(zerrors.ThrowInternal(nil, "id", "invalid key id")),
			},
			res{
				resp: nil,
				err:  zerrors.ThrowInternal(nil, "id", "invalid key id"),
			},
		}, {
			"successful oauth",
			args{
				intent: &command.IDPIntentWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:       "intentID",
						ProcessedSequence: 123,
						ResourceOwner:     "ro",
						InstanceID:        "instanceID",
						ChangeDate:        time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local),
					},
					IDPID:       "idpID",
					IDPUser:     []byte(`{"userID": "idpUserID", "username": "username"}`),
					IDPUserID:   "idpUserID",
					IDPUserName: "username",
					IDPAccessToken: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("accessToken"),
					},
					IDPIDToken: "idToken",
					UserID:     "",
					State:      domain.IDPIntentStateSucceeded,
				},
				alg: decryption(nil),
			},
			res{
				resp: &user.RetrieveIdentityProviderIntentResponse{
					Details: &object_pb.Details{
						Sequence:      123,
						ChangeDate:    timestamppb.New(time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local)),
						ResourceOwner: "ro",
					},
					IdpInformation: &user.IDPInformation{
						Access: &user.IDPInformation_Oauth{
							Oauth: &user.IDPOAuthAccessInformation{
								AccessToken: "accessToken",
								IdToken:     gu.Ptr("idToken"),
							},
						},
						IdpId:    "idpID",
						UserId:   "idpUserID",
						UserName: "username",
						RawInformation: func() *structpb.Struct {
							s, err := structpb.NewStruct(map[string]interface{}{
								"userID":   "idpUserID",
								"username": "username",
							})
							require.NoError(t, err)
							return s
						}(),
					},
				},
				err: nil,
			},
		},
		{
			"successful oauth with linked user",
			args{
				intent: &command.IDPIntentWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:       "intentID",
						ProcessedSequence: 123,
						ResourceOwner:     "ro",
						InstanceID:        "instanceID",
						ChangeDate:        time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local),
					},
					IDPID:       "idpID",
					IDPUser:     []byte(`{"userID": "idpUserID", "username": "username"}`),
					IDPUserID:   "idpUserID",
					IDPUserName: "username",
					IDPAccessToken: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("accessToken"),
					},
					IDPIDToken: "idToken",
					UserID:     "userID",
					State:      domain.IDPIntentStateSucceeded,
				},
				alg: decryption(nil),
			},
			res{
				resp: &user.RetrieveIdentityProviderIntentResponse{
					Details: &object_pb.Details{
						Sequence:      123,
						ChangeDate:    timestamppb.New(time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local)),
						ResourceOwner: "ro",
					},
					IdpInformation: &user.IDPInformation{
						Access: &user.IDPInformation_Oauth{
							Oauth: &user.IDPOAuthAccessInformation{
								AccessToken: "accessToken",
								IdToken:     gu.Ptr("idToken"),
							},
						},
						IdpId:    "idpID",
						UserId:   "idpUserID",
						UserName: "username",
						RawInformation: func() *structpb.Struct {
							s, err := structpb.NewStruct(map[string]interface{}{
								"userID":   "idpUserID",
								"username": "username",
							})
							require.NoError(t, err)
							return s
						}(),
					},
					UserId: "userID",
				},
				err: nil,
			},
		}, {
			"successful ldap",
			args{
				intent: &command.IDPIntentWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:       "intentID",
						ProcessedSequence: 123,
						ResourceOwner:     "ro",
						InstanceID:        "instanceID",
						ChangeDate:        time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local),
					},
					IDPID:       "idpID",
					IDPUser:     []byte(`{"userID": "idpUserID", "username": "username"}`),
					IDPUserID:   "idpUserID",
					IDPUserName: "username",
					IDPEntryAttributes: map[string][]string{
						"id":        {"idpUserID"},
						"firstName": {"firstname1", "firstname2"},
						"lastName":  {"lastname"},
					},
					UserID: "",
					State:  domain.IDPIntentStateSucceeded,
				},
			},
			res{
				resp: &user.RetrieveIdentityProviderIntentResponse{
					Details: &object_pb.Details{
						Sequence:      123,
						ChangeDate:    timestamppb.New(time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local)),
						ResourceOwner: "ro",
					},
					IdpInformation: &user.IDPInformation{
						Access: &user.IDPInformation_Ldap{
							Ldap: &user.IDPLDAPAccessInformation{
								Attributes: func() *structpb.Struct {
									s, err := structpb.NewStruct(map[string]interface{}{
										"id":        []interface{}{"idpUserID"},
										"firstName": []interface{}{"firstname1", "firstname2"},
										"lastName":  []interface{}{"lastname"},
									})
									require.NoError(t, err)
									return s
								}(),
							},
						},
						IdpId:    "idpID",
						UserId:   "idpUserID",
						UserName: "username",
						RawInformation: func() *structpb.Struct {
							s, err := structpb.NewStruct(map[string]interface{}{
								"userID":   "idpUserID",
								"username": "username",
							})
							require.NoError(t, err)
							return s
						}(),
					},
				},
				err: nil,
			},
		}, {
			"successful ldap with linked user",
			args{
				intent: &command.IDPIntentWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:       "intentID",
						ProcessedSequence: 123,
						ResourceOwner:     "ro",
						InstanceID:        "instanceID",
						ChangeDate:        time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local),
					},
					IDPID:       "idpID",
					IDPUser:     []byte(`{"userID": "idpUserID", "username": "username"}`),
					IDPUserID:   "idpUserID",
					IDPUserName: "username",
					IDPEntryAttributes: map[string][]string{
						"id":        {"idpUserID"},
						"firstName": {"firstname1", "firstname2"},
						"lastName":  {"lastname"},
					},
					UserID: "userID",
					State:  domain.IDPIntentStateSucceeded,
				},
			},
			res{
				resp: &user.RetrieveIdentityProviderIntentResponse{
					Details: &object_pb.Details{
						Sequence:      123,
						ChangeDate:    timestamppb.New(time.Date(2019, 4, 1, 1, 1, 1, 1, time.Local)),
						ResourceOwner: "ro",
					},
					IdpInformation: &user.IDPInformation{
						Access: &user.IDPInformation_Ldap{
							Ldap: &user.IDPLDAPAccessInformation{
								Attributes: func() *structpb.Struct {
									s, err := structpb.NewStruct(map[string]interface{}{
										"id":        []interface{}{"idpUserID"},
										"firstName": []interface{}{"firstname1", "firstname2"},
										"lastName":  []interface{}{"lastname"},
									})
									require.NoError(t, err)
									return s
								}(),
							},
						},
						IdpId:    "idpID",
						UserId:   "idpUserID",
						UserName: "username",
						RawInformation: func() *structpb.Struct {
							s, err := structpb.NewStruct(map[string]interface{}{
								"userID":   "idpUserID",
								"username": "username",
							})
							require.NoError(t, err)
							return s
						}(),
					},
					UserId: "userID",
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := idpIntentToIDPIntentPb(tt.args.intent, tt.args.alg)
			require.ErrorIs(t, err, tt.res.err)
			grpc.AllFieldsEqual(t, tt.res.resp.ProtoReflect(), got.ProtoReflect(), grpc.CustomMappers)
		})
	}
}

func Test_authMethodTypesToPb(t *testing.T) {
	tests := []struct {
		name        string
		methodTypes []domain.UserAuthMethodType
		want        []user.AuthenticationMethodType
	}{
		{
			"empty list",
			nil,
			[]user.AuthenticationMethodType{},
		},
		{
			"list",
			[]domain.UserAuthMethodType{
				domain.UserAuthMethodTypePasswordless,
			},
			[]user.AuthenticationMethodType{
				user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, authMethodTypesToPb(tt.methodTypes), "authMethodTypesToPb(%v)", tt.methodTypes)
		})
	}
}

func Test_authMethodTypeToPb(t *testing.T) {
	tests := []struct {
		name       string
		methodType domain.UserAuthMethodType
		want       user.AuthenticationMethodType
	}{
		{
			"uspecified",
			domain.UserAuthMethodTypeUnspecified,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED,
		},
		{
			"totp",
			domain.UserAuthMethodTypeTOTP,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_TOTP,
		},
		{
			"u2f",
			domain.UserAuthMethodTypeU2F,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_U2F,
		},
		{
			"passkey",
			domain.UserAuthMethodTypePasswordless,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
		},
		{
			"password",
			domain.UserAuthMethodTypePassword,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSWORD,
		},
		{
			"idp",
			domain.UserAuthMethodTypeIDP,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_IDP,
		},
		{
			"otp sms",
			domain.UserAuthMethodTypeOTPSMS,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_OTP_SMS,
		},
		{
			"otp email",
			domain.UserAuthMethodTypeOTPEmail,
			user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_OTP_EMAIL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, authMethodTypeToPb(tt.methodType), "authMethodTypeToPb(%v)", tt.methodType)
		})
	}
}
