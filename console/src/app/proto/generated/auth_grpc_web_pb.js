/**
 * @fileoverview gRPC-Web generated client stub for caos.citadel.auth.api.v1
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

var protoc$gen$swagger_options_annotations_pb = require('./protoc-gen-swagger/options/annotations_pb.js')

var validate_validate_pb = require('./validate/validate_pb.js')

var authoption_options_pb = require('./authoption/options_pb.js')
const proto = {};
proto.caos = {};
proto.caos.citadel = {};
proto.caos.citadel.auth = {};
proto.caos.citadel.auth.api = {};
proto.caos.citadel.auth.api.v1 = require('./auth_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient =
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
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient =
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
  '/caos.citadel.auth.api.v1.AuthService/Healthz',
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
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.healthz =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/Healthz',
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
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.healthz =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/Healthz',
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
  '/caos.citadel.auth.api.v1.AuthService/Ready',
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
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.ready =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/Ready',
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
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.ready =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/Ready',
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
  '/caos.citadel.auth.api.v1.AuthService/Validate',
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
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.validate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/Validate',
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
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.validate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/Validate',
      request,
      metadata || {},
      methodDescriptor_AuthService_Validate);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserAgentID,
 *   !proto.caos.citadel.auth.api.v1.UserAgent>}
 */
const methodDescriptor_AuthService_GetUserAgent = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetUserAgent',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserAgentID,
  proto.caos.citadel.auth.api.v1.UserAgent,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UserAgentID,
 *   !proto.caos.citadel.auth.api.v1.UserAgent>}
 */
const methodInfo_AuthService_GetUserAgent = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserAgent,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserAgent)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserAgent>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getUserAgent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserAgent',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserAgent,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserAgent>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getUserAgent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserAgent',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserAgent);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserAgentCreation,
 *   !proto.caos.citadel.auth.api.v1.UserAgent>}
 */
const methodDescriptor_AuthService_CreateUserAgent = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/CreateUserAgent',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserAgentCreation,
  proto.caos.citadel.auth.api.v1.UserAgent,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentCreation} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UserAgentCreation,
 *   !proto.caos.citadel.auth.api.v1.UserAgent>}
 */
const methodInfo_AuthService_CreateUserAgent = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserAgent,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentCreation} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentCreation} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserAgent)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserAgent>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.createUserAgent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/CreateUserAgent',
      request,
      metadata || {},
      methodDescriptor_AuthService_CreateUserAgent,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentCreation} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserAgent>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.createUserAgent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/CreateUserAgent',
      request,
      metadata || {},
      methodDescriptor_AuthService_CreateUserAgent);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserAgentID,
 *   !proto.caos.citadel.auth.api.v1.UserAgent>}
 */
const methodDescriptor_AuthService_RevokeUserAgent = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/RevokeUserAgent',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserAgentID,
  proto.caos.citadel.auth.api.v1.UserAgent,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UserAgentID,
 *   !proto.caos.citadel.auth.api.v1.UserAgent>}
 */
const methodInfo_AuthService_RevokeUserAgent = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserAgent,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAgent.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserAgent)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserAgent>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.revokeUserAgent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RevokeUserAgent',
      request,
      metadata || {},
      methodDescriptor_AuthService_RevokeUserAgent,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserAgent>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.revokeUserAgent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RevokeUserAgent',
      request,
      metadata || {},
      methodDescriptor_AuthService_RevokeUserAgent);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.AuthSessionCreation,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodDescriptor_AuthService_CreateAuthSession = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/CreateAuthSession',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.AuthSessionCreation,
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.AuthSessionCreation} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.AuthSessionCreation,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodInfo_AuthService_CreateAuthSession = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.AuthSessionCreation} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionCreation} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.AuthSessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.createAuthSession =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/CreateAuthSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_CreateAuthSession,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionCreation} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.createAuthSession =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/CreateAuthSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_CreateAuthSession);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.AuthSessionID,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodDescriptor_AuthService_GetAuthSession = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetAuthSession',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.AuthSessionID,
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.AuthSessionID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.AuthSessionID,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodInfo_AuthService_GetAuthSession = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.AuthSessionID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.AuthSessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getAuthSession =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetAuthSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetAuthSession,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.AuthSessionID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getAuthSession =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetAuthSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetAuthSession);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.TokenID,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionView>}
 */
