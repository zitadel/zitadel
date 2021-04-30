import { Component, OnInit } from '@angular/core';
import { MatButtonToggleChange } from '@angular/material/button-toggle';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { PolicyGridType } from 'src/app/modules/policy-grid/policy-grid.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Features } from 'src/app/proto/generated/zitadel/features_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Domain, Org, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddDomainDialogComponent } from './add-domain-dialog/add-domain-dialog.component';
import { DomainVerificationComponent } from './domain-verification/domain-verification.component';


@Component({
  selector: 'app-org-detail',
  templateUrl: './org-detail.component.html',
  styleUrls: ['./org-detail.component.scss'],
})
export class OrgDetailComponent implements OnInit {
  public org!: Org.AsObject;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public OrgState: any = OrgState;
  public ChangeType: any = ChangeType;

  public domains: Domain.AsObject[] = [];
  public primaryDomain: string = '';

  // members
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public totalMemberResult: number = 0;
  public membersSubject: BehaviorSubject<Member.AsObject[]>
    = new BehaviorSubject<Member.AsObject[]>([]);
  public PolicyGridType: any = PolicyGridType;

  public features!: Features.AsObject;

  constructor(
    private dialog: MatDialog,
    public translate: TranslateService,
    public mgmtService: ManagementService,
    private toast: ToastService,
    private router: Router,
  ) { }

  public ngOnInit(): void {
    this.getData();
  }

  private async getData(): Promise<void> {
    this.mgmtService.getMyOrg().then((resp) => {
      if (resp.org) {
        this.org = resp.org;
      }
    }).catch(error => {
      this.toast.showError(error);
    });
    this.loadMembers();
    this.loadDomains();
    this.loadFeatures();
  }

  public loadDomains(): void {
    this.mgmtService.listOrgDomains().then(result => {
      this.domains = result.resultList;
      this.primaryDomain = this.domains.find(domain => domain.isPrimary)?.domainName ?? '';
    });
  }

  public setPrimary(domain: Domain.AsObject): void {
    this.mgmtService.setPrimaryOrgDomain(domain.domainName).then(() => {
      this.toast.showInfo('ORG.TOAST.SETPRIMARY', true);
      this.loadDomains();
    }).catch((error) => {
      this.toast.showError(error);
    });
  }

  public changeState(event: MatButtonToggleChange | any): void {
    if (event.value === OrgState.ORG_STATE_ACTIVE) {
      this.mgmtService.reactivateOrg().then(() => {
        this.toast.showInfo('ORG.TOAST.REACTIVATED', true);
      }).catch((error) => {
        this.toast.showError(error);
      });
    } else if (event.value === OrgState.ORG_STATE_INACTIVE) {
      this.mgmtService.deactivateOrg().then(() => {
        this.toast.showInfo('ORG.TOAST.DEACTIVATED', true);
      }).catch((error) => {
        this.toast.showError(error);
      });
    }
  }

  public addNewDomain(): void {
    const dialogRef = this.dialog.open(AddDomainDialogComponent, {
      data: {},
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {
        this.mgmtService.addOrgDomain(resp).then(() => {
          this.toast.showInfo('ORG.TOAST.DOMAINADDED', true);

          setTimeout(() => {
            this.loadDomains();
          }, 1000);
        }).catch(error => {
          this.toast.showError(error);
        });
      }
    });
  }

  public removeDomain(domain: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'ORG.DOMAINS.DELETE.TITLE',
        descriptionKey: 'ORG.DOMAINS.DELETE.DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {
        this.mgmtService.removeOrgDomain(domain).then(() => {
          this.toast.showInfo('ORG.TOAST.DOMAINREMOVED', true);
          const index = this.domains.findIndex(d => d.domainName === domain);
          if (index > -1) {
            this.domains.splice(index, 1);
          }
        }).catch(error => {
          this.toast.showError(error);
        });
      }
    });
  }

  public openAddMember(): void {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.ORG,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {
        const users: User.AsObject[] = resp.users;
        const roles: string[] = resp.roles;

        if (users && users.length && roles && roles.length) {
          Promise.all(users.map(user => {
            return this.mgmtService.addOrgMember(user.id, roles);
          })).then(() => {
            this.toast.showInfo('ORG.TOAST.MEMBERADDED', true);
            setTimeout(() => {
              this.loadMembers();
            }, 1000);
          }).catch(error => {
            this.toast.showError(error);
          });
        }
      }
    });
  }

  public showDetail(): void {
    this.router.navigate(['org/members']);
  }

  public verifyDomain(domain: Domain.AsObject): void {
    const dialogRef = this.dialog.open(DomainVerificationComponent, {
      data: {
        domain: domain,
      },
      width: '500px',
    });

    dialogRef.afterClosed().subscribe((reload) => {
      if (reload) {
        this.loadDomains();
      }
    });
  }

  public loadMembers(): void {
    this.loadingSubject.next(true);
    from(this.mgmtService.listOrgMembers(100, 0)).pipe(
      map(resp => {
        if (resp.details?.totalResult) {
          this.totalMemberResult = resp.details?.totalResult;
        } else {
          this.totalMemberResult = 0;
        }

        return resp.resultList;
      }),
      catchError(() => of([])),
      finalize(() => this.loadingSubject.next(false)),
    ).subscribe(members => {
      this.membersSubject.next(members);
    });
  }

  public loadFeatures(): void {
    this.loadingSubject.next(true);
    this.mgmtService.getFeatures().then(resp => {
      if (resp.features) {
        this.features = resp.features;
      }
    });
  }
}
