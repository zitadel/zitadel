/**
 * @fileoverview gRPC-Web generated client stub for caos.zitadel.auth.api.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');


var google_api_annotations_pb = require('./google/api/annotations_pb.js')

var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')

var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var validate_validate_pb = require('./validate/validate_pb.js')

var protoc$gen$swagger_options_annotations_pb = require('./protoc-gen-swagger/options/annotations_pb.js')

var authoption_options_pb = require('./authoption/options_pb.js')
const proto = {};
proto.caos = {};
proto.caos.zitadel = {};
proto.caos.zitadel.auth = {};
proto.caos.zitadel.auth.api = {};
proto.caos.zitadel.auth.api.v1 = require('./auth_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_Healthz = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/Healthz',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_Healthz = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.healthz =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/Healthz',
      request,
      metadata || {},
      methodDescriptor_AuthService_Healthz,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.healthz =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/Healthz',
      request,
      metadata || {},
      methodDescriptor_AuthService_Healthz);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_Ready = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/Ready',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_Ready = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.ready =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/Ready',
      request,
      metadata || {},
      methodDescriptor_AuthService_Ready,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.ready =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/Ready',
      request,
      metadata || {},
      methodDescriptor_AuthService_Ready);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Struct>}
 */
const methodDescriptor_AuthService_Validate = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/Validate',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  google_protobuf_struct_pb.Struct,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_struct_pb.Struct.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Struct>}
 */
const methodInfo_AuthService_Validate = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_struct_pb.Struct,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_struct_pb.Struct.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Struct)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Struct>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.validate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/Validate',
      request,
      metadata || {},
      methodDescriptor_AuthService_Validate,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Struct>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.validate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/Validate',
      request,
      metadata || {},
      methodDescriptor_AuthService_Validate);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserSessionViews>}
 */
const methodDescriptor_AuthService_GetMyUserSessions = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/GetMyUserSessions',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.UserSessionViews,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserSessionViews.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserSessionViews>}
 */
const methodInfo_AuthService_GetMyUserSessions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserSessionViews,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserSessionViews.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserSessionViews)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserSessionViews>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.getMyUserSessions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserSessions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserSessions,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserSessionViews>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserSessions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserSessions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserSessions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserProfileView>}
 */
const methodDescriptor_AuthService_GetMyUserProfile = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/GetMyUserProfile',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.UserProfileView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserProfileView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserProfileView>}
 */
const methodInfo_AuthService_GetMyUserProfile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserProfileView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserProfileView.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserProfileView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserProfileView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.getMyUserProfile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserProfile',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserProfile,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserProfileView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserProfile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserProfile',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserProfile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserProfileRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserProfile>}
 */
const methodDescriptor_AuthService_UpdateMyUserProfile = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/UpdateMyUserProfile',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.UpdateUserProfileRequest,
  proto.caos.zitadel.auth.api.v1.UserProfile,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserProfileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserProfile.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserProfileRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserProfile>}
 */
const methodInfo_AuthService_UpdateMyUserProfile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserProfile,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserProfileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserProfile.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserProfileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserProfile)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserProfile>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.updateMyUserProfile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/UpdateMyUserProfile',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserProfile,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserProfileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserProfile>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.updateMyUserProfile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/UpdateMyUserProfile',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserProfile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserEmailView>}
 */
const methodDescriptor_AuthService_GetMyUserEmail = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/GetMyUserEmail',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.UserEmailView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserEmailView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserEmailView>}
 */
const methodInfo_AuthService_GetMyUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserEmailView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserEmailView.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserEmailView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserEmailView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.getMyUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserEmail,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserEmailView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserEmailRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserEmail>}
 */
const methodDescriptor_AuthService_ChangeMyUserEmail = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/ChangeMyUserEmail',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.UpdateUserEmailRequest,
  proto.caos.zitadel.auth.api.v1.UserEmail,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserEmail.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserEmailRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserEmail>}
 */
