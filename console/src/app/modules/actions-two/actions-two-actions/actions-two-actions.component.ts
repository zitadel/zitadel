import { ChangeDetectionStrategy, Component, DestroyRef, OnInit } from '@angular/core';
import { ActionService } from 'src/app/services/action.service';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { defer, firstValueFrom, lastValueFrom, Observable, of, shareReplay, Subject, TimeoutError } from 'rxjs';
import { catchError, map, startWith, switchMap, timeout } from 'rxjs/operators';
import { ToastService } from 'src/app/services/toast.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { ORGANIZATIONS } from '../../settings-list/settings';
import {
  ActionTwoAddActionDialogComponent,
  ActionTwoAddActionDialogData,
  ActionTwoAddActionDialogResult,
  correctlyTypeExecution,
} from '../actions-two-add-action/actions-two-add-action-dialog.component';
import { MatDialog } from '@angular/material/dialog';
import { MessageInitShape } from '@bufbuild/protobuf';
import { Execution } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';

@Component({
  selector: 'cnsl-actions-two-actions',
  templateUrl: './actions-two-actions.component.html',
  styleUrls: ['./actions-two-actions.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoActionsComponent implements OnInit {
  protected readonly refresh$ = new Subject<true>();
  private readonly actionsEnabled$: Observable<boolean>;
  protected readonly executions$: Observable<Execution[]>;
  protected readonly targets$: Observable<Target[]>;

  constructor(
    private readonly actionService: ActionService,
    private readonly featureService: NewFeatureService,
    private readonly toast: ToastService,
    private readonly destroyRef: DestroyRef,
    private readonly router: Router,
    private readonly route: ActivatedRoute,
    private readonly dialog: MatDialog,
  ) {
    this.actionsEnabled$ = this.getActionsEnabled$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.executions$ = this.getExecutions$(this.actionsEnabled$);
    this.targets$ = this.getTargets$(this.actionsEnabled$);
  }

  ngOnInit(): void {
    // this also preloads
    this.actionsEnabled$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe(async (enabled) => {
      if (enabled) {
        return;
      }
      await this.router.navigate([], {
        relativeTo: this.route,
        queryParams: {
          id: ORGANIZATIONS.id,
        },
        queryParamsHandling: 'merge',
      });
    });
  }

  private getExecutions$(actionsEnabled$: Observable<boolean>) {
    return this.refresh$.pipe(
      startWith(true),
      switchMap(() => {
        return this.actionService.listExecutions({});
      }),
      map(({ result }) => result),
      catchError(async (err) => {
        const actionsEnabled = await firstValueFrom(actionsEnabled$);
        if (actionsEnabled) {
          this.toast.showError(err);
        }
        return [];
      }),
    );
  }

  private getTargets$(actionsEnabled$: Observable<boolean>) {
    return this.refresh$.pipe(
      startWith(true),
      switchMap(() => {
        return this.actionService.listTargets({});
      }),
      map(({ result }) => result),
      catchError(async (err) => {
        const actionsEnabled = await firstValueFrom(actionsEnabled$);
        if (actionsEnabled) {
          this.toast.showError(err);
        }
        return [];
      }),
    );
  }

  private getActionsEnabled$() {
    return defer(() => this.featureService.getInstanceFeatures()).pipe(
      map(({ actions }) => actions?.enabled ?? false),
      timeout(1000),
      catchError((err) => {
        if (!(err instanceof TimeoutError)) {
          this.toast.showError(err);
        }
        return of(false);
      }),
    );
  }

  public async openDialog(execution?: Execution): Promise<void> {
    const request$ = this.dialog
      .open<ActionTwoAddActionDialogComponent, ActionTwoAddActionDialogData, ActionTwoAddActionDialogResult>(
        ActionTwoAddActionDialogComponent,
        {
          width: '400px',
          data: execution
            ? {
                execution: correctlyTypeExecution(execution),
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
      console.log(request);
      await this.actionService.setExecution(request);
      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);
    } catch (error) {
      console.error(error);
      this.toast.showError(error);
    }
  }

  public async deleteExecution(execution: Execution) {
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
}
