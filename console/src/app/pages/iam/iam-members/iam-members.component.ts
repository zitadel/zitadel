import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { IamMembersDataSource } from './iam-members-datasource';

@Component({
    selector: 'app-iam-members',
    templateUrl: './iam-members.component.html',
    styleUrls: ['./iam-members.component.scss'],
})
export class IamMembersComponent {
    public INITIALPAGESIZE: number = 25;
    public dataSource!: IamMembersDataSource;

    public memberRoleOptions: string[] = [];
    public changePageFactory!: Function;
    public changePage: EventEmitter<void> = new EventEmitter();
    public selection: Array<Member.AsObject> = [];

    constructor(private adminService: AdminService,
        private dialog: MatDialog,
        private toast: ToastService) {

        this.dataSource = new IamMembersDataSource(this.adminService);
        this.dataSource.loadMembers(0, 25);
        this.getRoleOptions();

        this.changePageFactory = (event?: PageEvent) => {
            return this.dataSource.loadMembers(
                event?.pageIndex ?? 0,
                event?.pageSize ?? this.INITIALPAGESIZE,
            );
        };
    }

    public getRoleOptions(): void {
        this.adminService.listIAMMemberRoles().then(resp => {
            this.memberRoleOptions = resp.rolesList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    updateRoles(member: Member.AsObject, selectionChange: MatSelectChange): void {
        this.adminService.updateIAMMember(member.userId, selectionChange.value)
            .then(() => {
                this.toast.showInfo('ORG.TOAST.MEMBERCHANGED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public removeMemberSelection(): void {
        Promise.all(this.selection.map(member => {
            return this.adminService.removeIAMMember(member.userId).then(() => {
                this.toast.showInfo('IAM.TOAST.MEMBERREMOVED', true);
                this.changePage.emit();
            }).catch(error => {
                this.toast.showError(error);
            });
        }));
    }

    public removeMember(member: Member.AsObject): void {
        this.adminService.removeIAMMember(member.userId).then(() => {
            this.toast.showInfo('IAM.TOAST.MEMBERREMOVED', true);
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
                creationType: CreationType.IAM,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.adminService.addIAMMember(user.id, roles);
                    })).then(() => {
                        this.toast.showInfo('IAM.TOAST.MEMBERADDED', true);
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
