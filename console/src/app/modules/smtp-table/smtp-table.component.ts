import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { Router, RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
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
import {
  ListSMSProvidersResponse,
  ListSMTPProvidersRequest,
  ListSMTPProvidersResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';

@Component({
  selector: 'cnsl-smtp-table',
  templateUrl: './smtp-table.component.html',
  styleUrls: ['./smtp-table.component.scss'],
})
export class SMTPTableComponent implements OnInit, OnDestroy {
  @Input() service!: AdminService | ManagementService;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<Provider.AsObject> = new MatTableDataSource<Provider.AsObject>();
  public selection: SelectionModel<Provider.AsObject> = new SelectionModel<Provider.AsObject>(true, []);
  public providersResult?: ListSMTPProvidersResponse.AsObject;
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

  // public deactivateSelectedIdps(): void {
  //   const map: Promise<any>[] = this.selection.selected.map((value) => {
  //     return (this.service as AdminService).deactivateIDP(value.id);
  //   });
  //   Promise.all(map)
  //     .then(() => {
  //       this.selection.clear();
  //       this.toast.showInfo('IDP.TOAST.SELECTEDDEACTIVATED', true);
  //       this.refreshPage();
  //     })
  //     .catch((error) => {
  //       this.toast.showError(error);
  //     });
  // }

  // public reactivateSelectedIdps(): void {
  //   const map: Promise<any>[] = this.selection.selected.map((value) => {
  //     return (this.service as AdminService).reactivateIDP(value.id);
  //   });
  //   Promise.all(map)
  //     .then(() => {
  //       this.selection.clear();
  //       this.toast.showInfo('IDP.TOAST.SELECTEDREACTIVATED', true);
  //       this.refreshPage();
  //     })
  //     .catch((error) => {
  //       this.toast.showError(error);
  //     });
  // }

  // public deleteIdp(idp: IDP.AsObject): void {
  //   const dialogRef = this.dialog.open(WarnDialogComponent, {
  //     data: {
  //       confirmKey: 'ACTIONS.DELETE',
  //       cancelKey: 'ACTIONS.CANCEL',
  //       titleKey: 'IDP.DELETE_TITLE',
  //       descriptionKey: 'IDP.DELETE_DESCRIPTION',
  //     },
  //     width: '400px',
  //   });

  //   dialogRef.afterClosed().subscribe((resp) => {
  //     if (resp) {
  //       this.service.deleteProvider(idp.id).then(
  //         () => {
  //           this.toast.showInfo('IDP.TOAST.DELETED', true);
  //           setTimeout(() => {
  //             this.refreshPage();
  //           }, 1000);
  //         },
  //         (error) => {
  //           this.toast.showError(error);
  //         },
  //       );
  //     }
  //   });
  // }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    const req = new ListSMTPProvidersRequest();
    const lq = new ListQuery();
    lq.setOffset(offset);
    lq.setLimit(limit);
    req.setQuery(lq);
    (this.service as AdminService)
      .listSMTPProviders(req)
      .then((resp) => {
        this.providersResult = resp;
        this.dataSource.data = resp.resultList;
        this.loadingSubject.next(false);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loadingSubject.next(false);
      });
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }

  public get createRouterLink(): RouterLink | any {
    return ['/instance', 'idp', 'create'];
  }

  // public routerLinkForRow(row: Provider.AsObject): any {
  //   if (row.id) {
  //     switch (row.type) {
  //       case ProviderType.PROVIDER_TYPE_AZURE_AD:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'azure-ad', row.id];
  //       case ProviderType.PROVIDER_TYPE_OIDC:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'oidc', row.id];
  //       case ProviderType.PROVIDER_TYPE_GITHUB_ES:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'github-es', row.id];
  //       case ProviderType.PROVIDER_TYPE_OAUTH:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'oauth', row.id];
  //       case ProviderType.PROVIDER_TYPE_JWT:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'jwt', row.id];
  //       case ProviderType.PROVIDER_TYPE_GOOGLE:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'google', row.id];
  //       case ProviderType.PROVIDER_TYPE_GITLAB:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'gitlab', row.id];
  //       case ProviderType.PROVIDER_TYPE_LDAP:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'ldap', row.id];
  //       case ProviderType.PROVIDER_TYPE_GITLAB_SELF_HOSTED:
  //         return [
  //           row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org',
  //           'provider',
  //           'gitlab-self-hosted',
  //           row.id,
  //         ];
  //       case ProviderType.PROVIDER_TYPE_GITHUB:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'github', row.id];
  //       case ProviderType.PROVIDER_TYPE_APPLE:
  //         return [row.owner === IDPOwnerType.IDP_OWNER_TYPE_SYSTEM ? '/instance' : '/org', 'provider', 'apple', row.id];
  //     }
  //   }
  // }

  private async getIdps(): Promise<IDPLoginPolicyLink.AsObject[]> {
    return (this.service as AdminService).getLoginPolicy().then((policyResp) => {
      if (policyResp.policy) {
        this.loginPolicy = policyResp.policy;
      }
      return policyResp.policy?.idpsList ?? [];
    });
  }

  // public addIdp(idp: Provider.AsObject): Promise<any> {
  //   return (this.service as AdminService)
  //     .addIDPToLoginPolicy(idp.id)
  //     .then(() => {
  //       this.toast.showInfo('IDP.TOAST.ADDED', true);
  //       setTimeout(() => {
  //         this.reloadIDPs$.next();
  //       }, 2000);
  //     })
  //     .catch((error) => {
  //       this.toast.showError(error);
  //     });
  // }

  // public removeIdp(idp: Provider.AsObject): void {
  //   const dialogRef = this.dialog.open(WarnDialogComponent, {
  //     data: {
  //       confirmKey: 'ACTIONS.CONTINUE',
  //       cancelKey: 'ACTIONS.CANCEL',
  //       titleKey: 'IDP.REMOVE_WARN_TITLE',
  //       descriptionKey: 'IDP.REMOVE_WARN_DESCRIPTION',
  //     },
  //     width: '400px',
  //   });

  //   dialogRef.afterClosed().subscribe((resp) => {
  //     if (resp) {
  //       (this.service as AdminService)
  //         .removeIDPFromLoginPolicy(idp.id)
  //         .then(() => {
  //           this.toast.showInfo('IDP.TOAST.REMOVED', true);
  //           setTimeout(() => {
  //             this.reloadIDPs$.next();
  //           }, 2000);
  //         })
  //         .catch((error) => {
  //           this.toast.showError(error);
  //         });
  //     }
  //   });
  // }

  public isEnabled(idp: Provider.AsObject): boolean {
    return this.idps.findIndex((i) => i.idpId === idp.id) > -1;
  }

  public get displayedColumnsWithActions(): string[] {
    return ['actions', ...this.displayedColumns];
  }
}
