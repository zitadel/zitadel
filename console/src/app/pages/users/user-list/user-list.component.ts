import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, OnDestroy, Output, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { User, UserSearchResponse } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.scss'],
})
export class UserListComponent implements OnDestroy {
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource: MatTableDataSource<User.AsObject> = new MatTableDataSource<User.AsObject>();
    public userResult!: UserSearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'state'];
    public selection: SelectionModel<User.AsObject> = new SelectionModel<User.AsObject>(true, []);
    @Output() public changedSelection: EventEmitter<Array<User.AsObject>> = new EventEmitter();

    private subscription?: Subscription;

    constructor(public translate: TranslateService, private route: ActivatedRoute, private userService: MgmtUserService,
        private toast: ToastService) {
        this.subscription = this.route.params.subscribe(() => this.getData(10, 0));

        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
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

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    public changePage(event: PageEvent): void {
        this.getData(event.pageSize, event.pageIndex * event.pageSize);
    }

    public deactivateSelectedUsers(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.userService.DeactivateUser(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
            this.getData(10, 0);
        });
    }

    public reactivateSelectedUsers(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.userService.ReactivateUser(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
            this.getData(10, 0);
        });
    }

    private async getData(limit: number, offset: number): Promise<void> {
        this.loadingSubject.next(true);
        this.userService.SearchUsers(limit, offset).then(resp => {
            this.userResult = resp.toObject();
            this.dataSource.data = this.userResult.resultList;
            this.loadingSubject.next(false);
        }).catch(error => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
        });
    }

    public refreshPage(): void {
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
    }
}
