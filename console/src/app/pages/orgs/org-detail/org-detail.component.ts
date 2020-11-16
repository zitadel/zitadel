import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatButtonToggleChange } from '@angular/material/button-toggle';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of, Subscription } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { PolicyGridType } from 'src/app/modules/policy-grid/policy-grid.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import {
    Org,
    OrgDomainView,
    OrgMember,
    OrgMemberSearchResponse,
    OrgMemberView,
    OrgState,
    UserView,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddDomainDialogComponent } from './add-domain-dialog/add-domain-dialog.component';
import { DomainVerificationComponent } from './domain-verification/domain-verification.component';


@Component({
    selector: 'app-org-detail',
    templateUrl: './org-detail.component.html',
    styleUrls: ['./org-detail.component.scss'],
})
export class OrgDetailComponent implements OnInit, OnDestroy {
    public org!: Org.AsObject;
    public PolicyComponentServiceType: any = PolicyComponentServiceType;

    public dataSource: MatTableDataSource<OrgMember.AsObject> = new MatTableDataSource<OrgMember.AsObject>();
    public memberResult!: OrgMemberSearchResponse.AsObject;
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];
    public selection: SelectionModel<OrgMember.AsObject> = new SelectionModel<OrgMember.AsObject>(true, []);
    public OrgState: any = OrgState;
    public ChangeType: any = ChangeType;

    private subscription: Subscription = new Subscription();

    public domains: OrgDomainView.AsObject[] = [];
    public primaryDomain: string = '';

    // members
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public totalMemberResult: number = 0;
    public membersSubject: BehaviorSubject<OrgMemberView.AsObject[]>
        = new BehaviorSubject<OrgMemberView.AsObject[]>([]);
    public PolicyGridType: any = PolicyGridType;

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

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

    private async getData(): Promise<void> {
        this.mgmtService.GetMyOrg().then((org: Org) => {
            this.org = org.toObject();
        }).catch(error => {
            this.toast.showError(error);
        });
        this.loadMembers();
        this.loadDomains();
    }

    public loadDomains(): void {
        this.mgmtService.SearchMyOrgDomains().then(result => {
            this.domains = result.toObject().resultList;
            this.primaryDomain = this.domains.find(domain => domain.primary)?.domain ?? '';
        });
    }

    public setPrimary(domain: OrgDomainView.AsObject): void {
        this.mgmtService.setMyPrimaryOrgDomain(domain.domain).then(() => {
            this.toast.showInfo('ORG.TOAST.SETPRIMARY', true);
            this.loadDomains();
        }).catch((error) => {
            this.toast.showError(error);
        });
    }

    public changeState(event: MatButtonToggleChange | any): void {
        if (event.value === OrgState.ORGSTATE_ACTIVE) {
            this.mgmtService.ReactivateMyOrg().then(() => {
                this.toast.showInfo('ORG.TOAST.REACTIVATED', true);
            }).catch((error) => {
                this.toast.showError(error);
            });
        } else if (event.value === OrgState.ORGSTATE_INACTIVE) {
            this.mgmtService.DeactivateMyOrg().then(() => {
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
                this.mgmtService.AddMyOrgDomain(resp).then(domain => {
                    const newDomain = domain;

                    const newDomainView = new OrgDomainView();
                    newDomainView.setChangeDate(newDomain.getChangeDate());
                    newDomainView.setCreationDate(newDomain.getCreationDate());
                    newDomainView.setDomain(newDomain.getDomain());
                    newDomainView.setOrgId(newDomain.getOrgId());
                    newDomainView.setPrimary(newDomain.getPrimary());
                    newDomainView.setSequence(newDomain.getSequence());
                    newDomainView.setVerified(newDomain.getVerified());

                    this.domains.push(newDomainView.toObject());

                    this.verifyDomain(newDomainView.toObject());
                    this.toast.showInfo('ORG.TOAST.DOMAINADDED', true);
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
                this.mgmtService.RemoveMyOrgDomain(domain).then(() => {
                    this.toast.showInfo('ORG.TOAST.DOMAINREMOVED', true);
                    const index = this.domains.findIndex(d => d.domain === domain);
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
                const users: UserView.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.mgmtService.AddMyOrgMember(user.id, roles);
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

    public verifyDomain(domain: OrgDomainView.AsObject): void {
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
        from(this.mgmtService.SearchMyOrgMembers(100, 0)).pipe(
            map(resp => {
                this.totalMemberResult = resp.toObject().totalResult;
                return resp.toObject().resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(members => {
            this.membersSubject.next(members);
        });
    }
}
