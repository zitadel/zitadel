import { Injectable } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { map, of, switchMap } from 'rxjs';
import { WarnDialogComponent } from '../modules/warn-dialog/warn-dialog.component';
import { EnvironmentService } from './environment.service';

@Injectable({
  providedIn: 'root',
})
export class ExhaustedService {
  private exhaustedKey = 'zitadel.quota.limiting';
  private isClosed = true;

  constructor(private envSvc: EnvironmentService, private dialog: MatDialog) {}

  public checkCookie() {
    const cookieIsThere = !!document.cookie
      .split(';')
      .map((c) => c.trim())
      .find((c) => c.startsWith(`${this.exhaustedKey}=`));

    if (cookieIsThere) {
      return this.showExhaustedDialog();
    }
    return of(undefined);
  }

  public showExhaustedDialog() {
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
      .pipe(switchMap(() => this.resolveExhausted()));
  }

  private resolveExhausted() {
    return this.envSvc.env.pipe(
      map((env) => {
        // Delete the cookie
        document.cookie = `${this.exhaustedKey}=; Path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT"`;
        // Just reload if there is no instance management url
        location.href = env.instance_management_url || location.href;
      }),
    );
  }
}
