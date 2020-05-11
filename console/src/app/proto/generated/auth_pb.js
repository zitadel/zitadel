/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var google_api_annotations_pb = require('./google/api/annotations_pb.js');
goog.object.extend(proto, google_api_annotations_pb);
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');
goog.object.extend(proto, google_protobuf_empty_pb);
var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js');
goog.object.extend(proto, google_protobuf_struct_pb);
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');
goog.object.extend(proto, google_protobuf_timestamp_pb);
var protoc$gen$swagger_options_annotations_pb = require('./protoc-gen-swagger/options/annotations_pb.js');
goog.object.extend(proto, protoc$gen$swagger_options_annotations_pb);
var validate_validate_pb = require('./validate/validate_pb.js');
goog.object.extend(proto, validate_validate_pb);
var authoption_options_pb = require('./authoption/options_pb.js');
goog.object.extend(proto, authoption_options_pb);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AppState', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.Application', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ApplicationID', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ApplicationSearchKey', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ApplicationSearchQuery', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ApplicationSearchRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ApplicationSearchResponse', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthRequestOIDC', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthSessionCreation', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthSessionID', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthSessionResponse', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthSessionType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthSessionView', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.AuthUser', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.BrowserInformation', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ChooseUser', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ChooseUserData', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.CodeChallenge', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.CodeChallengeMethod', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.CreateTokenRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.Gender', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.Grant', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.GrantSearchKey', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.GrantSearchQuery', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.GrantSearchRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.GrantSearchResponse', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.IDPProvider', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.IP', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.IsAdminResponse', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.LoginData', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MFAState', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MfaOtpResponse', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MfaPromptData', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MfaType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MfaVerifyData', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MultiFactor', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MultiFactors', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MyPermissions', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MyProjectOrgSearchKey', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.NextStep', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.NextStepType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.NotificationType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.OIDCApplicationType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.OIDCAuthMethodType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.OIDCClientAuth', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.OIDCConfig', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.OIDCGrantType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.OIDCResponseType', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.Org', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.PasswordChange', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.PasswordData', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.PasswordID', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.PasswordRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.Prompt', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.RegisterUserRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ResetPassword', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.ResetPasswordRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.SearchMethod', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.SelectUserRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.SessionRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.SkipMfaInitRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.Token', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.TokenID', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UniqueUserRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UniqueUserResponse', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.User', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserAddress', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserAgent', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserAgentCreation', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserAgentID', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserAgentState', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserAgents', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserEmail', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserID', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserPhone', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserProfile', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserSession', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserSessionID', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserSessionState', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserSessionView', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserSessionViews', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserSessions', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.UserState', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyMfaOtp', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyMfaRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyPasswordRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyUserInitRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest', null, global);
goog.exportSymbol('proto.caos.citadel.auth.api.v1.VerifyUserRequest', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.SessionRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.SessionRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.SessionRequest.displayName = 'proto.caos.citadel.auth.api.v1.SessionRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserAgent = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserAgent, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserAgent.displayName = 'proto.caos.citadel.auth.api.v1.UserAgent';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserAgentID = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserAgentID, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserAgentID.displayName = 'proto.caos.citadel.auth.api.v1.UserAgentID';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserAgentCreation, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserAgentCreation.displayName = 'proto.caos.citadel.auth.api.v1.UserAgentCreation';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserAgents = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.UserAgents.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserAgents, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserAgents.displayName = 'proto.caos.citadel.auth.api.v1.UserAgents';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.AuthSessionCreation.repeatedFields_, proto.caos.citadel.auth.api.v1.AuthSessionCreation.oneofGroups_);
};
goog.inherits(proto.caos.citadel.auth.api.v1.AuthSessionCreation, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.AuthSessionCreation.displayName = 'proto.caos.citadel.auth.api.v1.AuthSessionCreation';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.AuthSessionResponse.repeatedFields_, proto.caos.citadel.auth.api.v1.AuthSessionResponse.oneofGroups_);
};
goog.inherits(proto.caos.citadel.auth.api.v1.AuthSessionResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.displayName = 'proto.caos.citadel.auth.api.v1.AuthSessionResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.AuthSessionView = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.AuthSessionView.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.AuthSessionView, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.AuthSessionView.displayName = 'proto.caos.citadel.auth.api.v1.AuthSessionView';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.TokenID = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.TokenID, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.TokenID.displayName = 'proto.caos.citadel.auth.api.v1.TokenID';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserSessionID = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserSessionID, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserSessionID.displayName = 'proto.caos.citadel.auth.api.v1.UserSessionID';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserSessions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.UserSessions.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserSessions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserSessions.displayName = 'proto.caos.citadel.auth.api.v1.UserSessions';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserSession = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserSession, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserSession.displayName = 'proto.caos.citadel.auth.api.v1.UserSession';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserSessionViews = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.UserSessionViews.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserSessionViews, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserSessionViews.displayName = 'proto.caos.citadel.auth.api.v1.UserSessionViews';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserSessionView = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserSessionView, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserSessionView.displayName = 'proto.caos.citadel.auth.api.v1.UserSessionView';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.AuthUser = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.AuthUser, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.AuthUser.displayName = 'proto.caos.citadel.auth.api.v1.AuthUser';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.AuthSessionID = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.AuthSessionID, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.AuthSessionID.displayName = 'proto.caos.citadel.auth.api.v1.AuthSessionID';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.SelectUserRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.SelectUserRequest.displayName = 'proto.caos.citadel.auth.api.v1.SelectUserRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyUserRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyUserRequest.displayName = 'proto.caos.citadel.auth.api.v1.VerifyUserRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyPasswordRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.displayName = 'proto.caos.citadel.auth.api.v1.VerifyPasswordRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.caos.citadel.auth.api.v1.VerifyMfaRequest.oneofGroups_);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyMfaRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyMfaRequest.displayName = 'proto.caos.citadel.auth.api.v1.VerifyMfaRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.displayName = 'proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.NextStep = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_);
};
goog.inherits(proto.caos.citadel.auth.api.v1.NextStep, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.NextStep.displayName = 'proto.caos.citadel.auth.api.v1.NextStep';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.LoginData = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.LoginData, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.LoginData.displayName = 'proto.caos.citadel.auth.api.v1.LoginData';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.PasswordData = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.PasswordData, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.PasswordData.displayName = 'proto.caos.citadel.auth.api.v1.PasswordData';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.MfaVerifyData.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MfaVerifyData, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MfaVerifyData.displayName = 'proto.caos.citadel.auth.api.v1.MfaVerifyData';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MfaPromptData = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.MfaPromptData.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MfaPromptData, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MfaPromptData.displayName = 'proto.caos.citadel.auth.api.v1.MfaPromptData';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ChooseUserData = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.ChooseUserData.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ChooseUserData, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ChooseUserData.displayName = 'proto.caos.citadel.auth.api.v1.ChooseUserData';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ChooseUser = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ChooseUser, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ChooseUser.displayName = 'proto.caos.citadel.auth.api.v1.ChooseUser';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.SkipMfaInitRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.displayName = 'proto.caos.citadel.auth.api.v1.SkipMfaInitRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.BrowserInformation = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.BrowserInformation, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.BrowserInformation.displayName = 'proto.caos.citadel.auth.api.v1.BrowserInformation';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.IP = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.IP, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.IP.displayName = 'proto.caos.citadel.auth.api.v1.IP';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.AuthRequestOIDC.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.AuthRequestOIDC, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.AuthRequestOIDC.displayName = 'proto.caos.citadel.auth.api.v1.AuthRequestOIDC';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.CodeChallenge = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.CodeChallenge, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.CodeChallenge.displayName = 'proto.caos.citadel.auth.api.v1.CodeChallenge';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserID = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserID, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserID.displayName = 'proto.caos.citadel.auth.api.v1.UserID';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UniqueUserRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UniqueUserRequest.displayName = 'proto.caos.citadel.auth.api.v1.UniqueUserRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UniqueUserResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UniqueUserResponse.displayName = 'proto.caos.citadel.auth.api.v1.UniqueUserResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.RegisterUserRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.RegisterUserRequest.displayName = 'proto.caos.citadel.auth.api.v1.RegisterUserRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.displayName = 'proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.IDPProvider = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.IDPProvider, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.IDPProvider.displayName = 'proto.caos.citadel.auth.api.v1.IDPProvider';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.User = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.User, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.User.displayName = 'proto.caos.citadel.auth.api.v1.User';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserProfile = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserProfile, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserProfile.displayName = 'proto.caos.citadel.auth.api.v1.UserProfile';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.displayName = 'proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserEmail = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserEmail, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserEmail.displayName = 'proto.caos.citadel.auth.api.v1.UserEmail';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.displayName = 'proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.displayName = 'proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.displayName = 'proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserPhone = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserPhone, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserPhone.displayName = 'proto.caos.citadel.auth.api.v1.UserPhone';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.displayName = 'proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.displayName = 'proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UserAddress = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UserAddress, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UserAddress.displayName = 'proto.caos.citadel.auth.api.v1.UserAddress';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.displayName = 'proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.PasswordID = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.PasswordID, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.PasswordID.displayName = 'proto.caos.citadel.auth.api.v1.PasswordID';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.PasswordRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.PasswordRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.PasswordRequest.displayName = 'proto.caos.citadel.auth.api.v1.PasswordRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ResetPasswordRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ResetPasswordRequest.displayName = 'proto.caos.citadel.auth.api.v1.ResetPasswordRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ResetPassword = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ResetPassword, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ResetPassword.displayName = 'proto.caos.citadel.auth.api.v1.ResetPassword';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.displayName = 'proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.PasswordChange = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.PasswordChange, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.PasswordChange.displayName = 'proto.caos.citadel.auth.api.v1.PasswordChange';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyMfaOtp, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyMfaOtp.displayName = 'proto.caos.citadel.auth.api.v1.VerifyMfaOtp';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MultiFactors = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.MultiFactors.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MultiFactors, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MultiFactors.displayName = 'proto.caos.citadel.auth.api.v1.MultiFactors';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MultiFactor = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MultiFactor, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MultiFactor.displayName = 'proto.caos.citadel.auth.api.v1.MultiFactor';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MfaOtpResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MfaOtpResponse.displayName = 'proto.caos.citadel.auth.api.v1.MfaOtpResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ApplicationID = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ApplicationID, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ApplicationID.displayName = 'proto.caos.citadel.auth.api.v1.ApplicationID';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.Application = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.caos.citadel.auth.api.v1.Application.oneofGroups_);
};
goog.inherits(proto.caos.citadel.auth.api.v1.Application, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.Application.displayName = 'proto.caos.citadel.auth.api.v1.Application';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.OIDCConfig = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.OIDCConfig.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.OIDCConfig, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.OIDCConfig.displayName = 'proto.caos.citadel.auth.api.v1.OIDCConfig';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ApplicationSearchRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.displayName = 'proto.caos.citadel.auth.api.v1.ApplicationSearchRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ApplicationSearchQuery, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.displayName = 'proto.caos.citadel.auth.api.v1.ApplicationSearchQuery';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ApplicationSearchResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.displayName = 'proto.caos.citadel.auth.api.v1.ApplicationSearchResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.oneofGroups_);
};
goog.inherits(proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.displayName = 'proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.OIDCClientAuth, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.OIDCClientAuth.displayName = 'proto.caos.citadel.auth.api.v1.OIDCClientAuth';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.GrantSearchRequest.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.GrantSearchRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.GrantSearchRequest.displayName = 'proto.caos.citadel.auth.api.v1.GrantSearchRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.GrantSearchQuery, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.GrantSearchQuery.displayName = 'proto.caos.citadel.auth.api.v1.GrantSearchQuery';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.GrantSearchResponse.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.GrantSearchResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.GrantSearchResponse.displayName = 'proto.caos.citadel.auth.api.v1.GrantSearchResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.Grant = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.Grant.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.Grant, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.Grant.displayName = 'proto.caos.citadel.auth.api.v1.Grant';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.displayName = 'proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.displayName = 'proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.displayName = 'proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.IsAdminResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.IsAdminResponse.displayName = 'proto.caos.citadel.auth.api.v1.IsAdminResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.Org = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.Org, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.Org.displayName = 'proto.caos.citadel.auth.api.v1.Org';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.CreateTokenRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.CreateTokenRequest.displayName = 'proto.caos.citadel.auth.api.v1.CreateTokenRequest';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.Token = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.Token, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.Token.displayName = 'proto.caos.citadel.auth.api.v1.Token';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.MyPermissions = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.caos.citadel.auth.api.v1.MyPermissions.repeatedFields_, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.MyPermissions, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.MyPermissions.displayName = 'proto.caos.citadel.auth.api.v1.MyPermissions';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.caos.citadel.auth.api.v1.VerifyUserInitRequest, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.displayName = 'proto.caos.citadel.auth.api.v1.VerifyUserInitRequest';
}



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.SessionRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.SessionRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SessionRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    userId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.SessionRequest}
 */
proto.caos.citadel.auth.api.v1.SessionRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.SessionRequest;
  return proto.caos.citadel.auth.api.v1.SessionRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.SessionRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.SessionRequest}
 */
proto.caos.citadel.auth.api.v1.SessionRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 2:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.SessionRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.SessionRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SessionRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string user_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.setUserId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional BrowserInformation browser_info = 2;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 2));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.SessionRequest.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 2) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserAgent.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserAgent} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgent.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f),
    state: jspb.Message.getFieldWithDefault(msg, 3, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgent}
 */
proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserAgent;
  return proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgent} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgent}
 */
proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    case 3:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.UserAgentState} */ (reader.readEnum());
      msg.setState(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserAgent.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgent} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgent.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional BrowserInformation browser_info = 2;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 2));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional UserAgentState state = 3;
 * @return {!proto.caos.citadel.auth.api.v1.UserAgentState}
 */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.getState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.UserAgentState} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.UserAgentState} value */
proto.caos.citadel.auth.api.v1.UserAgent.prototype.setState = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserAgentID.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserAgentID.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgentID.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgentID}
 */
proto.caos.citadel.auth.api.v1.UserAgentID.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserAgentID;
  return proto.caos.citadel.auth.api.v1.UserAgentID.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgentID}
 */
proto.caos.citadel.auth.api.v1.UserAgentID.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserAgentID.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserAgentID.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgentID.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAgentID.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAgentID.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserAgentCreation.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentCreation} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.toObject = function(includeInstance, msg) {
  var f, obj = {
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgentCreation}
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserAgentCreation;
  return proto.caos.citadel.auth.api.v1.UserAgentCreation.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentCreation} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgentCreation}
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserAgentCreation.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentCreation} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
};


/**
 * optional BrowserInformation browser_info = 1;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 1));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.UserAgentCreation.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 1, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserAgentCreation.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 1) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.UserAgents.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserAgents.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserAgents.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserAgents} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgents.toObject = function(includeInstance, msg) {
  var f, obj = {
    sessionsList: jspb.Message.toObjectList(msg.getSessionsList(),
    proto.caos.citadel.auth.api.v1.UserAgent.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgents}
 */
proto.caos.citadel.auth.api.v1.UserAgents.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserAgents;
  return proto.caos.citadel.auth.api.v1.UserAgents.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgents} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserAgents}
 */
proto.caos.citadel.auth.api.v1.UserAgents.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.caos.citadel.auth.api.v1.UserAgent;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinaryFromReader);
      msg.addSessions(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserAgents.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserAgents.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserAgents} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAgents.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getSessionsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.caos.citadel.auth.api.v1.UserAgent.serializeBinaryToWriter
    );
  }
};


/**
 * repeated UserAgent sessions = 1;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.UserAgent>}
 */
proto.caos.citadel.auth.api.v1.UserAgents.prototype.getSessionsList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.UserAgent>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.UserAgent, 1));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.UserAgent>} value */
proto.caos.citadel.auth.api.v1.UserAgents.prototype.setSessionsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgent=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.UserAgent}
 */
