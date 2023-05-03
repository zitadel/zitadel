import { Injectable } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { BehaviorSubject, lastValueFrom, of, take } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';

const exhaustedKey = 'zitadel.quota.limiting';

@Injectable({
  providedIn: 'root',
})
export class ExhaustedService {
  private isClosed = new BehaviorSubject<boolean>(true);

  constructor(private dialog: MatDialog) {}

  public checkCookie() {
    const cookieValue = this.cookieIsPresent();
    console.log('cookie value', cookieValue);
    if (!cookieValue || !this.isClosed.value) {
      return lastValueFrom(of(true));
    }
    this.isClosed.next(false);
    let isURL = false;
    try {
      new URL(cookieValue).toString();
      isURL = true;
    } catch (_) {
      // is not url
    }
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: isURL ? 'ACTIONS.CONTINUE' : '',
        titleKey: 'ERRORS.EXHAUSTED.TITLE',
        descriptionKey: 'ERRORS.EXHAUSTED.DESCRIPTION',
      },
      disableClose: true,
      width: '400px',
      id: 'authenticated-requests-exhausted-dialog',
    });
    const newClosed = dialogRef.afterClosed();
    if (isURL) {
      newClosed.pipe(take(1)).subscribe((resp) => {
        if (resp && isURL) {
          this.deleteCookie();
          location.href = cookieValue;
        }
        this.isClosed.next(true);
      });
    }
    return newClosed;
  }

  private cookieIsPresent() {
    return document.cookie
      .split(';')
      .map((c) => c.trim())
      .find((c) => c.startsWith(`${exhaustedKey}=`))
      ?.replace(`${exhaustedKey}=`, '');
  }

  private deleteCookie() {
    document.cookie = `${exhaustedKey}=; Path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT"`;
  }
}
