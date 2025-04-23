import { ChangeDetectionStrategy, Component, effect, EventEmitter, Input, Output } from '@angular/core';
import { ReplaySubject } from 'rxjs';
import { filter, startWith } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { Target } from '@zitadel/proto/zitadel/action/v2beta/target_pb';
import { toSignal } from '@angular/core/rxjs-interop';

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
  public readonly selected = new EventEmitter<Target>();

  @Output()
  public readonly delete = new EventEmitter<Target>();

  @Input({ required: true })
  public set targets(targets: Target[] | null) {
    this.targets$.next(targets);
  }

  protected readonly targets$ = new ReplaySubject<Target[] | null>(1);
  protected readonly dataSource: MatTableDataSource<Target>;

  constructor() {
    this.dataSource = this.getDataSource();
  }

  private getDataSource() {
    const targets$ = this.targets$.pipe(filter(Boolean), startWith<Target[]>([]));
    const targetsSignal = toSignal(targets$, { requireSync: true });

    const dataSource = new MatTableDataSource(targetsSignal());
    effect(() => {
      const targets = targetsSignal();
      if (dataSource.data !== targets) {
        dataSource.data = targets;
      }
    });
    return dataSource;
  }
}
