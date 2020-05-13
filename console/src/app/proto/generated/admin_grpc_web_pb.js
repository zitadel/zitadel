/**
 * @fileoverview gRPC-Web generated client stub for zitadel.admin.api.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');


var google_api_annotations_pb = require('./google/api/annotations_pb.js')

var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js')

var validate_validate_pb = require('./validate/validate_pb.js')

var protoc$gen$swagger_options_annotations_pb = require('./protoc-gen-swagger/options/annotations_pb.js')

var authoption_options_pb = require('./authoption/options_pb.js')
const proto = {};
proto.zitadel = {};
proto.zitadel.admin = {};
proto.zitadel.admin.api = {};
proto.zitadel.admin.api.v1 = require('./admin_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.zitadel.admin.api.v1.AdminServiceClient =
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
proto.zitadel.admin.api.v1.AdminServicePromiseClient =
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
const methodDescriptor_AdminService_Healthz = new grpc.web.MethodDescriptor(
  '/zitadel.admin.api.v1.AdminService/Healthz',
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
const methodInfo_AdminService_Healthz = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.zitadel.admin.api.v1.AdminServiceClient.prototype.healthz =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/Healthz',
      request,
      metadata || {},
      methodDescriptor_AdminService_Healthz,
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
proto.zitadel.admin.api.v1.AdminServicePromiseClient.prototype.healthz =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/Healthz',
      request,
      metadata || {},
      methodDescriptor_AdminService_Healthz);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_AdminService_Ready = new grpc.web.MethodDescriptor(
  '/zitadel.admin.api.v1.AdminService/Ready',
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
const methodInfo_AdminService_Ready = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.zitadel.admin.api.v1.AdminServiceClient.prototype.ready =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/Ready',
      request,
      metadata || {},
      methodDescriptor_AdminService_Ready,
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
proto.zitadel.admin.api.v1.AdminServicePromiseClient.prototype.ready =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/Ready',
      request,
      metadata || {},
      methodDescriptor_AdminService_Ready);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Struct>}
 */
const methodDescriptor_AdminService_Validate = new grpc.web.MethodDescriptor(
  '/zitadel.admin.api.v1.AdminService/Validate',
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
const methodInfo_AdminService_Validate = new grpc.web.AbstractClientBase.MethodInfo(
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
proto.zitadel.admin.api.v1.AdminServiceClient.prototype.validate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/Validate',
      request,
      metadata || {},
      methodDescriptor_AdminService_Validate,
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
proto.zitadel.admin.api.v1.AdminServicePromiseClient.prototype.validate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/Validate',
      request,
      metadata || {},
      methodDescriptor_AdminService_Validate);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zitadel.admin.api.v1.UniqueOrgRequest,
 *   !proto.zitadel.admin.api.v1.UniqueOrgResponse>}
 */
const methodDescriptor_AdminService_IsOrgUnique = new grpc.web.MethodDescriptor(
  '/zitadel.admin.api.v1.AdminService/IsOrgUnique',
  grpc.web.MethodType.UNARY,
  proto.zitadel.admin.api.v1.UniqueOrgRequest,
  proto.zitadel.admin.api.v1.UniqueOrgResponse,
  /**
   * @param {!proto.zitadel.admin.api.v1.UniqueOrgRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.UniqueOrgResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zitadel.admin.api.v1.UniqueOrgRequest,
 *   !proto.zitadel.admin.api.v1.UniqueOrgResponse>}
 */
const methodInfo_AdminService_IsOrgUnique = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zitadel.admin.api.v1.UniqueOrgResponse,
  /**
   * @param {!proto.zitadel.admin.api.v1.UniqueOrgRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.UniqueOrgResponse.deserializeBinary
);


/**
 * @param {!proto.zitadel.admin.api.v1.UniqueOrgRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zitadel.admin.api.v1.UniqueOrgResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zitadel.admin.api.v1.UniqueOrgResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zitadel.admin.api.v1.AdminServiceClient.prototype.isOrgUnique =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/IsOrgUnique',
      request,
      metadata || {},
      methodDescriptor_AdminService_IsOrgUnique,
      callback);
};


/**
 * @param {!proto.zitadel.admin.api.v1.UniqueOrgRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zitadel.admin.api.v1.UniqueOrgResponse>}
 *     A native promise that resolves to the response
 */
proto.zitadel.admin.api.v1.AdminServicePromiseClient.prototype.isOrgUnique =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/IsOrgUnique',
      request,
      metadata || {},
      methodDescriptor_AdminService_IsOrgUnique);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zitadel.admin.api.v1.OrgID,
 *   !proto.zitadel.admin.api.v1.Org>}
 */
