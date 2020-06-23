/**
 * @fileoverview gRPC-Web generated client stub for caos.zitadel.management.api.v1
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

var google_protobuf_descriptor_pb = require('google-protobuf/google/protobuf/descriptor_pb.js')

var authoption_options_pb = require('./authoption/options_pb.js')
const proto = {};
proto.caos = {};
proto.caos.zitadel = {};
proto.caos.zitadel.management = {};
proto.caos.zitadel.management.api = {};
proto.caos.zitadel.management.api.v1 = require('./management_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient =
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
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient =
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
const methodDescriptor_ManagementService_Healthz = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/Healthz',
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
const methodInfo_ManagementService_Healthz = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.healthz =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/Healthz',
      request,
      metadata || {},
      methodDescriptor_ManagementService_Healthz,
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
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.healthz =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/Healthz',
      request,
      metadata || {},
      methodDescriptor_ManagementService_Healthz);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_Ready = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/Ready',
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
const methodInfo_ManagementService_Ready = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.ready =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/Ready',
      request,
      metadata || {},
      methodDescriptor_ManagementService_Ready,
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
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.ready =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/Ready',
      request,
      metadata || {},
      methodDescriptor_ManagementService_Ready);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Struct>}
 */
const methodDescriptor_ManagementService_Validate = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/Validate',
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
const methodInfo_ManagementService_Validate = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.validate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/Validate',
      request,
      metadata || {},
      methodDescriptor_ManagementService_Validate,
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
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.validate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/Validate',
      request,
      metadata || {},
      methodDescriptor_ManagementService_Validate);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.Iam>}
 */
const methodDescriptor_ManagementService_GetIam = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetIam',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.Iam,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Iam.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.Iam>}
 */
const methodInfo_ManagementService_GetIam = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Iam,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Iam.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Iam)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Iam>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getIam =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetIam',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetIam,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Iam>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getIam =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetIam',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetIam);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserView>}
 */
const methodDescriptor_ManagementService_GetUserByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetUserByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.UserView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserView>}
 */
const methodInfo_ManagementService_GetUserByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getUserByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getUserByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.Email,
 *   !proto.caos.zitadel.management.api.v1.UserView>}
 */
const methodDescriptor_ManagementService_GetUserByEmailGlobal = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetUserByEmailGlobal',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.Email,
  proto.caos.zitadel.management.api.v1.UserView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.Email} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.Email,
 *   !proto.caos.zitadel.management.api.v1.UserView>}
 */
const methodInfo_ManagementService_GetUserByEmailGlobal = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.Email} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.Email} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getUserByEmailGlobal =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserByEmailGlobal',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserByEmailGlobal,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.Email} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getUserByEmailGlobal =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserByEmailGlobal',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserByEmailGlobal);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchUsers = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchUsers',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserSearchRequest,
  proto.caos.zitadel.management.api.v1.UserSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserSearchResponse>}
 */
const methodInfo_ManagementService_SearchUsers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchUsers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchUsers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchUsers,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchUsers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchUsers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchUsers);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UniqueUserRequest,
 *   !proto.caos.zitadel.management.api.v1.UniqueUserResponse>}
 */
const methodDescriptor_ManagementService_IsUserUnique = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/IsUserUnique',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UniqueUserRequest,
  proto.caos.zitadel.management.api.v1.UniqueUserResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UniqueUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UniqueUserResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UniqueUserRequest,
 *   !proto.caos.zitadel.management.api.v1.UniqueUserResponse>}
 */
const methodInfo_ManagementService_IsUserUnique = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UniqueUserResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UniqueUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UniqueUserResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UniqueUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UniqueUserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UniqueUserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.isUserUnique =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/IsUserUnique',
      request,
      metadata || {},
      methodDescriptor_ManagementService_IsUserUnique,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UniqueUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UniqueUserResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.isUserUnique =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/IsUserUnique',
      request,
      metadata || {},
      methodDescriptor_ManagementService_IsUserUnique);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.CreateUserRequest,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodDescriptor_ManagementService_CreateUser = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreateUser',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.CreateUserRequest,
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.CreateUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.CreateUserRequest,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodInfo_ManagementService_CreateUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.CreateUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.CreateUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateUser,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.CreateUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodDescriptor_ManagementService_DeactivateUser = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateUser',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodInfo_ManagementService_DeactivateUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateUser,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodDescriptor_ManagementService_ReactivateUser = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateUser',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodInfo_ManagementService_ReactivateUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateUser,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodDescriptor_ManagementService_LockUser = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/LockUser',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodInfo_ManagementService_LockUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.lockUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/LockUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_LockUser,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.lockUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/LockUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_LockUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodDescriptor_ManagementService_UnlockUser = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UnlockUser',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.User>}
 */
const methodInfo_ManagementService_UnlockUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.User,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.User.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.unlockUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UnlockUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UnlockUser,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.User>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.unlockUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UnlockUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UnlockUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_DeleteUser = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeleteUser',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
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
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_DeleteUser = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deleteUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeleteUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeleteUser,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deleteUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeleteUser',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeleteUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodDescriptor_ManagementService_UserChanges = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UserChanges',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ChangeRequest,
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodInfo_ManagementService_UserChanges = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Changes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Changes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.userChanges =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UserChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UserChanges,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Changes>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.userChanges =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UserChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UserChanges);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodDescriptor_ManagementService_ApplicationChanges = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ApplicationChanges',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ChangeRequest,
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodInfo_ManagementService_ApplicationChanges = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Changes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Changes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.applicationChanges =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ApplicationChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ApplicationChanges,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Changes>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.applicationChanges =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ApplicationChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ApplicationChanges);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodDescriptor_ManagementService_OrgChanges = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/OrgChanges',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ChangeRequest,
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodInfo_ManagementService_OrgChanges = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Changes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Changes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.orgChanges =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/OrgChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_OrgChanges,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Changes>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.orgChanges =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/OrgChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_OrgChanges);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodDescriptor_ManagementService_ProjectChanges = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ProjectChanges',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ChangeRequest,
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ChangeRequest,
 *   !proto.caos.zitadel.management.api.v1.Changes>}
 */
