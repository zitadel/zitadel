import { ChangeDetectionStrategy, Component, DestroyRef } from '@angular/core';
import { lastValueFrom, Observable, of, ReplaySubject } from 'rxjs';
import { ActionService } from 'src/app/services/action.service';
import { ToastService } from 'src/app/services/toast.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { catchError, map, startWith, switchMap } from 'rxjs/operators';
import { MatDialog } from '@angular/material/dialog';
import { ActionTwoAddTargetDialogComponent } from '../actions-two-add-target/actions-two-add-target-dialog.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import {
  CreateTargetRequestSchema,
  UpdateTargetRequestSchema,
} from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';
import { InfoSectionType } from '../../info-section/info-section.component';

@Component({
  selector: 'cnsl-actions-two-targets',
  templateUrl: './actions-two-targets.component.html',
  styleUrls: ['./actions-two-targets.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoTargetsComponent {
  protected readonly targets$: Observable<Target[]>;
  protected readonly refresh$ = new ReplaySubject<true>(1);

  constructor(
    private readonly actionService: ActionService,
    private readonly toast: ToastService,
    private readonly destroyRef: DestroyRef,
    private readonly dialog: MatDialog,
  ) {
    this.targets$ = this.getTargets$();
  }

  private getTargets$() {
    return this.refresh$.pipe(
      startWith(true),
      switchMap(() => {
        return this.actionService.listTargets({});
      }),
      map(({ targets }) => targets),
      catchError((err) => {
        this.toast.showError(err);
        return of([]);
      }),
    );
  }

  public async deleteTarget(target: Target) {
    await this.actionService.deleteTarget({ id: target.id });
    await new Promise((res) => setTimeout(res, 1000));
    this.refresh$.next(true);
  }

  public async openDialog(target?: Target) {
    const request$ = this.dialog
      .open<
        ActionTwoAddTargetDialogComponent,
        { target?: Target },
        MessageInitShape<typeof UpdateTargetRequestSchema | typeof CreateTargetRequestSchema>
      >(ActionTwoAddTargetDialogComponent, {
        width: '550px',
        data: {
          target: target,
        },
      })
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef));

    const request = await lastValueFrom(request$);
    if (!request) {
      return;
    }

    try {
      if ('id' in request) {
        await this.actionService.updateTarget(request);
      } else {
        const resp = await this.actionService.createTarget(request);
        console.log(`Your singing key: ${resp.signingKey}`);
      }

      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);
    } catch (error) {
      console.error(error);
      this.toast.showError(error);
    }
  }

  protected readonly InfoSectionType = InfoSectionType;
}
