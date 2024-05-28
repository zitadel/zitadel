import { createPromiseClient, Transport } from "@connectrpc/connect";
import type { ServiceType } from "@bufbuild/protobuf";

export function createClientFor<TService extends ServiceType>(
  service: TService,
) {
  return (transport: Transport) => createPromiseClient(service, transport);
}