const methodInfo_ManagementService_ProjectChanges = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Changes,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Changes.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Changes)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Changes>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.projectChanges =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectChanges,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Changes>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.projectChanges =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectChanges',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectChanges);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserProfileView>}
 */
const methodDescriptor_ManagementService_GetUserProfile = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetUserProfile',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.UserProfileView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserProfileView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserProfileView>}
 */
const methodInfo_ManagementService_GetUserProfile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserProfileView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserProfileView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserProfileView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserProfileView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getUserProfile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserProfile',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserProfile,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserProfileView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getUserProfile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserProfile',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserProfile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserProfileRequest,
 *   !proto.caos.zitadel.management.api.v1.UserProfile>}
 */
const methodDescriptor_ManagementService_UpdateUserProfile = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateUserProfile',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UpdateUserProfileRequest,
  proto.caos.zitadel.management.api.v1.UserProfile,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserProfileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserProfile.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserProfileRequest,
 *   !proto.caos.zitadel.management.api.v1.UserProfile>}
 */
const methodInfo_ManagementService_UpdateUserProfile = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserProfile,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserProfileRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserProfile.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserProfileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserProfile)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserProfile>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateUserProfile =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateUserProfile',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateUserProfile,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserProfileRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserProfile>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateUserProfile =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateUserProfile',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateUserProfile);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserEmailView>}
 */
const methodDescriptor_ManagementService_GetUserEmail = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetUserEmail',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.UserEmailView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserEmailView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserEmailView>}
 */
const methodInfo_ManagementService_GetUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserEmailView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserEmailView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserEmailView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserEmailView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserEmail',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserEmail,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserEmailView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserEmail',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserEmailRequest,
 *   !proto.caos.zitadel.management.api.v1.UserEmail>}
 */
const methodDescriptor_ManagementService_ChangeUserEmail = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ChangeUserEmail',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UpdateUserEmailRequest,
  proto.caos.zitadel.management.api.v1.UserEmail,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserEmail.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserEmailRequest,
 *   !proto.caos.zitadel.management.api.v1.UserEmail>}
 */
const methodInfo_ManagementService_ChangeUserEmail = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserEmail,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserEmailRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserEmail.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserEmail)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserEmail>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.changeUserEmail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeUserEmail',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeUserEmail,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserEmailRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserEmail>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.changeUserEmail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeUserEmail',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeUserEmail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_ResendEmailVerificationMail = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ResendEmailVerificationMail',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
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
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_ResendEmailVerificationMail = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.resendEmailVerificationMail =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ResendEmailVerificationMail',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ResendEmailVerificationMail,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.resendEmailVerificationMail =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ResendEmailVerificationMail',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ResendEmailVerificationMail);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserPhoneView>}
 */
const methodDescriptor_ManagementService_GetUserPhone = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetUserPhone',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.UserPhoneView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserPhoneView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserPhoneView>}
 */
const methodInfo_ManagementService_GetUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserPhoneView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserPhoneView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserPhoneView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserPhoneView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserPhone',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserPhone,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserPhoneView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserPhone',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserPhone);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserPhoneRequest,
 *   !proto.caos.zitadel.management.api.v1.UserPhone>}
 */
const methodDescriptor_ManagementService_ChangeUserPhone = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ChangeUserPhone',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UpdateUserPhoneRequest,
  proto.caos.zitadel.management.api.v1.UserPhone,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserPhone.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserPhoneRequest,
 *   !proto.caos.zitadel.management.api.v1.UserPhone>}
 */
const methodInfo_ManagementService_ChangeUserPhone = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserPhone,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserPhoneRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserPhone.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserPhone)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserPhone>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.changeUserPhone =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeUserPhone',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeUserPhone,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserPhoneRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserPhone>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.changeUserPhone =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeUserPhone',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeUserPhone);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_ResendPhoneVerificationCode = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ResendPhoneVerificationCode',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
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
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_ResendPhoneVerificationCode = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.resendPhoneVerificationCode =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ResendPhoneVerificationCode',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ResendPhoneVerificationCode,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.resendPhoneVerificationCode =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ResendPhoneVerificationCode',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ResendPhoneVerificationCode);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserAddressView>}
 */
const methodDescriptor_ManagementService_GetUserAddress = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetUserAddress',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.UserAddressView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserAddressView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.UserAddressView>}
 */
const methodInfo_ManagementService_GetUserAddress = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserAddressView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserAddressView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserAddressView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserAddressView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getUserAddress =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserAddress',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserAddress,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserAddressView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getUserAddress =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserAddress',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserAddress);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserAddressRequest,
 *   !proto.caos.zitadel.management.api.v1.UserAddress>}
 */
const methodDescriptor_ManagementService_UpdateUserAddress = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateUserAddress',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UpdateUserAddressRequest,
  proto.caos.zitadel.management.api.v1.UserAddress,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserAddressRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserAddress.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UpdateUserAddressRequest,
 *   !proto.caos.zitadel.management.api.v1.UserAddress>}
 */
const methodInfo_ManagementService_UpdateUserAddress = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserAddress,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UpdateUserAddressRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserAddress.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserAddressRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserAddress)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserAddress>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateUserAddress =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateUserAddress',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateUserAddress,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UpdateUserAddressRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserAddress>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateUserAddress =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateUserAddress',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateUserAddress);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.MultiFactors>}
 */
