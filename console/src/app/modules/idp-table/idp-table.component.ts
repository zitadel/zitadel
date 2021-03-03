import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { BehaviorSubject, Observable } from 'rxjs';
import { IDP, IDPState, IDPStylingType, IDPType } from 'src/app/proto/generated/zitadel/idp_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';

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
    public dataSource: MatTableDataSource<IDP.AsObject>
        = new MatTableDataSource<IDP.AsObject>();
    public selection: SelectionModel<IDP.AsObject>
        = new SelectionModel<IDP.AsObject>(true, []);
    public idpResult!: AdminIdpSearchResponse.AsObject;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    public IDPType: any = IDPType;
    public IDPState: any = IDPState;
    public IdpStylingType: any = IDPStylingType;
    @Input() public displayedColumns: string[] = ['select', 'name', 'config', 'dates', 'state'];

    @Output() public changedSelection: EventEmitter<Array<IDP.AsObject>>
        = new EventEmitter();

    constructor(public translate: TranslateService, private toast: ToastService, private dialog: MatDialog) {
        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    ngOnInit(): void {
        this.getData(10, 0);
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            this.displayedColumns = ['select', 'name', 'config', 'dates', 'state', 'type'];
        }

        if (!this.disabled) {
            this.displayedColumns.push('actions');
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
        const map: Promise<any>[] = this.selection.selected.map(value => {
            if (this.serviceType === PolicyComponentServiceType.MGMT) {
                return (this.service as ManagementService).deactivateOrgIDP(value.id);
            } else {
                return (this.service as AdminService).deactivateIDP(value.id);
            }
        });
        Promise.all(map).then(() => {
            this.selection.clear();
            this.toast.showInfo('IDP.TOAST.SELECTEDDEACTIVATED', true);
            this.refreshPage();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public reactivateSelectedIdps(): void {
        const map: Promise<Empty>[] = this.selection.selected.map(value => {
            return this.service.ReactivateIdpConfig(value.id);
        });
        Promise.all(map).then(() => {
            this.selection.clear();
            this.toast.showInfo('IDP.TOAST.SELECTEDREACTIVATED', true);
            this.refreshPage();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public removeSelectedIdps(): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'IDP.DELETE_SELECTION_TITLE',
                descriptionKey: 'IDP.DELETE_SELECTION_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                this.selection.clear();

                Promise.all(this.selection.selected.map(value => {
                    return this.service.RemoveIdpConfig(value.id);
                })).then(() => {
                    this.toast.showInfo('IDP.TOAST.SELECTEDDEACTIVATED', true);
                    this.refreshPage();
                });
            }
        });
    }

    public removeIdp(idp: IDP.AsObject): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'IDP.DELETE_TITLE',
                descriptionKey: 'IDP.DELETE_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                this.service.RemoveIdpConfig(idp.id).then(() => {
                    this.toast.showInfo('IDP.TOAST.REMOVED', true);
                    setTimeout(() => {
                        this.refreshPage();
                    }, 1000);
                });
            }
        });
    }

    private async getData(limit: number, offset: number): Promise<void> {
        this.loadingSubject.next(true);

        this.service.SearchIdps(limit, offset).then(resp => {
            this.idpResult = resp.toObject();
            this.dataSource.data = this.idpResult.resultList;
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

    public routerLinkForRow(row: IDP.AsObject): any {
        if (row.id) {
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    switch ((row as IDP.AsObject).) {
                        case IdpProviderType.IDPPROVIDERTYPE_SYSTEM:
                            return ['/iam', 'idp', row.id];
                        case IdpProviderType.IDPPROVIDERTYPE_ORG:
                            return ['/org', 'idp', row.id];
                    }
                    break;
                case PolicyComponentServiceType.ADMIN:
                    return ['/iam', 'idp', row.id];
            }
        }
    }
}