proto.caos.citadel.auth.api.v1.UserAgents.prototype.addSessions = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.caos.citadel.auth.api.v1.UserAgent, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.UserAgents.prototype.clearSessionsList = function() {
  this.setSessionsList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.repeatedFields_ = [8,9];

/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.oneofGroups_ = [[12]];

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.TypeInfoCase = {
  TYPE_INFO_NOT_SET: 0,
  OIDC: 12
};

/**
 * @return {proto.caos.citadel.auth.api.v1.AuthSessionCreation.TypeInfoCase}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getTypeInfoCase = function() {
  return /** @type {proto.caos.citadel.auth.api.v1.AuthSessionCreation.TypeInfoCase} */(jspb.Message.computeOneofCase(this, proto.caos.citadel.auth.api.v1.AuthSessionCreation.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.AuthSessionCreation.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionCreation} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    type: jspb.Message.getFieldWithDefault(msg, 2, 0),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f),
    clientId: jspb.Message.getFieldWithDefault(msg, 4, ""),
    redirectUri: jspb.Message.getFieldWithDefault(msg, 5, ""),
    state: jspb.Message.getFieldWithDefault(msg, 6, ""),
    prompt: jspb.Message.getFieldWithDefault(msg, 7, 0),
    authContextClassReferenceList: jspb.Message.getRepeatedField(msg, 8),
    uiLocalesList: jspb.Message.getRepeatedField(msg, 9),
    loginHint: jspb.Message.getFieldWithDefault(msg, 10, ""),
    maxAge: jspb.Message.getFieldWithDefault(msg, 11, 0),
    oidc: (f = msg.getOidc()) && proto.caos.citadel.auth.api.v1.AuthRequestOIDC.toObject(includeInstance, f),
    preselectedUserId: jspb.Message.getFieldWithDefault(msg, 13, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionCreation}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.AuthSessionCreation;
  return proto.caos.citadel.auth.api.v1.AuthSessionCreation.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionCreation} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionCreation}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.AuthSessionType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 3:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setRedirectUri(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setState(value);
      break;
    case 7:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.Prompt} */ (reader.readEnum());
      msg.setPrompt(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.addAuthContextClassReference(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.addUiLocales(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setLoginHint(value);
      break;
    case 11:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMaxAge(value);
      break;
    case 12:
      var value = new proto.caos.citadel.auth.api.v1.AuthRequestOIDC;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.AuthRequestOIDC.deserializeBinaryFromReader);
      msg.setOidc(value);
      break;
    case 13:
      var value = /** @type {string} */ (reader.readString());
      msg.setPreselectedUserId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.AuthSessionCreation.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionCreation} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getRedirectUri();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getState();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getPrompt();
  if (f !== 0.0) {
    writer.writeEnum(
      7,
      f
    );
  }
  f = message.getAuthContextClassReferenceList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      8,
      f
    );
  }
  f = message.getUiLocalesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      9,
      f
    );
  }
  f = message.getLoginHint();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getMaxAge();
  if (f !== 0) {
    writer.writeUint32(
      11,
      f
    );
  }
  f = message.getOidc();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      proto.caos.citadel.auth.api.v1.AuthRequestOIDC.serializeBinaryToWriter
    );
  }
  f = message.getPreselectedUserId();
  if (f.length > 0) {
    writer.writeString(
      13,
      f
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional AuthSessionType type = 2;
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionType}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.AuthSessionType} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.AuthSessionType} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional BrowserInformation browser_info = 3;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 3));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional string client_id = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string redirect_uri = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getRedirectUri = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setRedirectUri = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string state = 6;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getState = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setState = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional Prompt prompt = 7;
 * @return {!proto.caos.citadel.auth.api.v1.Prompt}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getPrompt = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.Prompt} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.Prompt} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setPrompt = function(value) {
  jspb.Message.setProto3EnumField(this, 7, value);
};


/**
 * repeated string auth_context_class_reference = 8;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getAuthContextClassReferenceList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 8));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setAuthContextClassReferenceList = function(value) {
  jspb.Message.setField(this, 8, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.addAuthContextClassReference = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 8, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.clearAuthContextClassReferenceList = function() {
  this.setAuthContextClassReferenceList([]);
};


/**
 * repeated string ui_locales = 9;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getUiLocalesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 9));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setUiLocalesList = function(value) {
  jspb.Message.setField(this, 9, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.addUiLocales = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 9, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.clearUiLocalesList = function() {
  this.setUiLocalesList([]);
};


/**
 * optional string login_hint = 10;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getLoginHint = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setLoginHint = function(value) {
  jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional uint32 max_age = 11;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getMaxAge = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setMaxAge = function(value) {
  jspb.Message.setProto3IntField(this, 11, value);
};


/**
 * optional AuthRequestOIDC oidc = 12;
 * @return {?proto.caos.citadel.auth.api.v1.AuthRequestOIDC}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getOidc = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.AuthRequestOIDC} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.AuthRequestOIDC, 12));
};


/** @param {?proto.caos.citadel.auth.api.v1.AuthRequestOIDC|undefined} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setOidc = function(value) {
  jspb.Message.setOneofWrapperField(this, 12, proto.caos.citadel.auth.api.v1.AuthSessionCreation.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.clearOidc = function() {
  this.setOidc(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.hasOidc = function() {
  return jspb.Message.getField(this, 12) != null;
};


/**
 * optional string preselected_user_id = 13;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.getPreselectedUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 13, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionCreation.prototype.setPreselectedUserId = function(value) {
  jspb.Message.setProto3StringField(this, 13, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.repeatedFields_ = [8,9,13,14];

/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.oneofGroups_ = [[12]];

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.TypeInfoCase = {
  TYPE_INFO_NOT_SET: 0,
  OIDC: 12
};

/**
 * @return {proto.caos.citadel.auth.api.v1.AuthSessionResponse.TypeInfoCase}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getTypeInfoCase = function() {
  return /** @type {proto.caos.citadel.auth.api.v1.AuthSessionResponse.TypeInfoCase} */(jspb.Message.computeOneofCase(this, proto.caos.citadel.auth.api.v1.AuthSessionResponse.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.AuthSessionResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    id: jspb.Message.getFieldWithDefault(msg, 2, ""),
    type: jspb.Message.getFieldWithDefault(msg, 3, 0),
    clientId: jspb.Message.getFieldWithDefault(msg, 4, ""),
    redirectUri: jspb.Message.getFieldWithDefault(msg, 5, ""),
    state: jspb.Message.getFieldWithDefault(msg, 6, ""),
    prompt: jspb.Message.getFieldWithDefault(msg, 7, 0),
    authContextClassReferenceList: jspb.Message.getRepeatedField(msg, 8),
    uiLocalesList: jspb.Message.getRepeatedField(msg, 9),
    loginHint: jspb.Message.getFieldWithDefault(msg, 10, ""),
    maxAge: jspb.Message.getFieldWithDefault(msg, 11, 0),
    oidc: (f = msg.getOidc()) && proto.caos.citadel.auth.api.v1.AuthRequestOIDC.toObject(includeInstance, f),
    possibleStepsList: jspb.Message.toObjectList(msg.getPossibleStepsList(),
    proto.caos.citadel.auth.api.v1.NextStep.toObject, includeInstance),
    projectClientIdsList: jspb.Message.getRepeatedField(msg, 14),
    userSession: (f = msg.getUserSession()) && proto.caos.citadel.auth.api.v1.UserSession.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionResponse}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.AuthSessionResponse;
  return proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionResponse}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 3:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.AuthSessionType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setRedirectUri(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setState(value);
      break;
    case 7:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.Prompt} */ (reader.readEnum());
      msg.setPrompt(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.addAuthContextClassReference(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.addUiLocales(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setLoginHint(value);
      break;
    case 11:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setMaxAge(value);
      break;
    case 12:
      var value = new proto.caos.citadel.auth.api.v1.AuthRequestOIDC;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.AuthRequestOIDC.deserializeBinaryFromReader);
      msg.setOidc(value);
      break;
    case 13:
      var value = new proto.caos.citadel.auth.api.v1.NextStep;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.NextStep.deserializeBinaryFromReader);
      msg.addPossibleSteps(value);
      break;
    case 14:
      var value = /** @type {string} */ (reader.readString());
      msg.addProjectClientIds(value);
      break;
    case 15:
      var value = new proto.caos.citadel.auth.api.v1.UserSession;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.UserSession.deserializeBinaryFromReader);
      msg.setUserSession(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getRedirectUri();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getState();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getPrompt();
  if (f !== 0.0) {
    writer.writeEnum(
      7,
      f
    );
  }
  f = message.getAuthContextClassReferenceList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      8,
      f
    );
  }
  f = message.getUiLocalesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      9,
      f
    );
  }
  f = message.getLoginHint();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getMaxAge();
  if (f !== 0) {
    writer.writeUint32(
      11,
      f
    );
  }
  f = message.getOidc();
  if (f != null) {
    writer.writeMessage(
      12,
      f,
      proto.caos.citadel.auth.api.v1.AuthRequestOIDC.serializeBinaryToWriter
    );
  }
  f = message.getPossibleStepsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      13,
      f,
      proto.caos.citadel.auth.api.v1.NextStep.serializeBinaryToWriter
    );
  }
  f = message.getProjectClientIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      14,
      f
    );
  }
  f = message.getUserSession();
  if (f != null) {
    writer.writeMessage(
      15,
      f,
      proto.caos.citadel.auth.api.v1.UserSession.serializeBinaryToWriter
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional AuthSessionType type = 3;
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionType}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.AuthSessionType} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.AuthSessionType} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional string client_id = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string redirect_uri = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getRedirectUri = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setRedirectUri = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string state = 6;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getState = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setState = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional Prompt prompt = 7;
 * @return {!proto.caos.citadel.auth.api.v1.Prompt}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getPrompt = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.Prompt} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.Prompt} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setPrompt = function(value) {
  jspb.Message.setProto3EnumField(this, 7, value);
};


/**
 * repeated string auth_context_class_reference = 8;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getAuthContextClassReferenceList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 8));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setAuthContextClassReferenceList = function(value) {
  jspb.Message.setField(this, 8, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.addAuthContextClassReference = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 8, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.clearAuthContextClassReferenceList = function() {
  this.setAuthContextClassReferenceList([]);
};


/**
 * repeated string ui_locales = 9;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getUiLocalesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 9));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setUiLocalesList = function(value) {
  jspb.Message.setField(this, 9, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.addUiLocales = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 9, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.clearUiLocalesList = function() {
  this.setUiLocalesList([]);
};


/**
 * optional string login_hint = 10;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getLoginHint = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setLoginHint = function(value) {
  jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional uint32 max_age = 11;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getMaxAge = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 11, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setMaxAge = function(value) {
  jspb.Message.setProto3IntField(this, 11, value);
};


/**
 * optional AuthRequestOIDC oidc = 12;
 * @return {?proto.caos.citadel.auth.api.v1.AuthRequestOIDC}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getOidc = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.AuthRequestOIDC} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.AuthRequestOIDC, 12));
};


/** @param {?proto.caos.citadel.auth.api.v1.AuthRequestOIDC|undefined} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setOidc = function(value) {
  jspb.Message.setOneofWrapperField(this, 12, proto.caos.citadel.auth.api.v1.AuthSessionResponse.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.clearOidc = function() {
  this.setOidc(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.hasOidc = function() {
  return jspb.Message.getField(this, 12) != null;
};


/**
 * repeated NextStep possible_steps = 13;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.NextStep>}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getPossibleStepsList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.NextStep>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.NextStep, 13));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.NextStep>} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setPossibleStepsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 13, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.NextStep=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.NextStep}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.addPossibleSteps = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 13, opt_value, proto.caos.citadel.auth.api.v1.NextStep, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.clearPossibleStepsList = function() {
  this.setPossibleStepsList([]);
};


/**
 * repeated string project_client_ids = 14;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getProjectClientIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 14));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setProjectClientIdsList = function(value) {
  jspb.Message.setField(this, 14, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.addProjectClientIds = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 14, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.clearProjectClientIdsList = function() {
  this.setProjectClientIdsList([]);
};


/**
 * optional UserSession user_session = 15;
 * @return {?proto.caos.citadel.auth.api.v1.UserSession}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.getUserSession = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.UserSession} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.UserSession, 15));
};


/** @param {?proto.caos.citadel.auth.api.v1.UserSession|undefined} value */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.setUserSession = function(value) {
  jspb.Message.setWrapperField(this, 15, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.clearUserSession = function() {
  this.setUserSession(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.AuthSessionResponse.prototype.hasUserSession = function() {
  return jspb.Message.getField(this, 15) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.repeatedFields_ = [6];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.AuthSessionView.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionView} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    authSessionId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    type: jspb.Message.getFieldWithDefault(msg, 3, 0),
    clientId: jspb.Message.getFieldWithDefault(msg, 4, ""),
    userSessionId: jspb.Message.getFieldWithDefault(msg, 5, ""),
    projectClientIdsList: jspb.Message.getRepeatedField(msg, 6),
    tokenId: jspb.Message.getFieldWithDefault(msg, 7, ""),
    tokenExpiration: (f = msg.getTokenExpiration()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    userId: jspb.Message.getFieldWithDefault(msg, 9, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionView}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.AuthSessionView;
  return proto.caos.citadel.auth.api.v1.AuthSessionView.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionView} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionView}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthSessionId(value);
      break;
    case 3:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.AuthSessionType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserSessionId(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.addProjectClientIds(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setTokenId(value);
      break;
    case 8:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setTokenExpiration(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.AuthSessionView.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionView} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAuthSessionId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getUserSessionId();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getProjectClientIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      6,
      f
    );
  }
  f = message.getTokenId();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getTokenExpiration();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string auth_session_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getAuthSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setAuthSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional AuthSessionType type = 3;
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionType}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.AuthSessionType} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.AuthSessionType} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional string client_id = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string user_session_id = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getUserSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setUserSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * repeated string project_client_ids = 6;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getProjectClientIdsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 6));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setProjectClientIdsList = function(value) {
  jspb.Message.setField(this, 6, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.addProjectClientIds = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 6, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.clearProjectClientIdsList = function() {
  this.setProjectClientIdsList([]);
};


/**
 * optional string token_id = 7;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getTokenId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setTokenId = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional google.protobuf.Timestamp token_expiration = 8;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getTokenExpiration = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 8));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setTokenExpiration = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.clearTokenExpiration = function() {
  this.setTokenExpiration(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.hasTokenExpiration = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional string user_id = 9;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionView.prototype.setUserId = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.TokenID.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.TokenID.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.TokenID} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.TokenID.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.TokenID}
 */
proto.caos.citadel.auth.api.v1.TokenID.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.TokenID;
  return proto.caos.citadel.auth.api.v1.TokenID.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.TokenID} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.TokenID}
 */
proto.caos.citadel.auth.api.v1.TokenID.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.TokenID.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.TokenID.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.TokenID} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.TokenID.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.TokenID.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.TokenID.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserSessionID.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserSessionID.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessionID.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    agentId: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionID}
 */
proto.caos.citadel.auth.api.v1.UserSessionID.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserSessionID;
  return proto.caos.citadel.auth.api.v1.UserSessionID.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionID}
 */
proto.caos.citadel.auth.api.v1.UserSessionID.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserSessionID.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserSessionID.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessionID.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSessionID.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSessionID.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string agent_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSessionID.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSessionID.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.UserSessions.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserSessions.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserSessions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserSessions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessions.toObject = function(includeInstance, msg) {
  var f, obj = {
    userSessionsList: jspb.Message.toObjectList(msg.getUserSessionsList(),
    proto.caos.citadel.auth.api.v1.UserSession.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessions}
 */
proto.caos.citadel.auth.api.v1.UserSessions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserSessions;
  return proto.caos.citadel.auth.api.v1.UserSessions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessions}
 */
proto.caos.citadel.auth.api.v1.UserSessions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.caos.citadel.auth.api.v1.UserSession;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.UserSession.deserializeBinaryFromReader);
      msg.addUserSessions(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserSessions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserSessions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserSessionsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.caos.citadel.auth.api.v1.UserSession.serializeBinaryToWriter
    );
  }
};


/**
 * repeated UserSession user_sessions = 1;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.UserSession>}
 */
