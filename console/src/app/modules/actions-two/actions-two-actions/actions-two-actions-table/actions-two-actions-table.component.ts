import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { Execution } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';

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

  private readonly executions$ = new ReplaySubject<Execution[] | null>(1);
  protected readonly dataSource$ = this.executions$.pipe(
    filter(Boolean),
    map((keys) => new MatTableDataSource(keys)),
  );
}
