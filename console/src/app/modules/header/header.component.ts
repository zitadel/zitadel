import { ConnectedPosition, ConnectionPositionPair } from '@angular/cdk/overlay';
import { Component, EventEmitter, input, Input, linkedSignal, Output } from '@angular/core';
import { Router } from '@angular/router';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ActionKeysType } from '../action-keys/action-keys.component';
import { Organization } from '@zitadel/proto/zitadel/org/v2/org_pb';

@Component({
  selector: 'cnsl-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
  standalone: false,
})
export class HeaderComponent {
  @Input() public isDarkTheme: boolean = true;
  @Input({ required: true }) public user!: User.AsObject;
  public showOrgContext: boolean = false;

  public org = input<Organization | null>();

  @Output() public changedActiveOrg = new EventEmitter<void>();
  public showAccount: boolean = false;
  public BreadcrumbType: any = BreadcrumbType;
  public ActionKeysType: any = ActionKeysType;

  public positions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];

  public accountCardPositions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];
  constructor(
    public authService: GrpcAuthService,
    public mgmtService: ManagementService,
    public breadcrumbService: BreadcrumbService,
    public router: Router,
  ) {}

  public async setActiveOrg(orgId: string): Promise<void> {
    await this.authService.getActiveOrg(orgId);
    this.changedActiveOrg.emit();
  }

  public get isOnMe(): boolean {
    return this.router.url === '/users/me';
  }

  public errorHandler(event: any, fallbackSrc: string) {
    (event.target as HTMLImageElement).src = fallbackSrc;
  }

  public get isOnInstance(): boolean {
    const pages: string[] = [
      '/instance',
      '/settings',
      '/views',
      '/events',
      '/orgs',
      '/settings',
      '/failed-events',
      '/instance/members',
    ];

    return pages.findIndex((p) => this.router.url.includes(p)) > -1;
  }
}