proto.caos.citadel.auth.api.v1.UserSessions.prototype.getUserSessionsList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.UserSession>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.UserSession, 1));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.UserSession>} value */
proto.caos.citadel.auth.api.v1.UserSessions.prototype.setUserSessionsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserSession=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.UserSession}
 */
proto.caos.citadel.auth.api.v1.UserSessions.prototype.addUserSessions = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.caos.citadel.auth.api.v1.UserSession, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.UserSessions.prototype.clearUserSessionsList = function() {
  this.setUserSessionsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserSession.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserSession} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSession.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    agentId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    authState: jspb.Message.getFieldWithDefault(msg, 3, 0),
    user: (f = msg.getUser()) && proto.caos.citadel.auth.api.v1.AuthUser.toObject(includeInstance, f),
    passwordVerified: jspb.Message.getFieldWithDefault(msg, 5, false),
    mfa: jspb.Message.getFieldWithDefault(msg, 6, 0),
    mfaVerified: jspb.Message.getFieldWithDefault(msg, 7, false),
    authTime: (f = msg.getAuthTime()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserSession}
 */
proto.caos.citadel.auth.api.v1.UserSession.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserSession;
  return proto.caos.citadel.auth.api.v1.UserSession.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserSession} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserSession}
 */
proto.caos.citadel.auth.api.v1.UserSession.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 3:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.UserSessionState} */ (reader.readEnum());
      msg.setAuthState(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.AuthUser;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.AuthUser.deserializeBinaryFromReader);
      msg.setUser(value);
      break;
    case 5:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setPasswordVerified(value);
      break;
    case 6:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.MfaType} */ (reader.readEnum());
      msg.setMfa(value);
      break;
    case 7:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setMfaVerified(value);
      break;
    case 8:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setAuthTime(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserSession.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserSession} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSession.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getAuthState();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getUser();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.AuthUser.serializeBinaryToWriter
    );
  }
  f = message.getPasswordVerified();
  if (f) {
    writer.writeBool(
      5,
      f
    );
  }
  f = message.getMfa();
  if (f !== 0.0) {
    writer.writeEnum(
      6,
      f
    );
  }
  f = message.getMfaVerified();
  if (f) {
    writer.writeBool(
      7,
      f
    );
  }
  f = message.getAuthTime();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string agent_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional UserSessionState auth_state = 3;
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionState}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getAuthState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.UserSessionState} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.UserSessionState} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setAuthState = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional AuthUser user = 4;
 * @return {?proto.caos.citadel.auth.api.v1.AuthUser}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getUser = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.AuthUser} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.AuthUser, 4));
};


/** @param {?proto.caos.citadel.auth.api.v1.AuthUser|undefined} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setUser = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.clearUser = function() {
  this.setUser(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.hasUser = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional bool password_verified = 5;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getPasswordVerified = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 5, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setPasswordVerified = function(value) {
  jspb.Message.setProto3BooleanField(this, 5, value);
};


/**
 * optional MfaType mfa = 6;
 * @return {!proto.caos.citadel.auth.api.v1.MfaType}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getMfa = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.MfaType} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.MfaType} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setMfa = function(value) {
  jspb.Message.setProto3EnumField(this, 6, value);
};


/**
 * optional bool mfa_verified = 7;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getMfaVerified = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 7, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setMfaVerified = function(value) {
  jspb.Message.setProto3BooleanField(this, 7, value);
};


/**
 * optional google.protobuf.Timestamp auth_time = 8;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.getAuthTime = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 8));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.UserSession.prototype.setAuthTime = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.clearAuthTime = function() {
  this.setAuthTime(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserSession.prototype.hasAuthTime = function() {
  return jspb.Message.getField(this, 8) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserSessionViews.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionViews} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.toObject = function(includeInstance, msg) {
  var f, obj = {
    userSessionsList: jspb.Message.toObjectList(msg.getUserSessionsList(),
    proto.caos.citadel.auth.api.v1.UserSessionView.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionViews}
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserSessionViews;
  return proto.caos.citadel.auth.api.v1.UserSessionViews.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionViews} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionViews}
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.caos.citadel.auth.api.v1.UserSessionView;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.UserSessionView.deserializeBinaryFromReader);
      msg.addUserSessions(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserSessionViews.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionViews} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserSessionsList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.caos.citadel.auth.api.v1.UserSessionView.serializeBinaryToWriter
    );
  }
};


/**
 * repeated UserSessionView user_sessions = 1;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.UserSessionView>}
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.prototype.getUserSessionsList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.UserSessionView>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.UserSessionView, 1));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.UserSessionView>} value */
proto.caos.citadel.auth.api.v1.UserSessionViews.prototype.setUserSessionsList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionView=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionView}
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.prototype.addUserSessions = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.caos.citadel.auth.api.v1.UserSessionView, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.UserSessionViews.prototype.clearUserSessionsList = function() {
  this.setUserSessionsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserSessionView.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionView} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessionView.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    agentId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    authState: jspb.Message.getFieldWithDefault(msg, 3, 0),
    userId: jspb.Message.getFieldWithDefault(msg, 4, ""),
    userName: jspb.Message.getFieldWithDefault(msg, 5, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionView}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserSessionView;
  return proto.caos.citadel.auth.api.v1.UserSessionView.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionView} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionView}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 3:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.UserSessionState} */ (reader.readEnum());
      msg.setAuthState(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserSessionView.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionView} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserSessionView.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getAuthState();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string agent_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional UserSessionState auth_state = 3;
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionState}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.getAuthState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.UserSessionState} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.UserSessionState} value */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.setAuthState = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional string user_id = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.setUserId = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string user_name = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserSessionView.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.AuthUser.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.AuthUser.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.AuthUser} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthUser.toObject = function(includeInstance, msg) {
  var f, obj = {
    userId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userName: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.AuthUser}
 */
proto.caos.citadel.auth.api.v1.AuthUser.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.AuthUser;
  return proto.caos.citadel.auth.api.v1.AuthUser.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.AuthUser} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.AuthUser}
 */
proto.caos.citadel.auth.api.v1.AuthUser.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.AuthUser.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.AuthUser.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.AuthUser} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthUser.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string user_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthUser.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthUser.prototype.setUserId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string user_name = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthUser.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthUser.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.AuthSessionID.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionID} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    agentId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionID}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.AuthSessionID;
  return proto.caos.citadel.auth.api.v1.AuthSessionID.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionID} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionID}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 3:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.AuthSessionID.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionID} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string agent_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional BrowserInformation browser_info = 3;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 3));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.AuthSessionID.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 3) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.SelectUserRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.SelectUserRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    authSessionId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    userSessionId: jspb.Message.getFieldWithDefault(msg, 3, ""),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.SelectUserRequest}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.SelectUserRequest;
  return proto.caos.citadel.auth.api.v1.SelectUserRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.SelectUserRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.SelectUserRequest}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthSessionId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserSessionId(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.SelectUserRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.SelectUserRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAuthSessionId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUserSessionId();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string auth_session_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.getAuthSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.setAuthSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string user_session_id = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.getUserSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.setUserSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional BrowserInformation browser_info = 4;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 4));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.SelectUserRequest.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 4) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyUserRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    authSessionId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    userName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyUserRequest;
  return proto.caos.citadel.auth.api.v1.VerifyUserRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthSessionId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyUserRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAuthSessionId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string auth_session_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.getAuthSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.setAuthSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string user_name = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional BrowserInformation browser_info = 4;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 4));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.VerifyUserRequest.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 4) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    authSessionId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    password: jspb.Message.getFieldWithDefault(msg, 3, ""),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyPasswordRequest;
  return proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthSessionId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setPassword(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAuthSessionId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getPassword();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string auth_session_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.getAuthSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.setAuthSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string password = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.getPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.setPassword = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional BrowserInformation browser_info = 4;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 4));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.VerifyPasswordRequest.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 4) != null;
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.oneofGroups_ = [[4]];

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.MfaCase = {
  MFA_NOT_SET: 0,
  OTP: 4
};

/**
 * @return {proto.caos.citadel.auth.api.v1.VerifyMfaRequest.MfaCase}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.getMfaCase = function() {
  return /** @type {proto.caos.citadel.auth.api.v1.VerifyMfaRequest.MfaCase} */(jspb.Message.computeOneofCase(this, proto.caos.citadel.auth.api.v1.VerifyMfaRequest.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyMfaRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    authSessionId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    browserInfo: (f = msg.getBrowserInfo()) && proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(includeInstance, f),
    otp: (f = msg.getOtp()) && proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyMfaRequest;
  return proto.caos.citadel.auth.api.v1.VerifyMfaRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthSessionId(value);
      break;
    case 3:
      var value = new proto.caos.citadel.auth.api.v1.BrowserInformation;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader);
      msg.setBrowserInfo(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.deserializeBinaryFromReader);
      msg.setOtp(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyMfaRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAuthSessionId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getBrowserInfo();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter
    );
  }
  f = message.getOtp();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.serializeBinaryToWriter
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string auth_session_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.getAuthSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.setAuthSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional BrowserInformation browser_info = 3;
 * @return {?proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.getBrowserInfo = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.BrowserInformation} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.BrowserInformation, 3));
};


/** @param {?proto.caos.citadel.auth.api.v1.BrowserInformation|undefined} value */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.setBrowserInfo = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.clearBrowserInfo = function() {
  this.setBrowserInfo(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.hasBrowserInfo = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional AuthSessionMultiFactorOTP otp = 4;
 * @return {?proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.getOtp = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP, 4));
};


/** @param {?proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP|undefined} value */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.setOtp = function(value) {
  jspb.Message.setOneofWrapperField(this, 4, proto.caos.citadel.auth.api.v1.VerifyMfaRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.clearOtp = function() {
  this.setOtp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaRequest.prototype.hasOtp = function() {
  return jspb.Message.getField(this, 4) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.toObject = function(includeInstance, msg) {
  var f, obj = {
    code: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP}
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP;
  return proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP}
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCode();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string code = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.prototype.getCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthSessionMultiFactorOTP.prototype.setCode = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_ = [[2,3,4,5,6]];

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.NextStep.DataCase = {
  DATA_NOT_SET: 0,
  LOGIN: 2,
  PASSWORD: 3,
  MFA_VERIFY: 4,
  MFA_PROMPT: 5,
  CHOOSE_USER: 6
};

/**
 * @return {proto.caos.citadel.auth.api.v1.NextStep.DataCase}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.getDataCase = function() {
  return /** @type {proto.caos.citadel.auth.api.v1.NextStep.DataCase} */(jspb.Message.computeOneofCase(this, proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.NextStep.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.NextStep} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.NextStep.toObject = function(includeInstance, msg) {
  var f, obj = {
    type: jspb.Message.getFieldWithDefault(msg, 1, 0),
    login: (f = msg.getLogin()) && proto.caos.citadel.auth.api.v1.LoginData.toObject(includeInstance, f),
    password: (f = msg.getPassword()) && proto.caos.citadel.auth.api.v1.PasswordData.toObject(includeInstance, f),
    mfaVerify: (f = msg.getMfaVerify()) && proto.caos.citadel.auth.api.v1.MfaVerifyData.toObject(includeInstance, f),
    mfaPrompt: (f = msg.getMfaPrompt()) && proto.caos.citadel.auth.api.v1.MfaPromptData.toObject(includeInstance, f),
    chooseUser: (f = msg.getChooseUser()) && proto.caos.citadel.auth.api.v1.ChooseUserData.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.NextStep}
 */
proto.caos.citadel.auth.api.v1.NextStep.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.NextStep;
  return proto.caos.citadel.auth.api.v1.NextStep.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.NextStep} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.NextStep}
 */
proto.caos.citadel.auth.api.v1.NextStep.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.NextStepType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 2:
      var value = new proto.caos.citadel.auth.api.v1.LoginData;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.LoginData.deserializeBinaryFromReader);
      msg.setLogin(value);
      break;
    case 3:
      var value = new proto.caos.citadel.auth.api.v1.PasswordData;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.PasswordData.deserializeBinaryFromReader);
      msg.setPassword(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.MfaVerifyData;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.MfaVerifyData.deserializeBinaryFromReader);
      msg.setMfaVerify(value);
      break;
    case 5:
      var value = new proto.caos.citadel.auth.api.v1.MfaPromptData;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.MfaPromptData.deserializeBinaryFromReader);
      msg.setMfaPrompt(value);
      break;
    case 6:
      var value = new proto.caos.citadel.auth.api.v1.ChooseUserData;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.ChooseUserData.deserializeBinaryFromReader);
      msg.setChooseUser(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.NextStep.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.NextStep} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.NextStep.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getLogin();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.caos.citadel.auth.api.v1.LoginData.serializeBinaryToWriter
    );
  }
  f = message.getPassword();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      proto.caos.citadel.auth.api.v1.PasswordData.serializeBinaryToWriter
    );
  }
  f = message.getMfaVerify();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.MfaVerifyData.serializeBinaryToWriter
    );
  }
  f = message.getMfaPrompt();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      proto.caos.citadel.auth.api.v1.MfaPromptData.serializeBinaryToWriter
    );
  }
  f = message.getChooseUser();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      proto.caos.citadel.auth.api.v1.ChooseUserData.serializeBinaryToWriter
    );
  }
};


/**
 * optional NextStepType type = 1;
 * @return {!proto.caos.citadel.auth.api.v1.NextStepType}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.getType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.NextStepType} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.NextStepType} value */
proto.caos.citadel.auth.api.v1.NextStep.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional LoginData login = 2;
 * @return {?proto.caos.citadel.auth.api.v1.LoginData}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.getLogin = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.LoginData} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.LoginData, 2));
};


/** @param {?proto.caos.citadel.auth.api.v1.LoginData|undefined} value */
proto.caos.citadel.auth.api.v1.NextStep.prototype.setLogin = function(value) {
  jspb.Message.setOneofWrapperField(this, 2, proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.clearLogin = function() {
  this.setLogin(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.hasLogin = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional PasswordData password = 3;
 * @return {?proto.caos.citadel.auth.api.v1.PasswordData}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.getPassword = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.PasswordData} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.PasswordData, 3));
};


/** @param {?proto.caos.citadel.auth.api.v1.PasswordData|undefined} value */
proto.caos.citadel.auth.api.v1.NextStep.prototype.setPassword = function(value) {
  jspb.Message.setOneofWrapperField(this, 3, proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.clearPassword = function() {
  this.setPassword(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.hasPassword = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional MfaVerifyData mfa_verify = 4;
 * @return {?proto.caos.citadel.auth.api.v1.MfaVerifyData}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.getMfaVerify = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.MfaVerifyData} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.MfaVerifyData, 4));
};


/** @param {?proto.caos.citadel.auth.api.v1.MfaVerifyData|undefined} value */
proto.caos.citadel.auth.api.v1.NextStep.prototype.setMfaVerify = function(value) {
  jspb.Message.setOneofWrapperField(this, 4, proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.clearMfaVerify = function() {
  this.setMfaVerify(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.hasMfaVerify = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional MfaPromptData mfa_prompt = 5;
 * @return {?proto.caos.citadel.auth.api.v1.MfaPromptData}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.getMfaPrompt = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.MfaPromptData} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.MfaPromptData, 5));
};


/** @param {?proto.caos.citadel.auth.api.v1.MfaPromptData|undefined} value */
proto.caos.citadel.auth.api.v1.NextStep.prototype.setMfaPrompt = function(value) {
  jspb.Message.setOneofWrapperField(this, 5, proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.clearMfaPrompt = function() {
  this.setMfaPrompt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.hasMfaPrompt = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional ChooseUserData choose_user = 6;
 * @return {?proto.caos.citadel.auth.api.v1.ChooseUserData}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.getChooseUser = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.ChooseUserData} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.ChooseUserData, 6));
};


