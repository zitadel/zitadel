import { InjectionToken } from '@angular/core';
import { UnaryInterceptor } from 'grpc-web';

export const GRPC_INTERCEPTORS = new InjectionToken<Array<UnaryInterceptor<any, any>>>('GRPC_INTERCEPTORS');
