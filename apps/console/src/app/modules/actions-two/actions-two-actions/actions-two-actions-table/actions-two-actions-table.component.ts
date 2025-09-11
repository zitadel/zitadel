import { ChangeDetectionStrategy, Component, computed, effect, EventEmitter, Input, Output } from '@angular/core';
import { combineLatestWith, Observable, ReplaySubject } from 'rxjs';
import { filter, map, startWith } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import { toSignal } from '@angular/core/rxjs-interop';
import { CorrectlyTypedExecution } from '../../actions-two-add-action/actions-two-add-action-dialog.component';

@Component({
  selector: 'cnsl-actions-two-actions-table',
  templateUrl: './actions-two-actions-table.component.html',
  styleUrls: ['./actions-two-actions-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoActionsTableComponent {
  @Output()
  public readonly refresh = new EventEmitter<void>();

  @Output()
  public readonly selected = new EventEmitter<CorrectlyTypedExecution>();

  @Output()
  public readonly delete = new EventEmitter<CorrectlyTypedExecution>();

  @Input({ required: true })
  public set executions(executions: CorrectlyTypedExecution[] | null) {
    this.executions$.next(executions);
  }

  @Input({ required: true })
  public set targets(targets: Target[] | null) {
    this.targets$.next(targets);
  }

  private readonly executions$ = new ReplaySubject<CorrectlyTypedExecution[] | null>(1);

  private readonly targets$ = new ReplaySubject<Target[] | null>(1);

  protected readonly dataSource = this.getDataSource();

  protected readonly loading = this.getLoading();

  private getDataSource() {
    const executions$: Observable<CorrectlyTypedExecution[]> = this.executions$.pipe(filter(Boolean), startWith([]));
    const executionsSignal = toSignal(executions$, { requireSync: true });

    const targetsMapSignal = this.getTargetsMap();

    const dataSignal = computed(() => {
      const executions = executionsSignal();
      const targetsMap = targetsMapSignal();

      if (targetsMap.size === 0) {
        return [];
      }

      return executions.map((execution) => {
        const mappedTargets = execution.targets
          .map((target) => targetsMap.get(target))
          .filter((target): target is NonNullable<typeof target> => !!target);
        return { execution, mappedTargets };
      });
    });

    const dataSource = new MatTableDataSource(dataSignal());

    effect(() => {
      const data = dataSignal();
      if (dataSource.data !== data) {
        dataSource.data = data;
      }
    });

    return dataSource;
  }

  private getTargetsMap() {
    const targets$ = this.targets$.pipe(filter(Boolean), startWith([] as Target[]));
    const targetsSignal = toSignal(targets$, { requireSync: true });

    return computed(() => {
      const map = new Map<string, Target>();
      for (const target of targetsSignal()) {
        map.set(target.id, target);
      }
      return map;
    });
  }

  private getLoading() {
    const loading$ = this.executions$.pipe(
      combineLatestWith(this.targets$),
      map(([executions, targets]) => executions === null || targets === null),
      startWith(true),
    );

    return toSignal(loading$, { requireSync: true });
  }

  protected trackTarget(_: number, target: Target) {
    return target.id;
  }
}
