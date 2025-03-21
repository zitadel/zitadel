import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { GetTarget } from '@zitadel/proto/zitadel/resources/action/v3alpha/target_pb';

@Component({
  selector: 'cnsl-actions-two-targets-table',
  templateUrl: './actions-two-targets-table.component.html',
  styleUrls: ['./actions-two-targets-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ActionsTwoTargetsTableComponent {
  @Output()
  public readonly refresh = new EventEmitter<void>();

  @Output()
  public readonly delete = new EventEmitter<GetTarget>();

  @Input({ required: true })
  public set targets(targets: GetTarget[] | null) {
    this.targets$.next(targets);
  }

  private readonly targets$ = new ReplaySubject<GetTarget[] | null>(1);
  protected readonly dataSource$ = this.targets$.pipe(
    filter(Boolean),
    map((keys) => new MatTableDataSource(keys)),
  );
}
