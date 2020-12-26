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
import { UserView } from 'src/app/proto/generated/zitadel/auth_pb';
import {
    SearchMethod,
    UserSearchKey,
    UserSearchQuery,
    UserSearchResponse,
    UserState,
} from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { UserType } from '../user-list.component';

@Component({
    selector: 'app-user-table',
    templateUrl: './user-table.component.html',
    styleUrls: ['./user-table.component.scss'],
    animations: [
        enterAnimations,
    ],
})
export class UserTableComponent implements OnInit {
    public userSearchKey: UserSearchKey | undefined = undefined;
    public UserType: any = UserType;
    @Input() userType: UserType = UserType.HUMAN;
    @Input() refreshOnPreviousRoutes: string[] = [];
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild('input') public filter!: Input;
    public dataSource: MatTableDataSource<UserView.AsObject> = new MatTableDataSource<UserView.AsObject>();
    public selection: SelectionModel<UserView.AsObject> = new SelectionModel<UserView.AsObject>(true, []);
    public userResult!: UserSearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    @Input() public displayedColumns: string[] = ['select', 'displayName', 'username', 'email', 'state', 'actions'];

    @Output() public changedSelection: EventEmitter<Array<UserView.AsObject>> = new EventEmitter();
    UserSearchKey: any = UserSearchKey;

    public UserState: any = UserState;

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
            this.getData(10, 0, this.userType);
            if (params.deferredReload) {
                setTimeout(() => {
                    this.getData(10, 0, this.userType);
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
        this.getData(event.pageSize, event.pageIndex * event.pageSize, this.userType);
    }

    public deactivateSelectedUsers(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.userService.DeactivateUser(value.id);
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
            return this.userService.ReactivateUser(value.id);
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

    private async getData(limit: number, offset: number, filterTypeValue: UserType, filterName?: string): Promise<void> {
        this.loadingSubject.next(true);
        const query = new UserSearchQuery();
        query.setKey(UserSearchKey.USERSEARCHKEY_TYPE);
        query.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
        query.setValue(filterTypeValue);

        let namequery;
        if (filterName && this.userSearchKey !== undefined) {
            namequery = new UserSearchQuery();
            namequery.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS_IGNORE_CASE);
            namequery.setKey(this.userSearchKey);
            namequery.setValue(filterName.toLowerCase());
        }

        this.userService.SearchUsers(limit, offset, namequery ? [query, namequery] : [query]).then(resp => {
            this.userResult = resp.toObject();
            this.dataSource.data = this.userResult.resultList;
            this.loadingSubject.next(false);
        }).catch(error => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
        });
    }

    public refreshPage(): void {
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.userType);
    }

    public applyFilter(event: Event): void {
        this.selection.clear();
        const filterValue = (event.target as HTMLInputElement).value;

        this.getData(
            this.paginator.pageSize,
            this.paginator.pageIndex * this.paginator.pageSize,
            this.userType,
            filterValue,
        );
    }

    public setFilter(key: UserSearchKey): void {
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

    public deleteUser(user: UserView.AsObject): void {
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
                this.userService.DeleteUser(user.id).then(() => {
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
