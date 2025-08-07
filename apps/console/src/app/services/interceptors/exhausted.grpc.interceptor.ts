import { Injectable } from '@angular/core';
import { Request, RpcError, StatusCode, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { EnvironmentService } from '../environment.service';
import { ExhaustedService } from '../exhausted.service';
import { lastValueFrom } from 'rxjs';

/**
 * ExhaustedGrpcInterceptor shows the exhausted dialog after receiving a gRPC response status 8.
 */
@Injectable({ providedIn: 'root' })
export class ExhaustedGrpcInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
  constructor(
    private exhaustedSvc: ExhaustedService,
    private envSvc: EnvironmentService,
  ) {}

  public async intercept(
    request: Request<TReq, TResp>,
    invoker: (request: Request<TReq, TResp>) => Promise<UnaryResponse<TReq, TResp>>,
  ): Promise<UnaryResponse<TReq, TResp>> {
    try {
      return await invoker(request);
    } catch (error: any) {
      if (error instanceof RpcError && error.code === StatusCode.RESOURCE_EXHAUSTED) {
        await lastValueFrom(this.exhaustedSvc.showExhaustedDialog(this.envSvc.env));
      }
      throw error;
    }
  }
}