const methodDescriptor_AuthService_GetAuthSessionByTokenID = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetAuthSessionByTokenID',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.TokenID,
  proto.caos.citadel.auth.api.v1.AuthSessionView,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.TokenID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.TokenID,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionView>}
 */
const methodInfo_AuthService_GetAuthSessionByTokenID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.AuthSessionView,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.TokenID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionView.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.TokenID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.AuthSessionView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.AuthSessionView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getAuthSessionByTokenID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetAuthSessionByTokenID',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetAuthSessionByTokenID,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.TokenID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.AuthSessionView>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getAuthSessionByTokenID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetAuthSessionByTokenID',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetAuthSessionByTokenID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.SelectUserRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodDescriptor_AuthService_SelectUser = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/SelectUser',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.SelectUserRequest,
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.SelectUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.SelectUserRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodInfo_AuthService_SelectUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.SelectUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.SelectUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.AuthSessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.selectUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SelectUser',
      request,
      metadata || {},
      methodDescriptor_AuthService_SelectUser,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.SelectUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.selectUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SelectUser',
      request,
      metadata || {},
      methodDescriptor_AuthService_SelectUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyUserRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodDescriptor_AuthService_VerifyUser = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyUser',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyUserRequest,
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.VerifyUserRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodInfo_AuthService_VerifyUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.AuthSessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyUser',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyUser,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyUser',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyPasswordRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodDescriptor_AuthService_VerifyPassword = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyPassword',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyPasswordRequest,
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.VerifyPasswordRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodInfo_AuthService_VerifyPassword = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.AuthSessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyPassword,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyMfaRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodDescriptor_AuthService_VerifyMfa = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyMfa',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyMfaRequest,
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.VerifyMfaRequest,
 *   !proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 */
const methodInfo_AuthService_VerifyMfa = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.AuthSessionResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.AuthSessionResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.AuthSessionResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyMfa =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMfa',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMfa,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.AuthSessionResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyMfa =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMfa',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMfa);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserAgentID,
 *   !proto.caos.citadel.auth.api.v1.UserSessions>}
 */
const methodDescriptor_AuthService_GetUserAgentSessions = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetUserAgentSessions',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserAgentID,
  proto.caos.citadel.auth.api.v1.UserSessions,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserSessions.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UserAgentID,
 *   !proto.caos.citadel.auth.api.v1.UserSessions>}
 */
const methodInfo_AuthService_GetUserAgentSessions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserSessions,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserSessions.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserSessions)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserSessions>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getUserAgentSessions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserAgentSessions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserAgentSessions,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserAgentID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserSessions>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getUserAgentSessions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserAgentSessions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserAgentSessions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserSessionID,
 *   !proto.caos.citadel.auth.api.v1.UserSession>}
 */
const methodDescriptor_AuthService_GetUserSession = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetUserSession',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserSessionID,
  proto.caos.citadel.auth.api.v1.UserSession,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserSession.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UserSessionID,
 *   !proto.caos.citadel.auth.api.v1.UserSession>}
 */
const methodInfo_AuthService_GetUserSession = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserSession,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserSession.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserSession)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserSession>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getUserSession =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserSession,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserSession>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getUserSession =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserSession);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserSessionViews>}
 */
const methodDescriptor_AuthService_GetMyUserSessions = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetMyUserSessions',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.UserSessionViews,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserSessionViews.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserSessionViews>}
 */
