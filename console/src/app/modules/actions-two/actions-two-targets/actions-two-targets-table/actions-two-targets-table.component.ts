import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';

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
  public readonly delete = new EventEmitter<Target>();

  @Input({ required: true })
  public set targets(targets: Target[] | null) {
    this.targets$.next(targets);
  }

  @Output()
  public readonly selected = new EventEmitter<Target>();

  private readonly targets$ = new ReplaySubject<Target[] | null>(1);
  protected readonly dataSource$ = this.targets$.pipe(
    filter(Boolean),
    map((keys) => new MatTableDataSource(keys)),
  );
}