const methodInfo_AuthService_ChangeMyUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserEmail,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserEmail.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserEmail)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserEmail>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.changeMyUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ChangeMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserEmail,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserEmail>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.changeMyUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ChangeMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.VerifyMyUserEmailRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_VerifyMyUserEmail = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/VerifyMyUserEmail',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.VerifyMyUserEmailRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.VerifyMyUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.VerifyMyUserEmailRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_VerifyMyUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.VerifyMyUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.VerifyMyUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.verifyMyUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/VerifyMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMyUserEmail,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.VerifyMyUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyMyUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/VerifyMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMyUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_ResendMyEmailVerificationMail = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/ResendMyEmailVerificationMail',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_ResendMyEmailVerificationMail = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.resendMyEmailVerificationMail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ResendMyEmailVerificationMail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendMyEmailVerificationMail,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.resendMyEmailVerificationMail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ResendMyEmailVerificationMail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendMyEmailVerificationMail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserPhoneView>}
 */
const methodDescriptor_AuthService_GetMyUserPhone = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/GetMyUserPhone',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.UserPhoneView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserPhoneView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserPhoneView>}
 */
const methodInfo_AuthService_GetMyUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserPhoneView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserPhoneView.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserPhoneView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserPhoneView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.getMyUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserPhone,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserPhoneView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserPhone);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserPhoneRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserPhone>}
 */
const methodDescriptor_AuthService_ChangeMyUserPhone = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/ChangeMyUserPhone',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.UpdateUserPhoneRequest,
  proto.caos.zitadel.auth.api.v1.UserPhone,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserPhone.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserPhoneRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserPhone>}
 */
const methodInfo_AuthService_ChangeMyUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserPhone,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserPhone.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserPhone)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserPhone>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.changeMyUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ChangeMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserPhone,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserPhone>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.changeMyUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ChangeMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserPhone);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.VerifyUserPhoneRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_VerifyMyUserPhone = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/VerifyMyUserPhone',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.VerifyUserPhoneRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.VerifyUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.VerifyUserPhoneRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_VerifyMyUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.VerifyUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.VerifyUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.verifyMyUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/VerifyMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMyUserPhone,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.VerifyUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyMyUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/VerifyMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMyUserPhone);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_ResendMyPhoneVerificationCode = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/ResendMyPhoneVerificationCode',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_ResendMyPhoneVerificationCode = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.resendMyPhoneVerificationCode =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ResendMyPhoneVerificationCode',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendMyPhoneVerificationCode,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.resendMyPhoneVerificationCode =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ResendMyPhoneVerificationCode',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendMyPhoneVerificationCode);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserAddressView>}
 */
const methodDescriptor_AuthService_GetMyUserAddress = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/GetMyUserAddress',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.UserAddressView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserAddressView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.UserAddressView>}
 */
const methodInfo_AuthService_GetMyUserAddress = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserAddressView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserAddressView.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserAddressView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserAddressView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.getMyUserAddress =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserAddress',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserAddress,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserAddressView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserAddress =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyUserAddress',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserAddress);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserAddressRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserAddress>}
 */
const methodDescriptor_AuthService_UpdateMyUserAddress = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/UpdateMyUserAddress',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.UpdateUserAddressRequest,
  proto.caos.zitadel.auth.api.v1.UserAddress,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserAddressRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserAddress.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.UpdateUserAddressRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserAddress>}
 */
const methodInfo_AuthService_UpdateMyUserAddress = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserAddress,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserAddressRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserAddress.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserAddressRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserAddress)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserAddress>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.updateMyUserAddress =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/UpdateMyUserAddress',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserAddress,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UpdateUserAddressRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserAddress>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.updateMyUserAddress =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/UpdateMyUserAddress',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserAddress);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.MultiFactors>}
 */
const methodDescriptor_AuthService_GetMyMfas = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/GetMyMfas',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.MultiFactors,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MultiFactors.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.MultiFactors>}
 */
