import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';

@Injectable({
    providedIn: 'root',
})
export class ToastService {
    constructor(private snackBar: MatSnackBar) { }

    public showInfo(message: string): void {
        this.showMessage(message, 'close');
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
