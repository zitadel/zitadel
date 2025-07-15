"use strict";Object.defineProperty(exports, "__esModule", {value: true});// src/helpers.ts
var _wkt = require('@bufbuild/protobuf/wkt');
var _connect = require('@connectrpc/connect');
function createClientFor(service) {
  return (transport) => _connect.createClient.call(void 0, service, transport);
}
function toDate(timestamp) {
  return timestamp ? _wkt.timestampDate.call(void 0, timestamp) : void 0;
}




exports.createClientFor = createClientFor; exports.toDate = toDate;
//# sourceMappingURL=chunk-DUECDNWC.cjs.map