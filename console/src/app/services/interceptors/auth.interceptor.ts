import { DestroyRef, Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Request, RpcError, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { firstValueFrom, identity, lastValueFrom, Observable, Subject } from 'rxjs';
import { debounceTime, filter, map } from 'rxjs/operators';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';

import { AuthenticationService } from '../authentication.service';
import { StorageService } from '../storage.service';
import { AuthConfig } from 'angular-oauth2-oidc';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ConnectError, Interceptor } from '@connectrpc/connect';

const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer';
const accessTokenStorageKey = 'access_token';

@Injectable({ providedIn: 'root' })
export class AuthInterceptorProvider {
  private readonly triggerDialog: Subject<boolean> = new Subject();

  constructor(
    private readonly authenticationService: AuthenticationService,
    private readonly storageService: StorageService,
    private readonly dialog: MatDialog,
    destroyRef: DestroyRef,
  ) {
    this.triggerDialog.pipe(debounceTime(1000), takeUntilDestroyed(destroyRef)).subscribe(() => this.openDialog());
  }

  getToken(): Observable<string> {
    return this.authenticationService.authenticationChanged.pipe(
      filter(identity),
      map(() => this.storageService.getItem(accessTokenStorageKey)),
      map((token) => `${bearerPrefix} ${token}`),
    );
  }

  handleError = (error: any): never => {
    if (!(error instanceof RpcError) && !(error instanceof ConnectError)) {
      throw error;
    }

    if (error.code === 16 || (error.code === 7 && error.message === 'mfa required (AUTHZ-Kl3p0)')) {
      this.triggerDialog.next(true);
    }
    throw error;
  };

  private async openDialog(): Promise<void> {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.LOGIN',
        titleKey: 'ERRORS.TOKENINVALID.TITLE',
        descriptionKey: 'ERRORS.TOKENINVALID.DESCRIPTION',
      },
      disableClose: true,
      width: '400px',
    });

    const resp = await lastValueFrom(dialogRef.afterClosed());
    if (!resp) {
      return;
    }

    const idToken = this.authenticationService.getIdToken();
    const configWithPrompt: Partial<AuthConfig> = {
      customQueryParams: {
        id_token_hint: idToken,
      },
    };

    await this.authenticationService.authenticate(configWithPrompt, true);
  }
}

@Injectable({ providedIn: 'root' })
/**
 * Set the authentication token
 */
export class AuthInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
  constructor(private readonly authInterceptorProvider: AuthInterceptorProvider) {}

  public async intercept(
    request: Request<TReq, TResp>,
    invoker: (request: Request<TReq, TResp>) => Promise<UnaryResponse<TReq, TResp>>,
  ): Promise<UnaryResponse<TReq, TResp>> {
    const metadata = request.getMetadata();
    metadata[authorizationKey] = await firstValueFrom(this.authInterceptorProvider.getToken());

    return invoker(request).catch(this.authInterceptorProvider.handleError);
  }
}

export function NewConnectWebAuthInterceptor(authInterceptorProvider: AuthInterceptorProvider): Interceptor {
  return (next) => async (req) => {
    if (!req.header.get('Authorization')) {
      const token = await firstValueFrom(authInterceptorProvider.getToken());
      req.header.set('Authorization', token);
    }

    return next(req).catch(authInterceptorProvider.handleError);
  };
}
