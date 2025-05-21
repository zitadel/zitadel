import { Component, DestroyRef } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { BehaviorSubject, defer, from, Observable, of, shareReplay, TimeoutError } from 'rxjs';
import { catchError, finalize, map, timeout } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { InstanceDetail, State } from 'src/app/proto/generated/zitadel/instance_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import {
  BRANDING,
  COMPLEXITY,
  DOMAIN,
  LANGUAGES,
  IDP,
  LOCKOUT,
  AGE,
  LOGIN,
  LOGINTEXTS,
  MESSAGETEXTS,
  NOTIFICATIONS,
  OIDC,
  WEBKEYS,
  PRIVACYPOLICY,
  SECRETS,
  SECURITY,
  SMS_PROVIDER,
  SMTP_PROVIDER,
  VIEWS,
  FAILEDEVENTS,
  EVENTS,
  ORGANIZATIONS,
  FEATURESETTINGS,
  ACTIONS,
  ACTIONS_TARGETS,
} from 'src/app/modules/settings-list/settings';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { EnvironmentService } from 'src/app/services/environment.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
@Component({
  selector: 'cnsl-instance',
  templateUrl: './instance.component.html',
  styleUrls: ['./instance.component.scss'],
})
export class InstanceComponent {
  protected instance?: InstanceDetail.AsObject;
  protected readonly PolicyComponentServiceType = PolicyComponentServiceType;
  private readonly loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  protected readonly loading$: Observable<boolean> = this.loadingSubject.asObservable();
  protected totalMemberResult: number = 0;
  protected readonly membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);
  protected readonly State = State;

  protected id: string = '';
  protected readonly defaultSettingsList: SidenavSetting[] = [
    ORGANIZATIONS,
    FEATURESETTINGS,
    ACTIONS,
    ACTIONS_TARGETS,
    // notifications
    // { showWarn: true, ...NOTIFICATIONS },
    NOTIFICATIONS,
    SMTP_PROVIDER,
    SMS_PROVIDER,
    // login
    LOGIN,
    IDP,
    COMPLEXITY,
    AGE,
    LOCKOUT,
    DOMAIN,
    // appearance
    BRANDING,
    MESSAGETEXTS,
    LOGINTEXTS,
    // storage
    VIEWS,
    EVENTS,
    FAILEDEVENTS,
    // others
    PRIVACYPOLICY,
    LANGUAGES,
    OIDC,
    WEBKEYS,
    SECRETS,
    SECURITY,
  ];

  protected readonly settingsList$: Observable<SidenavSetting[]>;
  protected readonly customerPortalLink$ = this.envService.env.pipe(map((env) => env.customer_portal));

  constructor(
    protected readonly adminService: AdminService,
    private readonly dialog: MatDialog,
    private readonly toast: ToastService,
    breadcrumbService: BreadcrumbService,
    private readonly router: Router,
    private readonly authService: GrpcAuthService,
    private readonly envService: EnvironmentService,
    activatedRoute: ActivatedRoute,
    private readonly destroyRef: DestroyRef,
  ) {
    this.loadMembers();

    const instanceBread = new Breadcrumb({
      type: BreadcrumbType.INSTANCE,
      name: 'Instance',
      routerLink: ['/instance'],
    });
    breadcrumbService.setBreadcrumb([instanceBread]);

    this.adminService
      .getMyInstance()
      .then((instanceResp) => {
        if (instanceResp.instance) {
          this.instance = instanceResp.instance;
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });

    activatedRoute.queryParamMap.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((params) => {
      const id = params.get('id');
      if (id) {
        this.id = id;
      }
    });

    this.settingsList$ = this.getSettingsList();
  }

  private getSettingsList(): Observable<SidenavSetting[]> {
    return this.authService.isAllowedMapper(this.defaultSettingsList, (setting) => setting.requiredRoles.admin || []);
  }

  public loadMembers(): void {
    this.loadingSubject.next(true);
    from(this.adminService.listIAMMembers(100, 0))
      .pipe(
        map((resp) => {
          if (resp.details?.totalResult) {
            this.totalMemberResult = resp.details.totalResult;
          } else {
            this.totalMemberResult = 0;
          }
          return resp.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((members) => {
        this.membersSubject.next(members);
      });
  }

  public openAddMember(): void {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.IAM,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const users: User.AsObject[] = resp.users;
        const roles: string[] = resp.roles;

        if (users && users.length && roles && roles.length) {
          Promise.all(
            users.map((user) => {
              return this.adminService.addIAMMember(user.id, roles);
            }),
          )
            .then(() => {
              this.toast.showInfo('IAM.TOAST.MEMBERADDED', true);
              setTimeout(() => {
                this.loadMembers();
              }, 1000);
            })
            .catch((error) => {
              this.toast.showError(error);
              setTimeout(() => {
                this.loadMembers();
              }, 1000);
            });
        }
      }
    });
  }

  public showDetail(): void {
    this.router.navigate(['/instance', 'members']).then();
  }
}