const methodInfo_AuthService_GetMyUserSessions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserSessionViews,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserSessionViews.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserSessionViews)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserSessionViews>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getMyUserSessions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserSessions',
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
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserSessionViews>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserSessions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserSessions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserSessions);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserSessionID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_TerminateUserSession = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/TerminateUserSession',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserSessionID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request
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
 *   !proto.caos.citadel.auth.api.v1.UserSessionID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_TerminateUserSession = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.terminateUserSession =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/TerminateUserSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_TerminateUserSession,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserSessionID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.terminateUserSession =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/TerminateUserSession',
      request,
      metadata || {},
      methodDescriptor_AuthService_TerminateUserSession);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.CreateTokenRequest,
 *   !proto.caos.citadel.auth.api.v1.Token>}
 */
const methodDescriptor_AuthService_CreateToken = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/CreateToken',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.CreateTokenRequest,
  proto.caos.citadel.auth.api.v1.Token,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.CreateTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.Token.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.CreateTokenRequest,
 *   !proto.caos.citadel.auth.api.v1.Token>}
 */
const methodInfo_AuthService_CreateToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.Token,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.CreateTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.Token.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.CreateTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.Token)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.Token>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.createToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/CreateToken',
      request,
      metadata || {},
      methodDescriptor_AuthService_CreateToken,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.CreateTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.Token>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.createToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/CreateToken',
      request,
      metadata || {},
      methodDescriptor_AuthService_CreateToken);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UniqueUserRequest,
 *   !proto.caos.citadel.auth.api.v1.UniqueUserResponse>}
 */
const methodDescriptor_AuthService_IsUserUnique = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/IsUserUnique',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UniqueUserRequest,
  proto.caos.citadel.auth.api.v1.UniqueUserResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UniqueUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UniqueUserResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UniqueUserRequest,
 *   !proto.caos.citadel.auth.api.v1.UniqueUserResponse>}
 */
const methodInfo_AuthService_IsUserUnique = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UniqueUserResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UniqueUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UniqueUserResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UniqueUserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UniqueUserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.isUserUnique =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/IsUserUnique',
      request,
      metadata || {},
      methodDescriptor_AuthService_IsUserUnique,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UniqueUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UniqueUserResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.isUserUnique =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/IsUserUnique',
      request,
      metadata || {},
      methodDescriptor_AuthService_IsUserUnique);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.RegisterUserRequest,
 *   !proto.caos.citadel.auth.api.v1.User>}
 */
const methodDescriptor_AuthService_RegisterUser = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/RegisterUser',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.RegisterUserRequest,
  proto.caos.citadel.auth.api.v1.User,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.RegisterUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.RegisterUserRequest,
 *   !proto.caos.citadel.auth.api.v1.User>}
 */
const methodInfo_AuthService_RegisterUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.User,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.RegisterUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.registerUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RegisterUser',
      request,
      metadata || {},
      methodDescriptor_AuthService_RegisterUser,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.registerUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RegisterUser',
      request,
      metadata || {},
      methodDescriptor_AuthService_RegisterUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest,
 *   !proto.caos.citadel.auth.api.v1.User>}
 */
const methodDescriptor_AuthService_RegisterUserWithExternal = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/RegisterUserWithExternal',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest,
  proto.caos.citadel.auth.api.v1.User,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest,
 *   !proto.caos.citadel.auth.api.v1.User>}
 */
const methodInfo_AuthService_RegisterUserWithExternal = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.User,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.registerUserWithExternal =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RegisterUserWithExternal',
      request,
      metadata || {},
      methodDescriptor_AuthService_RegisterUserWithExternal,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.RegisterUserExternalIDPRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.registerUserWithExternal =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RegisterUserWithExternal',
      request,
      metadata || {},
      methodDescriptor_AuthService_RegisterUserWithExternal);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserProfile>}
 */
const methodDescriptor_AuthService_GetMyUserProfile = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetMyUserProfile',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.UserProfile,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserProfile.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserProfile>}
 */
const methodInfo_AuthService_GetMyUserProfile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserProfile,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserProfile.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserProfile)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserProfile>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getMyUserProfile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserProfile',
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
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserProfile>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserProfile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserProfile',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserProfile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest,
 *   !proto.caos.citadel.auth.api.v1.UserProfile>}
 */
