import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { Org, OrgMember, OrgMemberView, OrgState, User } from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    CreationType,
    MemberCreateDialogComponent,
} from '../../../modules/add-member-dialog/member-create-dialog.component';

@Component({
    selector: 'app-org-contributors',
    templateUrl: './org-contributors.component.html',
    styleUrls: ['./org-contributors.component.scss'],
})
export class OrgContributorsComponent implements OnInit {
    @Input() public org!: Org.AsObject;
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
    constructor(private orgService: OrgService, private dialog: MatDialog,
        private toast: ToastService,
        private router: Router) { }

    public ngOnInit(): void {
        this.loadMembers(0, 25, 'asc');
    }

    public loadMembers(pageIndex: number, pageSize: number, sortDirection?: string): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);
        from(this.orgService.SearchMyOrgMembers(pageSize, offset)).pipe(
            map(resp => {
                this.totalResult = resp.toObject().totalResult;
                return resp.toObject().resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(members => {
            console.log(members);
            this.membersSubject.next(members);
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
                        return this.orgService.AddMyOrgMember(user.id, roles);
                    })).then(() => {
                        this.toast.showError('members added');
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                }
            }
        });
    }

    public showDetail(): void {
        if (this.org?.state === OrgState.ORGSTATE_ACTIVE) {
            this.router.navigate(['orgs', this.org.id, 'members']);
        }
    }
}
