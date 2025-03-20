import { ChangeDetectionStrategy, Component, DestroyRef, OnInit } from '@angular/core';
import { ActionService } from 'src/app/services/action.service';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { defer, firstValueFrom, Observable, of, shareReplay, TimeoutError } from 'rxjs';
import { catchError, map, timeout } from 'rxjs/operators';
import { ToastService } from 'src/app/services/toast.service';
import { GetExecution } from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { ORGANIZATIONS } from '../../settings-list/settings';
import { ActionTwoAddActionDialogComponent } from '../actions-two-add-action/actions-two-add-action-dialog.component';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'cnsl-actions-two-actions',
  templateUrl: './actions-two-actions.component.html',
  styleUrls: ['./actions-two-actions.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoActionsComponent implements OnInit {
  private readonly actionsEnabled$: Observable<boolean>;
  protected readonly executions$: Observable<GetExecution[]>;

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
    return defer(() => this.actionService.searchExections({})).pipe(
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
    const ref = this.dialog.open(ActionTwoAddActionDialogComponent, {
      width: '400px',
      data: {},
    });

    ref.afterClosed().subscribe((resp) => {
      if (resp) {
        this.actionService.setExecution({});
      }
    });
  }
}
