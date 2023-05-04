import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable, of, switchMap } from 'rxjs';
import { ExhaustedService } from '../exhausted.service';

/**
 * ExhaustedHttpInterceptor shows the exhausted dialog before sending the request if the exhausted cookie is there.
 * Also, it shows the exhausted dialog after receiving an HTTP response status 429.
 */
@Injectable()
export class ExhaustedHttpInterceptor implements HttpInterceptor {
  constructor(private exhaustedSvc: ExhaustedService) {}

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return this.exhaustedSvc.checkCookie().pipe(
      switchMap(() =>
        next.handle(req).pipe(
          switchMap((event) => {
            if (!(event instanceof HttpResponse) || event.status != 429) {
              return of(event);
            }
            return this.exhaustedSvc.showExhaustedDialog().pipe(
              // This map just makes the compiler happy.
              // It should never be executed, as we expect a new page load now.
              map(() => event),
            );
          }),
        ),
      ),
    );
  }
}