const methodDescriptor_ManagementService_GetUserMfas = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetUserMfas',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserID,
  proto.caos.zitadel.management.api.v1.MultiFactors,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.MultiFactors.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserID,
 *   !proto.caos.zitadel.management.api.v1.MultiFactors>}
 */
const methodInfo_ManagementService_GetUserMfas = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.MultiFactors,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.MultiFactors.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.MultiFactors)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.MultiFactors>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getUserMfas =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserMfas',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserMfas,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.MultiFactors>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getUserMfas =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetUserMfas',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetUserMfas);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.SetPasswordNotificationRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_SendSetPasswordNotification = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SendSetPasswordNotification',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.SetPasswordNotificationRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.SetPasswordNotificationRequest} request
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
 *   !proto.caos.zitadel.management.api.v1.SetPasswordNotificationRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_SendSetPasswordNotification = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.SetPasswordNotificationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.SetPasswordNotificationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.sendSetPasswordNotification =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SendSetPasswordNotification',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SendSetPasswordNotification,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.SetPasswordNotificationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.sendSetPasswordNotification =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SendSetPasswordNotification',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SendSetPasswordNotification);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_SetInitialPassword = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SetInitialPassword',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordRequest} request
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
 *   !proto.caos.zitadel.management.api.v1.PasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_SetInitialPassword = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.setInitialPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SetInitialPassword',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SetInitialPassword,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.setInitialPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SetInitialPassword',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SetInitialPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 */
const methodDescriptor_ManagementService_GetPasswordComplexityPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetPasswordComplexityPolicy',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 */
const methodInfo_ManagementService_GetPasswordComplexityPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getPasswordComplexityPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetPasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetPasswordComplexityPolicy,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getPasswordComplexityPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetPasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetPasswordComplexityPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyCreate,
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 */
const methodDescriptor_ManagementService_CreatePasswordComplexityPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordComplexityPolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyCreate,
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyCreate,
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 */
const methodInfo_ManagementService_CreatePasswordComplexityPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createPasswordComplexityPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreatePasswordComplexityPolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createPasswordComplexityPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreatePasswordComplexityPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyUpdate,
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 */
const methodDescriptor_ManagementService_UpdatePasswordComplexityPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordComplexityPolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyUpdate,
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyUpdate,
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 */
const methodInfo_ManagementService_UpdatePasswordComplexityPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updatePasswordComplexityPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdatePasswordComplexityPolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updatePasswordComplexityPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdatePasswordComplexityPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_DeletePasswordComplexityPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordComplexityPolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyID} request
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
 *   !proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_DeletePasswordComplexityPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deletePasswordComplexityPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeletePasswordComplexityPolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordComplexityPolicyID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deletePasswordComplexityPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordComplexityPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeletePasswordComplexityPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 */
const methodDescriptor_ManagementService_GetPasswordAgePolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetPasswordAgePolicy',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 */
const methodInfo_ManagementService_GetPasswordAgePolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordAgePolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordAgePolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getPasswordAgePolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetPasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetPasswordAgePolicy,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getPasswordAgePolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetPasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetPasswordAgePolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicyCreate,
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 */
const methodDescriptor_ManagementService_CreatePasswordAgePolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordAgePolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordAgePolicyCreate,
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicyCreate,
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 */
const methodInfo_ManagementService_CreatePasswordAgePolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordAgePolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordAgePolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createPasswordAgePolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreatePasswordAgePolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createPasswordAgePolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreatePasswordAgePolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicyUpdate,
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 */
const methodDescriptor_ManagementService_UpdatePasswordAgePolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordAgePolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordAgePolicyUpdate,
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicyUpdate,
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 */
const methodInfo_ManagementService_UpdatePasswordAgePolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordAgePolicy.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordAgePolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordAgePolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updatePasswordAgePolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdatePasswordAgePolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordAgePolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updatePasswordAgePolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdatePasswordAgePolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicyID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_DeletePasswordAgePolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordAgePolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordAgePolicyID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyID} request
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
 *   !proto.caos.zitadel.management.api.v1.PasswordAgePolicyID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_DeletePasswordAgePolicy = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deletePasswordAgePolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeletePasswordAgePolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordAgePolicyID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deletePasswordAgePolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordAgePolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeletePasswordAgePolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 */
const methodDescriptor_ManagementService_GetPasswordLockoutPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetPasswordLockoutPolicy',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 */
const methodInfo_ManagementService_GetPasswordLockoutPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getPasswordLockoutPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetPasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetPasswordLockoutPolicy,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getPasswordLockoutPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetPasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetPasswordLockoutPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyCreate,
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 */
const methodDescriptor_ManagementService_CreatePasswordLockoutPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordLockoutPolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyCreate,
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyCreate,
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 */
const methodInfo_ManagementService_CreatePasswordLockoutPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createPasswordLockoutPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreatePasswordLockoutPolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createPasswordLockoutPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreatePasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreatePasswordLockoutPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyUpdate,
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 */
const methodDescriptor_ManagementService_UpdatePasswordLockoutPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordLockoutPolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyUpdate,
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyUpdate,
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 */
const methodInfo_ManagementService_UpdatePasswordLockoutPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updatePasswordLockoutPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdatePasswordLockoutPolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updatePasswordLockoutPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdatePasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdatePasswordLockoutPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_DeletePasswordLockoutPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordLockoutPolicy',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyID} request
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
 *   !proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_DeletePasswordLockoutPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deletePasswordLockoutPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeletePasswordLockoutPolicy,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.PasswordLockoutPolicyID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deletePasswordLockoutPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeletePasswordLockoutPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeletePasswordLockoutPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.OrgView>}
 */
