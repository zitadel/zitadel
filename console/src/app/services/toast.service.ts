import { Injectable } from '@angular/core';
import { MatSnackBar, MatSnackBarHorizontalPosition, MatSnackBarVerticalPosition } from '@angular/material/snack-bar';
import { TranslateService } from '@ngx-translate/core';
import { Observable } from 'rxjs';
import { take } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class ToastService {
  horizontalPosition: MatSnackBarHorizontalPosition = 'end';
  verticalPosition: MatSnackBarVerticalPosition = 'top';

  constructor(
    private snackBar: MatSnackBar,
    private translate: TranslateService,
  ) {}

  public showInfo(message: string, i18nkey: boolean = false): void {
    if (i18nkey) {
      this.translate.get(message).subscribe((data) => {
        this.translate
          .get('ACTIONS.CLOSE')
          .pipe(take(1))
          .subscribe((value) => {
            this.showMessage(data, value, true);
          });
      });
    } else {
      this.translate
        .get('ACTIONS.CLOSE')
        .pipe(take(1))
        .subscribe((value) => {
          this.showMessage(message, value, true);
        });
    }
  }

  public showError(error: any | string, isGrpc: boolean = true, i18nKey: boolean = false): void {
    if (isGrpc) {
      const { message, code, metadata } = error;
      if (code !== 16) {
        this.translate
          .get('ACTIONS.CLOSE')
          .pipe(take(1))
          .subscribe((value) => {
            this.showMessage(decodeURI(message), value, false);
          });
      }
    } else if (i18nKey) {
      this.translate
        .get(error)
        .pipe(take(1))
        .subscribe((value) => {
          this.showMessage(value, '', false);
        });
    } else {
      this.showMessage(error as string, '', false);
    }
  }

  private showMessage(message: string, action: string, success: boolean): Observable<void> {
    const ref = this.snackBar.open(message, action, {
      data: {
        message,
      },
      duration: success ? 4000 : 50000,
      panelClass: ['data-e2e-message', success ? 'data-e2e-success' : 'data-e2e-failure'],
      horizontalPosition: this.horizontalPosition,
      verticalPosition: this.verticalPosition,
    });

    return ref.onAction();
  }
}
