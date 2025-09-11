import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { map, Observable, of, switchMap, tap } from 'rxjs';
import { WarnDialogComponent } from '../modules/warn-dialog/warn-dialog.component';
import { Environment } from './environment.service';

@Injectable({
  providedIn: 'root',
})
export class ExhaustedService {
  private isClosed = true;

  constructor(private dialog: MatDialog) {}

  public showExhaustedDialog(env$: Observable<Environment>) {
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
        switchMap(() => env$),
        tap((env) => {
          // Just reload if there is no instance management url
          location.href = env.instance_management_url || location.href;
        }),
        map(() => undefined),
      );
  }
}
