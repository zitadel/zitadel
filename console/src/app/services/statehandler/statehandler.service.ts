import { Injectable, Injector, OnDestroy, Type } from '@angular/core';
import { GuardsCheckStart, Router, RouterEvent } from '@angular/router';
import { OAuthService } from 'angular-oauth2-oidc';
import { Observable, of, Subject, throwError } from 'rxjs';
import { filter, map, shareReplay, switchMap, take, takeUntil } from 'rxjs/operators';

import { StatehandlerProcessorService } from './statehandler-processor.service';

export abstract class StatehandlerService {
  public abstract createState(): Observable<string | undefined>;
  public abstract initStateHandler(): void;
}

@Injectable()
export class StatehandlerServiceImpl implements StatehandlerService, OnDestroy {
  private events?: Observable<string>;
  private unsubscribe$: Subject<void> = new Subject();

  constructor(
    oauthService: OAuthService,
    private injector: Injector,
    private processor: StatehandlerProcessorService,
  ) {
    oauthService.events
      .pipe(
        filter((event) => event.type === 'token_received'),
        map(() => oauthService.state),
        takeUntil(this.unsubscribe$),
      )
      .subscribe((state) => {
        processor.restoreState(state);
      });
  }

  public initStateHandler(): void {
    const router = this.injector.get(Router as Type<Router>);
    this.events = (router.events as Observable<RouterEvent>).pipe(
      filter((event) => event instanceof GuardsCheckStart),
      map((event) => event.url),
      shareReplay(1),
    );

    this.events.pipe(takeUntil(this.unsubscribe$)).subscribe();
  }

  public createState(): Observable<string | undefined> {
    if (this.events === undefined) {
      return throwError(() => new Error('no router events'));
    }

    return this.events.pipe(
      take(1),
      switchMap((url: string) => {
        if (url.includes('?login_hint=')) {
          const newUrl = this.removeParam('login_hint', url);
          const urlWithoutBasePath = newUrl.includes('/ui/console') ? newUrl.replace('/ui/console', '') : newUrl;
          return of(this.processor.createState(urlWithoutBasePath));
        } else if (url) {
          const urlWithoutBasePath = url.includes('/ui/console') ? url.replace('/ui/console', '') : url;
          return of(this.processor.createState(urlWithoutBasePath));
        } else {
          return of(undefined);
        }
      }),
    );
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

  removeParam(key: string, sourceURL: string) {
    var rtn = sourceURL.split('?')[0],
      param,
      params_arr = [],
      queryString = sourceURL.indexOf('?') !== -1 ? sourceURL.split('?')[1] : '';
    if (queryString !== '') {
      params_arr = queryString.split('&');
      for (var i = params_arr.length - 1; i >= 0; i -= 1) {
        param = params_arr[i].split('=')[0];
        if (param === key) {
          params_arr.splice(i, 1);
        }
      }
      if (params_arr.length) rtn = rtn + '?' + params_arr.join('&');
    }
    return rtn;
  }
}