const methodDescriptor_ManagementService_GetMyOrg = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetMyOrg',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.OrgView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.OrgView>}
 */
const methodInfo_ManagementService_GetMyOrg = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgView,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgView.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getMyOrg =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetMyOrg',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetMyOrg,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getMyOrg =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetMyOrg',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetMyOrg);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.Domain,
 *   !proto.caos.zitadel.management.api.v1.OrgView>}
 */
const methodDescriptor_ManagementService_GetOrgByDomainGlobal = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetOrgByDomainGlobal',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.Domain,
  proto.caos.zitadel.management.api.v1.OrgView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.Domain} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.Domain,
 *   !proto.caos.zitadel.management.api.v1.OrgView>}
 */
const methodInfo_ManagementService_GetOrgByDomainGlobal = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.Domain} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.Domain} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getOrgByDomainGlobal =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetOrgByDomainGlobal',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetOrgByDomainGlobal,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.Domain} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getOrgByDomainGlobal =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetOrgByDomainGlobal',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetOrgByDomainGlobal);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.Org>}
 */
const methodDescriptor_ManagementService_DeactivateMyOrg = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateMyOrg',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.Org,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Org.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.Org>}
 */
const methodInfo_ManagementService_DeactivateMyOrg = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Org,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Org.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Org)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Org>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateMyOrg =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateMyOrg',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateMyOrg,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Org>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateMyOrg =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateMyOrg',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateMyOrg);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.Org>}
 */
const methodDescriptor_ManagementService_ReactivateMyOrg = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateMyOrg',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.Org,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Org.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.Org>}
 */
const methodInfo_ManagementService_ReactivateMyOrg = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Org,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Org.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Org)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Org>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateMyOrg =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateMyOrg',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateMyOrg,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Org>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateMyOrg =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateMyOrg',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateMyOrg);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.OrgDomainSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchMyOrgDomains = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgDomains',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.OrgDomainSearchRequest,
  proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OrgDomainSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.OrgDomainSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse>}
 */
const methodInfo_ManagementService_SearchMyOrgDomains = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OrgDomainSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.OrgDomainSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchMyOrgDomains =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgDomains',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchMyOrgDomains,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.OrgDomainSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgDomainSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchMyOrgDomains =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgDomains',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchMyOrgDomains);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.AddOrgDomainRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgDomain>}
 */
const methodDescriptor_ManagementService_AddMyOrgDomain = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/AddMyOrgDomain',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.AddOrgDomainRequest,
  proto.caos.zitadel.management.api.v1.OrgDomain,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.AddOrgDomainRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgDomain.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.AddOrgDomainRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgDomain>}
 */
const methodInfo_ManagementService_AddMyOrgDomain = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgDomain,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.AddOrgDomainRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgDomain.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.AddOrgDomainRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgDomain)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgDomain>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.addMyOrgDomain =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddMyOrgDomain',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddMyOrgDomain,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.AddOrgDomainRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgDomain>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.addMyOrgDomain =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddMyOrgDomain',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddMyOrgDomain);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.RemoveOrgDomainRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveMyOrgDomain = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgDomain',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.RemoveOrgDomainRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgDomainRequest} request
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
 *   !proto.caos.zitadel.management.api.v1.RemoveOrgDomainRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveMyOrgDomain = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgDomainRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgDomainRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeMyOrgDomain =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgDomain',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveMyOrgDomain,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgDomainRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeMyOrgDomain =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgDomain',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveMyOrgDomain);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.OrgIamPolicy>}
 */
const methodDescriptor_ManagementService_GetMyOrgIamPolicy = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetMyOrgIamPolicy',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.OrgIamPolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgIamPolicy.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.OrgIamPolicy>}
 */
const methodInfo_ManagementService_GetMyOrgIamPolicy = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgIamPolicy,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgIamPolicy.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgIamPolicy)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgIamPolicy>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getMyOrgIamPolicy =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetMyOrgIamPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetMyOrgIamPolicy,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgIamPolicy>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getMyOrgIamPolicy =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetMyOrgIamPolicy',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetMyOrgIamPolicy);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.OrgMemberRoles>}
 */
const methodDescriptor_ManagementService_GetOrgMemberRoles = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetOrgMemberRoles',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.OrgMemberRoles,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMemberRoles.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.OrgMemberRoles>}
 */
const methodInfo_ManagementService_GetOrgMemberRoles = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgMemberRoles,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMemberRoles.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgMemberRoles)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgMemberRoles>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getOrgMemberRoles =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetOrgMemberRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetOrgMemberRoles,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgMemberRoles>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getOrgMemberRoles =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetOrgMemberRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetOrgMemberRoles);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.AddOrgMemberRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgMember>}
 */
const methodDescriptor_ManagementService_AddMyOrgMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/AddMyOrgMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.AddOrgMemberRequest,
  proto.caos.zitadel.management.api.v1.OrgMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.AddOrgMemberRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMember.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.AddOrgMemberRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgMember>}
 */
const methodInfo_ManagementService_AddMyOrgMember = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.AddOrgMemberRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMember.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.AddOrgMemberRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgMember)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgMember>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.addMyOrgMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddMyOrgMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddMyOrgMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.AddOrgMemberRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgMember>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.addMyOrgMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddMyOrgMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddMyOrgMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ChangeOrgMemberRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgMember>}
 */
const methodDescriptor_ManagementService_ChangeMyOrgMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ChangeMyOrgMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ChangeOrgMemberRequest,
  proto.caos.zitadel.management.api.v1.OrgMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeOrgMemberRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMember.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ChangeOrgMemberRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgMember>}
 */