const methodDescriptor_AuthService_UpdateMyUserProfile = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/UpdateMyUserProfile',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest,
  proto.caos.citadel.auth.api.v1.UserProfile,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserProfile.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest,
 *   !proto.caos.citadel.auth.api.v1.UserProfile>}
 */
const methodInfo_AuthService_UpdateMyUserProfile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserProfile,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserProfile.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserProfile)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserProfile>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.updateMyUserProfile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/UpdateMyUserProfile',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserProfile,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserProfileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserProfile>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.updateMyUserProfile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/UpdateMyUserProfile',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserProfile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserEmail>}
 */
const methodDescriptor_AuthService_GetMyUserEmail = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetMyUserEmail',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.UserEmail,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserEmail.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserEmail>}
 */
const methodInfo_AuthService_GetMyUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserEmail,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserEmail.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserEmail)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserEmail>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getMyUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserEmail',
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
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserEmail>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest,
 *   !proto.caos.citadel.auth.api.v1.UserEmail>}
 */
const methodDescriptor_AuthService_ChangeMyUserEmail = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/ChangeMyUserEmail',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest,
  proto.caos.citadel.auth.api.v1.UserEmail,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserEmail.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest,
 *   !proto.caos.citadel.auth.api.v1.UserEmail>}
 */
const methodInfo_AuthService_ChangeMyUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserEmail,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserEmail.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserEmail)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserEmail>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.changeMyUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ChangeMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserEmail,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserEmail>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.changeMyUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ChangeMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_VerifyMyUserEmail = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyMyUserEmail',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest} request
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
 *   !proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_VerifyMyUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyMyUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMyUserEmail,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMyUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyMyUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMyUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_VerifyUserEmail = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyUserEmail',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest} request
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
 *   !proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_VerifyUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyUserEmail,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyUserEmail',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_ResendMyEmailVerificationMail = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/ResendMyEmailVerificationMail',
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
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.resendMyEmailVerificationMail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendMyEmailVerificationMail',
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
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.resendMyEmailVerificationMail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendMyEmailVerificationMail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendMyEmailVerificationMail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_ResendEmailVerificationMail = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/ResendEmailVerificationMail',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserID} request
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
 *   !proto.caos.citadel.auth.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_ResendEmailVerificationMail = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.resendEmailVerificationMail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendEmailVerificationMail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendEmailVerificationMail,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.resendEmailVerificationMail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendEmailVerificationMail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendEmailVerificationMail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserPhone>}
 */
const methodDescriptor_AuthService_GetMyUserPhone = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetMyUserPhone',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.UserPhone,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserPhone.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserPhone>}
 */
const methodInfo_AuthService_GetMyUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserPhone,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserPhone.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserPhone)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserPhone>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getMyUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserPhone',
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
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserPhone>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserPhone);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest,
 *   !proto.caos.citadel.auth.api.v1.UserPhone>}
 */
const methodDescriptor_AuthService_ChangeMyUserPhone = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/ChangeMyUserPhone',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest,
  proto.caos.citadel.auth.api.v1.UserPhone,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserPhone.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest,
 *   !proto.caos.citadel.auth.api.v1.UserPhone>}
 */
const methodInfo_AuthService_ChangeMyUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserPhone,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserPhone.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserPhone)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserPhone>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.changeMyUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ChangeMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserPhone,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserPhone>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.changeMyUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ChangeMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyUserPhone);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_VerifyMyUserPhone = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyMyUserPhone',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest} request
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
 *   !proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_VerifyMyUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyMyUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMyUserPhone',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMyUserPhone,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyMyUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMyUserPhone',
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
  '/caos.citadel.auth.api.v1.AuthService/ResendMyPhoneVerificationCode',
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
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.resendMyPhoneVerificationCode =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendMyPhoneVerificationCode',
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
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.resendMyPhoneVerificationCode =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendMyPhoneVerificationCode',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendMyPhoneVerificationCode);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserAddress>}
 */
const methodDescriptor_AuthService_GetMyUserAddress = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetMyUserAddress',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.UserAddress,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAddress.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.UserAddress>}
 */