/** @param {?proto.caos.citadel.auth.api.v1.ChooseUserData|undefined} value */
proto.caos.citadel.auth.api.v1.NextStep.prototype.setChooseUser = function(value) {
  jspb.Message.setOneofWrapperField(this, 6, proto.caos.citadel.auth.api.v1.NextStep.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.clearChooseUser = function() {
  this.setChooseUser(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.NextStep.prototype.hasChooseUser = function() {
  return jspb.Message.getField(this, 6) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.LoginData.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.LoginData.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.LoginData} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.LoginData.toObject = function(includeInstance, msg) {
  var f, obj = {
    errMsg: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.LoginData}
 */
proto.caos.citadel.auth.api.v1.LoginData.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.LoginData;
  return proto.caos.citadel.auth.api.v1.LoginData.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.LoginData} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.LoginData}
 */
proto.caos.citadel.auth.api.v1.LoginData.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setErrMsg(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.LoginData.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.LoginData.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.LoginData} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.LoginData.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getErrMsg();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string err_msg = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.LoginData.prototype.getErrMsg = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.LoginData.prototype.setErrMsg = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.PasswordData.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.PasswordData.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.PasswordData} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordData.toObject = function(includeInstance, msg) {
  var f, obj = {
    errMsg: jspb.Message.getFieldWithDefault(msg, 1, ""),
    failureCount: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordData}
 */
proto.caos.citadel.auth.api.v1.PasswordData.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.PasswordData;
  return proto.caos.citadel.auth.api.v1.PasswordData.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordData} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordData}
 */
proto.caos.citadel.auth.api.v1.PasswordData.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setErrMsg(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setFailureCount(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.PasswordData.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.PasswordData.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordData} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordData.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getErrMsg();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getFailureCount();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
};


/**
 * optional string err_msg = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.PasswordData.prototype.getErrMsg = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.PasswordData.prototype.setErrMsg = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional uint32 failure_count = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.PasswordData.prototype.getFailureCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.PasswordData.prototype.setFailureCount = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MfaVerifyData.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MfaVerifyData} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.toObject = function(includeInstance, msg) {
  var f, obj = {
    errMsg: jspb.Message.getFieldWithDefault(msg, 1, ""),
    failureCount: jspb.Message.getFieldWithDefault(msg, 2, 0),
    mfaProvidersList: jspb.Message.getRepeatedField(msg, 3)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MfaVerifyData}
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MfaVerifyData;
  return proto.caos.citadel.auth.api.v1.MfaVerifyData.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MfaVerifyData} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MfaVerifyData}
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setErrMsg(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint32());
      msg.setFailureCount(value);
      break;
    case 3:
      var value = /** @type {!Array<!proto.caos.citadel.auth.api.v1.MfaType>} */ (reader.readPackedEnum());
      msg.setMfaProvidersList(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MfaVerifyData.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MfaVerifyData} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getErrMsg();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getFailureCount();
  if (f !== 0) {
    writer.writeUint32(
      2,
      f
    );
  }
  f = message.getMfaProvidersList();
  if (f.length > 0) {
    writer.writePackedEnum(
      3,
      f
    );
  }
};


/**
 * optional string err_msg = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.getErrMsg = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.setErrMsg = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional uint32 failure_count = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.getFailureCount = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.setFailureCount = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * repeated MfaType mfa_providers = 3;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.MfaType>}
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.getMfaProvidersList = function() {
  return /** @type {!Array<!proto.caos.citadel.auth.api.v1.MfaType>} */ (jspb.Message.getRepeatedField(this, 3));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.MfaType>} value */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.setMfaProvidersList = function(value) {
  jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.MfaType} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.addMfaProviders = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.MfaVerifyData.prototype.clearMfaProvidersList = function() {
  this.setMfaProvidersList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.repeatedFields_ = [2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MfaPromptData.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MfaPromptData} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.toObject = function(includeInstance, msg) {
  var f, obj = {
    required: jspb.Message.getFieldWithDefault(msg, 1, false),
    mfaProvidersList: jspb.Message.getRepeatedField(msg, 2)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MfaPromptData}
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MfaPromptData;
  return proto.caos.citadel.auth.api.v1.MfaPromptData.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MfaPromptData} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MfaPromptData}
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setRequired(value);
      break;
    case 2:
      var value = /** @type {!Array<!proto.caos.citadel.auth.api.v1.MfaType>} */ (reader.readPackedEnum());
      msg.setMfaProvidersList(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MfaPromptData.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MfaPromptData} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRequired();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
  f = message.getMfaProvidersList();
  if (f.length > 0) {
    writer.writePackedEnum(
      2,
      f
    );
  }
};


/**
 * optional bool required = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.getRequired = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.setRequired = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};


/**
 * repeated MfaType mfa_providers = 2;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.MfaType>}
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.getMfaProvidersList = function() {
  return /** @type {!Array<!proto.caos.citadel.auth.api.v1.MfaType>} */ (jspb.Message.getRepeatedField(this, 2));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.MfaType>} value */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.setMfaProvidersList = function(value) {
  jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.MfaType} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.addMfaProviders = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.MfaPromptData.prototype.clearMfaProvidersList = function() {
  this.setMfaProvidersList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ChooseUserData.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ChooseUserData} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.toObject = function(includeInstance, msg) {
  var f, obj = {
    usersList: jspb.Message.toObjectList(msg.getUsersList(),
    proto.caos.citadel.auth.api.v1.ChooseUser.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ChooseUserData}
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ChooseUserData;
  return proto.caos.citadel.auth.api.v1.ChooseUserData.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ChooseUserData} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ChooseUserData}
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.caos.citadel.auth.api.v1.ChooseUser;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.ChooseUser.deserializeBinaryFromReader);
      msg.addUsers(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ChooseUserData.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ChooseUserData} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUsersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.caos.citadel.auth.api.v1.ChooseUser.serializeBinaryToWriter
    );
  }
};


/**
 * repeated ChooseUser users = 1;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.ChooseUser>}
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.prototype.getUsersList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.ChooseUser>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.ChooseUser, 1));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.ChooseUser>} value */
proto.caos.citadel.auth.api.v1.ChooseUserData.prototype.setUsersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.ChooseUser=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.ChooseUser}
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.prototype.addUsers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.caos.citadel.auth.api.v1.ChooseUser, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.ChooseUserData.prototype.clearUsersList = function() {
  this.setUsersList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ChooseUser.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ChooseUser} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ChooseUser.toObject = function(includeInstance, msg) {
  var f, obj = {
    userSessionId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userId: jspb.Message.getFieldWithDefault(msg, 2, ""),
    userName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    userSessionState: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ChooseUser}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ChooseUser;
  return proto.caos.citadel.auth.api.v1.ChooseUser.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ChooseUser} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ChooseUser}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserSessionId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    case 4:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.UserSessionState} */ (reader.readEnum());
      msg.setUserSessionState(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ChooseUser.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ChooseUser} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ChooseUser.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserSessionId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getUserSessionState();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
};


/**
 * optional string user_session_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.getUserSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.setUserSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string user_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.setUserId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string user_name = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional UserSessionState user_session_state = 4;
 * @return {!proto.caos.citadel.auth.api.v1.UserSessionState}
 */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.getUserSessionState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.UserSessionState} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.UserSessionState} value */
proto.caos.citadel.auth.api.v1.ChooseUser.prototype.setUserSessionState = function(value) {
  jspb.Message.setProto3EnumField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    userId: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest}
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.SkipMfaInitRequest;
  return proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest}
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string user_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.SkipMfaInitRequest.prototype.setUserId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.BrowserInformation.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.BrowserInformation} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.toObject = function(includeInstance, msg) {
  var f, obj = {
    userAgent: jspb.Message.getFieldWithDefault(msg, 1, ""),
    remoteIp: (f = msg.getRemoteIp()) && proto.caos.citadel.auth.api.v1.IP.toObject(includeInstance, f),
    acceptLanguage: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.BrowserInformation;
  return proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.BrowserInformation} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.BrowserInformation}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserAgent(value);
      break;
    case 2:
      var value = new proto.caos.citadel.auth.api.v1.IP;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.IP.deserializeBinaryFromReader);
      msg.setRemoteIp(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setAcceptLanguage(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.BrowserInformation} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserAgent();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getRemoteIp();
  if (f != null) {
    writer.writeMessage(
      2,
      f,
      proto.caos.citadel.auth.api.v1.IP.serializeBinaryToWriter
    );
  }
  f = message.getAcceptLanguage();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string user_agent = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.getUserAgent = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.setUserAgent = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional IP remote_ip = 2;
 * @return {?proto.caos.citadel.auth.api.v1.IP}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.getRemoteIp = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.IP} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.IP, 2));
};


/** @param {?proto.caos.citadel.auth.api.v1.IP|undefined} value */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.setRemoteIp = function(value) {
  jspb.Message.setWrapperField(this, 2, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.clearRemoteIp = function() {
  this.setRemoteIp(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.hasRemoteIp = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * optional string accept_language = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.getAcceptLanguage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.BrowserInformation.prototype.setAcceptLanguage = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.IP.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.IP.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.IP} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.IP.toObject = function(includeInstance, msg) {
  var f, obj = {
    v4: jspb.Message.getFieldWithDefault(msg, 1, ""),
    v6: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.IP}
 */
proto.caos.citadel.auth.api.v1.IP.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.IP;
  return proto.caos.citadel.auth.api.v1.IP.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.IP} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.IP}
 */
proto.caos.citadel.auth.api.v1.IP.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setV4(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setV6(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.IP.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.IP.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.IP} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.IP.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getV4();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getV6();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string V4 = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.IP.prototype.getV4 = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.IP.prototype.setV4 = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string V6 = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.IP.prototype.getV6 = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.IP.prototype.setV6 = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.AuthRequestOIDC.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.AuthRequestOIDC} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.toObject = function(includeInstance, msg) {
  var f, obj = {
    scopeList: jspb.Message.getRepeatedField(msg, 1),
    responseType: jspb.Message.getFieldWithDefault(msg, 2, 0),
    nonce: jspb.Message.getFieldWithDefault(msg, 3, ""),
    codeChallenge: (f = msg.getCodeChallenge()) && proto.caos.citadel.auth.api.v1.CodeChallenge.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.AuthRequestOIDC}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.AuthRequestOIDC;
  return proto.caos.citadel.auth.api.v1.AuthRequestOIDC.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.AuthRequestOIDC} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.AuthRequestOIDC}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addScope(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.OIDCResponseType} */ (reader.readEnum());
      msg.setResponseType(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setNonce(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.CodeChallenge;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.CodeChallenge.deserializeBinaryFromReader);
      msg.setCodeChallenge(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.AuthRequestOIDC.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.AuthRequestOIDC} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getScopeList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
  f = message.getResponseType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getNonce();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getCodeChallenge();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.CodeChallenge.serializeBinaryToWriter
    );
  }
};


/**
 * repeated string scope = 1;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.getScopeList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.setScopeList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.addScope = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.clearScopeList = function() {
  this.setScopeList([]);
};


/**
 * optional OIDCResponseType response_type = 2;
 * @return {!proto.caos.citadel.auth.api.v1.OIDCResponseType}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.getResponseType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.OIDCResponseType} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.OIDCResponseType} value */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.setResponseType = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional string nonce = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.getNonce = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.setNonce = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional CodeChallenge code_challenge = 4;
 * @return {?proto.caos.citadel.auth.api.v1.CodeChallenge}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.getCodeChallenge = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.CodeChallenge} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.CodeChallenge, 4));
};


/** @param {?proto.caos.citadel.auth.api.v1.CodeChallenge|undefined} value */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.setCodeChallenge = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.clearCodeChallenge = function() {
  this.setCodeChallenge(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.AuthRequestOIDC.prototype.hasCodeChallenge = function() {
  return jspb.Message.getField(this, 4) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.CodeChallenge.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.CodeChallenge} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.toObject = function(includeInstance, msg) {
  var f, obj = {
    challenge: jspb.Message.getFieldWithDefault(msg, 1, ""),
    method: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.CodeChallenge}
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.CodeChallenge;
  return proto.caos.citadel.auth.api.v1.CodeChallenge.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.CodeChallenge} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.CodeChallenge}
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setChallenge(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.CodeChallengeMethod} */ (reader.readEnum());
      msg.setMethod(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.CodeChallenge.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.CodeChallenge} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getChallenge();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getMethod();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional string challenge = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.prototype.getChallenge = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.CodeChallenge.prototype.setChallenge = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional CodeChallengeMethod method = 2;
 * @return {!proto.caos.citadel.auth.api.v1.CodeChallengeMethod}
 */
proto.caos.citadel.auth.api.v1.CodeChallenge.prototype.getMethod = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.CodeChallengeMethod} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.CodeChallengeMethod} value */
proto.caos.citadel.auth.api.v1.CodeChallenge.prototype.setMethod = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserID.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserID.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserID} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserID.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserID}
 */
proto.caos.citadel.auth.api.v1.UserID.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserID;
  return proto.caos.citadel.auth.api.v1.UserID.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserID} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserID}
 */