const methodInfo_ManagementService_ChangeMyOrgMember = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ChangeOrgMemberRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMember.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeOrgMemberRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgMember)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgMember>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.changeMyOrgMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeMyOrgMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeMyOrgMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ChangeOrgMemberRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgMember>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.changeMyOrgMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeMyOrgMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeMyOrgMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.RemoveOrgMemberRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveMyOrgMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.RemoveOrgMemberRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgMemberRequest} request
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
 *   !proto.caos.zitadel.management.api.v1.RemoveOrgMemberRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveMyOrgMember = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgMemberRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgMemberRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeMyOrgMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveMyOrgMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.RemoveOrgMemberRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeMyOrgMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveMyOrgMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveMyOrgMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.OrgMemberSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchMyOrgMembers = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgMembers',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.OrgMemberSearchRequest,
  proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OrgMemberSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.OrgMemberSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse>}
 */
const methodInfo_ManagementService_SearchMyOrgMembers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OrgMemberSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.OrgMemberSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchMyOrgMembers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgMembers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchMyOrgMembers,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.OrgMemberSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OrgMemberSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchMyOrgMembers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchMyOrgMembers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchMyOrgMembers);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchProjects = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchProjects',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectSearchRequest,
  proto.caos.zitadel.management.api.v1.ProjectSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectSearchResponse>}
 */
const methodInfo_ManagementService_SearchProjects = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchProjects =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjects',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjects,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchProjects =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjects',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjects);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectID,
 *   !proto.caos.zitadel.management.api.v1.ProjectView>}
 */
const methodDescriptor_ManagementService_ProjectByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ProjectByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectID,
  proto.caos.zitadel.management.api.v1.ProjectView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectID,
 *   !proto.caos.zitadel.management.api.v1.ProjectView>}
 */
const methodInfo_ManagementService_ProjectByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.projectByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.projectByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectCreateRequest,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodDescriptor_ManagementService_CreateProject = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreateProject',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectCreateRequest,
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectCreateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectCreateRequest,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodInfo_ManagementService_CreateProject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectCreateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectCreateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Project)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Project>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createProject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProject,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectCreateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Project>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createProject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectUpdateRequest,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodDescriptor_ManagementService_UpdateProject = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateProject',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectUpdateRequest,
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectUpdateRequest,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodInfo_ManagementService_UpdateProject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Project)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Project>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateProject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProject,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Project>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateProject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectID,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodDescriptor_ManagementService_DeactivateProject = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateProject',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectID,
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectID,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodInfo_ManagementService_DeactivateProject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Project)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Project>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateProject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProject,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Project>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateProject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectID,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodDescriptor_ManagementService_ReactivateProject = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateProject',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectID,
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectID,
 *   !proto.caos.zitadel.management.api.v1.Project>}
 */
const methodInfo_ManagementService_ReactivateProject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Project,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Project.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Project)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Project>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateProject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProject,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Project>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateProject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProject',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.GrantedProjectSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchGrantedProjects = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchGrantedProjects',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.GrantedProjectSearchRequest,
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.GrantedProjectSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.GrantedProjectSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>}
 */
const methodInfo_ManagementService_SearchGrantedProjects = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.GrantedProjectSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.GrantedProjectSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchGrantedProjects =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchGrantedProjects',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchGrantedProjects,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.GrantedProjectSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchGrantedProjects =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchGrantedProjects',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchGrantedProjects);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantView>}
 */
const methodDescriptor_ManagementService_GetGrantedProjectByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetGrantedProjectByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantID,
  proto.caos.zitadel.management.api.v1.ProjectGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantView>}
 */
const methodInfo_ManagementService_GetGrantedProjectByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getGrantedProjectByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetGrantedProjectByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetGrantedProjectByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getGrantedProjectByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetGrantedProjectByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetGrantedProjectByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberRoles>}
 */
const methodDescriptor_ManagementService_GetProjectMemberRoles = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetProjectMemberRoles',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.ProjectMemberRoles,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMemberRoles.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberRoles>}
 */
const methodInfo_ManagementService_GetProjectMemberRoles = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectMemberRoles,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMemberRoles.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectMemberRoles)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectMemberRoles>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getProjectMemberRoles =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetProjectMemberRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetProjectMemberRoles,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectMemberRoles>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getProjectMemberRoles =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetProjectMemberRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetProjectMemberRoles);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchProjectMembers = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchProjectMembers',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectMemberSearchRequest,
  proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse>}
 */
const methodInfo_ManagementService_SearchProjectMembers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchProjectMembers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectMembers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectMembers,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectMemberSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchProjectMembers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectMembers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectMembers);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberAdd,
 *   !proto.caos.zitadel.management.api.v1.ProjectMember>}
 */
const methodDescriptor_ManagementService_AddProjectMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/AddProjectMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectMemberAdd,
  proto.caos.zitadel.management.api.v1.ProjectMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberAdd} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMember.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberAdd,
 *   !proto.caos.zitadel.management.api.v1.ProjectMember>}
 */
const methodInfo_ManagementService_AddProjectMember = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberAdd} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMember.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberAdd} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectMember)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectMember>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.addProjectMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddProjectMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddProjectMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberAdd} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectMember>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.addProjectMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddProjectMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddProjectMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberChange,
 *   !proto.caos.zitadel.management.api.v1.ProjectMember>}
 */
const methodDescriptor_ManagementService_ChangeProjectMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectMemberChange,
  proto.caos.zitadel.management.api.v1.ProjectMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMember.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberChange,
 *   !proto.caos.zitadel.management.api.v1.ProjectMember>}
 */
