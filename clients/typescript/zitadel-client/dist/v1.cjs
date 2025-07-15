"use strict";Object.defineProperty(exports, "__esModule", {value: true});

var _chunkDUECDNWCcjs = require('./chunk-DUECDNWC.cjs');

// src/v1.ts
var _admin_pbjs = require('@zitadel/proto/zitadel/admin_pb.js');
var _auth_pbjs = require('@zitadel/proto/zitadel/auth_pb.js');
var _management_pbjs = require('@zitadel/proto/zitadel/management_pb.js');
var _system_pbjs = require('@zitadel/proto/zitadel/system_pb.js');
var createAdminServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _admin_pbjs.AdminService);
var createAuthServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _auth_pbjs.AuthService);
var createManagementServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _management_pbjs.ManagementService);
var createSystemServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _system_pbjs.SystemService);





exports.createAdminServiceClient = createAdminServiceClient; exports.createAuthServiceClient = createAuthServiceClient; exports.createManagementServiceClient = createManagementServiceClient; exports.createSystemServiceClient = createSystemServiceClient;
//# sourceMappingURL=v1.cjs.map