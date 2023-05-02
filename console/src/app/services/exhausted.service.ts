import { Injectable } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { take } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';

const exhaustedKey = 'zitadel.quota.limiting';

@Injectable({
  providedIn: 'root',
})
export class ExhaustedService {
  constructor(private dialog: MatDialog) {}

  public checkCookie() {
    const cookieValue = this.cookieIsPresent()
    if(cookieValue){
      let isURL = false
      try {
        new URL(cookieValue).toString();
        isURL = true
      } catch (_) {
        // is not url
      }
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: isURL ? 'ACTIONS.CONTINUE' : "",
          titleKey: 'ERRORS.EXHAUSTED.TITLE',
          descriptionKey: 'ERRORS.EXHAUSTED.DESCRIPTION',
        },
        disableClose: false,
        width: '400px',
      });

      dialogRef
        .afterClosed()
        .pipe(take(1))
        .subscribe((resp) => {
          if (resp && isURL) {
            window.open(cookieValue, "_blank");
          }
        });
    }
  }

  private cookieIsPresent(){
    return document.cookie.split(";")
      .map(c => c.trim())
      .find(c => c.startsWith(`${exhaustedKey}=`))
      ?.replace(`${exhaustedKey}=`, "")
  }
}
