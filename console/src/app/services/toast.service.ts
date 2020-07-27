import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material/snack-bar';
import { TranslateService } from '@ngx-translate/core';
import { Observable } from 'rxjs';
import { take } from 'rxjs/operators';

import { WarnDialogComponent } from '../modules/warn-dialog/warn-dialog.component';
import { AuthService } from './auth.service';

@Injectable({
    providedIn: 'root',
})
export class ToastService {
    constructor(private dialog: MatDialog,
        private snackBar: MatSnackBar,
        private translate: TranslateService,
        private authService: AuthService,
    ) { }

    public showInfo(message: string, i18nkey: boolean = false): void {
        if (i18nkey) {
            this.translate
                .get(message)
                .subscribe(data => {
                    this.showMessage(data, 'close');
                });
        } else {
            this.showMessage(message, 'close');
        }
    }

    public showError(grpcError: any): void {
        const { message, code, metadata } = grpcError;
        // TODO: remove check for code === 7
        if (code === 16 || (code === 7 && message === 'invalid token')) {
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
                    this.authService.authenticate(undefined, true, true);
                }
            });
        } else {
            this.showMessage(decodeURI(message), 'close', { duration: 3000 });
        }
    }

    private showMessage(message: string, action: string, config?: MatSnackBarConfig): Observable<void> {
        const ref = this.snackBar.open(message, action, config);

        return ref.onAction();
    }
}
