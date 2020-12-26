import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { IamMember, IamMemberView } from 'src/app/proto/generated/zitadel/admin_pb';
import { ProjectMember, ProjectType, UserView } from 'src/app/proto/generated/zitadel/management_pb';
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
    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public dataSource!: IamMembersDataSource;

    public memberRoleOptions: string[] = [];
    public changePageFactory!: Function;
    public changePage: EventEmitter<void> = new EventEmitter();
    public selection: Array<IamMemberView.AsObject> = [];

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
        this.adminService.GetIamMemberRoles().then(resp => {
            this.memberRoleOptions = resp.toObject().rolesList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    updateRoles(member: IamMemberView.AsObject, selectionChange: MatSelectChange): void {
        this.adminService.ChangeIamMember(member.userId, selectionChange.value)
            .then((newmember: IamMember) => {
                this.toast.showInfo('ORG.TOAST.MEMBERCHANGED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public removeMemberSelection(): void {
        Promise.all(this.selection.map(member => {
            return this.adminService.RemoveIamMember(member.userId).then(() => {
                this.toast.showInfo('IAM.TOAST.MEMBERREMOVED', true);
                this.changePage.emit();
            }).catch(error => {
                this.toast.showError(error);
            });
        }));
    }

    public removeMember(member: ProjectMember.AsObject): void {
        this.adminService.RemoveIamMember(member.userId).then(() => {
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
                const users: UserView.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.adminService.AddIamMember(user.id, roles);
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
