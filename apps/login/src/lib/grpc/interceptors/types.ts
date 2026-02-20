export interface GrpcRequestHeader {
  readonly set: (key: string, value: string) => void;
}

export interface GrpcRequest {
  readonly service?: { readonly typeName: string };
  readonly method?: { readonly name: string };
  readonly url?: string;
  readonly header: GrpcRequestHeader;
}

export interface GrpcError extends Error {
  readonly code?: number;
}

export type Interceptor<T = unknown> = (
  next: (req: GrpcRequest) => Promise<T>,
) => (req: GrpcRequest) => Promise<T>;
