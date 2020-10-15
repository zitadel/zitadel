import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { IdpSearchResponse as AdminIdpSearchResponse, IdpView as AdminIdpView } from 'src/app/proto/generated/admin_pb';
import { IdpProviderType, IdpView as MgmtIdpView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

@Component({
    selector: 'app-idp-table',
    templateUrl: './idp-table.component.html',
    styleUrls: ['./idp-table.component.scss'],
})
export class IdpTableComponent implements OnInit {
    @Input() public serviceType!: PolicyComponentServiceType;
    @Input() service!: AdminService | ManagementService;
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource: MatTableDataSource<AdminIdpView.AsObject | MgmtIdpView.AsObject>
        = new MatTableDataSource<AdminIdpView.AsObject | MgmtIdpView.AsObject>();
    public selection: SelectionModel<AdminIdpView.AsObject | MgmtIdpView.AsObject>
        = new SelectionModel<AdminIdpView.AsObject | MgmtIdpView.AsObject>(true, []);
    public idpResult!: AdminIdpSearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    public IdpProviderType: any = IdpProviderType;
    @Input() public displayedColumns: string[] = ['select', 'name', 'config', 'creationDate', 'changeDate', 'state'];

    @Output() public changedSelection: EventEmitter<Array<AdminIdpView.AsObject | MgmtIdpView.AsObject>>
        = new EventEmitter();

    constructor(public translate: TranslateService, private toast: ToastService) {
        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    ngOnInit(): void {
        this.getData(10, 0);
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            this.displayedColumns = ['select', 'name', 'config', 'creationDate', 'changeDate', 'state', 'type'];
        }
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
            return this.service.DeactivateIdpConfig(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
            this.getData(10, 0);
        });
    }

    public reactivateSelectedIdps(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.service.ReactivateIdpConfig(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
            this.getData(10, 0);
        });
    }

    public removeSelectedIdps(): void {
        Promise.all(this.selection.selected.map(value => {
            return this.service.RemoveIdpConfig(value.id);
        })).then(() => {
            this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
            this.getData(10, 0);
        });
    }

    private async getData(limit: number, offset: number): Promise<void> {
        this.loadingSubject.next(true);

        // let query: AdminIdpSearchQuery | MgmtIdpSearchQuery;
        // if (this.service instanceof AdminService) {
        //     query = new AdminIdpSearchQuery();
        //     query.setKey(AdminIdpSearchKey.IDPSEARCHKEY_IDP_CONFIG_ID);
        // } else if (this.service instanceof ManagementService) {
        //     query = new MgmtIdpSearchQuery();
        //     query.setKey(MgmtIdpSearchKey.IDPSEARCHKEY_PROVIDER_TYPE);
        // }

        this.service.SearchIdps(limit, offset).then(resp => {
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

    public get createRouterLink(): RouterLink | any {
        if (this.service instanceof AdminService) {
            return ['/iam', 'idp', 'create'];
        } else if (this.service instanceof ManagementService) {
            return ['/org', 'idp', 'create'];
        }
    }

    public routerLinkForRow(row: MgmtIdpView.AsObject | AdminIdpView.AsObject): any {
        if (row.id) {
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    return ['/org', 'idp', row.id];
                case PolicyComponentServiceType.ADMIN:
                    return ['/iam', 'idp', row.id];
            }
        }
    }
}
