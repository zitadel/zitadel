import * as _connectrpc_connect from '@connectrpc/connect';
import { GrpcTransportOptions } from '@connectrpc/connect-node';

/**
 * Create a server transport using grpc with the given token and configuration options.
 * @param token
 * @param opts
 */
declare function createServerTransport(token: string, opts: GrpcTransportOptions): _connectrpc_connect.Transport;
declare function newSystemToken({ audience, subject, key, expirationTime, }: {
    audience: string;
    subject: string;
    key: string;
    expirationTime?: number | string | Date;
}): Promise<string>;

export { createServerTransport, newSystemToken };