proto.caos.citadel.auth.api.v1.UserID.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserID.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserID.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserID} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserID.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserID.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserID.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UniqueUserRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    userName: jspb.Message.getFieldWithDefault(msg, 1, ""),
    email: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UniqueUserRequest}
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UniqueUserRequest;
  return proto.caos.citadel.auth.api.v1.UniqueUserRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UniqueUserRequest}
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UniqueUserRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string user_name = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string email = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UniqueUserRequest.prototype.setEmail = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UniqueUserResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    isUnique: jspb.Message.getFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UniqueUserResponse}
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UniqueUserResponse;
  return proto.caos.citadel.auth.api.v1.UniqueUserResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UniqueUserResponse}
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsUnique(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UniqueUserResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getIsUnique();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool is_unique = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.prototype.getIsUnique = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.UniqueUserResponse.prototype.setIsUnique = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.RegisterUserRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    email: jspb.Message.getFieldWithDefault(msg, 1, ""),
    firstName: jspb.Message.getFieldWithDefault(msg, 2, ""),
    lastName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    nickName: jspb.Message.getFieldWithDefault(msg, 4, ""),
    displayName: jspb.Message.getFieldWithDefault(msg, 5, ""),
    preferredLanguage: jspb.Message.getFieldWithDefault(msg, 6, ""),
    gender: jspb.Message.getFieldWithDefault(msg, 7, 0),
    password: jspb.Message.getFieldWithDefault(msg, 8, ""),
    orgId: jspb.Message.getFieldWithDefault(msg, 9, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.RegisterUserRequest}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.RegisterUserRequest;
  return proto.caos.citadel.auth.api.v1.RegisterUserRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.RegisterUserRequest}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setFirstName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLastName(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setNickName(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setDisplayName(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setPreferredLanguage(value);
      break;
    case 7:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (reader.readEnum());
      msg.setGender(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setPassword(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setOrgId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.RegisterUserRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getFirstName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLastName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getNickName();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getDisplayName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getPreferredLanguage();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getGender();
  if (f !== 0.0) {
    writer.writeEnum(
      7,
      f
    );
  }
  f = message.getPassword();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getOrgId();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};


/**
 * optional string email = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setEmail = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string first_name = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getFirstName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setFirstName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string last_name = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getLastName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setLastName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string nick_name = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getNickName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setNickName = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string display_name = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getDisplayName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setDisplayName = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string preferred_language = 6;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getPreferredLanguage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setPreferredLanguage = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional Gender gender = 7;
 * @return {!proto.caos.citadel.auth.api.v1.Gender}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getGender = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.Gender} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setGender = function(value) {
  jspb.Message.setProto3EnumField(this, 7, value);
};


/**
 * optional string password = 8;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setPassword = function(value) {
  jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string org_id = 9;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.getOrgId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserRequest.prototype.setOrgId = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    email: jspb.Message.getFieldWithDefault(msg, 1, ""),
    firstName: jspb.Message.getFieldWithDefault(msg, 2, ""),
    lastName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    nickName: jspb.Message.getFieldWithDefault(msg, 4, ""),
    displayName: jspb.Message.getFieldWithDefault(msg, 5, ""),
    preferredLanguage: jspb.Message.getFieldWithDefault(msg, 6, ""),
    gender: jspb.Message.getFieldWithDefault(msg, 7, 0),
    idpProvider: (f = msg.getIdpProvider()) && proto.caos.citadel.auth.api.v1.IDPProvider.toObject(includeInstance, f),
    orgId: jspb.Message.getFieldWithDefault(msg, 9, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest;
  return proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setFirstName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLastName(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setNickName(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setDisplayName(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setPreferredLanguage(value);
      break;
    case 7:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (reader.readEnum());
      msg.setGender(value);
      break;
    case 8:
      var value = new proto.caos.citadel.auth.api.v1.IDPProvider;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.IDPProvider.deserializeBinaryFromReader);
      msg.setIdpProvider(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setOrgId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getFirstName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLastName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getNickName();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getDisplayName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getPreferredLanguage();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getGender();
  if (f !== 0.0) {
    writer.writeEnum(
      7,
      f
    );
  }
  f = message.getIdpProvider();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.caos.citadel.auth.api.v1.IDPProvider.serializeBinaryToWriter
    );
  }
  f = message.getOrgId();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};


/**
 * optional string email = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setEmail = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string first_name = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getFirstName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setFirstName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string last_name = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getLastName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setLastName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string nick_name = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getNickName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setNickName = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string display_name = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getDisplayName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setDisplayName = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string preferred_language = 6;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getPreferredLanguage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setPreferredLanguage = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional Gender gender = 7;
 * @return {!proto.caos.citadel.auth.api.v1.Gender}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getGender = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.Gender} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setGender = function(value) {
  jspb.Message.setProto3EnumField(this, 7, value);
};


/**
 * optional IDPProvider idp_provider = 8;
 * @return {?proto.caos.citadel.auth.api.v1.IDPProvider}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getIdpProvider = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.IDPProvider} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.IDPProvider, 8));
};


/** @param {?proto.caos.citadel.auth.api.v1.IDPProvider|undefined} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setIdpProvider = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.clearIdpProvider = function() {
  this.setIdpProvider(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.hasIdpProvider = function() {
  return jspb.Message.getField(this, 8) != null;
};


/**
 * optional string org_id = 9;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.getOrgId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest.prototype.setOrgId = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.IDPProvider.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.IDPProvider.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.IDPProvider} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.IDPProvider.toObject = function(includeInstance, msg) {
  var f, obj = {
    provider: jspb.Message.getFieldWithDefault(msg, 8, ""),
    externalidpid: jspb.Message.getFieldWithDefault(msg, 9, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.IDPProvider}
 */
proto.caos.citadel.auth.api.v1.IDPProvider.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.IDPProvider;
  return proto.caos.citadel.auth.api.v1.IDPProvider.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.IDPProvider} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.IDPProvider}
 */
proto.caos.citadel.auth.api.v1.IDPProvider.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setProvider(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setExternalidpid(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.IDPProvider.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.IDPProvider.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.IDPProvider} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.IDPProvider.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getProvider();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getExternalidpid();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
};


/**
 * optional string provider = 8;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.IDPProvider.prototype.getProvider = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.IDPProvider.prototype.setProvider = function(value) {
  jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string externalIdpID = 9;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.IDPProvider.prototype.getExternalidpid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.IDPProvider.prototype.setExternalidpid = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.User.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.User.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.User} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.User.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    state: jspb.Message.getFieldWithDefault(msg, 2, 0),
    creationDate: (f = msg.getCreationDate()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    activationDate: (f = msg.getActivationDate()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    changeDate: (f = msg.getChangeDate()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    lastLogin: (f = msg.getLastLogin()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    passwordChanged: (f = msg.getPasswordChanged()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    userName: jspb.Message.getFieldWithDefault(msg, 8, ""),
    firstName: jspb.Message.getFieldWithDefault(msg, 9, ""),
    lastName: jspb.Message.getFieldWithDefault(msg, 10, ""),
    nickName: jspb.Message.getFieldWithDefault(msg, 11, ""),
    displayName: jspb.Message.getFieldWithDefault(msg, 12, ""),
    preferredLanguage: jspb.Message.getFieldWithDefault(msg, 13, ""),
    gender: jspb.Message.getFieldWithDefault(msg, 14, 0),
    email: jspb.Message.getFieldWithDefault(msg, 15, ""),
    isemailverified: jspb.Message.getFieldWithDefault(msg, 16, false),
    phone: jspb.Message.getFieldWithDefault(msg, 17, ""),
    isphoneverified: jspb.Message.getFieldWithDefault(msg, 18, false),
    country: jspb.Message.getFieldWithDefault(msg, 19, ""),
    locality: jspb.Message.getFieldWithDefault(msg, 20, ""),
    postalCode: jspb.Message.getFieldWithDefault(msg, 21, ""),
    region: jspb.Message.getFieldWithDefault(msg, 22, ""),
    streetAddress: jspb.Message.getFieldWithDefault(msg, 23, ""),
    passwordChangeRequired: jspb.Message.getFieldWithDefault(msg, 24, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.User}
 */
proto.caos.citadel.auth.api.v1.User.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.User;
  return proto.caos.citadel.auth.api.v1.User.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.User} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.User}
 */
