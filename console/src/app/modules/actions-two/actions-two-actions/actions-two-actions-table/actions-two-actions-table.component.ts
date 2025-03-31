import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { Observable, ReplaySubject } from 'rxjs';
import { filter, map, switchMap } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import {
  Condition,
  EventExecution,
  Execution,
  ExecutionTargetType,
  FunctionExecution,
  RequestExecution,
  ResponseExecution,
} from '@zitadel/proto/zitadel/action/v2beta/execution_pb';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';

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

  private readonly executions$ = new ReplaySubject<Execution[] | null>(1);
  private readonly targets$ = new ReplaySubject<Target[] | null>(1);

  protected readonly dataSource$ = this.executions$.pipe(
    filter(Boolean),
    map((keys) => new MatTableDataSource(keys)),
  );

  protected filteredTargetTypes(targets: ExecutionTargetType[]): Observable<Target[]> {
    const targetIds = targets.filter((t) => t.type.case === 'target').map((t) => t.type.value as string);

    return this.targets$.pipe(
      filter(Boolean),
      map((alltargets) => alltargets!.filter((target) => targetIds.includes(target.id))),
    );
  }
}
