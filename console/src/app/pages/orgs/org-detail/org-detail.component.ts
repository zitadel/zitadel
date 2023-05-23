import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { Router } from '@angular/router';
import { Buffer } from 'buffer';
import { BehaviorSubject, from, Observable, of, Subject, takeUntil } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { MetadataDialogComponent } from 'src/app/modules/metadata/metadata-dialog/metadata-dialog.component';
import { NameDialogComponent } from 'src/app/modules/name-dialog/name-dialog.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { Org, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-org-detail',
  templateUrl: './org-detail.component.html',
  styleUrls: ['./org-detail.component.scss'],
})
export class OrgDetailComponent implements OnInit, OnDestroy {
  public org?: Org.AsObject;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public OrgState: any = OrgState;
  public ChangeType: any = ChangeType;

  public metadata: Metadata.AsObject[] = [];
  public loadingMetadata: boolean = true;

  // members
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public totalMemberResult: number = 0;
  public membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);
  private destroy$: Subject<void> = new Subject();

  public InfoSectionType: any = InfoSectionType;

  constructor(
    private auth: GrpcAuthService,
    private dialog: MatDialog,
    public mgmtService: ManagementService,
    private toast: ToastService,
    private router: Router,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);

    auth.activeOrgChanged.pipe(takeUntil(this.destroy$)).subscribe((org) => {
      if (this.org && org) {
        this.getData();
        this.loadMetadata();
      }
    });
  }

  public ngOnInit(): void {
    this.getData();
    this.loadMetadata();
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public changeState(newState: OrgState): void {
    if (newState === OrgState.ORG_STATE_ACTIVE) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.REACTIVATE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'ORG.DIALOG.REACTIVATE.TITLE',
          descriptionKey: 'ORG.DIALOG.REACTIVATE.DESCRIPTION',
        },
        width: '400px',
      });
      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.mgmtService
            .reactivateOrg()
            .then(() => {
              this.toast.showInfo('ORG.TOAST.REACTIVATED', true);
              this.org!.state = OrgState.ORG_STATE_ACTIVE;
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    } else if (newState === OrgState.ORG_STATE_INACTIVE) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.DEACTIVATE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'ORG.DIALOG.DEACTIVATE.TITLE',
          descriptionKey: 'ORG.DIALOG.DEACTIVATE.DESCRIPTION',
        },
        width: '400px',
      });
      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.mgmtService
            .deactivateOrg()
            .then(() => {
              this.toast.showInfo('ORG.TOAST.DEACTIVATED', true);
              this.org!.state = OrgState.ORG_STATE_INACTIVE;
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }

  public deleteOrg(): void {
    const mgmtUserData = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'ORG.DIALOG.DELETE.TITLE',
      warnSectionKey: 'ORG.DIALOG.DELETE.DESCRIPTION',
      hintKey: 'ORG.DIALOG.DELETE.TYPENAME',
      hintParam: 'ORG.DIALOG.DELETE.DESCRIPTION',
      confirmationKey: 'ORG.DIALOG.DELETE.ORGNAME',
      confirmation: this.org?.name,
    };

    if (this.org) {
      let dialogRef;

      dialogRef = this.dialog.open(WarnDialogComponent, {
        data: mgmtUserData,
        width: '400px',
      });

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.mgmtService
            .removeOrg()
            .then(() => {
              setTimeout(() => {
                this.router.navigate(['/orgs']);
              }, 1000);
              this.toast.showInfo('ORG.TOAST.DELETED', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }

  private async getData(): Promise<void> {
    this.mgmtService
      .getMyOrg()
      .then((resp) => {
        if (resp.org) {
          this.org = resp.org;
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
    this.loadMembers();
  }

  public openAddMember(): void {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.ORG,
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
              return this.mgmtService.addOrgMember(user.id, roles);
            }),
          )
            .then(() => {
              this.toast.showInfo('ORG.TOAST.MEMBERADDED', true);
              setTimeout(() => {
                this.loadMembers();
              }, 1000);
            })
            .catch((error) => {
              setTimeout(() => {
                this.loadMembers();
              }, 1000);
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public showDetail(): void {
    this.router.navigate(['org/members']);
  }

  public loadMembers(): void {
    this.loadingSubject.next(true);
    from(this.mgmtService.listOrgMembers(100, 0))
      .pipe(
        map((resp) => {
          if (resp.details?.totalResult) {
            this.totalMemberResult = resp.details?.totalResult;
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

  public loadMetadata(): Promise<any> | void {
    this.loadingMetadata = true;
    return this.mgmtService
      .listOrgMetadata()
      .then((resp) => {
        this.loadingMetadata = false;
        this.metadata = resp.resultList.map((md) => {
          return {
            key: md.key,
            value: Buffer.from(md.value as string, 'base64').toString('ascii'),
          };
        });
      })
      .catch((error) => {
        this.loadingMetadata = false;
        this.toast.showError(error);
      });
  }

  public editMetadata(): void {
    const setFcn = (key: string, value: string): Promise<any> => this.mgmtService.setOrgMetadata(key, btoa(value));
    const removeFcn = (key: string): Promise<any> => this.mgmtService.removeOrgMetadata(key);

    const dialogRef = this.dialog.open(MetadataDialogComponent, {
      data: {
        metadata: this.metadata,
        setFcn: setFcn,
        removeFcn: removeFcn,
      },
    });

    dialogRef.afterClosed().subscribe(() => {
      this.loadMetadata();
    });
  }

  public renameOrg(): void {
    const dialogRef = this.dialog.open(NameDialogComponent, {
      data: {
        name: this.org?.name,
        titleKey: 'ORG.PAGES.RENAME.TITLE',
        descKey: 'ORG.PAGES.RENAME.DESCRIPTION',
        labelKey: 'ORG.PAGES.NAME',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((name) => {
      if (name) {
        this.updateOrg(name);
      }
    });
  }

  public updateOrg(name: string): void {
    if (this.org) {
      this.mgmtService
        .updateOrg(name)
        .then(() => {
          this.toast.showInfo('ORG.TOAST.UPDATED', true);
          if (this.org) {
            this.org.name = name;
          }
          this.mgmtService
            .getMyOrg()
            .then((resp) => {
              if (resp.org) {
                this.org = resp.org;
                this.auth.setActiveOrg(resp.org);
              }
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }
}
