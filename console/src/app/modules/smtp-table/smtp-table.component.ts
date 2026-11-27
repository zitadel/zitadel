import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { ListQuery } from 'src/app/proto/generated/zitadel/object_pb';
import { LoginPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { ListEmailProvidersResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { EmailProvider, EmailProviderState } from 'src/app/proto/generated/zitadel/settings_pb';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { MatTableDataSource } from '@angular/material/table';
import { MatDialog } from '@angular/material/dialog';
import { SmtpTestDialogComponent } from '../smtp-test-dialog/smtp-test-dialog.component';
import { HttpProviderDialogComponent, HttpProviderDialogData } from '../http-provider-dialog/http-provider-dialog.component';

@Component({
  selector: 'cnsl-smtp-table',
  templateUrl: './smtp-table.component.html',
  styleUrls: ['./smtp-table.component.scss'],
  standalone: false,
})
export class SMTPTableComponent implements OnInit {
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<EmailProvider.AsObject> = new MatTableDataSource<EmailProvider.AsObject>();
  public selection: SelectionModel<EmailProvider.AsObject> = new SelectionModel<EmailProvider.AsObject>(true, []);
  public configsResult?: ListEmailProvidersResponse.AsObject;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public displayedColumns: string[] = ['activated', 'description', 'type', 'tls', 'host', 'senderAddress', 'actions'];
  @Output() public changedSelection: EventEmitter<Array<EmailProvider.AsObject>> = new EventEmitter();

  public loginPolicy!: LoginPolicy.AsObject;

  constructor(
    private adminService: AdminService,
    public translate: TranslateService,
    private toast: ToastService,
    private dialog: MatDialog,
    private router: Router,
  ) {
    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  ngOnInit(): void {
    this.getData(10, 0);
  }

  public isActive(state: number) {
    return state === EmailProviderState.EMAIL_PROVIDER_ACTIVE;
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

  public activateSMTPConfig(id: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CONTINUE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'SMTP.LIST.DIALOG.ACTIVATE_WARN_TITLE',
        descriptionKey: 'SMTP.LIST.DIALOG.ACTIVATE_WARN_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.adminService
          .activateSMTPConfig(id)
          .then(() => {
            this.toast.showInfo('SMTP.LIST.DIALOG.ACTIVATED', true);
            this.refreshPage();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public deactivateSMTPConfig(id: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CONTINUE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'SMTP.LIST.DIALOG.DEACTIVATE_WARN_TITLE',
        descriptionKey: 'SMTP.LIST.DIALOG.DEACTIVATE_WARN_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.adminService
          .deactivateSMTPConfig(id)
          .then(() => {
            this.toast.showInfo('SMTP.LIST.DIALOG.DEACTIVATED', true);
            this.refreshPage();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public deleteSMTPConfig(config: EmailProvider.AsObject): void {
    // an HTTP provider has no sender name, so its endpoint is typed to confirm instead
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: config.http
        ? {
            confirmKey: 'ACTIONS.DELETE',
            cancelKey: 'ACTIONS.CANCEL',
            titleKey: 'SETTING.HTTPPROVIDER.DELETE_TITLE',
            descriptionKey: 'SETTING.HTTPPROVIDER.DELETE_DESCRIPTION',
            confirmationKey: 'SETTING.HTTPPROVIDER.DELETE_CONFIRM',
            confirmation: config.http.endpoint,
          }
        : {
            confirmKey: 'ACTIONS.DELETE',
            cancelKey: 'ACTIONS.CANCEL',
            titleKey: 'SMTP.LIST.DIALOG.DELETE_TITLE',
            descriptionKey: 'SMTP.LIST.DIALOG.DELETE_DESCRIPTION',
            confirmationKey: 'SMTP.LIST.DIALOG.SENDER',
            confirmation: config.smtp?.senderName,
          },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.adminService
          .removeSMTPConfig(config.id)
          .then(() => {
            this.toast.showInfo(config.http ? 'SETTING.HTTPPROVIDER.REMOVED' : 'SMTP.LIST.DIALOG.DELETED', true);
            this.refreshPage();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public testSMTPConfig(id: string): void {
    this.dialog.open(SmtpTestDialogComponent, {
      data: {
        id: id,
        confirmKey: 'ACTIONS.TEST',
        cancelKey: 'ACTIONS.CLOSE',
        titleKey: 'SMTP.LIST.DIALOG.TEST_TITLE',
        descriptionKey: 'SMTP.LIST.DIALOG.TEST_DESCRIPTION',
        emailKey: 'SMTP.LIST.DIALOG.TEST_EMAIL',
        testResultKey: 'SMTP.LIST.DIALOG.TEST_RESULT',
      },
      width: '500px',
    });
  }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    const lq = new ListQuery();
    lq.setOffset(offset);
    lq.setLimit(limit);
    this.adminService
      .listEmailProviders(lq)
      .then((resp) => {
        this.configsResult = resp;
        if (resp.resultList) {
          this.dataSource.data = resp.resultList;
        }
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

  public routerLinkForRow(row: EmailProvider.AsObject): any {
    return ['/instance', 'smtpprovider', row.id];
  }

  public navigateToProvider(row: EmailProvider.AsObject) {
    if (row.http) {
      this.editHttpProvider(row.id, row.description, row.http);
      return;
    }
    if (!row.smtp?.senderAddress) {
      return;
    }
    this.router.navigate(this.routerLinkForRow(row));
  }

  private editHttpProvider(id: string, description: string, http: { endpoint: string; signingKey: string }): void {
    this.dialog
      .open<HttpProviderDialogComponent, HttpProviderDialogData, boolean>(HttpProviderDialogComponent, {
        data: { type: 'email', provider: { id, description, endpoint: http.endpoint, signingKey: http.signingKey } },
        width: '500px',
      })
      .afterClosed()
      .subscribe((saved) => {
        if (saved) {
          this.refreshPage();
        }
      });
  }
}
