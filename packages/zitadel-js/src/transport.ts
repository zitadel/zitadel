import { createConnectTransport as createWebTransport } from "@connectrpc/connect-web";
import { createConnectTransport as createNodeTransport } from "@connectrpc/connect-node";

export type WebTransportOptions = Parameters<typeof createWebTransport>[0];
export type NodeTransportOptions = Parameters<typeof createNodeTransport>[0];

/**
 * Creates a Connect transport suitable for browser environments.
 */
export function createConnectTransport(options: WebTransportOptions) {
  return createWebTransport(options);
}

/**
 * Creates a gRPC transport suitable for Node.js environments.
 */
export function createGrpcTransport(options: NodeTransportOptions) {
  return createNodeTransport(options);
}