const methodInfo_ManagementService_ChangeProjectMember = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectMember.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectMember)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectMember>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.changeProjectMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeProjectMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectMember>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.changeProjectMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeProjectMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberRemove,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveProjectMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectMemberRemove,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberRemove} request
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
 *   !proto.caos.zitadel.management.api.v1.ProjectMemberRemove,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveProjectMember = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberRemove} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberRemove} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeProjectMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectMemberRemove} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeProjectMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchProjectRoles = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchProjectRoles',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectRoleSearchRequest,
  proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse>}
 */
const methodInfo_ManagementService_SearchProjectRoles = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchProjectRoles =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectRoles,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectRoleSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchProjectRoles =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectRoles);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleAdd,
 *   !proto.caos.zitadel.management.api.v1.ProjectRole>}
 */
const methodDescriptor_ManagementService_AddProjectRole = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/AddProjectRole',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectRoleAdd,
  proto.caos.zitadel.management.api.v1.ProjectRole,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAdd} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectRole.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleAdd,
 *   !proto.caos.zitadel.management.api.v1.ProjectRole>}
 */
const methodInfo_ManagementService_AddProjectRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectRole,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAdd} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectRole.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAdd} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectRole)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectRole>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.addProjectRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddProjectRole,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAdd} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectRole>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.addProjectRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddProjectRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleAddBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_BulkAddProjectRole = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/BulkAddProjectRole',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectRoleAddBulk,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAddBulk} request
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
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleAddBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_BulkAddProjectRole = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAddBulk} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAddBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.bulkAddProjectRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkAddProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkAddProjectRole,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleAddBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.bulkAddProjectRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkAddProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkAddProjectRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleChange,
 *   !proto.caos.zitadel.management.api.v1.ProjectRole>}
 */
const methodDescriptor_ManagementService_ChangeProjectRole = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectRole',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectRoleChange,
  proto.caos.zitadel.management.api.v1.ProjectRole,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectRole.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleChange,
 *   !proto.caos.zitadel.management.api.v1.ProjectRole>}
 */
const methodInfo_ManagementService_ChangeProjectRole = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectRole,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectRole.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectRole)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectRole>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.changeProjectRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeProjectRole,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectRole>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.changeProjectRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeProjectRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleRemove,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveProjectRole = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectRole',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectRoleRemove,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleRemove} request
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
 *   !proto.caos.zitadel.management.api.v1.ProjectRoleRemove,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveProjectRole = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleRemove} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleRemove} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeProjectRole =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectRole,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectRoleRemove} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeProjectRole =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectRole',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectRole);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ApplicationSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ApplicationSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchApplications = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchApplications',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ApplicationSearchRequest,
  proto.caos.zitadel.management.api.v1.ApplicationSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ApplicationSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ApplicationSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ApplicationSearchResponse>}
 */
const methodInfo_ManagementService_SearchApplications = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ApplicationSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ApplicationSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ApplicationSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ApplicationSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchApplications =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchApplications',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchApplications,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ApplicationSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchApplications =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchApplications',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchApplications);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.ApplicationView>}
 */
const methodDescriptor_ManagementService_ApplicationByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ApplicationByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ApplicationID,
  proto.caos.zitadel.management.api.v1.ApplicationView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ApplicationView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.ApplicationView>}
 */
const methodInfo_ManagementService_ApplicationByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ApplicationView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ApplicationView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ApplicationView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ApplicationView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.applicationByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ApplicationByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ApplicationByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ApplicationView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.applicationByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ApplicationByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ApplicationByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.OIDCApplicationCreate,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodDescriptor_ManagementService_CreateOIDCApplication = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreateOIDCApplication',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.OIDCApplicationCreate,
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OIDCApplicationCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.OIDCApplicationCreate,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodInfo_ManagementService_CreateOIDCApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OIDCApplicationCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.OIDCApplicationCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Application)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Application>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createOIDCApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateOIDCApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateOIDCApplication,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.OIDCApplicationCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Application>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createOIDCApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateOIDCApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateOIDCApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ApplicationUpdate,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodDescriptor_ManagementService_UpdateApplication = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateApplication',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ApplicationUpdate,
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ApplicationUpdate,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodInfo_ManagementService_UpdateApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Application)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Application>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateApplication,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Application>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodDescriptor_ManagementService_DeactivateApplication = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateApplication',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ApplicationID,
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodInfo_ManagementService_DeactivateApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Application)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Application>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateApplication,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Application>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodDescriptor_ManagementService_ReactivateApplication = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateApplication',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ApplicationID,
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.Application>}
 */
const methodInfo_ManagementService_ReactivateApplication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.Application,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.Application.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.Application)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.Application>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateApplication,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.Application>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveApplication = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveApplication',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ApplicationID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
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
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveApplication = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeApplication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveApplication,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeApplication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveApplication',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveApplication);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.OIDCConfigUpdate,
 *   !proto.caos.zitadel.management.api.v1.OIDCConfig>}
 */
const methodDescriptor_ManagementService_UpdateApplicationOIDCConfig = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateApplicationOIDCConfig',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.OIDCConfigUpdate,
  proto.caos.zitadel.management.api.v1.OIDCConfig,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OIDCConfigUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OIDCConfig.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.OIDCConfigUpdate,
 *   !proto.caos.zitadel.management.api.v1.OIDCConfig>}
 */
const methodInfo_ManagementService_UpdateApplicationOIDCConfig = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.OIDCConfig,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.OIDCConfigUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.OIDCConfig.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.OIDCConfigUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.OIDCConfig)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.OIDCConfig>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateApplicationOIDCConfig =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateApplicationOIDCConfig',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateApplicationOIDCConfig,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.OIDCConfigUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.OIDCConfig>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateApplicationOIDCConfig =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateApplicationOIDCConfig',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateApplicationOIDCConfig);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.ClientSecret>}
 */
