import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { TranslateService } from '@ngx-translate/core';

@Injectable({
    providedIn: 'root',
})
export class ToastService {
    constructor(private snackBar: MatSnackBar, private translate: TranslateService) { }

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

    public showError(message: string): void {
        this.showMessage(decodeURI(message), 'close');
    }

    private showMessage(message: string, action: string): void {
        this.snackBar.open(message, action, {
            duration: 5000,
        });
    }
}
