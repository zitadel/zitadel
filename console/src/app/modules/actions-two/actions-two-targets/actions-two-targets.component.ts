import { ChangeDetectionStrategy, Component, DestroyRef, OnInit } from '@angular/core';
import { defer, firstValueFrom, merge, Observable, of, shareReplay, Subject, TimeoutError } from 'rxjs';
import { ActionService } from 'src/app/services/action.service';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { ToastService } from 'src/app/services/toast.service';
import { ActivatedRoute, Router } from '@angular/router';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ORGANIZATIONS } from '../../settings-list/settings';
import { catchError, map, startWith, switchMap, tap, timeout } from 'rxjs/operators';
import { GetTarget } from '@zitadel/proto/zitadel/resources/action/v3alpha/target_pb';
import { MatDialog } from '@angular/material/dialog';
import { ActionTwoAddTargetDialogComponent } from '../actions-two-add-target/actions-two-add-target-dialog.component';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  CreateTargetRequestSchema,
  PatchTargetRequestSchema,
} from '@zitadel/proto/zitadel/resources/action/v3alpha/action_service_pb';

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
    return this.refresh.pipe(
      startWith(true),
      switchMap(() => {
        return this.actionService.searchTargets({});
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

  public async deleteTarget(target: GetTarget) {
    this.actionService.deleteTarget({ id: target.details?.id });
    setTimeout(() => {
      this.refresh.next(true);
    }, 1000);
  }

  public openDialog(target?: GetTarget): void {
    const ref = this.dialog.open(ActionTwoAddTargetDialogComponent, {
      width: '550px',
      data: target
        ? {
            target: target,
          }
        : {},
    });

    ref.afterClosed().subscribe((dialogResponse) => {
      if (target?.details?.id && dialogResponse) {
        const req: MessageInitShape<typeof PatchTargetRequestSchema> = dialogResponse;

        this.actionService.patchTarget({ ...req, id: target.details?.id });
        setTimeout(() => {
          this.refresh.next(true);
        }, 1000);
      }
      if (dialogResponse) {
        const req: MessageInitShape<typeof CreateTargetRequestSchema> = dialogResponse;

        this.actionService.createTarget(req);
        setTimeout(() => {
          this.refresh.next(true);
        }, 1000);
      }
    });
  }
}
