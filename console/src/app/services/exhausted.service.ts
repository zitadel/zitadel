import { Injectable } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { of, tap } from 'rxjs';
import { WarnDialogComponent } from '../modules/warn-dialog/warn-dialog.component';

@Injectable({
  providedIn: 'root',
})
export class ExhaustedService {
  private isClosed = true;

  constructor(private dialog: MatDialog) {}

  public showExhaustedDialog(instanceManagementUrl?: string) {
    if (!this.isClosed) {
      return of(undefined);
    }
    this.isClosed = false;
    return this.dialog
      .open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.CONTINUE',
          titleKey: 'ERRORS.EXHAUSTED.TITLE',
          descriptionKey: 'ERRORS.EXHAUSTED.DESCRIPTION',
        },
        disableClose: true,
        width: '400px',
        id: 'authenticated-requests-exhausted-dialog',
      })
      .afterClosed()
      .pipe(
        tap(() => {
          // Just reload if there is no instance management url
          location.href = instanceManagementUrl || location.href;
        }),
      );
  }
}
