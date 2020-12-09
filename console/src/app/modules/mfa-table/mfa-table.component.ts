import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { MultiFactorsResult as AdminMultiFactorsResult, , MultiFactorType as AdminMultiFactorType } from 'src/app/proto/generated/admin_pb';
import { MultiFactorsResult as MgmtMultiFactorsResult, MultiFactorType as MgmtMultiFactorType } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';

@Component({
    selector: 'app-mfa-table',
    templateUrl: './mfa-table.component.html',
    styleUrls: ['./mfa-table.component.scss'],
})
export class MfaTableComponent implements OnInit {
    @Input() public serviceType!: PolicyComponentServiceType;
    @Input() service!: AdminService | ManagementService;
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource: MatTableDataSource<AdminMultiFactorsResult.AsObject | MgmtMultiFactorsResult.AsObject>
        = new MatTableDataSource<AdminMultiFactorsResult.AsObject | MgmtMultiFactorsResult.AsObject>();
    public selection: SelectionModel<AdminMultiFactorsResult.AsObject | MgmtMultiFactorsResult.AsObject>
        = new SelectionModel<MgmtMultiFactorType.AsObject | AdminMultiFactorType.AsObject>(true, []);
    public mfaResult!: AdminMultiFactorsResult.AsObject;

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    @Input() public displayedColumns: string[] = ['select', 'name', 'config', 'creationDate', 'changeDate', 'state'];

    @Output() public changedSelection: EventEmitter<Array<AdminMultiFactorsResult.AsObject | MgmtMultiFactorsResult.AsObject>>
        = new EventEmitter();

    constructor(public translate: TranslateService, private toast: ToastService, private dialog: MatDialog) {
        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    ngOnInit(): void {
        this.getData(10, 0);
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            this.displayedColumns = ['select', 'name', 'config', 'creationDate', 'changeDate', 'state', 'type'];
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
        this.selection.clear();
        Promise.all(this.selection.selected.map(value => {
            return this.service.DeactivateIdpConfig(value.id);
        })).then(() => {
            this.toast.showInfo('IDP.TOAST.SELECTEDDEACTIVATED', true);
            this.refreshPage();
        });
    }

    public reactivateSelectedIdps(): void {
        this.selection.clear();
        Promise.all(this.selection.selected.map(value => {
            return this.service.ReactivateIdpConfig(value.id);
        })).then(() => {
            this.toast.showInfo('IDP.TOAST.SELECTEDREACTIVATED', true);
            this.refreshPage();
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

    public removeIdp(idp: AdminIdpView.AsObject | MgmtIdpView.AsObject): void {
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

        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            (this.service as ManagementService).GetLoginPolicyMultiFactors().then(resp => {
                this.idpResult = resp.toObject();
                this.dataSource.data = this.idpResult.resultList;
                this.loadingSubject.next(false);
            }).catch(error => {
                this.toast.showError(error);
                this.loadingSubject.next(false);
            });
        } else {

        }
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
                    switch ((row as MgmtIdpView.AsObject).providerType) {
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
