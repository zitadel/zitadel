import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { Subject } from 'rxjs';
import { debounceTime, filter, first, take } from 'rxjs/operators';
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
    public triggerDialog: Subject<boolean> = new Subject();
    constructor(
        private authenticationService: AuthenticationService,
        private storageService: StorageService,
        private dialog: MatDialog,
    ) {
        this.triggerDialog.pipe(debounceTime(1000)).subscribe(() => {
            this.openDialog();
        });
    }

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
                this.triggerDialog.next(true);
            }
            return Promise.reject(error);
        });
    }

    openDialog() {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.LOGIN',
                titleKey: 'ERRORS.TOKENINVALID.TITLE',
                descriptionKey: 'ERRORS.TOKENINVALID.DESCRIPTION',
            },
            disableClose: true,
            width: '400px',
        });

        dialogRef.afterClosed().pipe(take(1)).subscribe(resp => {
            if (resp) {
                this.authenticationService.authenticate(undefined, true, true);
            }
        });
    }
}
