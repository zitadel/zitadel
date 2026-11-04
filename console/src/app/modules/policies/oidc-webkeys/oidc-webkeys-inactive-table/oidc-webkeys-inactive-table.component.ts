import { ChangeDetectionStrategy, Component, effect, input, output } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { WebKey } from '@zitadel/proto/zitadel/webkey/v2/key_pb';

@Component({
  selector: 'cnsl-oidc-webkeys-inactive-table',
  templateUrl: './oidc-webkeys-inactive-table.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false,
})
export class OidcWebKeysInactiveTableComponent {
  public readonly delete = output<WebKey>();

  public inactiveWebKeys = input.required<WebKey[] | null>();

  protected dataSource = new MatTableDataSource<WebKey>([]);

  constructor() {
    effect(() => {
      const inactiveWebKeys = this.inactiveWebKeys();
      if (inactiveWebKeys) {
        this.dataSource.data = inactiveWebKeys;
      }
    });
  }
}