const methodInfo_AuthService_GetMyMfas = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.MultiFactors,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MultiFactors.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.MultiFactors)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.MultiFactors>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.getMyMfas =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyMfas',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyMfas,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.MultiFactors>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyMfas =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyMfas',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyMfas);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.PasswordChange,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_ChangeMyPassword = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/ChangeMyPassword',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.PasswordChange,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.PasswordChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.PasswordChange,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_ChangeMyPassword = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.PasswordChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.PasswordChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.changeMyPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ChangeMyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyPassword,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.PasswordChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.changeMyPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/ChangeMyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.MfaOtpResponse>}
 */
const methodDescriptor_AuthService_AddMfaOTP = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/AddMfaOTP',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.MfaOtpResponse,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MfaOtpResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.MfaOtpResponse>}
 */
const methodInfo_AuthService_AddMfaOTP = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.MfaOtpResponse,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MfaOtpResponse.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.MfaOtpResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.MfaOtpResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.addMfaOTP =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/AddMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_AddMfaOTP,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.MfaOtpResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.addMfaOTP =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/AddMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_AddMfaOTP);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.VerifyMfaOtp,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_VerifyMfaOTP = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/VerifyMfaOTP',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.VerifyMfaOtp,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.VerifyMfaOtp} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.VerifyMfaOtp,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_VerifyMfaOTP = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.VerifyMfaOtp} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.VerifyMfaOtp} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.verifyMfaOTP =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/VerifyMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMfaOTP,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.VerifyMfaOtp} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyMfaOTP =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/VerifyMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMfaOTP);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_RemoveMfaOTP = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/RemoveMfaOTP',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_RemoveMfaOTP = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.removeMfaOTP =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/RemoveMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_RemoveMfaOTP,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.removeMfaOTP =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/RemoveMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_RemoveMfaOTP);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.UserGrantSearchRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse>}
 */
const methodDescriptor_AuthService_SearchMyUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/SearchMyUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.UserGrantSearchRequest,
  proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.UserGrantSearchRequest,
 *   !proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse>}
 */
const methodInfo_AuthService_SearchMyUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.UserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.searchMyUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/SearchMyUserGrant',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchMyUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.UserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.UserGrantSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.searchMyUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/SearchMyUserGrant',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchMyUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchRequest,
 *   !proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse>}
 */
const methodDescriptor_AuthService_SearchMyProjectOrgs = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/SearchMyProjectOrgs',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchRequest,
  proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchRequest,
 *   !proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse>}
 */
const methodInfo_AuthService_SearchMyProjectOrgs = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse,
  /**
   * @param {!proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.searchMyProjectOrgs =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/SearchMyProjectOrgs',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchMyProjectOrgs,
      callback);
};


/**
 * @param {!proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.MyProjectOrgSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.searchMyProjectOrgs =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/SearchMyProjectOrgs',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchMyProjectOrgs);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.MyPermissions>}
 */
const methodDescriptor_AuthService_GetMyZitadelPermissions = new grpc.web.MethodDescriptor(
  '/caos.zitadel.auth.api.v1.AuthService/GetMyZitadelPermissions',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.auth.api.v1.MyPermissions,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MyPermissions.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.auth.api.v1.MyPermissions>}
 */
const methodInfo_AuthService_GetMyZitadelPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.auth.api.v1.MyPermissions,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.auth.api.v1.MyPermissions.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.auth.api.v1.MyPermissions)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.auth.api.v1.MyPermissions>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.auth.api.v1.AuthServiceClient.prototype.getMyZitadelPermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyZitadelPermissions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyZitadelPermissions,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.auth.api.v1.MyPermissions>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyZitadelPermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.auth.api.v1.AuthService/GetMyZitadelPermissions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyZitadelPermissions);
};


module.exports = proto.caos.zitadel.auth.api.v1;

