import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { GetWebKey } from '@zitadel/proto/zitadel/resources/webkey/v3alpha/key_pb';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';

@Component({
  selector: 'cnsl-oidc-webkeys-inactive-table',
  templateUrl: './oidc-webkeys-inactive-table.component.html',
  styleUrls: ['./oidc-webkeys-inactive-table.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OidcWebKeysInactiveTableComponent {
  @Input({ required: true })
  public set InactiveWebKeys(webKeys: GetWebKey[] | null) {
    this.inactiveWebKeys$.next(webKeys);
  }

  private inactiveWebKeys$ = new ReplaySubject<GetWebKey[] | null>(1);
  protected dataSource$ = this.inactiveWebKeys$.pipe(
    filter(Boolean),
    map((webKeys) => new MatTableDataSource(webKeys)),
  );
}
