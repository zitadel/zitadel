import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { UserView } from 'src/app/proto/generated/auth_pb';
import { SearchMethod, UserSearchKey, UserSearchQuery, UserSearchResponse } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { UserType } from '../user-list.component';

@Component({
    selector: 'app-user-table',
    templateUrl: './user-table.component.html',
    styleUrls: ['./user-table.component.scss'],
})
export class UserTableComponent implements OnInit {
    public UserType: any = UserType;
    @Input() userType: UserType = UserType.HUMAN;
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource: MatTableDataSource<UserView.AsObject> = new MatTableDataSource<UserView.AsObject>();
    public selection: SelectionModel<UserView.AsObject> = new SelectionModel<UserView.AsObject>(true, []);
    public userResult!: UserSearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    @Input() public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'state'];

    @Output() public changedSelection: EventEmitter<Array<UserView.AsObject>> = new EventEmitter();

    constructor(public translate: TranslateService, private userService: ManagementService,
        private toast: ToastService) {
        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    ngOnInit(): void {
        this.getData(10, 0, this.userType);
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
            this.getData(10, 0, this.userType);
        });
    }

    public reactivateSelectedUsers(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.userService.ReactivateUser(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
            this.getData(10, 0, this.userType);
        });
    }

    private async getData(limit: number, offset: number, filterTypeValue: UserType): Promise<void> {
        this.loadingSubject.next(true);
        const query = new UserSearchQuery();
        query.setKey(UserSearchKey.USERSEARCHKEY_TYPE);
        query.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
        query.setValue(filterTypeValue);
        console.log(filterTypeValue);

        this.userService.SearchUsers(limit, offset, [query]).then(resp => {
            this.userResult = resp.toObject();
            this.dataSource.data = this.userResult.resultList;
            console.log(this.userResult.resultList);
            this.loadingSubject.next(false);
        }).catch(error => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
        });
    }

    public refreshPage(): void {
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.userType);
    }
}
