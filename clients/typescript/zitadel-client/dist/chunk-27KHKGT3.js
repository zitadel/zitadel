// src/helpers.ts
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { createClient } from "@connectrpc/connect";
function createClientFor(service) {
  return (transport) => createClient(service, transport);
}
function toDate(timestamp) {
  return timestamp ? timestampDate(timestamp) : void 0;
}

export {
  createClientFor,
  toDate
};
//# sourceMappingURL=chunk-27KHKGT3.js.map