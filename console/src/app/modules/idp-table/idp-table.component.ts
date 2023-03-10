import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { BehaviorSubject, Observable } from 'rxjs';
import {
  ListProvidersRequest as AdminListProvidersRequest,
  ListProvidersResponse as AdminListProvidersResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  IDP,
  IDPLoginPolicyLink,
  IDPOwnerType,
  IDPState,
  IDPStylingType,
  Provider,
  ProviderType,
} from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddCustomLoginPolicyRequest,
  AddCustomLoginPolicyResponse,
  ListProvidersRequest as MgmtListProvidersRequest,
  ListProvidersResponse as MgmtListProvidersResponse,
} from 'src/app/proto/generated/zitadel/management_pb';
import { ListQuery } from 'src/app/proto/generated/zitadel/object_pb';
import { LoginPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
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
  public dataSource: MatTableDataSource<Provider.AsObject> = new MatTableDataSource<Provider.AsObject>();
  public selection: SelectionModel<Provider.AsObject> = new SelectionModel<Provider.AsObject>(true, []);
  public idpResult?: MgmtListProvidersResponse.AsObject | AdminListProvidersResponse.AsObject;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public IDPOwnerType: any = IDPOwnerType;
  public IDPState: any = IDPState;
  public ProviderType: any = ProviderType;
  public displayedColumns: string[] = ['availability', 'name', 'type', 'creationDate', 'changeDate', 'actions'];
  @Output() public changedSelection: EventEmitter<Array<Provider.AsObject>> = new EventEmitter();

  public idps: IDPLoginPolicyLink.AsObject[] = [];
  public IDPStylingType: any = IDPStylingType;
  public loginPolicy!: LoginPolicy.AsObject;

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
      this.displayedColumns = ['availability', 'name', 'type', 'owner', 'creationDate', 'changeDate', 'actions'];
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
      const req = new MgmtListProvidersRequest();
      const lq = new ListQuery();
      lq.setOffset(offset);
      lq.setLimit(limit);
      req.setQuery(lq);
      (this.service as ManagementService)
        .listProviders(req)
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
      const req = new AdminListProvidersRequest();
      const lq = new ListQuery();
      lq.setOffset(offset);
      lq.setLimit(limit);
      req.setQuery(lq);
      (this.service as AdminService)
        .listProviders(req)
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

  public routerLinkForRow(row: Provider.AsObject): any {
    if (row.id) {
      switch (row.type) {
        case ProviderType.PROVIDER_TYPE_OIDC:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'oidc', row.id];
        case ProviderType.PROVIDER_TYPE_OAUTH:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'oauth', row.id];
        case ProviderType.PROVIDER_TYPE_JWT:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'jwt', row.id];
        case ProviderType.PROVIDER_TYPE_GOOGLE:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'google', row.id];
      }
    }
  }

  private async getIdps(): Promise<IDPLoginPolicyLink.AsObject[]> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getLoginPolicy().then((policyResp) => {
          if (policyResp.policy) {
            this.loginPolicy = policyResp.policy;
          }
          return policyResp.policy?.idpsList ?? [];
        });
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getLoginPolicy().then((policyResp) => {
          if (policyResp.policy) {
            this.loginPolicy = policyResp.policy;
          }
          return policyResp.policy?.idpsList ?? [];
        });
    }
  }

  private addLoginPolicy(): Promise<AddCustomLoginPolicyResponse.AsObject> {
    const mgmtreq = new AddCustomLoginPolicyRequest();
    mgmtreq.setAllowExternalIdp(this.loginPolicy.allowExternalIdp);
    mgmtreq.setAllowRegister(this.loginPolicy.allowRegister);
    mgmtreq.setAllowUsernamePassword(this.loginPolicy.allowUsernamePassword);
    mgmtreq.setForceMfa(this.loginPolicy.forceMfa);
    mgmtreq.setPasswordlessType(this.loginPolicy.passwordlessType);
    mgmtreq.setHidePasswordReset(this.loginPolicy.hidePasswordReset);
    mgmtreq.setMultiFactorsList(this.loginPolicy.multiFactorsList);
    mgmtreq.setSecondFactorsList(this.loginPolicy.secondFactorsList);

    const pcl = new Duration()
      .setSeconds(this.loginPolicy.passwordCheckLifetime?.seconds ?? 0)
      .setNanos(this.loginPolicy.passwordCheckLifetime?.nanos ?? 0);
    mgmtreq.setPasswordCheckLifetime(pcl);

    const elcl = new Duration()
      .setSeconds(this.loginPolicy.externalLoginCheckLifetime?.seconds ?? 0)
      .setNanos(this.loginPolicy.externalLoginCheckLifetime?.nanos ?? 0);
    mgmtreq.setExternalLoginCheckLifetime(elcl);

    const misl = new Duration()
      .setSeconds(this.loginPolicy.mfaInitSkipLifetime?.seconds ?? 0)
      .setNanos(this.loginPolicy.mfaInitSkipLifetime?.nanos ?? 0);
    mgmtreq.setMfaInitSkipLifetime(misl);

    const sfcl = new Duration()
      .setSeconds(this.loginPolicy.secondFactorCheckLifetime?.seconds ?? 0)
      .setNanos(this.loginPolicy.secondFactorCheckLifetime?.nanos ?? 0);
    mgmtreq.setSecondFactorCheckLifetime(sfcl);

    const mficl = new Duration()
      .setSeconds(this.loginPolicy.multiFactorCheckLifetime?.seconds ?? 0)
      .setNanos(this.loginPolicy.multiFactorCheckLifetime?.nanos ?? 0);
    mgmtreq.setMultiFactorCheckLifetime(mficl);

    mgmtreq.setAllowDomainDiscovery(this.loginPolicy.allowDomainDiscovery);
    mgmtreq.setIgnoreUnknownUsernames(this.loginPolicy.ignoreUnknownUsernames);
    mgmtreq.setDefaultRedirectUri(this.loginPolicy.defaultRedirectUri);

    return (this.service as ManagementService).addCustomLoginPolicy(mgmtreq);
  }

  public addIdp(idp: Provider.AsObject): Promise<any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if (this.isDefault) {
          return this.addLoginPolicy()
            .then(() => {
              this.loginPolicy.isDefault = false;
              return (this.service as ManagementService).addIDPToLoginPolicy(idp.id, idp.owner).then(() => {
                this.toast.showInfo('IDP.TOAST.ADDED', true);

                setTimeout(() => {
                  this.getIdps()
                    .then((resp) => {
                      this.idps = resp;
                    })
                    .catch((error) => {
                      this.toast.showError(error);
                    });
                }, 2000);
              });
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else {
          return (this.service as ManagementService)
            .addIDPToLoginPolicy(idp.id, idp.owner)
            .then(() => {
              this.toast.showInfo('IDP.TOAST.ADDED', true);
              setTimeout(() => {
                this.getIdps()
                  .then((resp) => {
                    this.idps = resp;
                  })
                  .catch((error) => {
                    this.toast.showError(error);
                  });
              }, 2000);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService)
          .addIDPToLoginPolicy(idp.id)
          .then(() => {
            this.toast.showInfo('IDP.TOAST.ADDED', true);
            setTimeout(() => {
              this.getIdps()
                .then((resp) => {
                  this.idps = resp;
                })
                .catch((error) => {
                  this.toast.showError(error);
                });
            }, 2000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
    }
  }

  public removeIdp(idp: Provider.AsObject): Promise<any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if (this.isDefault) {
          return this.addLoginPolicy()
            .then(() => {
              this.loginPolicy.isDefault = false;
              return (this.service as ManagementService)
                .removeIDPFromLoginPolicy(idp.id)
                .then(() => {
                  this.toast.showInfo('IDP.TOAST.REMOVED', true);
                  setTimeout(() => {
                    this.getIdps()
                      .then((resp) => {
                        this.idps = resp;
                      })
                      .catch((error) => {
                        this.toast.showError(error);
                      });
                  }, 2000);
                })
                .catch((error) => {
                  this.toast.showError(error);
                });
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else {
          return (this.service as ManagementService)
            .removeIDPFromLoginPolicy(idp.id)
            .then(() => {
              this.toast.showInfo('IDP.TOAST.REMOVED', true);
              setTimeout(() => {
                this.getIdps()
                  .then((resp) => {
                    this.idps = resp;
                  })
                  .catch((error) => {
                    this.toast.showError(error);
                  });
              }, 2000);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService)
          .removeIDPFromLoginPolicy(idp.id)
          .then(() => {
            this.toast.showInfo('IDP.TOAST.REMOVED', true);
            setTimeout(() => {
              this.getIdps()
                .then((resp) => {
                  this.idps = resp;
                })
                .catch((error) => {
                  this.toast.showError(error);
                });
            }, 2000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
    }
  }

  public isEnabled(idp: Provider.AsObject): boolean {
    return this.idps.findIndex((i) => i.idpId === idp.id) > -1;
  }

  public get displayedColumnsWithActions(): string[] {
    return ['actions', ...this.displayedColumns];
  }

  public get isDefault(): boolean {
    if (this.loginPolicy && this.serviceType === PolicyComponentServiceType.MGMT) {
      return this.loginPolicy.isDefault;
    } else {
      return false;
    }
  }
}
