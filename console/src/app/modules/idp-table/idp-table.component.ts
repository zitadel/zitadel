import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ListIDPsResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { IDP, IDPLoginPolicyLink, IDPOwnerType, IDPState, IDPStylingType } from 'src/app/proto/generated/zitadel/idp_pb';
import { ListOrgIDPsResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';

@Component({
  selector: 'cnsl-idp-table',
  templateUrl: './idp-table.component.html',
  styleUrls: ['./idp-table.component.scss'],
})
export class IdpTableComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  @Input() service!: AdminService | ManagementService;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<IDP.AsObject> = new MatTableDataSource<IDP.AsObject>();
  public selection: SelectionModel<IDP.AsObject> = new SelectionModel<IDP.AsObject>(true, []);
  public idpResult!: ListIDPsResponse.AsObject | ListOrgIDPsResponse.AsObject;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public IDPOwnerType: any = IDPOwnerType;
  public IDPState: any = IDPState;
  public displayedColumns: string[] = ['availability', 'name', 'type', 'creationDate', 'changeDate', 'state'];
  @Output() public changedSelection: EventEmitter<Array<IDP.AsObject>> = new EventEmitter();

  public idps: IDPLoginPolicyLink.AsObject[] = [];
  public IDPStylingType: any = IDPStylingType;

  constructor(public translate: TranslateService, private toast: ToastService, private dialog: MatDialog) {
    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  ngOnInit(): void {
    this.getData(10, 0);
    this.getIdps().then((resp) => {
      this.idps = resp;
    });

    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      this.displayedColumns = ['availability', 'name', 'type', 'owner', 'creationDate', 'changeDate', 'state'];
    }
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected() ? this.selection.clear() : this.dataSource.data.forEach((row) => this.selection.select(row));
  }

  public changePage(event: PageEvent): void {
    this.getData(event.pageSize, event.pageIndex * event.pageSize);
  }

  public deactivateSelectedIdps(): void {
    const map: Promise<any>[] = this.selection.selected.map((value) => {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        return (this.service as ManagementService).deactivateOrgIDP(value.id);
      } else {
        return (this.service as AdminService).deactivateIDP(value.id);
      }
    });
    Promise.all(map)
      .then(() => {
        this.selection.clear();
        this.toast.showInfo('IDP.TOAST.SELECTEDDEACTIVATED', true);
        this.refreshPage();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public reactivateSelectedIdps(): void {
    const map: Promise<any>[] = this.selection.selected.map((value) => {
      if (this.serviceType === PolicyComponentServiceType.MGMT) {
        return (this.service as ManagementService).reactivateOrgIDP(value.id);
      } else {
        return (this.service as AdminService).reactivateIDP(value.id);
      }
    });
    Promise.all(map)
      .then(() => {
        this.selection.clear();
        this.toast.showInfo('IDP.TOAST.SELECTEDREACTIVATED', true);
        this.refreshPage();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deleteIdp(idp: IDP.AsObject): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'IDP.DELETE_TITLE',
        descriptionKey: 'IDP.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
          (this.service as ManagementService).removeOrgIDP(idp.id).then(
            () => {
              this.toast.showInfo('IDP.TOAST.DELETED', true);
              setTimeout(() => {
                this.refreshPage();
              }, 1000);
            },
            (error) => {
              this.toast.showError(error);
            },
          );
        } else {
          (this.service as AdminService).removeIDP(idp.id).then(
            () => {
              this.toast.showInfo('IDP.TOAST.DELETED', true);
              setTimeout(() => {
                this.refreshPage();
              }, 1000);
            },
            (error) => {
              this.toast.showError(error);
            },
          );
        }
      }
    });
  }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      // const query: IDPQuery = new IDPQuery();
      // const otQuery: IDPOwnerTypeQuery = new IDPOwnerTypeQuery();
      // otQuery.setOwnerType(IDPOwnerType.IDP_OWNER_TYPE_ORG);
      // query.setOwnerTypeQuery(otQuery);
      (this.service as ManagementService)
        .listOrgIDPs(limit, offset)
        .then((resp) => {
          this.idpResult = resp;
          this.dataSource.data = resp.resultList;
          this.loadingSubject.next(false);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loadingSubject.next(false);
        });
    } else {
      (this.service as AdminService)
        .listIDPs(limit, offset)
        .then((resp) => {
          this.idpResult = resp;
          this.dataSource.data = resp.resultList;
          this.loadingSubject.next(false);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.loadingSubject.next(false);
        });
    }
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }

  public get createRouterLink(): RouterLink | any {
    if (this.service instanceof AdminService) {
      return ['/instance', 'idp', 'create'];
    } else if (this.service instanceof ManagementService) {
      return ['/org', 'idp', 'create'];
    }
  }

  public routerLinkForRow(row: IDP.AsObject): any {
    if (row.id) {
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          switch (row.owner) {
            case IDPOwnerType.IDP_OWNER_TYPE_SYSTEM:
              return ['/instance', 'idp', row.id];
            case IDPOwnerType.IDP_OWNER_TYPE_ORG:
              return ['/org', 'idp', row.id];
          }
          break;
        case PolicyComponentServiceType.ADMIN:
          return ['/instance', 'idp', row.id];
      }
    }
  }

  private async getIdps(): Promise<IDPLoginPolicyLink.AsObject[]> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).listLoginPolicyIDPs().then((resp) => {
          return resp.resultList;
        });
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).listLoginPolicyIDPs().then((providers) => {
          return providers.resultList;
        });
    }
  }

  public addIdp(idp: IDP.AsObject | IDP.AsObject): Promise<any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).addIDPToLoginPolicy(idp.id, idp.owner).then(() => {
          this.getIdps()
            .then((resp) => {
              this.idps = resp;
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        });
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService)
          .addIDPToLoginPolicy(idp.id)
          .then(() => {
            this.getIdps().then((resp) => {
              this.idps = resp;
            });
          })
          .catch((error) => {
            this.toast.showError(error);
          });
    }
  }

  public removeIdp(idp: IDP.AsObject): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        (this.service as ManagementService).removeIDPFromLoginPolicy(idp.id).then(
          () => {
            this.getIdps().then((resp) => {
              this.idps = resp;
            });
          },
          (error) => {
            this.toast.showError(error);
          },
        );
        break;
      case PolicyComponentServiceType.ADMIN:
        (this.service as AdminService).removeIDPFromLoginPolicy(idp.id).then(
          () => {
            this.getIdps().then((resp) => {
              this.idps = resp;
            });
          },
          (error) => {
            this.toast.showError(error);
          },
        );
        break;
    }
  }

  public isEnabled(idp: IDP.AsObject): boolean {
    return this.idps.findIndex((i) => i.idpId === idp.id) > -1;
  }

  public get displayedColumnsWithActions(): string[] {
    return ['actions', ...this.displayedColumns];
  }
}