const methodInfo_AuthService_GetMyUserAddress = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserAddress,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAddress.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserAddress)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserAddress>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getMyUserAddress =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserAddress',
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
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserAddress>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyUserAddress =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyUserAddress',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyUserAddress);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest,
 *   !proto.caos.citadel.auth.api.v1.UserAddress>}
 */
const methodDescriptor_AuthService_UpdateMyUserAddress = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/UpdateMyUserAddress',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest,
  proto.caos.citadel.auth.api.v1.UserAddress,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAddress.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest,
 *   !proto.caos.citadel.auth.api.v1.UserAddress>}
 */
const methodInfo_AuthService_UpdateMyUserAddress = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.UserAddress,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.UserAddress.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.UserAddress)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.UserAddress>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.updateMyUserAddress =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/UpdateMyUserAddress',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserAddress,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UpdateUserAddressRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.UserAddress>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.updateMyUserAddress =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/UpdateMyUserAddress',
      request,
      metadata || {},
      methodDescriptor_AuthService_UpdateMyUserAddress);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.MultiFactors>}
 */
const methodDescriptor_AuthService_GetMyMfas = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetMyMfas',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.MultiFactors,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MultiFactors.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.MultiFactors>}
 */
const methodInfo_AuthService_GetMyMfas = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.MultiFactors,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MultiFactors.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.MultiFactors)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.MultiFactors>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getMyMfas =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyMfas',
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
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.MultiFactors>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyMfas =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyMfas',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyMfas);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.PasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_SetMyPassword = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/SetMyPassword',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.PasswordRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.PasswordRequest} request
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
 *   !proto.caos.citadel.auth.api.v1.PasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_SetMyPassword = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.PasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.PasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.setMyPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SetMyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_SetMyPassword,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.PasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.setMyPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SetMyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_SetMyPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.ResetPasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_RequestPasswordReset = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/RequestPasswordReset',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.ResetPasswordRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest} request
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
 *   !proto.caos.citadel.auth.api.v1.ResetPasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_RequestPasswordReset = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.requestPasswordReset =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RequestPasswordReset',
      request,
      metadata || {},
      methodDescriptor_AuthService_RequestPasswordReset,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.ResetPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.requestPasswordReset =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RequestPasswordReset',
      request,
      metadata || {},
      methodDescriptor_AuthService_RequestPasswordReset);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.ResetPassword,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_PasswordReset = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/PasswordReset',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.ResetPassword,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ResetPassword} request
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
 *   !proto.caos.citadel.auth.api.v1.ResetPassword,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_PasswordReset = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ResetPassword} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.ResetPassword} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.passwordReset =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/PasswordReset',
      request,
      metadata || {},
      methodDescriptor_AuthService_PasswordReset,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.ResetPassword} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.passwordReset =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/PasswordReset',
      request,
      metadata || {},
      methodDescriptor_AuthService_PasswordReset);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.PasswordChange,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_ChangeMyPassword = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/ChangeMyPassword',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.PasswordChange,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.PasswordChange} request
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
 *   !proto.caos.citadel.auth.api.v1.PasswordChange,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_ChangeMyPassword = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.PasswordChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.PasswordChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.changeMyPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ChangeMyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyPassword,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.PasswordChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.changeMyPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ChangeMyPassword',
      request,
      metadata || {},
      methodDescriptor_AuthService_ChangeMyPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.MfaOtpResponse>}
 */
const methodDescriptor_AuthService_AddMfaOTP = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/AddMfaOTP',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.MfaOtpResponse,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MfaOtpResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.MfaOtpResponse>}
 */
const methodInfo_AuthService_AddMfaOTP = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.MfaOtpResponse,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MfaOtpResponse.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.MfaOtpResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.MfaOtpResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.addMfaOTP =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/AddMfaOTP',
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
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.MfaOtpResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.addMfaOTP =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/AddMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_AddMfaOTP);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyMfaOtp,
 *   !proto.caos.citadel.auth.api.v1.MfaOtpResponse>}
 */
