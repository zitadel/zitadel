import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { OrgMember, OrgMemberView, OrgState, User } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    CreationType,
    MemberCreateDialogComponent,
} from '../../../modules/add-member-dialog/member-create-dialog.component';

@Component({
    selector: 'app-iam-contributors',
    templateUrl: './iam-contributors.component.html',
    styleUrls: ['./iam-contributors.component.scss'],
})
export class IamContributorsComponent implements OnInit {
    @Input() public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<OrgMember.AsObject>;
    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public totalResult: number = 0;
    public membersSubject: BehaviorSubject<OrgMemberView.AsObject[]>
        = new BehaviorSubject<OrgMemberView.AsObject[]>([]);

    public OrgState: any = OrgState;
    constructor(private adminService: AdminService, private dialog: MatDialog,
        private toast: ToastService,
        private router: Router) { }

    public ngOnInit(): void {
        this.loadMembers(0, 25, 'asc');
    }

    public loadMembers(pageIndex: number, pageSize: number, sortDirection?: string): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);
        from(this.adminService.SearchIamMembers(pageSize, offset)).pipe(
            map(resp => {
                this.totalResult = resp.toObject().totalResult;
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
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.adminService.AddIamMember(user.id, roles);
                    })).then(() => {
                        this.toast.showError('members added');
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
