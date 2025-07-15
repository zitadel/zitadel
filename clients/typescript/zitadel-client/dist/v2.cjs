"use strict";Object.defineProperty(exports, "__esModule", {value: true});

var _chunkDUECDNWCcjs = require('./chunk-DUECDNWC.cjs');

// src/v2.ts
var _protobuf = require('@bufbuild/protobuf');
var _feature_service_pbjs = require('@zitadel/proto/zitadel/feature/v2/feature_service_pb.js');
var _idp_service_pbjs = require('@zitadel/proto/zitadel/idp/v2/idp_service_pb.js');
var _object_pbjs = require('@zitadel/proto/zitadel/object/v2/object_pb.js');
var _oidc_service_pbjs = require('@zitadel/proto/zitadel/oidc/v2/oidc_service_pb.js');
var _org_service_pbjs = require('@zitadel/proto/zitadel/org/v2/org_service_pb.js');
var _saml_service_pbjs = require('@zitadel/proto/zitadel/saml/v2/saml_service_pb.js');
var _session_service_pbjs = require('@zitadel/proto/zitadel/session/v2/session_service_pb.js');
var _settings_service_pbjs = require('@zitadel/proto/zitadel/settings/v2/settings_service_pb.js');
var _user_service_pbjs = require('@zitadel/proto/zitadel/user/v2/user_service_pb.js');
var createUserServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _user_service_pbjs.UserService);
var createSettingsServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _settings_service_pbjs.SettingsService);
var createSessionServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _session_service_pbjs.SessionService);
var createOIDCServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _oidc_service_pbjs.OIDCService);
var createSAMLServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _saml_service_pbjs.SAMLService);
var createOrganizationServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _org_service_pbjs.OrganizationService);
var createFeatureServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _feature_service_pbjs.FeatureService);
var createIdpServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _idp_service_pbjs.IdentityProviderService);
function makeReqCtx(orgId) {
  return _protobuf.create.call(void 0, _object_pbjs.RequestContextSchema, {
    resourceOwner: orgId ? { case: "orgId", value: orgId } : { case: "instance", value: true }
  });
}










exports.createFeatureServiceClient = createFeatureServiceClient; exports.createIdpServiceClient = createIdpServiceClient; exports.createOIDCServiceClient = createOIDCServiceClient; exports.createOrganizationServiceClient = createOrganizationServiceClient; exports.createSAMLServiceClient = createSAMLServiceClient; exports.createSessionServiceClient = createSessionServiceClient; exports.createSettingsServiceClient = createSettingsServiceClient; exports.createUserServiceClient = createUserServiceClient; exports.makeReqCtx = makeReqCtx;
//# sourceMappingURL=v2.cjs.map