const methodDescriptor_AuthService_VerifyMfaOTP = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyMfaOTP',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyMfaOtp,
  proto.caos.citadel.auth.api.v1.MfaOtpResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MfaOtpResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.VerifyMfaOtp,
 *   !proto.caos.citadel.auth.api.v1.MfaOtpResponse>}
 */
const methodInfo_AuthService_VerifyMfaOTP = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.MfaOtpResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MfaOtpResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.MfaOtpResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.MfaOtpResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyMfaOTP =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyMfaOTP,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyMfaOtp} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.MfaOtpResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyMfaOTP =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyMfaOTP',
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
  '/caos.citadel.auth.api.v1.AuthService/RemoveMfaOTP',
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
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.removeMfaOTP =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RemoveMfaOTP',
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
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.removeMfaOTP =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/RemoveMfaOTP',
      request,
      metadata || {},
      methodDescriptor_AuthService_RemoveMfaOTP);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.SkipMfaInitRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_SkipMfaInit = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/SkipMfaInit',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.SkipMfaInitRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest} request
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
 *   !proto.caos.citadel.auth.api.v1.SkipMfaInitRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_SkipMfaInit = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.skipMfaInit =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SkipMfaInit',
      request,
      metadata || {},
      methodDescriptor_AuthService_SkipMfaInit,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.SkipMfaInitRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.skipMfaInit =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SkipMfaInit',
      request,
      metadata || {},
      methodDescriptor_AuthService_SkipMfaInit);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.VerifyUserInitRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_VerifyUserInit = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/VerifyUserInit',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.VerifyUserInitRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest} request
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
 *   !proto.caos.citadel.auth.api.v1.VerifyUserInitRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_VerifyUserInit = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.verifyUserInit =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyUserInit',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyUserInit,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.VerifyUserInitRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.verifyUserInit =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/VerifyUserInit',
      request,
      metadata || {},
      methodDescriptor_AuthService_VerifyUserInit);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AuthService_ResendUserInitMail = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/ResendUserInitMail',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserID} request
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
 *   !proto.caos.citadel.auth.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_AuthService_ResendUserInitMail = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.resendUserInitMail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendUserInitMail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendUserInitMail,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.resendUserInitMail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/ResendUserInitMail',
      request,
      metadata || {},
      methodDescriptor_AuthService_ResendUserInitMail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.UserID,
 *   !proto.caos.citadel.auth.api.v1.User>}
 */
const methodDescriptor_AuthService_GetUserByID = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetUserByID',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.UserID,
  proto.caos.citadel.auth.api.v1.User,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.UserID,
 *   !proto.caos.citadel.auth.api.v1.User>}
 */
const methodInfo_AuthService_GetUserByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.User,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getUserByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserByID',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserByID,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getUserByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetUserByID',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetUserByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.ApplicationID,
 *   !proto.caos.citadel.auth.api.v1.Application>}
 */
const methodDescriptor_AuthService_GetApplicationByID = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetApplicationByID',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.ApplicationID,
  proto.caos.citadel.auth.api.v1.Application,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.Application.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.ApplicationID,
 *   !proto.caos.citadel.auth.api.v1.Application>}
 */
const methodInfo_AuthService_GetApplicationByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.Application,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.Application.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.Application)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.Application>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getApplicationByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetApplicationByID',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetApplicationByID,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.Application>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getApplicationByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetApplicationByID',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetApplicationByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.ApplicationSearchRequest,
 *   !proto.caos.citadel.auth.api.v1.ApplicationSearchResponse>}
 */
const methodDescriptor_AuthService_SearchApplications = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/SearchApplications',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.ApplicationSearchRequest,
  proto.caos.citadel.auth.api.v1.ApplicationSearchResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.ApplicationSearchRequest,
 *   !proto.caos.citadel.auth.api.v1.ApplicationSearchResponse>}
 */
const methodInfo_AuthService_SearchApplications = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.ApplicationSearchResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.ApplicationSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.ApplicationSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.ApplicationSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.searchApplications =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SearchApplications',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchApplications,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.ApplicationSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.searchApplications =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SearchApplications',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchApplications);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest,
 *   !proto.caos.citadel.auth.api.v1.Application>}
 */
