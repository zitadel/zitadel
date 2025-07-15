import * as _connectrpc_connect from '@connectrpc/connect';
import { Transport, Interceptor } from '@connectrpc/connect';
export { Client, Code, ConnectError } from '@connectrpc/connect';
import { DescService } from '@bufbuild/protobuf';
export { JsonObject, create, fromJson, toJson } from '@bufbuild/protobuf';
import { Timestamp } from '@bufbuild/protobuf/wkt';
export { Duration, Timestamp, TimestampSchema, timestampDate, timestampFromDate, timestampFromMs, timestampMs } from '@bufbuild/protobuf/wkt';
export { GenService } from '@bufbuild/protobuf/codegenv1';

declare function createClientFor<TService extends DescService>(service: TService): (transport: Transport) => _connectrpc_connect.Client<TService>;
declare function toDate(timestamp: Timestamp | undefined): Date | undefined;

/**
 * Creates an interceptor that adds an Authorization header with a Bearer token.
 * @param token
 */
declare function NewAuthorizationBearerInterceptor(token: string): Interceptor;

export { NewAuthorizationBearerInterceptor, createClientFor, toDate };
