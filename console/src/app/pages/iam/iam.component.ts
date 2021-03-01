import { Component } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { PolicyGridType } from 'src/app/modules/policy-grid/policy-grid.component';
import { OrgMemberView, UserView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-iam',
    templateUrl: './iam.component.html',
    styleUrls: ['./iam.component.scss'],
})
export class IamComponent {
    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public totalMemberResult: number = 0;
    public membersSubject: BehaviorSubject<OrgMemberView.AsObject[]>
        = new BehaviorSubject<OrgMemberView.AsObject[]>([]);

    public PolicyGridType: any = PolicyGridType;

    constructor(public adminService: AdminService, private dialog: MatDialog, private toast: ToastService,
        private router: Router) {
        this.loadMembers();
    }

    public loadMembers(): void {
        this.loadingSubject.next(true);
        from(this.adminService.listIAMMembers(100, 0)).pipe(
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
                        this.toast.showInfo('IAM.TOAST.MEMBERADDED');
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
        this.router.navigate(['iam/members']);
    }
}
