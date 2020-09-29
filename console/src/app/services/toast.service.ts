import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material/snack-bar';
import { TranslateService } from '@ngx-translate/core';
import { Observable } from 'rxjs';

import { AuthenticationService } from './authentication.service';

@Injectable({
    providedIn: 'root',
})
export class ToastService {
    constructor(private dialog: MatDialog,
        private snackBar: MatSnackBar,
        private translate: TranslateService,
        private authService: AuthenticationService,
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
        if (code !== 16) {
            this.showMessage(decodeURI(message), 'close', { duration: 4000 });
        }
    }

    private showMessage(message: string, action: string, config?: MatSnackBarConfig): Observable<void> {
        const ref = this.snackBar.open(message, action, config);

        return ref.onAction();
    }
}