const methodDescriptor_AuthService_AuthorizeApplication = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/AuthorizeApplication',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest,
  proto.caos.citadel.auth.api.v1.Application,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.Application.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest,
 *   !proto.caos.citadel.auth.api.v1.Application>}
 */
const methodInfo_AuthService_AuthorizeApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.Application,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.Application.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.Application)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.Application>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.authorizeApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/AuthorizeApplication',
      request,
      metadata || {},
      methodDescriptor_AuthService_AuthorizeApplication,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.ApplicationAuthorizeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.Application>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.authorizeApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/AuthorizeApplication',
      request,
      metadata || {},
      methodDescriptor_AuthService_AuthorizeApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.GrantSearchRequest,
 *   !proto.caos.citadel.auth.api.v1.GrantSearchResponse>}
 */
const methodDescriptor_AuthService_SearchGrant = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/SearchGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.GrantSearchRequest,
  proto.caos.citadel.auth.api.v1.GrantSearchResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.GrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.GrantSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.GrantSearchRequest,
 *   !proto.caos.citadel.auth.api.v1.GrantSearchResponse>}
 */
const methodInfo_AuthService_SearchGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.GrantSearchResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.GrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.GrantSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.GrantSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.GrantSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.searchGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SearchGrant',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchGrant,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.GrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.GrantSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.searchGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SearchGrant',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest,
 *   !proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse>}
 */
const methodDescriptor_AuthService_SearchMyProjectOrgs = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/SearchMyProjectOrgs',
  grpc.web.MethodType.UNARY,
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest,
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest,
 *   !proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse>}
 */
const methodInfo_AuthService_SearchMyProjectOrgs = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse,
  /**
   * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.searchMyProjectOrgs =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SearchMyProjectOrgs',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchMyProjectOrgs,
      callback);
};


/**
 * @param {!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.MyProjectOrgSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.searchMyProjectOrgs =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/SearchMyProjectOrgs',
      request,
      metadata || {},
      methodDescriptor_AuthService_SearchMyProjectOrgs);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.IsAdminResponse>}
 */
const methodDescriptor_AuthService_IsIamAdmin = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/IsIamAdmin',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.IsAdminResponse,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.IsAdminResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.IsAdminResponse>}
 */
const methodInfo_AuthService_IsIamAdmin = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.IsAdminResponse,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.IsAdminResponse.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.IsAdminResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.IsAdminResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.isIamAdmin =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/IsIamAdmin',
      request,
      metadata || {},
      methodDescriptor_AuthService_IsIamAdmin,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.IsAdminResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.isIamAdmin =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/IsIamAdmin',
      request,
      metadata || {},
      methodDescriptor_AuthService_IsIamAdmin);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.MyPermissions>}
 */
const methodDescriptor_AuthService_GetMyCitadelPermissions = new grpc.web.MethodDescriptor(
  '/caos.citadel.auth.api.v1.AuthService/GetMyCitadelPermissions',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.citadel.auth.api.v1.MyPermissions,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MyPermissions.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.citadel.auth.api.v1.MyPermissions>}
 */
const methodInfo_AuthService_GetMyCitadelPermissions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.citadel.auth.api.v1.MyPermissions,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.citadel.auth.api.v1.MyPermissions.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.citadel.auth.api.v1.MyPermissions)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.citadel.auth.api.v1.MyPermissions>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.citadel.auth.api.v1.AuthServiceClient.prototype.getMyCitadelPermissions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyCitadelPermissions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyCitadelPermissions,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.citadel.auth.api.v1.MyPermissions>}
 *     A native promise that resolves to the response
 */
proto.caos.citadel.auth.api.v1.AuthServicePromiseClient.prototype.getMyCitadelPermissions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.citadel.auth.api.v1.AuthService/GetMyCitadelPermissions',
      request,
      metadata || {},
      methodDescriptor_AuthService_GetMyCitadelPermissions);
};


module.exports = proto.caos.citadel.auth.api.v1;

