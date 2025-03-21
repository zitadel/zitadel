import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { GetExecution } from '@zitadel/proto/zitadel/resources/action/v3alpha/execution_pb';

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
  public readonly delete = new EventEmitter<GetExecution>();

  @Input({ required: true })
  public set executions(executions: GetExecution[] | null) {
    this.executions$.next(executions);
  }

  private readonly executions$ = new ReplaySubject<GetExecution[] | null>(1);
  protected readonly dataSource$ = this.executions$.pipe(
    filter(Boolean),
    map((keys) => new MatTableDataSource(keys)),
  );
}
