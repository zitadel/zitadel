import { ChangeDetectionStrategy, Component, computed, effect, EventEmitter, Input, Output } from '@angular/core';
import { Observable, ReplaySubject, shareReplay } from 'rxjs';
import { filter, map, startWith } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { Execution, ExecutionTargetType } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import { toSignal } from '@angular/core/rxjs-interop';

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
  public readonly delete = new EventEmitter<Execution>();

  @Input({ required: true })
  public set executions(executions: Execution[] | null) {
    this.executions$.next(executions);
  }

  @Input({ required: true })
  public set targets(targets: Target[] | null) {
    this.targets$.next(targets);
  }

  @Output()
  public readonly selected = new EventEmitter<Execution>();

  protected readonly executions$ = new ReplaySubject<Execution[] | null>(1);

  private readonly targets$ = new ReplaySubject<Target[] | null>(1);
  private readonly targetsMap = this.getTargetsMap();

  protected readonly dataSource = this.getDataSource();

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

  private getDataSource() {
    const executions$: Observable<Execution[]> = this.executions$.pipe(filter(Boolean), startWith([]));
    const executionsSignal = toSignal(executions$, { requireSync: true });

    const dataSource = new MatTableDataSource(executionsSignal());

    effect(() => {
      const executions = executionsSignal();
      if (dataSource.data !== executions) {
        dataSource.data = executions;
      }
    });

    return dataSource;
  }

  protected filteredTargetTypes(targets: ExecutionTargetType[]): Target[] {
    return targets
      .map((t) => t.type)
      .filter((t): t is Extract<ExecutionTargetType['type'], { case: 'target' }> => t.case === 'target')
      .map((t) => this.targetsMap().get(t.value))
      .filter((target): target is Target => !!target);
  }

  protected trackTarget(_: number, target: Target) {
    return target.id;
  }
}
