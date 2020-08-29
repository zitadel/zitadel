import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, OnDestroy, Output, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { UserSearchKey, UserSearchQuery, UserSearchResponse, UserView } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ToastService } from 'src/app/services/toast.service';

export enum UserType {
    HUMAN = 'human',
    MACHINE = 'machine',
}
@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.scss'],
})
export class UserListComponent implements OnDestroy {
    public UserType: any = UserType;

    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSources: {
        [type: string]: MatTableDataSource<UserView.AsObject>;
    } = {
            [UserType.HUMAN]: new MatTableDataSource<UserView.AsObject>(),
            [UserType.MACHINE]: new MatTableDataSource<UserView.AsObject>(),
        };
    public selections: {
        [type: string]: SelectionModel<UserView.AsObject>;
    } = {
            [UserType.HUMAN]: new SelectionModel<UserView.AsObject>(true, []),
            [UserType.MACHINE]: new SelectionModel<UserView.AsObject>(true, []),
        };

    public userResults: {
        [type: string]: UserSearchResponse.AsObject;
    } = {};
    private loadingSubjects: {
        [type: string]: BehaviorSubject<boolean>;
    } = {
            [UserType.HUMAN]: new BehaviorSubject<boolean>(false),
            [UserType.MACHINE]: new BehaviorSubject<boolean>(false),
        };
    public loading$: {
        [type: string]: Observable<boolean>;
    } = {
            [UserType.HUMAN]: this.loadingSubjects[UserType.HUMAN].asObservable(),
            [UserType.MACHINE]: this.loadingSubjects[UserType.MACHINE].asObservable(),
        };
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'state'];

    @Output() public changedSelection: EventEmitter<Array<UserView.AsObject>> = new EventEmitter();

    private subscription?: Subscription;

    constructor(public translate: TranslateService, private route: ActivatedRoute, private userService: MgmtUserService,
        private toast: ToastService) {
        this.subscription = this.route.params.subscribe(() => this.getData(10, 0, UserType.HUMAN));

        this.selections[UserType.HUMAN].changed.subscribe(() => {
            this.changedSelection.emit(this.selections[UserType.HUMAN].selected);
        });

        this.selections[UserType.MACHINE].changed.subscribe(() => {
            this.changedSelection.emit(this.selections[UserType.MACHINE].selected);
        });
    }

    public isAllSelected(usertype: UserType): boolean {
        const numSelected = this.selections[usertype].selected.length;
        const numRows = this.dataSources[usertype].data.length;
        return numSelected === numRows;
    }

    public masterToggle(usertype: UserType): void {
        this.isAllSelected(usertype) ?
            this.selections[usertype].clear() :
            this.dataSources[usertype].data.forEach(row => this.selections[usertype].select(row));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    public changePage(event: PageEvent): void {
        this.getData(event.pageSize, event.pageIndex * event.pageSize, UserType.HUMAN);
        this.getData(event.pageSize, event.pageIndex * event.pageSize, UserType.MACHINE);
    }

    public deactivateSelectedUsers(usertype: UserType): void {
        Promise.all(this.selections[usertype].selected.map(value => {
            return this.userService.DeactivateUser(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
            this.getData(10, 0, UserType.HUMAN);
            this.getData(10, 0, UserType.MACHINE);
        });
    }

    public reactivateSelectedUsers(usertype: UserType): void {
        Promise.all(this.selections[usertype].selected.map(value => {
            return this.userService.ReactivateUser(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
            this.getData(10, 0, UserType.HUMAN);
            this.getData(10, 0, UserType.MACHINE);
        });
    }

    private async getData(limit: number, offset: number, filterTypeValue: UserType): Promise<void> {
        this.loadingSubjects[filterTypeValue].next(true);
        const query = new UserSearchQuery();
        query.setKey(UserSearchKey.USERSEARCHKEY_TYPE);
        query.setValue(filterTypeValue);

        this.userService.SearchUsers(limit, offset).then(resp => {
            this.userResults[filterTypeValue] = resp.toObject();
            this.dataSources[filterTypeValue].data = this.userResults[filterTypeValue].resultList;
            console.log(this.userResults[filterTypeValue].resultList);
            this.loadingSubjects[filterTypeValue].next(false);
        }).catch(error => {
            this.toast.showError(error);
            this.loadingSubjects[filterTypeValue].next(false);
        });
    }

    public refreshPage(): void {
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, UserType.HUMAN);
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, UserType.MACHINE);
    }
}