const methodDescriptor_ManagementService_RegenerateOIDCClientSecret = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RegenerateOIDCClientSecret',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ApplicationID,
  proto.caos.zitadel.management.api.v1.ClientSecret,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ClientSecret.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ApplicationID,
 *   !proto.caos.zitadel.management.api.v1.ClientSecret>}
 */
const methodInfo_ManagementService_RegenerateOIDCClientSecret = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ClientSecret,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ClientSecret.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ClientSecret)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ClientSecret>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.regenerateOIDCClientSecret =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RegenerateOIDCClientSecret',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RegenerateOIDCClientSecret,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ApplicationID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ClientSecret>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.regenerateOIDCClientSecret =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RegenerateOIDCClientSecret',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RegenerateOIDCClientSecret);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchProjectGrants = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrants',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchRequest,
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>}
 */
const methodInfo_ManagementService_SearchProjectGrants = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchProjectGrants =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectGrants,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchProjectGrants =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectGrants);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantView>}
 */
const methodDescriptor_ManagementService_ProjectGrantByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ProjectGrantByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantID,
  proto.caos.zitadel.management.api.v1.ProjectGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantView>}
 */
const methodInfo_ManagementService_ProjectGrantByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.projectGrantByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectGrantByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.projectGrantByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectGrantByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodDescriptor_ManagementService_CreateProjectGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreateProjectGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantCreate,
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodInfo_ManagementService_CreateProjectGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createProjectGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProjectGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createProjectGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProjectGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodDescriptor_ManagementService_UpdateProjectGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantUpdate,
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodInfo_ManagementService_UpdateProjectGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateProjectGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProjectGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateProjectGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProjectGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodDescriptor_ManagementService_DeactivateProjectGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantID,
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodInfo_ManagementService_DeactivateProjectGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateProjectGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProjectGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateProjectGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProjectGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodDescriptor_ManagementService_ReactivateProjectGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantID,
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrant>}
 */
const methodInfo_ManagementService_ReactivateProjectGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateProjectGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProjectGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateProjectGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProjectGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveProjectGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
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
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveProjectGrant = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeProjectGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeProjectGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles>}
 */
const methodDescriptor_ManagementService_GetProjectGrantMemberRoles = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/GetProjectGrantMemberRoles',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles>}
 */
const methodInfo_ManagementService_GetProjectGrantMemberRoles = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.getProjectGrantMemberRoles =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetProjectGrantMemberRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetProjectGrantMemberRoles,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantMemberRoles>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.getProjectGrantMemberRoles =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/GetProjectGrantMemberRoles',
      request,
      metadata || {},
      methodDescriptor_ManagementService_GetProjectGrantMemberRoles);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchProjectGrantMembers = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrantMembers',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchRequest,
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse>}
 */
const methodInfo_ManagementService_SearchProjectGrantMembers = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchProjectGrantMembers =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrantMembers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectGrantMembers,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantMemberSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchProjectGrantMembers =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrantMembers',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectGrantMembers);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberAdd,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMember>}
 */
const methodDescriptor_ManagementService_AddProjectGrantMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/AddProjectGrantMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberAdd,
  proto.caos.zitadel.management.api.v1.ProjectGrantMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberAdd} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMember.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberAdd,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMember>}
 */
const methodInfo_ManagementService_AddProjectGrantMember = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberAdd} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMember.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberAdd} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantMember)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantMember>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.addProjectGrantMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddProjectGrantMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddProjectGrantMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberAdd} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantMember>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.addProjectGrantMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/AddProjectGrantMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_AddProjectGrantMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberChange,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMember>}
 */
const methodDescriptor_ManagementService_ChangeProjectGrantMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectGrantMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberChange,
  proto.caos.zitadel.management.api.v1.ProjectGrantMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMember.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberChange,
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMember>}
 */
const methodInfo_ManagementService_ChangeProjectGrantMember = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.ProjectGrantMember,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberChange} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.ProjectGrantMember.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.ProjectGrantMember)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.ProjectGrantMember>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.changeProjectGrantMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectGrantMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeProjectGrantMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberChange} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.ProjectGrantMember>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.changeProjectGrantMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ChangeProjectGrantMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ChangeProjectGrantMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberRemove,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveProjectGrantMember = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrantMember',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantMemberRemove,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberRemove} request
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
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantMemberRemove,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveProjectGrantMember = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberRemove} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberRemove} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeProjectGrantMember =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrantMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectGrantMember,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantMemberRemove} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeProjectGrantMember =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveProjectGrantMember',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveProjectGrantMember);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchUserGrants = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchUserGrants',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantSearchRequest,
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 */
const methodInfo_ManagementService_SearchUserGrants = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrantSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchUserGrants =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchUserGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchUserGrants,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchUserGrants =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchUserGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchUserGrants);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrantView>}
 */
const methodDescriptor_ManagementService_UserGrantByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UserGrantByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrantView>}
 */