proto.caos.citadel.auth.api.v1.User.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.UserState} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 3:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setCreationDate(value);
      break;
    case 4:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setActivationDate(value);
      break;
    case 5:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setChangeDate(value);
      break;
    case 6:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setLastLogin(value);
      break;
    case 7:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setPasswordChanged(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    case 9:
      var value = /** @type {string} */ (reader.readString());
      msg.setFirstName(value);
      break;
    case 10:
      var value = /** @type {string} */ (reader.readString());
      msg.setLastName(value);
      break;
    case 11:
      var value = /** @type {string} */ (reader.readString());
      msg.setNickName(value);
      break;
    case 12:
      var value = /** @type {string} */ (reader.readString());
      msg.setDisplayName(value);
      break;
    case 13:
      var value = /** @type {string} */ (reader.readString());
      msg.setPreferredLanguage(value);
      break;
    case 14:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (reader.readEnum());
      msg.setGender(value);
      break;
    case 15:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    case 16:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsemailverified(value);
      break;
    case 17:
      var value = /** @type {string} */ (reader.readString());
      msg.setPhone(value);
      break;
    case 18:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsphoneverified(value);
      break;
    case 19:
      var value = /** @type {string} */ (reader.readString());
      msg.setCountry(value);
      break;
    case 20:
      var value = /** @type {string} */ (reader.readString());
      msg.setLocality(value);
      break;
    case 21:
      var value = /** @type {string} */ (reader.readString());
      msg.setPostalCode(value);
      break;
    case 22:
      var value = /** @type {string} */ (reader.readString());
      msg.setRegion(value);
      break;
    case 23:
      var value = /** @type {string} */ (reader.readString());
      msg.setStreetAddress(value);
      break;
    case 24:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setPasswordChangeRequired(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.User.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.User.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.User} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.User.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getCreationDate();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getActivationDate();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getChangeDate();
  if (f != null) {
    writer.writeMessage(
      5,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getLastLogin();
  if (f != null) {
    writer.writeMessage(
      6,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getPasswordChanged();
  if (f != null) {
    writer.writeMessage(
      7,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getFirstName();
  if (f.length > 0) {
    writer.writeString(
      9,
      f
    );
  }
  f = message.getLastName();
  if (f.length > 0) {
    writer.writeString(
      10,
      f
    );
  }
  f = message.getNickName();
  if (f.length > 0) {
    writer.writeString(
      11,
      f
    );
  }
  f = message.getDisplayName();
  if (f.length > 0) {
    writer.writeString(
      12,
      f
    );
  }
  f = message.getPreferredLanguage();
  if (f.length > 0) {
    writer.writeString(
      13,
      f
    );
  }
  f = message.getGender();
  if (f !== 0.0) {
    writer.writeEnum(
      14,
      f
    );
  }
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      15,
      f
    );
  }
  f = message.getIsemailverified();
  if (f) {
    writer.writeBool(
      16,
      f
    );
  }
  f = message.getPhone();
  if (f.length > 0) {
    writer.writeString(
      17,
      f
    );
  }
  f = message.getIsphoneverified();
  if (f) {
    writer.writeBool(
      18,
      f
    );
  }
  f = message.getCountry();
  if (f.length > 0) {
    writer.writeString(
      19,
      f
    );
  }
  f = message.getLocality();
  if (f.length > 0) {
    writer.writeString(
      20,
      f
    );
  }
  f = message.getPostalCode();
  if (f.length > 0) {
    writer.writeString(
      21,
      f
    );
  }
  f = message.getRegion();
  if (f.length > 0) {
    writer.writeString(
      22,
      f
    );
  }
  f = message.getStreetAddress();
  if (f.length > 0) {
    writer.writeString(
      23,
      f
    );
  }
  f = message.getPasswordChangeRequired();
  if (f) {
    writer.writeBool(
      24,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional UserState state = 2;
 * @return {!proto.caos.citadel.auth.api.v1.UserState}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.UserState} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.UserState} value */
proto.caos.citadel.auth.api.v1.User.prototype.setState = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional google.protobuf.Timestamp creation_date = 3;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getCreationDate = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 3));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.User.prototype.setCreationDate = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.User.prototype.clearCreationDate = function() {
  this.setCreationDate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.hasCreationDate = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional google.protobuf.Timestamp activation_date = 4;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getActivationDate = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 4));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.User.prototype.setActivationDate = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.User.prototype.clearActivationDate = function() {
  this.setActivationDate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.hasActivationDate = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional google.protobuf.Timestamp change_date = 5;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getChangeDate = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 5));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.User.prototype.setChangeDate = function(value) {
  jspb.Message.setWrapperField(this, 5, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.User.prototype.clearChangeDate = function() {
  this.setChangeDate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.hasChangeDate = function() {
  return jspb.Message.getField(this, 5) != null;
};


/**
 * optional google.protobuf.Timestamp last_login = 6;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getLastLogin = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 6));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.User.prototype.setLastLogin = function(value) {
  jspb.Message.setWrapperField(this, 6, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.User.prototype.clearLastLogin = function() {
  this.setLastLogin(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.hasLastLogin = function() {
  return jspb.Message.getField(this, 6) != null;
};


/**
 * optional google.protobuf.Timestamp password_changed = 7;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getPasswordChanged = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 7));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.User.prototype.setPasswordChanged = function(value) {
  jspb.Message.setWrapperField(this, 7, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.User.prototype.clearPasswordChanged = function() {
  this.setPasswordChanged(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.hasPasswordChanged = function() {
  return jspb.Message.getField(this, 7) != null;
};


/**
 * optional string user_name = 8;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional string first_name = 9;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getFirstName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 9, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setFirstName = function(value) {
  jspb.Message.setProto3StringField(this, 9, value);
};


/**
 * optional string last_name = 10;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getLastName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 10, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setLastName = function(value) {
  jspb.Message.setProto3StringField(this, 10, value);
};


/**
 * optional string nick_name = 11;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getNickName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 11, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setNickName = function(value) {
  jspb.Message.setProto3StringField(this, 11, value);
};


/**
 * optional string display_name = 12;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getDisplayName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 12, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setDisplayName = function(value) {
  jspb.Message.setProto3StringField(this, 12, value);
};


/**
 * optional string preferred_language = 13;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getPreferredLanguage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 13, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setPreferredLanguage = function(value) {
  jspb.Message.setProto3StringField(this, 13, value);
};


/**
 * optional Gender gender = 14;
 * @return {!proto.caos.citadel.auth.api.v1.Gender}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getGender = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (jspb.Message.getFieldWithDefault(this, 14, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.Gender} value */
proto.caos.citadel.auth.api.v1.User.prototype.setGender = function(value) {
  jspb.Message.setProto3EnumField(this, 14, value);
};


/**
 * optional string email = 15;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 15, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setEmail = function(value) {
  jspb.Message.setProto3StringField(this, 15, value);
};


/**
 * optional bool isEmailVerified = 16;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getIsemailverified = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 16, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.User.prototype.setIsemailverified = function(value) {
  jspb.Message.setProto3BooleanField(this, 16, value);
};


/**
 * optional string phone = 17;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getPhone = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 17, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setPhone = function(value) {
  jspb.Message.setProto3StringField(this, 17, value);
};


/**
 * optional bool isPhoneVerified = 18;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getIsphoneverified = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 18, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.User.prototype.setIsphoneverified = function(value) {
  jspb.Message.setProto3BooleanField(this, 18, value);
};


/**
 * optional string country = 19;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getCountry = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 19, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setCountry = function(value) {
  jspb.Message.setProto3StringField(this, 19, value);
};


/**
 * optional string locality = 20;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getLocality = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 20, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setLocality = function(value) {
  jspb.Message.setProto3StringField(this, 20, value);
};


/**
 * optional string postal_code = 21;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getPostalCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 21, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setPostalCode = function(value) {
  jspb.Message.setProto3StringField(this, 21, value);
};


/**
 * optional string region = 22;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getRegion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 22, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setRegion = function(value) {
  jspb.Message.setProto3StringField(this, 22, value);
};


/**
 * optional string street_address = 23;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getStreetAddress = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 23, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.User.prototype.setStreetAddress = function(value) {
  jspb.Message.setProto3StringField(this, 23, value);
};


/**
 * optional bool password_change_required = 24;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.User.prototype.getPasswordChangeRequired = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 24, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.User.prototype.setPasswordChangeRequired = function(value) {
  jspb.Message.setProto3BooleanField(this, 24, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserProfile.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserProfile} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserProfile.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    userName: jspb.Message.getFieldWithDefault(msg, 2, ""),
    firstName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    lastName: jspb.Message.getFieldWithDefault(msg, 4, ""),
    nickName: jspb.Message.getFieldWithDefault(msg, 5, ""),
    displayName: jspb.Message.getFieldWithDefault(msg, 6, ""),
    preferredLanguage: jspb.Message.getFieldWithDefault(msg, 7, ""),
    gender: jspb.Message.getFieldWithDefault(msg, 8, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserProfile}
 */
proto.caos.citadel.auth.api.v1.UserProfile.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserProfile;
  return proto.caos.citadel.auth.api.v1.UserProfile.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserProfile} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserProfile}
 */
proto.caos.citadel.auth.api.v1.UserProfile.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setFirstName(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setLastName(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setNickName(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setDisplayName(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setPreferredLanguage(value);
      break;
    case 8:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (reader.readEnum());
      msg.setGender(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserProfile.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserProfile} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserProfile.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getFirstName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getLastName();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getNickName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getDisplayName();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getPreferredLanguage();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getGender();
  if (f !== 0.0) {
    writer.writeEnum(
      8,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string user_name = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string first_name = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getFirstName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setFirstName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string last_name = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getLastName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setLastName = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string nick_name = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getNickName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setNickName = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string display_name = 6;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getDisplayName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setDisplayName = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string preferred_language = 7;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getPreferredLanguage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setPreferredLanguage = function(value) {
  jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional Gender gender = 8;
 * @return {!proto.caos.citadel.auth.api.v1.Gender}
 */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.getGender = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (jspb.Message.getFieldWithDefault(this, 8, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.Gender} value */
proto.caos.citadel.auth.api.v1.UserProfile.prototype.setGender = function(value) {
  jspb.Message.setProto3EnumField(this, 8, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    firstName: jspb.Message.getFieldWithDefault(msg, 1, ""),
    lastName: jspb.Message.getFieldWithDefault(msg, 2, ""),
    nickName: jspb.Message.getFieldWithDefault(msg, 3, ""),
    displayName: jspb.Message.getFieldWithDefault(msg, 4, ""),
    preferredLanguage: jspb.Message.getFieldWithDefault(msg, 5, ""),
    gender: jspb.Message.getFieldWithDefault(msg, 6, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest;
  return proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setFirstName(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setLastName(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setNickName(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setDisplayName(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setPreferredLanguage(value);
      break;
    case 6:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (reader.readEnum());
      msg.setGender(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getFirstName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLastName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getNickName();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getDisplayName();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getPreferredLanguage();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getGender();
  if (f !== 0.0) {
    writer.writeEnum(
      6,
      f
    );
  }
};


/**
 * optional string first_name = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.getFirstName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.setFirstName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string last_name = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.getLastName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.setLastName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string nick_name = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.getNickName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.setNickName = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string display_name = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.getDisplayName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.setDisplayName = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string preferred_language = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.getPreferredLanguage = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.setPreferredLanguage = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional Gender gender = 6;
 * @return {!proto.caos.citadel.auth.api.v1.Gender}
 */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.getGender = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.Gender} */ (jspb.Message.getFieldWithDefault(this, 6, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.Gender} value */
proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest.prototype.setGender = function(value) {
  jspb.Message.setProto3EnumField(this, 6, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserEmail.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserEmail} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserEmail.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    email: jspb.Message.getFieldWithDefault(msg, 2, ""),
    isemailverified: jspb.Message.getFieldWithDefault(msg, 3, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserEmail}
 */
proto.caos.citadel.auth.api.v1.UserEmail.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserEmail;
  return proto.caos.citadel.auth.api.v1.UserEmail.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserEmail} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserEmail}
 */
proto.caos.citadel.auth.api.v1.UserEmail.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsemailverified(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserEmail.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserEmail} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserEmail.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getIsemailverified();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string email = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.setEmail = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool isEmailVerified = 3;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.getIsemailverified = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 3, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.UserEmail.prototype.setIsemailverified = function(value) {
  jspb.Message.setProto3BooleanField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    code: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest;
  return proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCode();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string code = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.prototype.getCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest.prototype.setCode = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    code: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest;
  return proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCode();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string code = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.prototype.getCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest.prototype.setCode = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    email: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest;
  return proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setEmail(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getEmail();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string email = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.prototype.getEmail = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest.prototype.setEmail = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserPhone.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserPhone} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserPhone.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    phone: jspb.Message.getFieldWithDefault(msg, 2, ""),
    isphoneverified: jspb.Message.getFieldWithDefault(msg, 3, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserPhone}
 */
proto.caos.citadel.auth.api.v1.UserPhone.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserPhone;
  return proto.caos.citadel.auth.api.v1.UserPhone.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserPhone} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserPhone}
 */
proto.caos.citadel.auth.api.v1.UserPhone.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setPhone(value);
      break;
    case 3:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsphoneverified(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserPhone.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserPhone} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserPhone.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getPhone();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getIsphoneverified();
  if (f) {
    writer.writeBool(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string phone = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.getPhone = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.setPhone = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional bool isPhoneVerified = 3;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.getIsphoneverified = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 3, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.UserPhone.prototype.setIsphoneverified = function(value) {
  jspb.Message.setProto3BooleanField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    phone: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest;
  return proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPhone(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPhone();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string phone = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.prototype.getPhone = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest.prototype.setPhone = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    code: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest;
  return proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCode();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string code = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.prototype.getCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest.prototype.setCode = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UserAddress.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UserAddress} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAddress.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    country: jspb.Message.getFieldWithDefault(msg, 2, ""),
    locality: jspb.Message.getFieldWithDefault(msg, 3, ""),
    postalCode: jspb.Message.getFieldWithDefault(msg, 4, ""),
    region: jspb.Message.getFieldWithDefault(msg, 5, ""),
    streetAddress: jspb.Message.getFieldWithDefault(msg, 6, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UserAddress}
 */
proto.caos.citadel.auth.api.v1.UserAddress.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UserAddress;
  return proto.caos.citadel.auth.api.v1.UserAddress.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UserAddress} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UserAddress}
 */
proto.caos.citadel.auth.api.v1.UserAddress.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCountry(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setLocality(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setPostalCode(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setRegion(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setStreetAddress(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UserAddress.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UserAddress} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UserAddress.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCountry();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getLocality();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getPostalCode();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getRegion();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getStreetAddress();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string country = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.getCountry = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.setCountry = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string locality = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.getLocality = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.setLocality = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string postal_code = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.getPostalCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.setPostalCode = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string region = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.getRegion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.setRegion = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string street_address = 6;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.getStreetAddress = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UserAddress.prototype.setStreetAddress = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    country: jspb.Message.getFieldWithDefault(msg, 1, ""),
    locality: jspb.Message.getFieldWithDefault(msg, 2, ""),
    postalCode: jspb.Message.getFieldWithDefault(msg, 3, ""),
    region: jspb.Message.getFieldWithDefault(msg, 4, ""),
    streetAddress: jspb.Message.getFieldWithDefault(msg, 5, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest;
  return proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCountry(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setLocality(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setPostalCode(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setRegion(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setStreetAddress(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCountry();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getLocality();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getPostalCode();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getRegion();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getStreetAddress();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
};


/**
 * optional string country = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.getCountry = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.setCountry = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string locality = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.getLocality = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.setLocality = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string postal_code = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.getPostalCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.setPostalCode = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string region = 4;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.getRegion = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.setRegion = function(value) {
  jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string street_address = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.getStreetAddress = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest.prototype.setStreetAddress = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.PasswordID.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.PasswordID.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.PasswordID} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordID.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordID}
 */
proto.caos.citadel.auth.api.v1.PasswordID.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.PasswordID;
  return proto.caos.citadel.auth.api.v1.PasswordID.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordID} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordID}
 */
proto.caos.citadel.auth.api.v1.PasswordID.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.PasswordID.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.PasswordID.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordID} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordID.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.PasswordID.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.PasswordID.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.PasswordRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.PasswordRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.PasswordRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    password: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordRequest}
 */
proto.caos.citadel.auth.api.v1.PasswordRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.PasswordRequest;
  return proto.caos.citadel.auth.api.v1.PasswordRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordRequest}
 */
proto.caos.citadel.auth.api.v1.PasswordRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setPassword(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.PasswordRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.PasswordRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPassword();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string password = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.PasswordRequest.prototype.getPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.PasswordRequest.prototype.setPassword = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ResetPasswordRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    userName: jspb.Message.getFieldWithDefault(msg, 1, ""),
    type: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest}
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ResetPasswordRequest;
  return proto.caos.citadel.auth.api.v1.ResetPasswordRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest}
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserName(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.NotificationType} */ (reader.readEnum());
      msg.setType(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ResetPasswordRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserName();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional string user_name = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.prototype.getUserName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.prototype.setUserName = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional NotificationType type = 2;
 * @return {!proto.caos.citadel.auth.api.v1.NotificationType}
 */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.prototype.getType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.NotificationType} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.NotificationType} value */
proto.caos.citadel.auth.api.v1.ResetPasswordRequest.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ResetPassword.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ResetPassword} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ResetPassword.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    code: jspb.Message.getFieldWithDefault(msg, 2, ""),
    newPassword: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ResetPassword}
 */
proto.caos.citadel.auth.api.v1.ResetPassword.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ResetPassword;
  return proto.caos.citadel.auth.api.v1.ResetPassword.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ResetPassword} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ResetPassword}
 */
proto.caos.citadel.auth.api.v1.ResetPassword.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCode(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setNewPassword(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ResetPassword.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ResetPassword} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ResetPassword.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCode();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getNewPassword();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string code = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.getCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.setCode = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string new_password = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.getNewPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ResetPassword.prototype.setNewPassword = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    type: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest}
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest;
  return proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest}
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.NotificationType} */ (reader.readEnum());
      msg.setType(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional NotificationType type = 2;
 * @return {!proto.caos.citadel.auth.api.v1.NotificationType}
 */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.prototype.getType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.NotificationType} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.NotificationType} value */
proto.caos.citadel.auth.api.v1.SetPasswordNotificationRequest.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.PasswordChange.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.PasswordChange.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.PasswordChange} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordChange.toObject = function(includeInstance, msg) {
  var f, obj = {
    oldPassword: jspb.Message.getFieldWithDefault(msg, 1, ""),
    newPassword: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordChange}
 */
proto.caos.citadel.auth.api.v1.PasswordChange.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.PasswordChange;
  return proto.caos.citadel.auth.api.v1.PasswordChange.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordChange} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.PasswordChange}
 */
proto.caos.citadel.auth.api.v1.PasswordChange.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setOldPassword(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setNewPassword(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.PasswordChange.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.PasswordChange.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.PasswordChange} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.PasswordChange.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOldPassword();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getNewPassword();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string old_password = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.PasswordChange.prototype.getOldPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.PasswordChange.prototype.setOldPassword = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string new_password = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.PasswordChange.prototype.getNewPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.PasswordChange.prototype.setNewPassword = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyMfaOtp.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.toObject = function(includeInstance, msg) {
  var f, obj = {
    code: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyMfaOtp;
  return proto.caos.citadel.auth.api.v1.VerifyMfaOtp.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setCode(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyMfaOtp.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCode();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string code = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.prototype.getCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyMfaOtp.prototype.setCode = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.MultiFactors.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MultiFactors.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MultiFactors.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MultiFactors} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MultiFactors.toObject = function(includeInstance, msg) {
  var f, obj = {
    mfasList: jspb.Message.toObjectList(msg.getMfasList(),
    proto.caos.citadel.auth.api.v1.MultiFactor.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MultiFactors}
 */
proto.caos.citadel.auth.api.v1.MultiFactors.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MultiFactors;
  return proto.caos.citadel.auth.api.v1.MultiFactors.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MultiFactors} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MultiFactors}
 */
proto.caos.citadel.auth.api.v1.MultiFactors.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.caos.citadel.auth.api.v1.MultiFactor;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.MultiFactor.deserializeBinaryFromReader);
      msg.addMfas(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MultiFactors.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MultiFactors.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MultiFactors} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MultiFactors.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getMfasList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      1,
      f,
      proto.caos.citadel.auth.api.v1.MultiFactor.serializeBinaryToWriter
    );
  }
};


/**
 * repeated MultiFactor mfas = 1;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.MultiFactor>}
 */
proto.caos.citadel.auth.api.v1.MultiFactors.prototype.getMfasList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.MultiFactor>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.MultiFactor, 1));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.MultiFactor>} value */
proto.caos.citadel.auth.api.v1.MultiFactors.prototype.setMfasList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 1, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.MultiFactor=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.MultiFactor}
 */
proto.caos.citadel.auth.api.v1.MultiFactors.prototype.addMfas = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 1, opt_value, proto.caos.citadel.auth.api.v1.MultiFactor, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.MultiFactors.prototype.clearMfasList = function() {
  this.setMfasList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MultiFactor.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MultiFactor.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MultiFactor} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MultiFactor.toObject = function(includeInstance, msg) {
  var f, obj = {
    type: jspb.Message.getFieldWithDefault(msg, 1, 0),
    state: jspb.Message.getFieldWithDefault(msg, 2, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MultiFactor}
 */
proto.caos.citadel.auth.api.v1.MultiFactor.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MultiFactor;
  return proto.caos.citadel.auth.api.v1.MultiFactor.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MultiFactor} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MultiFactor}
 */
proto.caos.citadel.auth.api.v1.MultiFactor.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.MfaType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.MFAState} */ (reader.readEnum());
      msg.setState(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MultiFactor.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MultiFactor.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MultiFactor} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MultiFactor.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getType();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
};


/**
 * optional MfaType type = 1;
 * @return {!proto.caos.citadel.auth.api.v1.MfaType}
 */
proto.caos.citadel.auth.api.v1.MultiFactor.prototype.getType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.MfaType} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.MfaType} value */
proto.caos.citadel.auth.api.v1.MultiFactor.prototype.setType = function(value) {
  jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional MFAState state = 2;
 * @return {!proto.caos.citadel.auth.api.v1.MFAState}
 */
proto.caos.citadel.auth.api.v1.MultiFactor.prototype.getState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.MFAState} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.MFAState} value */
proto.caos.citadel.auth.api.v1.MultiFactor.prototype.setState = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MfaOtpResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MfaOtpResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    userId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    url: jspb.Message.getFieldWithDefault(msg, 2, ""),
    secret: jspb.Message.getFieldWithDefault(msg, 3, ""),
    state: jspb.Message.getFieldWithDefault(msg, 4, 0)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MfaOtpResponse}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MfaOtpResponse;
  return proto.caos.citadel.auth.api.v1.MfaOtpResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MfaOtpResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MfaOtpResponse}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setUrl(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setSecret(value);
      break;
    case 4:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.MFAState} */ (reader.readEnum());
      msg.setState(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MfaOtpResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MfaOtpResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getUserId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUrl();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getSecret();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
};


/**
 * optional string user_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.getUserId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.setUserId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string url = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.getUrl = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.setUrl = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string secret = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.getSecret = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.setSecret = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional MFAState state = 4;
 * @return {!proto.caos.citadel.auth.api.v1.MFAState}
 */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.getState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.MFAState} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.MFAState} value */
proto.caos.citadel.auth.api.v1.MfaOtpResponse.prototype.setState = function(value) {
  jspb.Message.setProto3EnumField(this, 4, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ApplicationID.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ApplicationID.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationID} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationID.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationID}
 */
proto.caos.citadel.auth.api.v1.ApplicationID.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ApplicationID;
  return proto.caos.citadel.auth.api.v1.ApplicationID.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationID} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationID}
 */
proto.caos.citadel.auth.api.v1.ApplicationID.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ApplicationID.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ApplicationID.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationID} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationID.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ApplicationID.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ApplicationID.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.caos.citadel.auth.api.v1.Application.oneofGroups_ = [[8]];

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.Application.AppConfigCase = {
  APP_CONFIG_NOT_SET: 0,
  OIDC_CONFIG: 8
};

/**
 * @return {proto.caos.citadel.auth.api.v1.Application.AppConfigCase}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.getAppConfigCase = function() {
  return /** @type {proto.caos.citadel.auth.api.v1.Application.AppConfigCase} */(jspb.Message.computeOneofCase(this, proto.caos.citadel.auth.api.v1.Application.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.Application.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.Application} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Application.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    state: jspb.Message.getFieldWithDefault(msg, 2, 0),
    creationDate: (f = msg.getCreationDate()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    changeDate: (f = msg.getChangeDate()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f),
    name: jspb.Message.getFieldWithDefault(msg, 5, ""),
    oidcConfig: (f = msg.getOidcConfig()) && proto.caos.citadel.auth.api.v1.OIDCConfig.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.Application}
 */
proto.caos.citadel.auth.api.v1.Application.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.Application;
  return proto.caos.citadel.auth.api.v1.Application.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.Application} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.Application}
 */
