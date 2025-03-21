import { ChangeDetectionStrategy, Component, DestroyRef, OnInit } from '@angular/core';
import { defer, firstValueFrom, Observable, of, shareReplay, Subject, TimeoutError } from 'rxjs';
import { ActionService } from 'src/app/services/action.service';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { ToastService } from 'src/app/services/toast.service';
import { ActivatedRoute, Router } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ORGANIZATIONS } from '../../settings-list/settings';
import { catchError, map, tap, timeout } from 'rxjs/operators';
import { GetTarget } from '@zitadel/proto/zitadel/resources/action/v3alpha/target_pb';
import { MatDialog } from '@angular/material/dialog';
import { ActionTwoAddTargetDialogComponent } from '../actions-two-add-target/actions-two-add-target-dialog.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import { CreateTargetRequestSchema } from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';

@Component({
  selector: 'cnsl-actions-two-targets',
  templateUrl: './actions-two-targets.component.html',
  styleUrls: ['./actions-two-targets.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoTargetsComponent implements OnInit {
  private readonly actionsEnabled$: Observable<boolean>;
  protected readonly targets$: Observable<GetTarget[]>;
  protected readonly refresh = new Subject<true>();

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

  private getTargets$(actionsEnabled$: Observable<boolean>) {
    return defer(() => this.actionService.searchTargets({})).pipe(
      map(({ result }) => result),
      tap((targets) => {
        console.log(targets);
      }),
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
    const ref = this.dialog.open(ActionTwoAddTargetDialogComponent, {
      width: '400px',
      data: {},
    });

    ref.afterClosed().subscribe((resp: MessageInitShape<typeof CreateTargetRequestSchema>) => {
      if (resp) {
        this.actionService.createTarget(resp);
      }
    });
  }

  public deleteTarget(target: GetTarget) {}
}
