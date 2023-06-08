package main

import (
	"strings"

	"github.com/checkr/openmock"
	"github.com/golang/protobuf/proto"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var serviceMap = map[string]openmock.GRPCService{
	settings.SettingsService_ServiceDesc.ServiceName: {
		splitServiceMethod(settings.SettingsService_GetBrandingSettings_FullMethodName): pair(&settings.GetBrandingSettingsRequest{}, &settings.GetBrandingSettingsResponse{}),
	},
	user.UserService_ServiceDesc.ServiceName: {
		splitServiceMethod(user.UserService_VerifyEmail_FullMethodName): pair(&user.VerifyEmailRequest{}, &user.VerifyEmailResponse{}),
	},
}

func splitServiceMethod(fullMethodName string) string {
	return strings.Split(fullMethodName, "/")[2]
}

func pair(request, response proto.Message) openmock.GRPCRequestResponsePair {
	return openmock.GRPCRequestResponsePair{Request: proto.MessageV2(request), Response: proto.MessageV2(response)}
}
