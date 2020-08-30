import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { MachineKeySearchResponse, MachineKeyView } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-machine-keys',
    templateUrl: './machine-keys.component.html',
    styleUrls: ['./machine-keys.component.scss'],
})
export class MachineKeysComponent implements OnInit {
    @Input() userId!: string;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource: MatTableDataSource<MachineKeyView.AsObject> = new MatTableDataSource<MachineKeyView.AsObject>();
    public selection: SelectionModel<MachineKeyView.AsObject> = new SelectionModel<MachineKeyView.AsObject>(true, []);
    public keyResult!: MachineKeySearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    @Input() public displayedColumns: string[] = ['select', 'userId', 'type', 'expiry'];

    @Output() public changedSelection: EventEmitter<Array<MachineKeyView.AsObject>> = new EventEmitter();

    constructor(public translate: TranslateService, private userService: MgmtUserService,
        private toast: ToastService) {
        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    public ngOnInit(): void {
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

        this.userService.SearchMachineKeys(this.userId, limit, offset).then(resp => {
            this.keyResult = resp.toObject();
            this.dataSource.data = this.keyResult.resultList;
            console.log(this.keyResult.resultList);
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
