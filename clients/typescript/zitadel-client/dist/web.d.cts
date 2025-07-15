import * as _connectrpc_connect from '@connectrpc/connect';
import { GrpcTransportOptions } from '@connectrpc/connect-node';

/**
 * Create a client transport using grpc web with the given token and configuration options.
 * @param token
 * @param opts
 */
declare function createClientTransport(token: string, opts: GrpcTransportOptions): _connectrpc_connect.Transport;

export { createClientTransport };
