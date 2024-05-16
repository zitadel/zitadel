import { ConnectedPosition, ConnectionPositionPair } from '@angular/cdk/overlay';
import { Component, ElementRef, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, of, Subject } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ActionKeysType } from '../action-keys/action-keys.component';
import { GetPrivacyPolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';

@Component({
  selector: 'cnsl-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent implements OnDestroy {
  @ViewChild('input', { static: false }) input!: ElementRef;

  @Input() public isDarkTheme: boolean = true;
  @Input() public user?: User.AsObject;
  public showOrgContext: boolean = false;

  public orgs$: Observable<Org.AsObject[]> = of([]);
  @Input() public org!: Org.AsObject;
  @Output() public changedActiveOrg: EventEmitter<Org.AsObject> = new EventEmitter();
  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public showAccount: boolean = false;
  private destroy$: Subject<void> = new Subject();
  public BreadcrumbType: any = BreadcrumbType;
  public ActionKeysType: any = ActionKeysType;
  public docsLink = 'https://zitadel.com/docs';
  public customLink = '';
  public customLinkText = '';

  public positions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];

  public accountCardPositions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];
  constructor(
    public authenticationService: AuthenticationService,
    public authService: GrpcAuthService,
    public mgmtService: ManagementService,
    public breadcrumbService: BreadcrumbService,
    public router: Router,
  ) {
    this.loadData();
  }

  public async loadData(): Promise<any> {
    const getData = (): Promise<GetPrivacyPolicyResponse.AsObject> => {
      return this.mgmtService.getPrivacyPolicy();
    };

    getData()
      .then((resp) => {
        if (resp.policy) {
          this.docsLink = resp.policy.docsLink;
          this.customLink = resp.policy.customLink;
          this.customLinkText = resp.policy.customLinkText;
        }
      })
      .catch(() => {});
  }

  public ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public closeAccountCard(): void {
    if (this.showAccount) {
      this.showAccount = false;
    }
  }

  public setActiveOrg(org: Org.AsObject): void {
    this.org = org;
    this.authService.setActiveOrg(org);
    this.changedActiveOrg.emit(org);
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
