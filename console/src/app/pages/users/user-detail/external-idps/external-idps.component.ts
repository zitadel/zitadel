import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, firstValueFrom, Observable } from 'rxjs';
import { PageEvent, PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { IDPUserLink } from 'src/app/proto/generated/zitadel/idp_pb';

import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

const PROJECTION_SETTLE_DELAY_MS = 1000;

@Component({
  selector: 'cnsl-external-idps',
  templateUrl: './external-idps.component.html',
  styleUrls: ['./external-idps.component.scss'],
  standalone: false,
})
export class ExternalIdpsComponent implements OnInit, OnDestroy {
  @Input({ required: true }) service!: GrpcAuthService | ManagementService;
  @Input() userId!: string;
  @Output() linkedIdpsChanged = new EventEmitter<void>();
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;
  public dataSource: MatTableDataSource<IDPUserLink.AsObject> = new MatTableDataSource<IDPUserLink.AsObject>();
  public selection: SelectionModel<IDPUserLink.AsObject> = new SelectionModel<IDPUserLink.AsObject>(true, []);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumns: string[] = [
    'idpConfigId',
    'idpName',
    'externalUserId',
    'externalUserDisplayName',
    'actions',
  ];

  constructor(
    private toast: ToastService,
    private dialog: MatDialog,
  ) {}

  ngOnInit(): void {
    this.getData(10, 0).then();
  }

  ngOnDestroy(): void {
    this.loadingSubject.complete();
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

  private listLinkedIdps(limit: number, offset: number) {
    if (this.service instanceof ManagementService) {
      return this.service.listHumanLinkedIDPs(this.userId, limit, offset);
    }
    return this.service.listMyLinkedIDPs(limit, offset);
  }

  private removeLinkedIdp(idp: IDPUserLink.AsObject) {
    if (this.service instanceof ManagementService) {
      return this.service.removeHumanLinkedIDP(idp.idpId, idp.providedUserId, idp.userId);
    }
    return this.service.removeMyLinkedIDP(idp.idpId, idp.providedUserId);
  }

  private waitForProjection(): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, PROJECTION_SETTLE_DELAY_MS));
  }

  private async refreshAfterMutation(): Promise<void> {
    await this.waitForProjection();
    await this.refreshPage();
    this.linkedIdpsChanged.emit();
  }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    let resp;
    try {
      resp = await this.listLinkedIdps(limit, offset);
    } catch (error) {
      this.toast.showError(error);
      this.loadingSubject.next(false);
      return;
    }

    this.dataSource.data = resp.resultList;
    if (resp.details?.viewTimestamp) {
      this.viewTimestamp = resp.details.viewTimestamp;
    }
    if (resp.details?.totalResult) {
      this.totalResult = resp.details?.totalResult;
    } else {
      this.totalResult = 0;
    }
    this.loadingSubject.next(false);
  }

  public refreshPage(): Promise<void> {
    return this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }

  public async removeExternalIdp(idp: IDPUserLink.AsObject): Promise<void> {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.REMOVE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.EXTERNALIDP.DIALOG.REMOVE_TITLE',
        descriptionKey: 'USER.EXTERNALIDP.DIALOG.REMOVE_DESCRIPTION',
      },
      width: '400px',
    });

    const resp = await firstValueFrom(dialogRef.afterClosed());
    if (!resp) {
      return;
    }

    try {
      await this.removeLinkedIdp(idp);
      await this.refreshAfterMutation();
    } catch (error) {
      this.toast.showError(error);
    }
  }
}
