import { ChangeDetectionStrategy, Component, DestroyRef } from '@angular/core';
import { ActionService } from 'src/app/services/action.service';
import { lastValueFrom, Observable, of, Subject } from 'rxjs';
import { catchError, map, startWith, switchMap } from 'rxjs/operators';
import { ToastService } from 'src/app/services/toast.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import {
  ActionTwoAddActionDialogComponent,
  ActionTwoAddActionDialogData,
  ActionTwoAddActionDialogResult,
  CorrectlyTypedExecution,
  correctlyTypeExecution,
} from '../actions-two-add-action/actions-two-add-action-dialog.component';
import { MatDialog } from '@angular/material/dialog';
import { MessageInitShape } from '@bufbuild/protobuf';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import { InfoSectionType } from '../../info-section/info-section.component';
import { ExecutionFieldName } from '@zitadel/proto/zitadel/action/v2beta/query_pb';

@Component({
  selector: 'cnsl-actions-two-actions',
  templateUrl: './actions-two-actions.component.html',
  styleUrls: ['./actions-two-actions.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoActionsComponent {
  protected readonly refresh$ = new Subject<true>();
  protected readonly executions$: Observable<CorrectlyTypedExecution[]>;
  protected readonly targets$: Observable<Target[]>;

  constructor(
    private readonly actionService: ActionService,
    private readonly toast: ToastService,
    private readonly destroyRef: DestroyRef,
    private readonly dialog: MatDialog,
  ) {
    this.executions$ = this.getExecutions$();
    this.targets$ = this.getTargets$();
  }

  private getExecutions$() {
    return this.refresh$.pipe(
      startWith(true),
      switchMap(() => {
        return this.actionService.listExecutions({ sortingColumn: ExecutionFieldName.ID, pagination: { asc: true } });
      }),
      map(({ executions }) => executions.map(correctlyTypeExecution)),
      catchError((err) => {
        this.toast.showError(err);
        return of([]);
      }),
    );
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

  public async openDialog(execution?: CorrectlyTypedExecution): Promise<void> {
    const request$ = this.dialog
      .open<ActionTwoAddActionDialogComponent, ActionTwoAddActionDialogData, ActionTwoAddActionDialogResult>(
        ActionTwoAddActionDialogComponent,
        {
          width: '500px',
          data: execution
            ? {
                execution,
              }
            : {},
        },
      )
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef));

    const request = await lastValueFrom(request$);
    if (!request) {
      return;
    }

    try {
      await this.actionService.setExecution(request);
      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);
    } catch (error) {
      console.error(error);
      this.toast.showError(error);
    }
  }

  public async deleteExecution(execution: CorrectlyTypedExecution) {
    const deleteReq: MessageInitShape<typeof SetExecutionRequestSchema> = {
      condition: execution.condition,
      targets: [],
    };
    try {
      await this.actionService.setExecution(deleteReq);
      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);
    } catch (error) {
      console.error(error);
      this.toast.showError(error);
    }
  }

  protected readonly InfoSectionType = InfoSectionType;
}
