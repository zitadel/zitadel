import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { SwUpdate } from '@angular/service-worker';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { NEVER, switchMap } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class UpdateService {
  constructor(
    private swUpdate: SwUpdate,
    snackbar: MatSnackBar,
  ) {
    this.swUpdate.versionUpdates
      .pipe(
        switchMap((event) => {
          if (event.type !== 'VERSION_DETECTED') {
            return NEVER;
          }

          return snackbar.open('Update Available', 'Reload').onAction();
        }),
        takeUntilDestroyed(),
      )
      .subscribe(() => {
        window.location.reload();
      });
  }
}
