import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { State, WebKey } from '@zitadel/proto/zitadel/webkey/v2beta/key_pb';

@Component({
  selector: 'cnsl-oidc-webkeys-table',
  templateUrl: './oidc-webkeys-table.component.html',
  styleUrls: ['./oidc-webkeys-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OidcWebKeysTableComponent {
  @Output()
  public readonly refresh = new EventEmitter<void>();

  @Output()
  public readonly delete = new EventEmitter<WebKey>();

  @Input({ required: true })
  public set webKeys(webKeys: WebKey[] | null) {
    this.webKeys$.next(webKeys);
  }

  private readonly webKeys$ = new ReplaySubject<WebKey[] | null>(1);
  protected readonly dataSource$ = this.webKeys$.pipe(
    filter(Boolean),
    map((keys) => new MatTableDataSource(keys)),
  );
  protected readonly State = State;
}
