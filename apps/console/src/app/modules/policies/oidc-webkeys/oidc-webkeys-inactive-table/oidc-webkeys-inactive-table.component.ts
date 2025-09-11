import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';
import { WebKey } from '@zitadel/proto/zitadel/webkey/v2beta/key_pb';

@Component({
  selector: 'cnsl-oidc-webkeys-inactive-table',
  templateUrl: './oidc-webkeys-inactive-table.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OidcWebKeysInactiveTableComponent {
  @Input({ required: true })
  public set inactiveWebKeys(webKeys: WebKey[] | null) {
    this.inactiveWebKeys$.next(webKeys);
  }

  private inactiveWebKeys$ = new ReplaySubject<WebKey[] | null>(1);
  protected dataSource$ = this.inactiveWebKeys$.pipe(
    filter(Boolean),
    map((webKeys) => new MatTableDataSource(webKeys)),
  );
}
