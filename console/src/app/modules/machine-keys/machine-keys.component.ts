import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Moment } from 'moment';
import { BehaviorSubject, Observable } from 'rxjs';
import { Key, KeyType } from 'src/app/proto/generated/zitadel/auth_n_key_pb';
import { ListMachineTokensResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { MachineToken } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddTokenDialogComponent } from '../add-token-dialog/add-token-dialog.component';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { ShowTokenDialogComponent } from '../show-token-dialog/show-token-dialog.component';

@Component({
  selector: 'cnsl-machine-keys',
  templateUrl: './machine-keys.component.html',
  styleUrls: ['./machine-keys.component.scss'],
})
export class MachineKeysComponent implements OnInit {
  @Input() userId!: string;

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<MachineToken.AsObject> = new MatTableDataSource<MachineToken.AsObject>();
  public selection: SelectionModel<MachineToken.AsObject> = new SelectionModel<MachineToken.AsObject>(true, []);
  public tokenResult!: ListMachineTokensResponse.AsObject;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumns: string[] = ['select', 'id', 'type', 'creationDate', 'expirationDate', 'actions'];

  @Output() public changedSelection: EventEmitter<Array<MachineToken.AsObject>> = new EventEmitter();

  constructor(
    public translate: TranslateService,
    private mgmtService: ManagementService,
    private dialog: MatDialog,
    private toast: ToastService,
  ) {
    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  public ngOnInit(): void {
    this.getData(10, 0);
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

  public deleteKey(key: Key.AsObject): void {
    this.mgmtService
      .removeMachineKey(key.id, this.userId)
      .then(() => {
        this.selection.clear();
        this.toast.showInfo('USER.TOAST.SELECTEDKEYSDELETED', true);
        this.getData(10, 0);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public openAddKey(): void {
    const dialogRef = this.dialog.open(AddTokenDialogComponent, {
      data: {},
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const type: KeyType = resp.type;

        let date: Timestamp | undefined;

        if (resp.date as Moment) {
          const ts = new Timestamp();
          const milliseconds = resp.date.toDate().getTime();
          const seconds = Math.abs(milliseconds / 1000);
          const nanos = (milliseconds - seconds * 1000) * 1000 * 1000;
          ts.setSeconds(seconds);
          ts.setNanos(nanos);
          date = ts;
        }

        if (type) {
          this.mgmtService
            .addMachineKey(this.userId, type, date)
            .then((response) => {
              if (response) {
                setTimeout(() => {
                  this.refreshPage();
                }, 1000);

                this.dialog.open(ShowTokenDialogComponent, {
                  data: {
                    key: response,
                  },
                  width: '400px',
                });
              }
            })
            .catch((error: any) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    if (this.userId) {
      this.mgmtService
        .listMachineTokens(this.userId, limit, offset)
        .then((resp) => {
          this.tokenResult = resp;
          if (resp.resultList) {
            this.dataSource.data = resp.resultList;
          }
          this.loadingSubject.next(false);
        })
        .catch((error: any) => {
          this.toast.showError(error);
          this.loadingSubject.next(false);
        });
    }
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }
}