proto.caos.citadel.auth.api.v1.Application.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.AppState} */ (reader.readEnum());
      msg.setState(value);
      break;
    case 3:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setCreationDate(value);
      break;
    case 4:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setChangeDate(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 8:
      var value = new proto.caos.citadel.auth.api.v1.OIDCConfig;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.OIDCConfig.deserializeBinaryFromReader);
      msg.setOidcConfig(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.Application.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.Application} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Application.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getState();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getCreationDate();
  if (f != null) {
    writer.writeMessage(
      3,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getChangeDate();
  if (f != null) {
    writer.writeMessage(
      4,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getOidcConfig();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      proto.caos.citadel.auth.api.v1.OIDCConfig.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Application.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional AppState state = 2;
 * @return {!proto.caos.citadel.auth.api.v1.AppState}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.getState = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.AppState} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.AppState} value */
proto.caos.citadel.auth.api.v1.Application.prototype.setState = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional google.protobuf.Timestamp creation_date = 3;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.getCreationDate = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 3));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.Application.prototype.setCreationDate = function(value) {
  jspb.Message.setWrapperField(this, 3, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.Application.prototype.clearCreationDate = function() {
  this.setCreationDate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.hasCreationDate = function() {
  return jspb.Message.getField(this, 3) != null;
};


/**
 * optional google.protobuf.Timestamp change_date = 4;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.getChangeDate = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 4));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.Application.prototype.setChangeDate = function(value) {
  jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.Application.prototype.clearChangeDate = function() {
  this.setChangeDate(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.hasChangeDate = function() {
  return jspb.Message.getField(this, 4) != null;
};


/**
 * optional string name = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Application.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional OIDCConfig oidc_config = 8;
 * @return {?proto.caos.citadel.auth.api.v1.OIDCConfig}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.getOidcConfig = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.OIDCConfig} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.OIDCConfig, 8));
};


/** @param {?proto.caos.citadel.auth.api.v1.OIDCConfig|undefined} value */
proto.caos.citadel.auth.api.v1.Application.prototype.setOidcConfig = function(value) {
  jspb.Message.setOneofWrapperField(this, 8, proto.caos.citadel.auth.api.v1.Application.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.Application.prototype.clearOidcConfig = function() {
  this.setOidcConfig(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.Application.prototype.hasOidcConfig = function() {
  return jspb.Message.getField(this, 8) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.repeatedFields_ = [1,2,3,8];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.OIDCConfig.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.OIDCConfig} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.toObject = function(includeInstance, msg) {
  var f, obj = {
    redirectUrisList: jspb.Message.getRepeatedField(msg, 1),
    responseTypesList: jspb.Message.getRepeatedField(msg, 2),
    grantTypesList: jspb.Message.getRepeatedField(msg, 3),
    applicationType: jspb.Message.getFieldWithDefault(msg, 4, 0),
    clientSecret: jspb.Message.getFieldWithDefault(msg, 5, ""),
    clientId: jspb.Message.getFieldWithDefault(msg, 6, ""),
    authMethodType: jspb.Message.getFieldWithDefault(msg, 7, 0),
    postLogoutRedirectUrisList: jspb.Message.getRepeatedField(msg, 8)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.OIDCConfig}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.OIDCConfig;
  return proto.caos.citadel.auth.api.v1.OIDCConfig.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.OIDCConfig} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.OIDCConfig}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addRedirectUris(value);
      break;
    case 2:
      var value = /** @type {!Array<!proto.caos.citadel.auth.api.v1.OIDCResponseType>} */ (reader.readPackedEnum());
      msg.setResponseTypesList(value);
      break;
    case 3:
      var value = /** @type {!Array<!proto.caos.citadel.auth.api.v1.OIDCGrantType>} */ (reader.readPackedEnum());
      msg.setGrantTypesList(value);
      break;
    case 4:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.OIDCApplicationType} */ (reader.readEnum());
      msg.setApplicationType(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientSecret(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 7:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.OIDCAuthMethodType} */ (reader.readEnum());
      msg.setAuthMethodType(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.addPostLogoutRedirectUris(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.OIDCConfig.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.OIDCConfig} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getRedirectUrisList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
  f = message.getResponseTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      2,
      f
    );
  }
  f = message.getGrantTypesList();
  if (f.length > 0) {
    writer.writePackedEnum(
      3,
      f
    );
  }
  f = message.getApplicationType();
  if (f !== 0.0) {
    writer.writeEnum(
      4,
      f
    );
  }
  f = message.getClientSecret();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getAuthMethodType();
  if (f !== 0.0) {
    writer.writeEnum(
      7,
      f
    );
  }
  f = message.getPostLogoutRedirectUrisList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      8,
      f
    );
  }
};


/**
 * repeated string redirect_uris = 1;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getRedirectUrisList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setRedirectUrisList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.addRedirectUris = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.clearRedirectUrisList = function() {
  this.setRedirectUrisList([]);
};


/**
 * repeated OIDCResponseType response_types = 2;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.OIDCResponseType>}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getResponseTypesList = function() {
  return /** @type {!Array<!proto.caos.citadel.auth.api.v1.OIDCResponseType>} */ (jspb.Message.getRepeatedField(this, 2));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.OIDCResponseType>} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setResponseTypesList = function(value) {
  jspb.Message.setField(this, 2, value || []);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.OIDCResponseType} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.addResponseTypes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 2, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.clearResponseTypesList = function() {
  this.setResponseTypesList([]);
};


/**
 * repeated OIDCGrantType grant_types = 3;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.OIDCGrantType>}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getGrantTypesList = function() {
  return /** @type {!Array<!proto.caos.citadel.auth.api.v1.OIDCGrantType>} */ (jspb.Message.getRepeatedField(this, 3));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.OIDCGrantType>} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setGrantTypesList = function(value) {
  jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.OIDCGrantType} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.addGrantTypes = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.clearGrantTypesList = function() {
  this.setGrantTypesList([]);
};


/**
 * optional OIDCApplicationType application_type = 4;
 * @return {!proto.caos.citadel.auth.api.v1.OIDCApplicationType}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getApplicationType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.OIDCApplicationType} */ (jspb.Message.getFieldWithDefault(this, 4, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.OIDCApplicationType} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setApplicationType = function(value) {
  jspb.Message.setProto3EnumField(this, 4, value);
};


/**
 * optional string client_secret = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getClientSecret = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setClientSecret = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string client_id = 6;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional OIDCAuthMethodType auth_method_type = 7;
 * @return {!proto.caos.citadel.auth.api.v1.OIDCAuthMethodType}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getAuthMethodType = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.OIDCAuthMethodType} */ (jspb.Message.getFieldWithDefault(this, 7, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.OIDCAuthMethodType} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setAuthMethodType = function(value) {
  jspb.Message.setProto3EnumField(this, 7, value);
};


/**
 * repeated string post_logout_redirect_uris = 8;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.getPostLogoutRedirectUrisList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 8));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.setPostLogoutRedirectUrisList = function(value) {
  jspb.Message.setField(this, 8, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.addPostLogoutRedirectUris = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 8, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.OIDCConfig.prototype.clearPostLogoutRedirectUrisList = function() {
  this.setPostLogoutRedirectUrisList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.repeatedFields_ = [5];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    offset: jspb.Message.getFieldWithDefault(msg, 1, 0),
    limit: jspb.Message.getFieldWithDefault(msg, 2, 0),
    sortingColumn: jspb.Message.getFieldWithDefault(msg, 3, 0),
    asc: jspb.Message.getFieldWithDefault(msg, 4, false),
    queriesList: jspb.Message.toObjectList(msg.getQueriesList(),
    proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ApplicationSearchRequest;
  return proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setOffset(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setLimit(value);
      break;
    case 3:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey} */ (reader.readEnum());
      msg.setSortingColumn(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setAsc(value);
      break;
    case 5:
      var value = new proto.caos.citadel.auth.api.v1.ApplicationSearchQuery;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.deserializeBinaryFromReader);
      msg.addQueries(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOffset();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getLimit();
  if (f !== 0) {
    writer.writeUint64(
      2,
      f
    );
  }
  f = message.getSortingColumn();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getAsc();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getQueriesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 offset = 1;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.getOffset = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.setOffset = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint64 limit = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.getLimit = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.setLimit = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional ApplicationSearchKey sorting_column = 3;
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.getSortingColumn = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.setSortingColumn = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional bool asc = 4;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.getAsc = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 4, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.setAsc = function(value) {
  jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * repeated ApplicationSearchQuery queries = 5;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery>}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.getQueriesList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.ApplicationSearchQuery, 5));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery>} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.setQueriesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.addQueries = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.caos.citadel.auth.api.v1.ApplicationSearchQuery, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchRequest.prototype.clearQueriesList = function() {
  this.setQueriesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.toObject = function(includeInstance, msg) {
  var f, obj = {
    key: jspb.Message.getFieldWithDefault(msg, 1, 0),
    method: jspb.Message.getFieldWithDefault(msg, 2, 0),
    value: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ApplicationSearchQuery;
  return proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey} */ (reader.readEnum());
      msg.setKey(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.SearchMethod} */ (reader.readEnum());
      msg.setMethod(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchQuery} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKey();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getMethod();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getValue();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional ApplicationSearchKey key = 1;
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.getKey = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchKey} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.setKey = function(value) {
  jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional SearchMethod method = 2;
 * @return {!proto.caos.citadel.auth.api.v1.SearchMethod}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.getMethod = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.SearchMethod} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.SearchMethod} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.setMethod = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional string value = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchQuery.prototype.setValue = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    offset: jspb.Message.getFieldWithDefault(msg, 1, 0),
    limit: jspb.Message.getFieldWithDefault(msg, 2, 0),
    totalResult: jspb.Message.getFieldWithDefault(msg, 3, 0),
    resultList: jspb.Message.toObjectList(msg.getResultList(),
    proto.caos.citadel.auth.api.v1.Application.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchResponse}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ApplicationSearchResponse;
  return proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationSearchResponse}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setOffset(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setLimit(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setTotalResult(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.Application;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.Application.deserializeBinaryFromReader);
      msg.addResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOffset();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getLimit();
  if (f !== 0) {
    writer.writeUint64(
      2,
      f
    );
  }
  f = message.getTotalResult();
  if (f !== 0) {
    writer.writeUint64(
      3,
      f
    );
  }
  f = message.getResultList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.Application.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 offset = 1;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.getOffset = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.setOffset = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint64 limit = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.getLimit = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.setLimit = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint64 total_result = 3;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.getTotalResult = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.setTotalResult = function(value) {
  jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * repeated Application result = 4;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.Application>}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.getResultList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.Application>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.Application, 4));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.Application>} value */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.setResultList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.Application=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.Application}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.addResult = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.caos.citadel.auth.api.v1.Application, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.prototype.clearResultList = function() {
  this.setResultList([]);
};



/**
 * Oneof group definitions for this message. Each group defines the field
 * numbers belonging to that group. When of these fields' value is set, all
 * other fields in the group are cleared. During deserialization, if multiple
 * fields are encountered for a group, only the last value seen will be kept.
 * @private {!Array<!Array<number>>}
 * @const
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.oneofGroups_ = [[1]];

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.AuthCase = {
  AUTH_NOT_SET: 0,
  OIDC_CLIENT_AUTH: 1
};

/**
 * @return {proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.AuthCase}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.prototype.getAuthCase = function() {
  return /** @type {proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.AuthCase} */(jspb.Message.computeOneofCase(this, proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.oneofGroups_[0]));
};



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    oidcClientAuth: (f = msg.getOidcClientAuth()) && proto.caos.citadel.auth.api.v1.OIDCClientAuth.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest;
  return proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = new proto.caos.citadel.auth.api.v1.OIDCClientAuth;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.OIDCClientAuth.deserializeBinaryFromReader);
      msg.setOidcClientAuth(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOidcClientAuth();
  if (f != null) {
    writer.writeMessage(
      1,
      f,
      proto.caos.citadel.auth.api.v1.OIDCClientAuth.serializeBinaryToWriter
    );
  }
};


/**
 * optional OIDCClientAuth oidc_client_auth = 1;
 * @return {?proto.caos.citadel.auth.api.v1.OIDCClientAuth}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.prototype.getOidcClientAuth = function() {
  return /** @type{?proto.caos.citadel.auth.api.v1.OIDCClientAuth} */ (
    jspb.Message.getWrapperField(this, proto.caos.citadel.auth.api.v1.OIDCClientAuth, 1));
};


/** @param {?proto.caos.citadel.auth.api.v1.OIDCClientAuth|undefined} value */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.prototype.setOidcClientAuth = function(value) {
  jspb.Message.setOneofWrapperField(this, 1, proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.oneofGroups_[0], value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.prototype.clearOidcClientAuth = function() {
  this.setOidcClientAuth(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest.prototype.hasOidcClientAuth = function() {
  return jspb.Message.getField(this, 1) != null;
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.OIDCClientAuth.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.OIDCClientAuth} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.toObject = function(includeInstance, msg) {
  var f, obj = {
    clientId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    clientSecret: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.OIDCClientAuth}
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.OIDCClientAuth;
  return proto.caos.citadel.auth.api.v1.OIDCClientAuth.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.OIDCClientAuth} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.OIDCClientAuth}
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setClientSecret(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.OIDCClientAuth.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.OIDCClientAuth} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getClientId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getClientSecret();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string client_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.prototype.getClientId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.prototype.setClientId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string client_secret = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.prototype.getClientSecret = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.OIDCClientAuth.prototype.setClientSecret = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.repeatedFields_ = [5];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.GrantSearchRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    offset: jspb.Message.getFieldWithDefault(msg, 1, 0),
    limit: jspb.Message.getFieldWithDefault(msg, 2, 0),
    sortingColumn: jspb.Message.getFieldWithDefault(msg, 3, 0),
    asc: jspb.Message.getFieldWithDefault(msg, 4, false),
    queriesList: jspb.Message.toObjectList(msg.getQueriesList(),
    proto.caos.citadel.auth.api.v1.GrantSearchQuery.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchRequest}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.GrantSearchRequest;
  return proto.caos.citadel.auth.api.v1.GrantSearchRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchRequest}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setOffset(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setLimit(value);
      break;
    case 3:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.GrantSearchKey} */ (reader.readEnum());
      msg.setSortingColumn(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setAsc(value);
      break;
    case 5:
      var value = new proto.caos.citadel.auth.api.v1.GrantSearchQuery;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.GrantSearchQuery.deserializeBinaryFromReader);
      msg.addQueries(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.GrantSearchRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOffset();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getLimit();
  if (f !== 0) {
    writer.writeUint64(
      2,
      f
    );
  }
  f = message.getSortingColumn();
  if (f !== 0.0) {
    writer.writeEnum(
      3,
      f
    );
  }
  f = message.getAsc();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getQueriesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.caos.citadel.auth.api.v1.GrantSearchQuery.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 offset = 1;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.getOffset = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.setOffset = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint64 limit = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.getLimit = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.setLimit = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional GrantSearchKey sorting_column = 3;
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchKey}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.getSortingColumn = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.GrantSearchKey} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.GrantSearchKey} value */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.setSortingColumn = function(value) {
  jspb.Message.setProto3EnumField(this, 3, value);
};


/**
 * optional bool asc = 4;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.getAsc = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 4, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.setAsc = function(value) {
  jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * repeated GrantSearchQuery queries = 5;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.GrantSearchQuery>}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.getQueriesList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.GrantSearchQuery>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.GrantSearchQuery, 5));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.GrantSearchQuery>} value */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.setQueriesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchQuery=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchQuery}
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.addQueries = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.caos.citadel.auth.api.v1.GrantSearchQuery, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.GrantSearchRequest.prototype.clearQueriesList = function() {
  this.setQueriesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.GrantSearchQuery.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchQuery} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.toObject = function(includeInstance, msg) {
  var f, obj = {
    key: jspb.Message.getFieldWithDefault(msg, 1, 0),
    method: jspb.Message.getFieldWithDefault(msg, 2, 0),
    value: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchQuery}
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.GrantSearchQuery;
  return proto.caos.citadel.auth.api.v1.GrantSearchQuery.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchQuery} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchQuery}
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.GrantSearchKey} */ (reader.readEnum());
      msg.setKey(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.SearchMethod} */ (reader.readEnum());
      msg.setMethod(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.GrantSearchQuery.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchQuery} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKey();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getMethod();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getValue();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional GrantSearchKey key = 1;
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchKey}
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.getKey = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.GrantSearchKey} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.GrantSearchKey} value */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.setKey = function(value) {
  jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional SearchMethod method = 2;
 * @return {!proto.caos.citadel.auth.api.v1.SearchMethod}
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.getMethod = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.SearchMethod} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.SearchMethod} value */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.setMethod = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional string value = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.GrantSearchQuery.prototype.setValue = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.GrantSearchResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    offset: jspb.Message.getFieldWithDefault(msg, 1, 0),
    limit: jspb.Message.getFieldWithDefault(msg, 2, 0),
    totalResult: jspb.Message.getFieldWithDefault(msg, 3, 0),
    resultList: jspb.Message.toObjectList(msg.getResultList(),
    proto.caos.citadel.auth.api.v1.Grant.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchResponse}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.GrantSearchResponse;
  return proto.caos.citadel.auth.api.v1.GrantSearchResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.GrantSearchResponse}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setOffset(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setLimit(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setTotalResult(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.Grant;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.Grant.deserializeBinaryFromReader);
      msg.addResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.GrantSearchResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOffset();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getLimit();
  if (f !== 0) {
    writer.writeUint64(
      2,
      f
    );
  }
  f = message.getTotalResult();
  if (f !== 0) {
    writer.writeUint64(
      3,
      f
    );
  }
  f = message.getResultList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.Grant.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 offset = 1;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.getOffset = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.setOffset = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint64 limit = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.getLimit = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.setLimit = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint64 total_result = 3;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.getTotalResult = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.setTotalResult = function(value) {
  jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * repeated Grant result = 4;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.Grant>}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.getResultList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.Grant>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.Grant, 4));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.Grant>} value */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.setResultList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.Grant=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.Grant}
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.addResult = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.caos.citadel.auth.api.v1.Grant, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.GrantSearchResponse.prototype.clearResultList = function() {
  this.setResultList([]);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.Grant.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.Grant.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.Grant} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Grant.toObject = function(includeInstance, msg) {
  var f, obj = {
    orgid: jspb.Message.getFieldWithDefault(msg, 1, ""),
    projectid: jspb.Message.getFieldWithDefault(msg, 2, ""),
    userid: jspb.Message.getFieldWithDefault(msg, 3, ""),
    rolesList: jspb.Message.getRepeatedField(msg, 4),
    orgname: jspb.Message.getFieldWithDefault(msg, 5, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.Grant}
 */
proto.caos.citadel.auth.api.v1.Grant.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.Grant;
  return proto.caos.citadel.auth.api.v1.Grant.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.Grant} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.Grant}
 */
