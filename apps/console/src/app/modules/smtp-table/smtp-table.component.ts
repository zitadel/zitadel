import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { SMTPConfigState } from 'src/app/proto/generated/zitadel/settings_pb';
import { ListQuery } from 'src/app/proto/generated/zitadel/object_pb';
import { LoginPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { ListSMTPConfigsRequest, ListSMTPConfigsResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { SMTPConfig } from 'src/app/proto/generated/zitadel/settings_pb';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { MatTableDataSource } from '@angular/material/table';
import { MatDialog } from '@angular/material/dialog';
import { SmtpTestDialogComponent } from '../smtp-test-dialog/smtp-test-dialog.component';

@Component({
  selector: 'cnsl-smtp-table',
  templateUrl: './smtp-table.component.html',
  styleUrls: ['./smtp-table.component.scss'],
})
export class SMTPTableComponent implements OnInit {
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<SMTPConfig.AsObject> = new MatTableDataSource<SMTPConfig.AsObject>();
  public selection: SelectionModel<SMTPConfig.AsObject> = new SelectionModel<SMTPConfig.AsObject>(true, []);
  public configsResult?: ListSMTPConfigsResponse.AsObject;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public displayedColumns: string[] = ['activated', 'description', 'tls', 'host', 'senderAddress', 'actions'];
  @Output() public changedSelection: EventEmitter<Array<SMTPConfig.AsObject>> = new EventEmitter();

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
    return state === SMTPConfigState.SMTP_CONFIG_ACTIVE;
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

  public deleteSMTPConfig(id: string, senderName: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'SMTP.LIST.DIALOG.DELETE_TITLE',
        descriptionKey: 'SMTP.LIST.DIALOG.DELETE_DESCRIPTION',
        confirmationKey: 'SMTP.LIST.DIALOG.SENDER',
        confirmation: senderName,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.adminService
          .removeSMTPConfig(id)
          .then(() => {
            this.toast.showInfo('SMTP.LIST.DIALOG.DELETED', true);
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

    const req = new ListSMTPConfigsRequest();
    const lq = new ListQuery();
    lq.setOffset(offset);
    lq.setLimit(limit);
    req.setQuery(lq);
    this.adminService
      .listSMTPConfigs()
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

  public get createRouterLink(): RouterLink | any {
    return ['/instance', 'idp', 'create'];
  }

  public routerLinkForRow(row: SMTPConfig.AsObject): any {
    return ['/instance', 'smtpprovider', row.id];
  }

  public get displayedColumnsWithActions(): string[] {
    return ['actions', ...this.displayedColumns];
  }

  public navigateToProvider(row: SMTPConfig.AsObject) {
    if (!row.senderAddress) {
      return;
    }
    this.router.navigate(this.routerLinkForRow(row));
  }
}
