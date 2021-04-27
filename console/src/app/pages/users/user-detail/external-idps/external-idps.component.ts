import { SelectionModel } from '@angular/cdk/collections';
import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable } from 'rxjs';
import { PageEvent, PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { IDPUserLink } from 'src/app/proto/generated/zitadel/idp_pb';

import { GrpcAuthService } from '../../../../services/grpc-auth.service';
import { ManagementService } from '../../../../services/mgmt.service';
import { ToastService } from '../../../../services/toast.service';

@Component({
  selector: 'app-external-idps',
  templateUrl: './external-idps.component.html',
  styleUrls: ['./external-idps.component.scss'],
})
export class ExternalIdpsComponent implements OnInit {
  @Input() service!: GrpcAuthService | ManagementService;
  @Input() userId!: string;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;
  public dataSource: MatTableDataSource<IDPUserLink.AsObject>
    = new MatTableDataSource<IDPUserLink.AsObject>();
  public selection: SelectionModel<IDPUserLink.AsObject>
    = new SelectionModel<IDPUserLink.AsObject>(true, []);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumns: string[] = ['idpConfigId', 'idpName', 'externalUserId', 'externalUserDisplayName', 'actions'];

  constructor(private toast: ToastService, private dialog: MatDialog) { }

  ngOnInit(): void {
    this.getData(10, 0);
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

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    let promise;
    if (this.service instanceof ManagementService) {
      promise = (this.service as ManagementService).listHumanLinkedIDPs(this.userId, limit, offset);
    } else if (this.service instanceof GrpcAuthService) {
      promise = (this.service as GrpcAuthService).listMyLinkedIDPs(limit, offset);
    }

    if (promise) {
      promise.then(resp => {
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
      }).catch((error: any) => {
        this.toast.showError(error);
        this.loadingSubject.next(false);
      });
    }
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }

  public removeExternalIdp(idp: IDPUserLink.AsObject): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.REMOVE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.EXTERNALIDP.DIALOG.DELETE_TITLE',
        descriptionKey: 'USER.EXTERNALIDP.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {
        let promise;
        if (this.service instanceof ManagementService) {
          promise = (this.service as ManagementService)
            .removeHumanLinkedIDP(idp.providedUserId, idp.idpId, idp.userId);
        } else if (this.service instanceof GrpcAuthService) {
          promise = (this.service as GrpcAuthService)
            .removeMyLinkedIDP(idp.providedUserId, idp.idpId);
        }

        if (promise) {
          promise.then(_ => {
            setTimeout(() => {
              this.refreshPage();
            }, 1000);
          }).catch((error: any) => {
            this.toast.showError(error);
          });
        }
      }
    });
  }
}
