import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { Action, ActionState } from 'src/app/proto/generated/zitadel/action_pb';
import {
  CreateActionRequest,
  ListActionsResponse,
  UpdateActionRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddActionDialogComponent } from '../add-action-dialog/add-action-dialog.component';

@Component({
  selector: 'cnsl-action-table',
  templateUrl: './action-table.component.html',
  styleUrls: ['./action-table.component.scss']
})
export class ActionTableComponent implements OnInit {
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<Action.AsObject> = new MatTableDataSource<Action.AsObject>();
  public selection: SelectionModel<Action.AsObject> = new SelectionModel<Action.AsObject>(true, []);
  public actionsResult!: ListActionsResponse.AsObject;
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumns: string[] = ['select', 'id', 'name', 'state', 'timeout', 'allowedToFail'];

  @Output() public changedSelection: EventEmitter<Array<Action.AsObject>> = new EventEmitter();

  public ActionState: any = ActionState;
  constructor(public translate: TranslateService, private mgmtService: ManagementService, private dialog: MatDialog,
    private toast: ToastService) {
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
    this.isAllSelected() ?
      this.selection.clear() :
      this.dataSource.data.forEach(row => this.selection.select(row));
  }


  public changePage(event: PageEvent): void {
    this.getData(event.pageSize, event.pageIndex * event.pageSize);
  }

  public deleteKey(action: Action.AsObject): void {
    this.mgmtService.deleteAction(action.id).then(() => {
      this.selection.clear();
      this.toast.showInfo('FLOWS.TOAST.SELECTEDKEYSDELETED', true);
      this.getData(10, 0);
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public openAddAction(): void {
    const dialogRef = this.dialog.open(AddActionDialogComponent, {
      data: {},
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: CreateActionRequest) => {
      if (req) {
        this.mgmtService.createAction(req).then(resp => {
          this.refreshPage();
        }).catch((error: any) => {
          this.toast.showError(error);
        });
      }
    });
  }

  public openDialog(action: Action.AsObject): void {
    const dialogRef = this.dialog.open(AddActionDialogComponent, {
      data: {
        action: action,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: UpdateActionRequest) => {
      if (req) {
        this.mgmtService.updateAction(req).then(resp => {
          this.refreshPage();
        }).catch((error: any) => {
          this.toast.showError(error);
        });
      }
    });
  }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);
    
    this.mgmtService.listActions(limit, offset).then(resp => {
      this.actionsResult = resp;
      this.dataSource.data = this.actionsResult.resultList;
      this.loadingSubject.next(false);
    }).catch((error: any) => {
      this.toast.showError(error);
      this.loadingSubject.next(false);
    });
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }
}
