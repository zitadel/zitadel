import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { filter, first, take } from 'rxjs/operators';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';

import { AuthenticationService } from '../authentication.service';
import { StorageService } from '../storage.service';


const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer';
const accessTokenStorageKey = 'access_token';
@Injectable({ providedIn: 'root' })
/**
 * Set the authentication token
 */
export class AuthInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
    constructor(
        private authenticationService: AuthenticationService,
        private storageService: StorageService,
        private dialog: MatDialog,
    ) { }

    public async intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {
        await this.authenticationService.authenticationChanged.pipe(
            filter((authed) => !!authed),
            first(),
        ).toPromise();

        const metadata = request.getMetadata();
        const accessToken = this.storageService.getItem(accessTokenStorageKey);
        metadata[authorizationKey] = `${bearerPrefix} ${accessToken}`;

        return invoker(request).then((response: any) => {
            return response;
        }).catch((error: any) => {
            if (error.code === 16) {
                const dialogRef = this.dialog.open(WarnDialogComponent, {
                    data: {
                        confirmKey: 'ACTIONS.LOGIN',
                        titleKey: 'ERRORS.TOKENINVALID.TITLE',
                        descriptionKey: 'ERRORS.TOKENINVALID.DESCRIPTION',
                    },
                    width: '400px',
                });

                dialogRef.afterClosed().pipe(take(1)).subscribe(resp => {
                    if (resp) {
                        this.authenticationService.authenticate(undefined, true, true);
                    }
                });
            }
            return Promise.reject(error);
        });
    }
}
