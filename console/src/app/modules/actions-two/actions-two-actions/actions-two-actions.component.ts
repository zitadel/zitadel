import { ChangeDetectionStrategy, Component, DestroyRef, OnInit } from '@angular/core';
import { ActionService } from 'src/app/services/action.service';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { defer, firstValueFrom, Observable, of, shareReplay, Subject, TimeoutError } from 'rxjs';
import { catchError, filter, map, startWith, switchMap, tap, timeout } from 'rxjs/operators';
import { ToastService } from 'src/app/services/toast.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { ORGANIZATIONS } from '../../settings-list/settings';
import { ActionTwoAddActionDialogComponent } from '../actions-two-add-action/actions-two-add-action-dialog.component';
import { MatDialog } from '@angular/material/dialog';
import { MessageInitShape } from '@bufbuild/protobuf';
import { Execution, ExecutionSchema } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { SetExecutionRequestSchema } from '@zitadel/proto/zitadel/action/v2beta/action_service_pb';

@Component({
  selector: 'cnsl-actions-two-actions',
  templateUrl: './actions-two-actions.component.html',
  styleUrls: ['./actions-two-actions.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoActionsComponent implements OnInit {
  protected readonly refresh = new Subject<true>();
  private readonly actionsEnabled$: Observable<boolean>;
  protected readonly executions$: Observable<Execution[]>;

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
    return this.refresh.pipe(
      startWith(true),
      switchMap(() => {
        return this.actionService.listExecutions({});
      }),
      tap(console.log),
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

  public openDialog(): void {
    const ref = this.dialog.open<ActionTwoAddActionDialogComponent>(ActionTwoAddActionDialogComponent, {
      width: '400px',
    });

    ref.afterClosed().subscribe((request?: MessageInitShape<typeof SetExecutionRequestSchema>) => {
      console.log('request', request);
      // if (request) {
      //   this.actionService
      //     .setExecution(request)
      //     .then(() => {
      //       setTimeout(() => {
      //         this.refresh.next(true);
      //       }, 1000);
      //     })
      //     .catch((error) => {
      //       this.toast.showError(error);
      //     });
      // }
    });
  }

  public deleteExecution(execution: Execution) {}
}
