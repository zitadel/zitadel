"use strict";Object.defineProperty(exports, "__esModule", {value: true});

var _chunkDUECDNWCcjs = require('./chunk-DUECDNWC.cjs');

// src/v3alpha.ts
var _user_service_pbjs = require('@zitadel/proto/zitadel/resources/user/v3alpha/user_service_pb.js');
var _user_schema_service_pbjs = require('@zitadel/proto/zitadel/resources/userschema/v3alpha/user_schema_service_pb.js');
var createUserSchemaServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _user_schema_service_pbjs.ZITADELUserSchemas);
var createUserServiceClient = _chunkDUECDNWCcjs.createClientFor.call(void 0, _user_service_pbjs.ZITADELUsers);



exports.createUserSchemaServiceClient = createUserSchemaServiceClient; exports.createUserServiceClient = createUserServiceClient;
//# sourceMappingURL=v3alpha.cjs.map