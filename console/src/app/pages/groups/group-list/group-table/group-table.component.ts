import { LiveAnnouncer } from '@angular/cdk/a11y';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort, Sort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable, of } from 'rxjs';
import { take } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { PageEvent, PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';
import { Group, GroupState, GroupQuery, GroupFieldName } from 'src/app/proto/generated/zitadel/group_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-group-table',
  templateUrl: './group-table.component.html',
  styleUrls: ['./group-table.component.scss'],
  animations: [enterAnimations],
})
export class GroupTableComponent implements OnInit {
  @Input() refreshOnPreviousRoutes: string[] = [];
  @Input() public canWrite$: Observable<boolean> = of(false);
  @Input() public canDelete$: Observable<boolean> = of(false);

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatSort) public sort!: MatSort;
  public INITIAL_PAGE_SIZE: number = 20;

  public viewTimestamp!: Timestamp.AsObject;
  public totalResult: number = 0;
  public dataSource: MatTableDataSource<Group.AsObject> = new MatTableDataSource<Group.AsObject>();
  public selection: SelectionModel<Group.AsObject> = new SelectionModel<Group.AsObject>(true, []);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumnsHuman: string[] = [
    'select',
    'name',
    'state',
    'creationDate',
    'changeDate',
    'actions',
  ];

  @Output() public changedSelection: EventEmitter<Array<Group.AsObject>> = new EventEmitter();

  public GroupState: any = GroupState;
  public ActionKeysType: any = ActionKeysType;
  public filterOpen: boolean = false;

  private searchQueries: GroupQuery[] = [];
  constructor(
    private router: Router,
    public translate: TranslateService,
    private authService: GrpcAuthService,
    private groupService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    private route: ActivatedRoute,
    private _liveAnnouncer: LiveAnnouncer,
  ) {
    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  ngOnInit(): void {
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      if (!params['filter']) {
        this.getData(this.INITIAL_PAGE_SIZE, 0, this.searchQueries);
      }

      if (params['deferredReload']) {
        setTimeout(() => {
          this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
        }, 2000);
      }
    });
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
    this.selection.clear();
    this.getData(event.pageSize, event.pageIndex * event.pageSize, this.searchQueries);
  }

  public deactivateSelectedGroups(): void {
    Promise.all(
      this.selection.selected
        .filter((u) => u.state === GroupState.GROUP_STATE_ACTIVE)
        .map((value) => {
          return this.groupService.deactivateGroup(value.id);
        }),
    )
      .then(() => {
        this.toast.showInfo('GROUP.TOAST.SELECTEDDEACTIVATED', true);
        this.selection.clear();
        setTimeout(() => {
          this.refreshPage();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public reactivateSelectedGroups(): void {
    Promise.all(
      this.selection.selected
        .filter((u) => u.state === GroupState.GROUP_STATE_INACTIVE)
        .map((value) => {
          return this.groupService.reactivateGroup(value.id);
        }),
    )
      .then(() => {
        this.toast.showInfo('GROUP.TOAST.SELECTEDREACTIVATED', true);
        this.selection.clear();
        setTimeout(() => {
          this.refreshPage();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public gotoRouterLink(rL: any): void {
    this.router.navigate(rL);
  }

  private async getData(limit: number, offset: number, searchQueries?: GroupQuery[]): Promise<void> {
    this.loadingSubject.next(true);

    let queryT = new GroupQuery();
    let sortingField: GroupFieldName | undefined = undefined;
    if (this.sort?.active && this.sort?.direction)
      switch (this.sort.active) {
        case 'name':
          sortingField = GroupFieldName.GROUP_FIELD_NAME_NAME;
          break;
        case 'creationDate':
          sortingField = GroupFieldName.GROUP_FIELD_NAME_CREATION_DATE;
          break;
      }

    this.groupService
      .listGroups(
        limit,
        offset,
        searchQueries,
        sortingField,
        this.sort?.direction,
      )
      .then((resp) => {
        if (resp.details?.totalResult) {
          this.totalResult = resp.details?.totalResult;
        } else {
          this.totalResult = 0;
        }
        if (resp.details?.viewTimestamp) {
          this.viewTimestamp = resp.details?.viewTimestamp;
        }
        this.dataSource.data = resp.resultList;
        this.loadingSubject.next(false);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loadingSubject.next(false);
      });
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.searchQueries);
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction} ending`);
      this.refreshPage();
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  public applySearchQuery(searchQueries: GroupQuery[]): void {
    this.selection.clear();
    this.searchQueries = searchQueries;
    this.getData(
      this.paginator ? this.paginator.pageSize : this.INITIAL_PAGE_SIZE,
      this.paginator ? this.paginator.pageIndex * this.paginator.pageSize : 0,
      searchQueries,
    );
  }

  public deleteGroup(group: Group.AsObject): void {
    const mgmtGroupData = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'GROUP.DIALOG.DELETE_TITLE',
      warnSectionKey: 'GROUP.DIALOG.DELETE_DESCRIPTION',
      hintKey: 'GROUP.DIALOG.TYPEUSERNAME',
      hintParam: 'GROUP.DIALOG.DELETE_DESCRIPTION',
      confirmationKey: 'GROUP.DIALOG.GROUPNAME',
      confirmation: group.name,
    };

    if (group && group.id) {
      let dialogRef;
      dialogRef = this.dialog.open(WarnDialogComponent, {
        data: mgmtGroupData,
        width: '400px',
      });

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.groupService
            .removeGroup(group.id)
            .then(() => {
              setTimeout(() => {
                this.refreshPage();
              }, 1000);
              this.selection.clear();
              this.toast.showInfo('GROUP.TOAST.DELETED', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }

  public get multipleActivatePossible(): boolean {
    const selected = this.selection.selected;
    return selected ? selected.findIndex((group) => group.state !== GroupState.GROUP_STATE_ACTIVE) > -1 : false;
  }

  public get multipleDeactivatePossible(): boolean {
    const selected = this.selection.selected;
    return selected ? selected.findIndex((group) => group.state !== GroupState.GROUP_STATE_INACTIVE) > -1 : false;
  }

}
