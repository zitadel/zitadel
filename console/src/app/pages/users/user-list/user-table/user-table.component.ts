import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { take } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import {
    DisplayNameQuery,
    EmailQuery,
    FirstNameQuery,
    LastNameQuery,
    SearchQuery,
    Type,
    TypeQuery,
    User,
    UserNameQuery,
    UserState,
} from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

enum UserListSearchKey {
    FIRST_NAME,
    LAST_NAME,
    DISPLAY_NAME,
    USER_NAME,
    EMAIL,
}

@Component({
    selector: 'app-user-table',
    templateUrl: './user-table.component.html',
    styleUrls: ['./user-table.component.scss'],
    animations: [
        enterAnimations,
    ],
})
export class UserTableComponent implements OnInit {
    public userSearchKey: UserListSearchKey | undefined = undefined;
    public Type: any = Type;
    @Input() type: Type = Type.TYPE_HUMAN;
    @Input() refreshOnPreviousRoutes: string[] = [];
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild('input') public filter!: Input;

    public viewTimestamp!: Timestamp.AsObject;
    public totalResult: number = 0;
    public dataSource: MatTableDataSource<User.AsObject> = new MatTableDataSource<User.AsObject>();
    public selection: SelectionModel<User.AsObject> = new SelectionModel<User.AsObject>(true, []);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    @Input() public displayedColumns: string[] = ['select', 'displayName', 'username', 'email', 'state', 'actions'];

    @Output() public changedSelection: EventEmitter<Array<User.AsObject>> = new EventEmitter();

    public UserState: any = UserState;
    public UserListSearchKey: any = UserListSearchKey;

    constructor(
        public translate: TranslateService,
        private userService: ManagementService,
        private toast: ToastService,
        private dialog: MatDialog,
        private route: ActivatedRoute,
    ) {
        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    ngOnInit(): void {
        this.route.queryParams.pipe(take(1)).subscribe(params => {
            this.getData(10, 0, this.type);
            if (params.deferredReload) {
                setTimeout(() => {
                    this.getData(10, 0, this.type);
                }, 2000);
            }
        });
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.data.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.data.forEach(row => this.selection.select(row));
    }


    public changePage(event: PageEvent): void {
        this.getData(event.pageSize, event.pageIndex * event.pageSize, this.type);
    }

    public deactivateSelectedUsers(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.userService.deactivateUser(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
            this.selection.clear();
            setTimeout(() => {
                this.refreshPage();
            }, 1000);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public reactivateSelectedUsers(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.userService.reactivateUser(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
            this.selection.clear();
            setTimeout(() => {
                this.refreshPage();
            }, 1000);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    private async getData(limit: number, offset: number, type: Type, searchValue?: string): Promise<void> {
        this.loadingSubject.next(true);
        const query = new SearchQuery();
        const typeQuery = new TypeQuery();
        typeQuery.setType(type);
        query.setTypeQuery(typeQuery);

        if (searchValue && this.userSearchKey !== undefined) {
            switch (this.userSearchKey) {
                case UserListSearchKey.DISPLAY_NAME:
                    const dNQuery = new DisplayNameQuery();
                    dNQuery.setDisplayName(searchValue);
                    dNQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);

                    query.setDisplayNameQuery(dNQuery);
                    break;
                case UserListSearchKey.USER_NAME:
                    const uNQuery = new UserNameQuery();
                    uNQuery.setUserName(searchValue);
                    uNQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);

                    query.setUserNameQuery(uNQuery);
                    break;
                case UserListSearchKey.FIRST_NAME:
                    const fNQuery = new FirstNameQuery();
                    fNQuery.setFirstName(searchValue);
                    fNQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);

                    query.setFirstNameQuery(fNQuery);
                    break;
                case UserListSearchKey.FIRST_NAME:
                    const lNQuery = new LastNameQuery();
                    lNQuery.setLastName(searchValue);
                    lNQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);

                    query.setLastNameQuery(lNQuery);
                    break;
                case UserListSearchKey.EMAIL:
                    const eQuery = new EmailQuery();
                    eQuery.setEmailAddress(searchValue);
                    eQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);

                    query.setEmailQuery(eQuery);
                    break;
            }
        }

        this.userService.listUsers(limit, offset, [query]).then(resp => {
            if (resp.details?.totalResult) {
                this.totalResult = resp.details?.totalResult;
            }
            if (resp.details?.viewTimestamp) {
                this.viewTimestamp = resp.details?.viewTimestamp;
            }
            this.dataSource.data = resp.resultList;
            this.loadingSubject.next(false);
        }).catch(error => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
        });
    }

    public refreshPage(): void {
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.type);
    }

    public applyFilter(event: Event): void {
        this.selection.clear();
        const filterValue = (event.target as HTMLInputElement).value;

        this.getData(
            this.paginator.pageSize,
            this.paginator.pageIndex * this.paginator.pageSize,
            this.type,
            filterValue,
        );
    }

    public setFilter(key: UserListSearchKey): void {
        setTimeout(() => {
            if (this.filter) {
                (this.filter as any).nativeElement.focus();
            }
        }, 100);

        if (this.userSearchKey !== key) {
            this.userSearchKey = key;
        } else {
            this.userSearchKey = undefined;
            this.refreshPage();
        }
    }

    public deleteUser(user: User.AsObject): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.DIALOG.DELETE_TITLE',
                descriptionKey: 'USER.DIALOG.DELETE_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                this.userService.removeUser(user.id).then(() => {
                    setTimeout(() => {
                        this.refreshPage();
                    }, 1000);
                    this.toast.showInfo('USER.TOAST.DELETED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }
}