const methodInfo_ManagementService_UserGrantByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrantView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrantView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.userGrantByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UserGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UserGrantByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrantView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.userGrantByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UserGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UserGrantByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_CreateUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreateUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantCreate,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_CreateUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_UpdateUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantUpdate,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_UpdateUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_DeactivateUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_DeactivateUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_ReactivateUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_ReactivateUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_RemoveUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/RemoveUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantID,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
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
 *   !proto.caos.zitadel.management.api.v1.UserGrantID,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_RemoveUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.removeUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.removeUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/RemoveUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_RemoveUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantCreateBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_BulkCreateUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/BulkCreateUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantCreateBulk,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreateBulk} request
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
 *   !proto.caos.zitadel.management.api.v1.UserGrantCreateBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_BulkCreateUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreateBulk} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreateBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.bulkCreateUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkCreateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkCreateUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreateBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.bulkCreateUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkCreateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkCreateUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantUpdateBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_BulkUpdateUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/BulkUpdateUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantUpdateBulk,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdateBulk} request
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
 *   !proto.caos.zitadel.management.api.v1.UserGrantUpdateBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_BulkUpdateUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdateBulk} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdateBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.bulkUpdateUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkUpdateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkUpdateUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantUpdateBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.bulkUpdateUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkUpdateUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkUpdateUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantRemoveBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ManagementService_BulkRemoveUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/BulkRemoveUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantRemoveBulk,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantRemoveBulk} request
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
 *   !proto.caos.zitadel.management.api.v1.UserGrantRemoveBulk,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ManagementService_BulkRemoveUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantRemoveBulk} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantRemoveBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.bulkRemoveUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkRemoveUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkRemoveUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantRemoveBulk} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.bulkRemoveUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/BulkRemoveUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_BulkRemoveUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchProjectUserGrants = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchProjectUserGrants',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectUserGrantSearchRequest,
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 */
const methodInfo_ManagementService_SearchProjectUserGrants = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrantSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchProjectUserGrants =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectUserGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectUserGrants,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchProjectUserGrants =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectUserGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectUserGrants);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrantView>}
 */
const methodDescriptor_ManagementService_ProjectUserGrantByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ProjectUserGrantByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrantView>}
 */
const methodInfo_ManagementService_ProjectUserGrantByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrantView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrantView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.projectUserGrantByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectUserGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectUserGrantByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrantView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.projectUserGrantByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectUserGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectUserGrantByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.UserGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_CreateProjectUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreateProjectUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.UserGrantCreate,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.UserGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_CreateProjectUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createProjectUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProjectUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.UserGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createProjectUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProjectUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_UpdateProjectUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectUserGrantUpdate,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_UpdateProjectUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateProjectUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProjectUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateProjectUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProjectUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_DeactivateProjectUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_DeactivateProjectUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateProjectUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProjectUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateProjectUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProjectUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_ReactivateProjectUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_ReactivateProjectUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateProjectUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProjectUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateProjectUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProjectUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 */
const methodDescriptor_ManagementService_SearchProjectGrantUserGrants = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrantUserGrants',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantSearchRequest,
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantSearchRequest,
 *   !proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 */
const methodInfo_ManagementService_SearchProjectGrantUserGrants = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantSearchResponse.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrantSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.searchProjectGrantUserGrants =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrantUserGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectGrantUserGrants,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrantSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.searchProjectGrantUserGrants =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/SearchProjectGrantUserGrants',
      request,
      metadata || {},
      methodDescriptor_ManagementService_SearchProjectGrantUserGrants);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrantView>}
 */
const methodDescriptor_ManagementService_ProjectGrantUserGrantByID = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ProjectGrantUserGrantByID',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantView.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrantView>}
 */
const methodInfo_ManagementService_ProjectGrantUserGrantByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrantView,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrantView.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrantView)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrantView>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.projectGrantUserGrantByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectGrantUserGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectGrantUserGrantByID,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrantView>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.projectGrantUserGrantByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ProjectGrantUserGrantByID',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ProjectGrantUserGrantByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_CreateProjectGrantUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/CreateProjectGrantUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantCreate,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantCreate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_CreateProjectGrantUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantCreate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.createProjectGrantUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProjectGrantUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantCreate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.createProjectGrantUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/CreateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_CreateProjectGrantUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_UpdateProjectGrantUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectGrantUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantUpdate,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantUpdate,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_UpdateProjectGrantUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantUpdate} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.updateProjectGrantUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProjectGrantUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantUpdate} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.updateProjectGrantUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/UpdateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_UpdateProjectGrantUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_DeactivateProjectGrantUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectGrantUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_DeactivateProjectGrantUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.deactivateProjectGrantUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProjectGrantUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.deactivateProjectGrantUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/DeactivateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_DeactivateProjectGrantUserGrant);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodDescriptor_ManagementService_ReactivateProjectGrantUserGrant = new grpc.web.MethodDescriptor(
  '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectGrantUserGrant',
  grpc.web.MethodType.UNARY,
  proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID,
 *   !proto.caos.zitadel.management.api.v1.UserGrant>}
 */
const methodInfo_ManagementService_ReactivateProjectGrantUserGrant = new grpc.web.AbstractClientBase.MethodInfo(
  proto.caos.zitadel.management.api.v1.UserGrant,
  /**
   * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.caos.zitadel.management.api.v1.UserGrant.deserializeBinary
);


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.caos.zitadel.management.api.v1.UserGrant)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.caos.zitadel.management.api.v1.UserGrant>|undefined}
 *     The XHR Node Readable Stream
 */
proto.caos.zitadel.management.api.v1.ManagementServiceClient.prototype.reactivateProjectGrantUserGrant =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProjectGrantUserGrant,
      callback);
};


/**
 * @param {!proto.caos.zitadel.management.api.v1.ProjectGrantUserGrantID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.caos.zitadel.management.api.v1.UserGrant>}
 *     A native promise that resolves to the response
 */
proto.caos.zitadel.management.api.v1.ManagementServicePromiseClient.prototype.reactivateProjectGrantUserGrant =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/caos.zitadel.management.api.v1.ManagementService/ReactivateProjectGrantUserGrant',
      request,
      metadata || {},
      methodDescriptor_ManagementService_ReactivateProjectGrantUserGrant);
};


module.exports = proto.caos.zitadel.management.api.v1;

