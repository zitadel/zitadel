import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { TranslateService } from '@ngx-translate/core';
import { Observable } from 'rxjs';
import { take } from 'rxjs/operators';

@Injectable({
    providedIn: 'root',
})
export class ToastService {
    constructor(
        private snackBar: MatSnackBar,
        private translate: TranslateService,
    ) { }

    public showInfo(message: string, i18nkey: boolean = false): void {
        if (i18nkey) {
            this.translate
                .get(message)
                .subscribe(data => {
                    this.translate.get('ACTIONS.CLOSE').pipe(take(1)).subscribe(value => {
                        this.showMessage(data, value);
                    });
                });
        } else {
            this.translate.get('ACTIONS.CLOSE').pipe(take(1)).subscribe(value => {
                this.showMessage(message, value);
            });
        }
    }

    public showError(grpcError: any): void {
        const { message, code, metadata } = grpcError;
        if (code !== 16) {
            this.translate.get('ACTIONS.CLOSE').pipe(take(1)).subscribe(value => {
                this.showMessage(decodeURI(message), value);
            });
        }
    }

    private showMessage(message: string, action: string): Observable<void> {
        const ref = this.snackBar.open(message, action, { duration: 4000 });

        return ref.onAction();
    }
}
