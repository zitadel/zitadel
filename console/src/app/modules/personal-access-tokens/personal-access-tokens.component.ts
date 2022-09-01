import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Moment } from 'moment';
import { BehaviorSubject, Observable } from 'rxjs';
import { Key } from 'src/app/proto/generated/zitadel/auth_n_key_pb';
import { ListPersonalAccessTokensResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { PersonalAccessToken } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddTokenDialogComponent } from '../add-token-dialog/add-token-dialog.component';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { ShowTokenDialogComponent } from '../show-token-dialog/show-token-dialog.component';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';

@Component({
  selector: 'cnsl-personal-access-tokens',
  templateUrl: './personal-access-tokens.component.html',
  styleUrls: ['./personal-access-tokens.component.scss'],
})
export class PersonalAccessTokensComponent implements OnInit {
  @Input() userId!: string;

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<PersonalAccessToken.AsObject> = new MatTableDataSource<PersonalAccessToken.AsObject>(
    [],
  );
  public selection: SelectionModel<PersonalAccessToken.AsObject> = new SelectionModel<PersonalAccessToken.AsObject>(
    true,
    [],
  );
  public keyResult: ListPersonalAccessTokensResponse.AsObject | undefined = undefined;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumns: string[] = ['id', 'creationDate', 'expirationDate', 'actions'];

  @Output() public changedSelection: EventEmitter<Array<PersonalAccessToken.AsObject>> = new EventEmitter();

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
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.PERSONALACCESSTOKEN.DELETE.TITLE',
        descriptionKey: 'USER.PERSONALACCESSTOKEN.DELETE.DESCRIPTION',
      },
      width: '400px',
    });
    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .removePersonalAccessToken(key.id, this.userId)
          .then(() => {
            this.selection.clear();
            this.toast.showInfo('USER.PERSONALACCESSTOKEN.DELETED', true);
            this.getData(10, 0);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public openAddKey(): void {
    const dialogRef = this.dialog.open(AddTokenDialogComponent, {
      data: {},
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
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

        this.mgmtService
          .addPersonalAccessToken(this.userId, date)
          .then((response) => {
            if (response) {
              setTimeout(() => {
                this.refreshPage();
              }, 1000);

              this.dialog.open(ShowTokenDialogComponent, {
                data: {
                  token: response,
                },
                width: '400px',
              });
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });
  }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    if (this.userId) {
      this.mgmtService
        .listPersonalAccessTokens(this.userId, limit, offset)
        .then((resp) => {
          this.keyResult = resp;
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
