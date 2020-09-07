import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { IdpSearchResponse, IdpView } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';


@Component({
    selector: 'app-idp-table',
    templateUrl: './idp-table.component.html',
    styleUrls: ['./idp-table.component.scss'],
})
export class IdpTableComponent implements OnInit {
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource: MatTableDataSource<IdpView.AsObject> = new MatTableDataSource<IdpView.AsObject>();
    public selection: SelectionModel<IdpView.AsObject> = new SelectionModel<IdpView.AsObject>(true, []);
    public idpResult!: IdpSearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    @Input() public displayedColumns: string[] = ['select', 'name', 'config', 'creationDate', 'changeDate', 'state'];

    @Output() public changedSelection: EventEmitter<Array<IdpView.AsObject>> = new EventEmitter();

    constructor(public translate: TranslateService, private adminService: AdminService,
        private toast: ToastService) {
        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    ngOnInit(): void {
        this.getData(10, 0);
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
        this.getData(event.pageSize, event.pageIndex * event.pageSize);
    }

    public deactivateSelectedIdps(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.adminService.DeactivateIdpConfig(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
            this.getData(10, 0);
        });
    }

    public reactivateSelectedIdps(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.adminService.ReactivateIdpConfig(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
            this.getData(10, 0);
        });
    }

    private async getData(limit: number, offset: number): Promise<void> {
        this.loadingSubject.next(true);
        // const query = new UserSearchQuery();
        // query.setKey(UserSearchKey.USERSEARCHKEY_TYPE);
        // query.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
        // query.setValue(filterTypeValue);
        // console.log(filterTypeValue);

        this.adminService.SearchIdps(limit, offset).then(resp => {
            this.idpResult = resp.toObject();
            this.dataSource.data = this.idpResult.resultList;
            console.log(this.idpResult.resultList);
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
