import { Injectable } from '@angular/core';
import { Request, StatusCode, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { from, lastValueFrom, map, of, switchMap } from 'rxjs';
import { ExhaustedService } from '../exhausted.service';

/**
 * ExhaustedGrpcInterceptor shows the exhausted dialog before sending the request if the exhausted cookie is there.
 * Also, it shows the exhausted dialog after receiving a gRPC response status 8.
 */
@Injectable({ providedIn: 'root' })
export class ExhaustedGrpcInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
  constructor(private exhaustedSvc: ExhaustedService) {}

  public async intercept(
    request: Request<TReq, TResp>,
    invoker: (request: Request<TReq, TResp>) => Promise<UnaryResponse<TReq, TResp>>,
  ): Promise<UnaryResponse<TReq, TResp>> {
    return lastValueFrom(
      this.exhaustedSvc.checkCookie().pipe(
        switchMap(() =>
          from(invoker(request)).pipe(
            switchMap((response) => {
              if (response.getStatus().code != StatusCode.RESOURCE_EXHAUSTED) {
                return of(response);
              }
              return this.exhaustedSvc.showExhaustedDialog().pipe(
                // This map just makes the compiler happy.
                // It should never be executed, as we expect a new page load now.
                map(() => response),
              );
            }),
          ),
        ),
      ),
    );
  }
}
