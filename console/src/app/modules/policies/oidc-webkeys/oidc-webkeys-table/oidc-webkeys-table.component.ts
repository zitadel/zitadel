import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { GetWebKey, WebKeyState } from '@zitadel/proto/zitadel/resources/webkey/v3alpha/key_pb';
import { ReplaySubject } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { MatTableDataSource } from '@angular/material/table';

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
  public readonly delete = new EventEmitter<GetWebKey>();

  @Input({ required: true })
  public set webKeys(webKeys: GetWebKey[] | null) {
    this.webKeys$.next(webKeys);
  }

  private readonly webKeys$ = new ReplaySubject<GetWebKey[] | null>(1);
  protected readonly dataSource$ = this.webKeys$.pipe(
    filter(Boolean),
    map((keys) => new MatTableDataSource(keys)),
  );
  protected readonly WebKeyState = WebKeyState;
}
