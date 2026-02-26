import type { DescService } from "@bufbuild/protobuf";
import { createClient, Transport } from "@connectrpc/connect";

export type { Client } from "@connectrpc/connect";

/**
 * Creates a typed client factory for a given protobuf service descriptor.
 */
export function createClientFor<TService extends DescService>(service: TService) {
  return (transport: Transport) => createClient(service, transport);
}
