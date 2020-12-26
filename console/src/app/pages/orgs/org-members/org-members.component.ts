import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { Org, OrgMemberView, UserView } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { OrgMembersDataSource } from './org-members-datasource';

@Component({
    selector: 'app-org-members',
    templateUrl: './org-members.component.html',
    styleUrls: ['./org-members.component.scss'],
})
export class OrgMembersComponent {
    public INITIALPAGESIZE: number = 25;
    public org!: Org.AsObject;
    public disableWrite: boolean = false;
    public dataSource!: OrgMembersDataSource;

    public memberRoleOptions: string[] = [];
    public changePageFactory!: Function;
    public changePage: EventEmitter<void> = new EventEmitter();
    public selection: Array<OrgMemberView.AsObject> = [];

    constructor(
        private mgmtService: ManagementService,
        private dialog: MatDialog,
        private toast: ToastService,
    ) {
        this.mgmtService.GetMyOrg().then(org => {
            this.org = org.toObject();
            this.dataSource = new OrgMembersDataSource(this.mgmtService);
            this.dataSource.loadMembers(0, this.INITIALPAGESIZE);
        });

        this.getRoleOptions();

        this.changePageFactory = (event?: PageEvent) => {
            return this.dataSource.loadMembers(
                event?.pageIndex ?? 0,
                event?.pageSize ?? this.INITIALPAGESIZE,
            );
        };
    }

    public getRoleOptions(): void {
        this.mgmtService.GetOrgMemberRoles().then(resp => {
            this.memberRoleOptions = resp.toObject().rolesList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    updateRoles(member: OrgMemberView.AsObject, selectionChange: MatSelectChange): void {
        this.mgmtService.ChangeMyOrgMember(member.userId, selectionChange.value)
            .then(() => {
                this.toast.showInfo('ORG.TOAST.MEMBERCHANGED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public removeOrgMemberSelection(): void {
        Promise.all(this.selection.map(member => {
            return this.mgmtService.RemoveMyOrgMember(member.userId).then(() => {
                this.toast.showInfo('ORG.TOAST.MEMBERREMOVED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
        })).then(() => {
            setTimeout(() => {
                this.changePage.emit();
            }, 1000);
        });
    }

    public removeOrgMember(member: OrgMemberView.AsObject): void {
        this.mgmtService.RemoveMyOrgMember(member.userId).then(() => {
            this.toast.showInfo('ORG.TOAST.MEMBERREMOVED', true);

            setTimeout(() => {
                this.changePage.emit();
            }, 1000);
        }).catch(error => {
            this.toast.showError(error);
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
                            this.changePage.emit();
                        }, 1000);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }
}
