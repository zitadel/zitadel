import type { DescService } from "@bufbuild/protobuf";
import { Timestamp, timestampDate } from "@bufbuild/protobuf/wkt";
import { createClient, Transport } from "@connectrpc/connect";

export function createClientFor<TService extends DescService>(service: TService) {
  return (transport: Transport) => createClient(service, transport);
}

export function toDate(timestamp: Timestamp | undefined): Date | undefined {
  return timestamp ? timestampDate(timestamp) : undefined;
}
