import { createConnectTransport as createWebTransport } from "@connectrpc/connect-web";
import { createConnectTransport as createNodeTransport } from "@connectrpc/connect-node";

export interface TransportOptions {
  baseUrl: string;
  interceptors?: Parameters<typeof createWebTransport>[0]["interceptors"];
}

/**
 * Creates a Connect transport suitable for browser environments.
 */
export function createConnectTransport(options: TransportOptions) {
  return createWebTransport(options);
}

/**
 * Creates a gRPC transport suitable for Node.js environments.
 */
export function createGrpcTransport(options: TransportOptions) {
  return createNodeTransport(options);
}