const methodDescriptor_AdminService_GetOrgByID = new grpc.web.MethodDescriptor(
  '/zitadel.admin.api.v1.AdminService/GetOrgByID',
  grpc.web.MethodType.UNARY,
  proto.zitadel.admin.api.v1.OrgID,
  proto.zitadel.admin.api.v1.Org,
  /**
   * @param {!proto.zitadel.admin.api.v1.OrgID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.Org.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zitadel.admin.api.v1.OrgID,
 *   !proto.zitadel.admin.api.v1.Org>}
 */
const methodInfo_AdminService_GetOrgByID = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zitadel.admin.api.v1.Org,
  /**
   * @param {!proto.zitadel.admin.api.v1.OrgID} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.Org.deserializeBinary
);


/**
 * @param {!proto.zitadel.admin.api.v1.OrgID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zitadel.admin.api.v1.Org)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zitadel.admin.api.v1.Org>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zitadel.admin.api.v1.AdminServiceClient.prototype.getOrgByID =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/GetOrgByID',
      request,
      metadata || {},
      methodDescriptor_AdminService_GetOrgByID,
      callback);
};


/**
 * @param {!proto.zitadel.admin.api.v1.OrgID} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zitadel.admin.api.v1.Org>}
 *     A native promise that resolves to the response
 */
proto.zitadel.admin.api.v1.AdminServicePromiseClient.prototype.getOrgByID =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/GetOrgByID',
      request,
      metadata || {},
      methodDescriptor_AdminService_GetOrgByID);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zitadel.admin.api.v1.OrgSearchRequest,
 *   !proto.zitadel.admin.api.v1.OrgSearchResponse>}
 */
const methodDescriptor_AdminService_SearchOrgs = new grpc.web.MethodDescriptor(
  '/zitadel.admin.api.v1.AdminService/SearchOrgs',
  grpc.web.MethodType.UNARY,
  proto.zitadel.admin.api.v1.OrgSearchRequest,
  proto.zitadel.admin.api.v1.OrgSearchResponse,
  /**
   * @param {!proto.zitadel.admin.api.v1.OrgSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.OrgSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zitadel.admin.api.v1.OrgSearchRequest,
 *   !proto.zitadel.admin.api.v1.OrgSearchResponse>}
 */
const methodInfo_AdminService_SearchOrgs = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zitadel.admin.api.v1.OrgSearchResponse,
  /**
   * @param {!proto.zitadel.admin.api.v1.OrgSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.OrgSearchResponse.deserializeBinary
);


/**
 * @param {!proto.zitadel.admin.api.v1.OrgSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zitadel.admin.api.v1.OrgSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zitadel.admin.api.v1.OrgSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zitadel.admin.api.v1.AdminServiceClient.prototype.searchOrgs =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/SearchOrgs',
      request,
      metadata || {},
      methodDescriptor_AdminService_SearchOrgs,
      callback);
};


/**
 * @param {!proto.zitadel.admin.api.v1.OrgSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zitadel.admin.api.v1.OrgSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.zitadel.admin.api.v1.AdminServicePromiseClient.prototype.searchOrgs =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/SearchOrgs',
      request,
      metadata || {},
      methodDescriptor_AdminService_SearchOrgs);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.zitadel.admin.api.v1.OrgSetUpRequest,
 *   !proto.zitadel.admin.api.v1.OrgSetUpResponse>}
 */
const methodDescriptor_AdminService_SetUpOrg = new grpc.web.MethodDescriptor(
  '/zitadel.admin.api.v1.AdminService/SetUpOrg',
  grpc.web.MethodType.UNARY,
  proto.zitadel.admin.api.v1.OrgSetUpRequest,
  proto.zitadel.admin.api.v1.OrgSetUpResponse,
  /**
   * @param {!proto.zitadel.admin.api.v1.OrgSetUpRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.OrgSetUpResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.zitadel.admin.api.v1.OrgSetUpRequest,
 *   !proto.zitadel.admin.api.v1.OrgSetUpResponse>}
 */
const methodInfo_AdminService_SetUpOrg = new grpc.web.AbstractClientBase.MethodInfo(
  proto.zitadel.admin.api.v1.OrgSetUpResponse,
  /**
   * @param {!proto.zitadel.admin.api.v1.OrgSetUpRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.zitadel.admin.api.v1.OrgSetUpResponse.deserializeBinary
);


/**
 * @param {!proto.zitadel.admin.api.v1.OrgSetUpRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.zitadel.admin.api.v1.OrgSetUpResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.zitadel.admin.api.v1.OrgSetUpResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.zitadel.admin.api.v1.AdminServiceClient.prototype.setUpOrg =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/SetUpOrg',
      request,
      metadata || {},
      methodDescriptor_AdminService_SetUpOrg,
      callback);
};


/**
 * @param {!proto.zitadel.admin.api.v1.OrgSetUpRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.zitadel.admin.api.v1.OrgSetUpResponse>}
 *     A native promise that resolves to the response
 */
proto.zitadel.admin.api.v1.AdminServicePromiseClient.prototype.setUpOrg =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/zitadel.admin.api.v1.AdminService/SetUpOrg',
      request,
      metadata || {},
      methodDescriptor_AdminService_SetUpOrg);
};


module.exports = proto.zitadel.admin.api.v1;

