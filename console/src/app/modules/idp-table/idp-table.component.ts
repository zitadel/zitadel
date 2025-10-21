import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { Router, RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { BehaviorSubject, firstValueFrom, Observable, Subject } from 'rxjs';
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

import { OverlayWorkflowService } from 'src/app/services/overlay/overlay-workflow.service';
import { ContextChangedWorkflowOverlays } from 'src/app/services/overlay/workflows';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { LoginPolicyService } from '../../services/login-policy.service';
import { first } from 'rxjs/operators';

@Component({
  selector: 'cnsl-idp-table',
  templateUrl: './idp-table.component.html',
  styleUrls: ['./idp-table.component.scss'],
  standalone: false,
})
export class IdpTableComponent implements OnInit, OnDestroy {
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

  private reloadIDPs$: Subject<void> = new Subject();

  constructor(
    private workflowService: OverlayWorkflowService,
    public translate: TranslateService,
    private toast: ToastService,
    private dialog: MatDialog,
    private router: Router,
    private loginPolicySvc: LoginPolicyService,
  ) {
    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });

    this.reloadIDPs$.subscribe(() => {
      this.getIdps()
        .then((resp) => {
          this.idps = resp;
        })
        .catch((error) => {
          this.toast.showError(error);
        });
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

  ngOnDestroy(): void {
    this.reloadIDPs$.complete();
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
        this.service.deleteProvider(idp.id).then(
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
        case ProviderType.PROVIDER_TYPE_AZURE_AD:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'azure-ad', row.id];
        case ProviderType.PROVIDER_TYPE_OIDC:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'oidc', row.id];
        case ProviderType.PROVIDER_TYPE_GITHUB_ES:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'github-es', row.id];
        case ProviderType.PROVIDER_TYPE_OAUTH:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'oauth', row.id];
        case ProviderType.PROVIDER_TYPE_JWT:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'jwt', row.id];
        case ProviderType.PROVIDER_TYPE_GOOGLE:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'google', row.id];
        case ProviderType.PROVIDER_TYPE_GITLAB:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'gitlab', row.id];
        case ProviderType.PROVIDER_TYPE_LDAP:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'ldap', row.id];
        case ProviderType.PROVIDER_TYPE_GITLAB_SELF_HOSTED:
          return [
            row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org',
            'provider',
            'gitlab-self-hosted',
            row.id,
          ];
        case ProviderType.PROVIDER_TYPE_GITHUB:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'github', row.id];
        case ProviderType.PROVIDER_TYPE_APPLE:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'apple', row.id];
        case ProviderType.PROVIDER_TYPE_SAML:
          return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'saml', row.id];
      }
    }
  }

  navigateToIDP(row: Provider.AsObject) {
    this.router.navigate(this.routerLinkForRow(row)).then(() => {
      if (this.serviceType === PolicyComponentServiceType.MGMT && row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM) {
        setTimeout(() => {
          this.workflowService.startWorkflow(ContextChangedWorkflowOverlays, null);
        }, 1000);
      }
    });
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

  public addIdp(idp: Provider.AsObject): Promise<any> {
    return firstValueFrom(this.loginPolicySvc.activateIdp(this.service, idp.id, idp.owner, this.loginPolicy))
      .then(() => {
        this.toast.showInfo('IDP.TOAST.ADDED', true);
        setTimeout(() => {
          this.reloadIDPs$.next();
        }, 2000);
      })
      .catch(this.toast.showError);
  }

  public removeIdp(idp: Provider.AsObject): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CONTINUE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'IDP.REMOVE_WARN_TITLE',
        descriptionKey: 'IDP.REMOVE_WARN_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        switch (this.serviceType) {
          case PolicyComponentServiceType.MGMT:
            if (this.isDefault) {
              this.loginPolicySvc
                .createCustomLoginPolicy(this.service as ManagementService, this.loginPolicy)
                .then(() => {
                  this.loginPolicy.isDefault = false;
                  return (this.service as ManagementService)
                    .removeIDPFromLoginPolicy(idp.id)
                    .then(() => {
                      this.toast.showInfo('IDP.TOAST.REMOVED', true);
                      setTimeout(() => {
                        this.reloadIDPs$.next();
                      }, 2000);
                    })
                    .catch((error) => {
                      this.toast.showError(error);
                    });
                })
                .catch((error) => {
                  this.toast.showError(error);
                });
              break;
            } else {
              (this.service as ManagementService)
                .removeIDPFromLoginPolicy(idp.id)
                .then(() => {
                  this.toast.showInfo('IDP.TOAST.REMOVED', true);
                  setTimeout(() => {
                    this.reloadIDPs$.next();
                  }, 2000);
                })
                .catch((error) => {
                  this.toast.showError(error);
                });
              break;
            }
          case PolicyComponentServiceType.ADMIN:
            (this.service as AdminService)
              .removeIDPFromLoginPolicy(idp.id)
              .then(() => {
                this.toast.showInfo('IDP.TOAST.REMOVED', true);
                setTimeout(() => {
                  this.reloadIDPs$.next();
                }, 2000);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
            break;
        }
      }
    });
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