proto.caos.citadel.auth.api.v1.Grant.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setOrgid(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setProjectid(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setUserid(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.addRoles(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setOrgname(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.Grant.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.Grant} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Grant.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOrgid();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getProjectid();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getUserid();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
  f = message.getRolesList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      4,
      f
    );
  }
  f = message.getOrgname();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
};


/**
 * optional string OrgId = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.getOrgid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Grant.prototype.setOrgid = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string ProjectId = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.getProjectid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Grant.prototype.setProjectid = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string UserId = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.getUserid = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Grant.prototype.setUserid = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * repeated string Roles = 4;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.getRolesList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 4));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.Grant.prototype.setRolesList = function(value) {
  jspb.Message.setField(this, 4, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.addRoles = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 4, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.clearRolesList = function() {
  this.setRolesList([]);
};


/**
 * optional string OrgName = 5;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Grant.prototype.getOrgname = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Grant.prototype.setOrgname = function(value) {
  jspb.Message.setProto3StringField(this, 5, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.repeatedFields_ = [5];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    offset: jspb.Message.getFieldWithDefault(msg, 1, 0),
    limit: jspb.Message.getFieldWithDefault(msg, 2, 0),
    asc: jspb.Message.getFieldWithDefault(msg, 4, false),
    queriesList: jspb.Message.toObjectList(msg.getQueriesList(),
    proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest;
  return proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setOffset(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setLimit(value);
      break;
    case 4:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setAsc(value);
      break;
    case 5:
      var value = new proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.deserializeBinaryFromReader);
      msg.addQueries(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOffset();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getLimit();
  if (f !== 0) {
    writer.writeUint64(
      2,
      f
    );
  }
  f = message.getAsc();
  if (f) {
    writer.writeBool(
      4,
      f
    );
  }
  f = message.getQueriesList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      5,
      f,
      proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 offset = 1;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.getOffset = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.setOffset = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint64 limit = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.getLimit = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.setLimit = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional bool asc = 4;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.getAsc = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 4, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.setAsc = function(value) {
  jspb.Message.setProto3BooleanField(this, 4, value);
};


/**
 * repeated MyProjectOrgSearchQuery queries = 5;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery>}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.getQueriesList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery, 5));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery>} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.setQueriesList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 5, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.addQueries = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 5, opt_value, proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest.prototype.clearQueriesList = function() {
  this.setQueriesList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.toObject = function(includeInstance, msg) {
  var f, obj = {
    key: jspb.Message.getFieldWithDefault(msg, 1, 0),
    method: jspb.Message.getFieldWithDefault(msg, 2, 0),
    value: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery;
  return proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchKey} */ (reader.readEnum());
      msg.setKey(value);
      break;
    case 2:
      var value = /** @type {!proto.caos.citadel.auth.api.v1.SearchMethod} */ (reader.readEnum());
      msg.setMethod(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setValue(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getKey();
  if (f !== 0.0) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = message.getMethod();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getValue();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional MyProjectOrgSearchKey key = 1;
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchKey}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.getKey = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchKey} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchKey} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.setKey = function(value) {
  jspb.Message.setProto3EnumField(this, 1, value);
};


/**
 * optional SearchMethod method = 2;
 * @return {!proto.caos.citadel.auth.api.v1.SearchMethod}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.getMethod = function() {
  return /** @type {!proto.caos.citadel.auth.api.v1.SearchMethod} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {!proto.caos.citadel.auth.api.v1.SearchMethod} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.setMethod = function(value) {
  jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional string value = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.getValue = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchQuery.prototype.setValue = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.repeatedFields_ = [4];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    offset: jspb.Message.getFieldWithDefault(msg, 1, 0),
    limit: jspb.Message.getFieldWithDefault(msg, 2, 0),
    totalResult: jspb.Message.getFieldWithDefault(msg, 3, 0),
    resultList: jspb.Message.toObjectList(msg.getResultList(),
    proto.caos.citadel.auth.api.v1.Org.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse;
  return proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setOffset(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setLimit(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readUint64());
      msg.setTotalResult(value);
      break;
    case 4:
      var value = new proto.caos.citadel.auth.api.v1.Org;
      reader.readMessage(value,proto.caos.citadel.auth.api.v1.Org.deserializeBinaryFromReader);
      msg.addResult(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getOffset();
  if (f !== 0) {
    writer.writeUint64(
      1,
      f
    );
  }
  f = message.getLimit();
  if (f !== 0) {
    writer.writeUint64(
      2,
      f
    );
  }
  f = message.getTotalResult();
  if (f !== 0) {
    writer.writeUint64(
      3,
      f
    );
  }
  f = message.getResultList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      4,
      f,
      proto.caos.citadel.auth.api.v1.Org.serializeBinaryToWriter
    );
  }
};


/**
 * optional uint64 offset = 1;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.getOffset = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.setOffset = function(value) {
  jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional uint64 limit = 2;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.getLimit = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.setLimit = function(value) {
  jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * optional uint64 total_result = 3;
 * @return {number}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.getTotalResult = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/** @param {number} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.setTotalResult = function(value) {
  jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * repeated Org result = 4;
 * @return {!Array<!proto.caos.citadel.auth.api.v1.Org>}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.getResultList = function() {
  return /** @type{!Array<!proto.caos.citadel.auth.api.v1.Org>} */ (
    jspb.Message.getRepeatedWrapperField(this, proto.caos.citadel.auth.api.v1.Org, 4));
};


/** @param {!Array<!proto.caos.citadel.auth.api.v1.Org>} value */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.setResultList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 4, value);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.Org=} opt_value
 * @param {number=} opt_index
 * @return {!proto.caos.citadel.auth.api.v1.Org}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.addResult = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 4, opt_value, proto.caos.citadel.auth.api.v1.Org, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.prototype.clearResultList = function() {
  this.setResultList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.IsAdminResponse.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.IsAdminResponse} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse.toObject = function(includeInstance, msg) {
  var f, obj = {
    isAdmin: jspb.Message.getFieldWithDefault(msg, 1, false)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.IsAdminResponse}
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.IsAdminResponse;
  return proto.caos.citadel.auth.api.v1.IsAdminResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.IsAdminResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.IsAdminResponse}
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {boolean} */ (reader.readBool());
      msg.setIsAdmin(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.IsAdminResponse.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.IsAdminResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getIsAdmin();
  if (f) {
    writer.writeBool(
      1,
      f
    );
  }
};


/**
 * optional bool is_admin = 1;
 * Note that Boolean fields may be set to 0/1 when serialized from a Java server.
 * You should avoid comparisons like {@code val === true/false} in those cases.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.IsAdminResponse.prototype.getIsAdmin = function() {
  return /** @type {boolean} */ (jspb.Message.getFieldWithDefault(this, 1, false));
};


/** @param {boolean} value */
proto.caos.citadel.auth.api.v1.IsAdminResponse.prototype.setIsAdmin = function(value) {
  jspb.Message.setProto3BooleanField(this, 1, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.Org.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.Org.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.Org} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Org.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    name: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.Org}
 */
proto.caos.citadel.auth.api.v1.Org.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.Org;
  return proto.caos.citadel.auth.api.v1.Org.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.Org} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.Org}
 */
proto.caos.citadel.auth.api.v1.Org.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.Org.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.Org.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.Org} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Org.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string Id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Org.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Org.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string Name = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Org.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Org.prototype.setName = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.CreateTokenRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.CreateTokenRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    agentId: jspb.Message.getFieldWithDefault(msg, 1, ""),
    authSessionId: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.CreateTokenRequest}
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.CreateTokenRequest;
  return proto.caos.citadel.auth.api.v1.CreateTokenRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.CreateTokenRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.CreateTokenRequest}
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setAgentId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setAuthSessionId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.CreateTokenRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.CreateTokenRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getAgentId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getAuthSessionId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional string agent_id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.prototype.getAgentId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.prototype.setAgentId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string auth_session_id = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.prototype.getAuthSessionId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.CreateTokenRequest.prototype.setAuthSessionId = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.Token.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.Token.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.Token} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Token.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    expiration: (f = msg.getExpiration()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.Token}
 */
proto.caos.citadel.auth.api.v1.Token.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.Token;
  return proto.caos.citadel.auth.api.v1.Token.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.Token} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.Token}
 */
proto.caos.citadel.auth.api.v1.Token.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 8:
      var value = new google_protobuf_timestamp_pb.Timestamp;
      reader.readMessage(value,google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
      msg.setExpiration(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.Token.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.Token.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.Token} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.Token.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getExpiration();
  if (f != null) {
    writer.writeMessage(
      8,
      f,
      google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.Token.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.Token.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional google.protobuf.Timestamp expiration = 8;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.caos.citadel.auth.api.v1.Token.prototype.getExpiration = function() {
  return /** @type{?proto.google.protobuf.Timestamp} */ (
    jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 8));
};


/** @param {?proto.google.protobuf.Timestamp|undefined} value */
proto.caos.citadel.auth.api.v1.Token.prototype.setExpiration = function(value) {
  jspb.Message.setWrapperField(this, 8, value);
};


/**
 * Clears the message field making it undefined.
 */
proto.caos.citadel.auth.api.v1.Token.prototype.clearExpiration = function() {
  this.setExpiration(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.caos.citadel.auth.api.v1.Token.prototype.hasExpiration = function() {
  return jspb.Message.getField(this, 8) != null;
};



/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.caos.citadel.auth.api.v1.MyPermissions.repeatedFields_ = [1];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.MyPermissions.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.MyPermissions.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.MyPermissions} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyPermissions.toObject = function(includeInstance, msg) {
  var f, obj = {
    permissionsList: jspb.Message.getRepeatedField(msg, 1)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.MyPermissions}
 */
proto.caos.citadel.auth.api.v1.MyPermissions.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.MyPermissions;
  return proto.caos.citadel.auth.api.v1.MyPermissions.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.MyPermissions} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.MyPermissions}
 */
proto.caos.citadel.auth.api.v1.MyPermissions.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.addPermissions(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.MyPermissions.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.MyPermissions.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.MyPermissions} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.MyPermissions.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getPermissionsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      1,
      f
    );
  }
};


/**
 * repeated string permissions = 1;
 * @return {!Array<string>}
 */
proto.caos.citadel.auth.api.v1.MyPermissions.prototype.getPermissionsList = function() {
  return /** @type {!Array<string>} */ (jspb.Message.getRepeatedField(this, 1));
};


/** @param {!Array<string>} value */
proto.caos.citadel.auth.api.v1.MyPermissions.prototype.setPermissionsList = function(value) {
  jspb.Message.setField(this, 1, value || []);
};


/**
 * @param {string} value
 * @param {number=} opt_index
 */
proto.caos.citadel.auth.api.v1.MyPermissions.prototype.addPermissions = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 1, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 */
proto.caos.citadel.auth.api.v1.MyPermissions.prototype.clearPermissionsList = function() {
  this.setPermissionsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.toObject = function(opt_includeInstance) {
  return proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, ""),
    code: jspb.Message.getFieldWithDefault(msg, 2, ""),
    password: jspb.Message.getFieldWithDefault(msg, 3, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.caos.citadel.auth.api.v1.VerifyUserInitRequest;
  return proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest}
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setCode(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.setPassword(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f.length > 0) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getCode();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
  f = message.getPassword();
  if (f.length > 0) {
    writer.writeString(
      3,
      f
    );
  }
};


/**
 * optional string id = 1;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.getId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.setId = function(value) {
  jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string code = 2;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.getCode = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.setCode = function(value) {
  jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string password = 3;
 * @return {string}
 */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.getPassword = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/** @param {string} value */
proto.caos.citadel.auth.api.v1.VerifyUserInitRequest.prototype.setPassword = function(value) {
  jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.UserAgentState = {
  NO_STATE: 0,
  ACTIVE_SESSION: 1,
  TERMINATED_SESSION: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.UserSessionState = {
  USER_SESSION_STATE_UNKNOWN: 0,
  USER_SESSION_STATE_ACTIVE: 1,
  USER_SESSION_STATE_TERMINATED: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.NextStepType = {
  NEXT_STEP_UNSPECIFIED: 0,
  NEXT_STEP_LOGIN: 1,
  NEXT_STEP_PASSWORD: 2,
  NEXT_STEP_CHANGE_PASSWORD: 3,
  NEXT_STEP_MFA_PROMPT: 4,
  NEXT_STEP_MFA_INIT_CHOICE: 5,
  NEXT_STEP_MFA_INIT_CREATE: 6,
  NEXT_STEP_MFA_INIT_VERIFY: 7,
  NEXT_STEP_MFA_INIT_DONE: 8,
  NEXT_STEP_MFA_VERIFY: 9,
  NEXT_STEP_MFA_VERIFY_ASYNC: 10,
  NEXT_STEP_VERIFY_EMAIL: 11,
  NEXT_STEP_REDIRECT_TO_CALLBACK: 12,
  NEXT_STEP_INIT_PASSWORD: 13,
  NEXT_STEP_CHOOSE_USER: 14
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.AuthSessionType = {
  TYPE_UNKNOWN: 0,
  TYPE_OIDC: 1,
  TYPE_SAML: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.Prompt = {
  NO_PROMPT: 0,
  PROMPT_NONE: 1,
  PROMPT_LOGIN: 2,
  PROMPT_CONSENT: 3,
  PROMPT_SELECT_ACCOUNT: 4
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.OIDCResponseType = {
  CODE: 0,
  ID_TOKEN: 1,
  ID_TOKEN_TOKEN: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.CodeChallengeMethod = {
  PLAIN: 0,
  S256: 1
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.UserState = {
  NONE: 0,
  ACTIVE: 1,
  INACTIVE: 2,
  DELETED: 3,
  LOCKED: 4,
  SUSPEND: 5,
  INITIAL: 6
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.Gender = {
  UNKNOWN_GENDER: 0,
  FEMALE: 1,
  MALE: 2,
  DIVERSE: 3
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.NotificationType = {
  EMAIL: 0,
  SMS: 1
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.MfaType = {
  NO_MFA: 0,
  MFA_SMS: 1,
  MFA_OTP: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.MFAState = {
  MFASTATE_NO: 0,
  NOT_READY: 1,
  READY: 2,
  REMOVED: 3
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.AppState = {
  NONE_APP: 0,
  ACTIVE_APP: 1,
  INACTIVE_APP: 2,
  DELETED_APP: 3
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.OIDCGrantType = {
  AUTHORIZATION_CODE: 0,
  GRANT_TYPE_NONE: 1,
  REFRESH_TOKEN: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.OIDCApplicationType = {
  WEB: 0,
  USER_AGENT: 1,
  NATIVE: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.OIDCAuthMethodType = {
  AUTH_TYPE_BASIC: 0,
  AUTH_TYPE_POST: 1,
  AUTH_TYPE_NONE: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.ApplicationSearchKey = {
  UNKNOWN: 0,
  APP_TYPE: 1,
  STATE: 2,
  CLIENT_ID: 3,
  APP_NAME: 4,
  PROJECT_ID: 5
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.SearchMethod = {
  EQUALS: 0,
  STARTS_WITH: 1,
  CONTAINS: 2
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.GrantSearchKey = {
  GRANTSEARCHKEY_UNKNOWN: 0,
  GRANTSEARCHKEY_ORG_ID: 1,
  GRANTSEARCHKEY_PROJECT_ID: 2,
  GRANTSEARCHKEY_USER_ID: 3
};

/**
 * @enum {number}
 */
proto.caos.citadel.auth.api.v1.MyProjectOrgSearchKey = {
  MYPROJECTORGKEY_UNKNOWN: 0,
  MYPROJECTORGKEY_ORG_NAME: 1
};

goog.object.extend(exports, proto.caos.citadel.auth.api.v1);
