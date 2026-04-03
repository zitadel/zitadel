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

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
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
			name: "decryption invalid key id error",
			args: args{
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
			res: res{
				resp: nil,
				err:  zerrors.ThrowInternal(nil, "id", "invalid key id"),
			},
		}, {
			name: "successful oauth",
			args: args{
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
			res: res{
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
			name: "successful oauth with linked user",
			args: args{
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
			res: res{
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
			name: "successful ldap",
			args: args{
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
			res: res{
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
			name: "successful ldap with linked user",
			args: args{
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
			res: res{
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
			assert.EqualExportedValues(t, tt.res.resp, got)
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
			name:        "empty list",
			methodTypes: nil,
			want:        []user.AuthenticationMethodType{},
		},
		{
			name: "list",
			methodTypes: []domain.UserAuthMethodType{
				domain.UserAuthMethodTypePasswordless,
			},
			want: []user.AuthenticationMethodType{
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
			name:       "uspecified",
			methodType: domain.UserAuthMethodTypeUnspecified,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED,
		},
		{
			name:       "totp",
			methodType: domain.UserAuthMethodTypeTOTP,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_TOTP,
		},
		{
			name:       "u2f",
			methodType: domain.UserAuthMethodTypeU2F,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_U2F,
		},
		{
			name:       "passkey",
			methodType: domain.UserAuthMethodTypePasswordless,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
		},
		{
			name:       "password",
			methodType: domain.UserAuthMethodTypePassword,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSWORD,
		},
		{
			name:       "idp",
			methodType: domain.UserAuthMethodTypeIDP,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_IDP,
		},
		{
			name:       "otp sms",
			methodType: domain.UserAuthMethodTypeOTPSMS,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_OTP_SMS,
		},
		{
			name:       "otp email",
			methodType: domain.UserAuthMethodTypeOTPEmail,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_OTP_EMAIL,
		},
		{
			name:       "recovery code",
			methodType: domain.UserAuthMethodTypeRecoveryCode,
			want:       user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_RECOVERY_CODE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, authMethodTypeToPb(tt.methodType), "authMethodTypeToPb(%v)", tt.methodType)
		})
	}
}
