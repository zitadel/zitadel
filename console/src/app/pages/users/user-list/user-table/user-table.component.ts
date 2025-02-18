import { LiveAnnouncer } from '@angular/cdk/a11y';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, Signal, ViewChild } from '@angular/core';
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
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { UserService } from 'src/app/services/user.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { SearchQuery, SearchQuerySchema, Type, UserFieldName } from '@zitadel/proto/zitadel/user/v2/query_pb';
import { UserState, User as UserV2 } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { create } from '@bufbuild/protobuf';
import { Timestamp } from '@bufbuild/protobuf/wkt';

enum UserListSearchKey {
  FIRST_NAME,
  LAST_NAME,
  DISPLAY_NAME,
  USER_NAME,
  EMAIL,
}

@Component({
  selector: 'cnsl-user-table',
  templateUrl: './user-table.component.html',
  styleUrls: ['./user-table.component.scss'],
  animations: [enterAnimations],
})
export class UserTableComponent implements OnInit {
  public userSearchKey: UserListSearchKey | undefined = undefined;
  public Type = Type;
  @Input() public type: Type = Type.HUMAN;
  @Input() refreshOnPreviousRoutes: string[] = [];
  @Input() public canWrite$: Observable<boolean> = of(false);
  @Input() public canDelete$: Observable<boolean> = of(false);

  private user: Signal<User.AsObject | undefined> = toSignal(this.authService.user, { requireSync: true });

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatSort) public sort!: MatSort;
  public INITIAL_PAGE_SIZE: number = 20;

  public viewTimestamp!: Timestamp;
  public totalResult: number = 0;
  public dataSource: MatTableDataSource<UserV2> = new MatTableDataSource<UserV2>();
  public selection: SelectionModel<UserV2> = new SelectionModel<UserV2>(true, []);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumnsHuman: string[] = [
    'select',
    'displayName',
    'preferredLoginName',
    'email',
    'state',
    'creationDate',
    'changeDate',
    'actions',
  ];
  @Input() public displayedColumnsMachine: string[] = [
    'select',
    'displayName',
    'username',
    'creationDate',
    'changeDate',
    'state',
    'actions',
  ];

  @Output() public changedSelection: EventEmitter<Array<UserV2>> = new EventEmitter();

  public UserState: any = UserState;
  public UserListSearchKey: any = UserListSearchKey;

  public ActionKeysType: any = ActionKeysType;
  public filterOpen: boolean = false;

  private searchQueries: SearchQuery[] = [];
  constructor(
    private router: Router,
    public translate: TranslateService,
    private authService: GrpcAuthService,
    private userService: UserService,
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
        this.getData(this.INITIAL_PAGE_SIZE, 0, this.type, this.searchQueries).then();
      }

      if (params['deferredReload']) {
        setTimeout(() => {
          this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize, this.type).then();
        }, 2000);
      }
    });
  }

  public setType(type: Type): void {
    this.type = type;
    this.router
      .navigate([], {
        relativeTo: this.route,
        queryParams: {
          type: type === Type.HUMAN ? 'human' : type === Type.MACHINE ? 'machine' : 'human',
        },
        replaceUrl: true,
        queryParamsHandling: 'merge',
        skipLocationChange: false,
      })
      .then();
    this.getData(
      this.paginator.pageSize,
      this.paginator.pageIndex * this.paginator.pageSize,
      this.type,
      this.searchQueries,
    ).then();
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
    this.getData(event.pageSize, event.pageIndex * event.pageSize, this.type, this.searchQueries).then();
  }

  public deactivateSelectedUsers(): void {
    const usersToDeactivate = this.selection.selected
      .filter((u) => u.state === UserState.ACTIVE)
      .map((value) => {
        return this.userService.deactivateUser(value.userId);
      });

    Promise.all(usersToDeactivate)
      .then(() => {
        this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
        this.selection.clear();
        setTimeout(() => {
          this.refreshPage();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public reactivateSelectedUsers(): void {
    const usersToReactivate = this.selection.selected
      .filter((u) => u.state === UserState.INACTIVE)
      .map((value) => {
        return this.userService.reactivateUser(value.userId);
      });

    Promise.all(usersToReactivate)
      .then(() => {
        this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
        this.selection.clear();
        setTimeout(() => {
          this.refreshPage();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public gotoRouterLink(rL: any): Promise<boolean> {
    return this.router.navigate(rL);
  }

  private async getData(limit: number, offset: number, type: Type, searchQueries?: SearchQuery[]): Promise<void> {
    this.loadingSubject.next(true);

    let queryT = create(SearchQuerySchema, {
      query: {
        case: 'typeQuery',
        value: {
          type,
        },
      },
    });

    let sortingField: UserFieldName | undefined = undefined;
    if (this.sort?.active && this.sort?.direction)
      switch (this.sort.active) {
        case 'displayName':
          sortingField = UserFieldName.DISPLAY_NAME;
          break;
        case 'username':
          sortingField = UserFieldName.USER_NAME;
          break;
        case 'preferredLoginName':
          // TODO: replace with preferred username sorting once implemented
          sortingField = UserFieldName.USER_NAME;
          break;
        case 'email':
          sortingField = UserFieldName.EMAIL;
          break;
        case 'state':
          sortingField = UserFieldName.STATE;
          break;
        case 'creationDate':
          sortingField = UserFieldName.CREATION_DATE;
          break;
      }

    this.userService
      .listUsers(
        limit,
        offset,
        searchQueries?.length ? [queryT, ...searchQueries] : [queryT],
        sortingField,
        this.sort?.direction,
      )
      .then((resp) => {
        if (resp.details?.totalResult) {
          this.totalResult = Number(resp.details.totalResult);
        } else {
          this.totalResult = 0;
        }
        if (resp.details?.timestamp) {
          this.viewTimestamp = resp.details?.timestamp;
        }
        this.dataSource.data = resp.result;
        this.loadingSubject.next(false);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loadingSubject.next(false);
      });
  }

  public refreshPage(): void {
    this.getData(
      this.paginator.pageSize,
      this.paginator.pageIndex * this.paginator.pageSize,
      this.type,
      this.searchQueries,
    ).then();
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction} ending`).then();
      this.refreshPage();
    } else {
      this._liveAnnouncer.announce('Sorting cleared').then();
    }
  }

  public applySearchQuery(searchQueries: SearchQuery[]): void {
    this.selection.clear();
    this.searchQueries = searchQueries;
    this.getData(
      this.paginator ? this.paginator.pageSize : this.INITIAL_PAGE_SIZE,
      this.paginator ? this.paginator.pageIndex * this.paginator.pageSize : 0,
      this.type,
      searchQueries,
    ).then();
  }

  public deleteUser(user: UserV2): void {
    const authUserData = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'USER.DIALOG.DELETE_SELF_TITLE',
      warnSectionKey: 'USER.DIALOG.DELETE_SELF_DESCRIPTION',
      hintKey: 'USER.DIALOG.TYPEUSERNAME',
      hintParam: 'USER.DIALOG.DELETE_DESCRIPTION',
      confirmationKey: 'USER.DIALOG.USERNAME',
      confirmation: user.preferredLoginName,
    };

    const mgmtUserData = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'USER.DIALOG.DELETE_TITLE',
      warnSectionKey: 'USER.DIALOG.DELETE_DESCRIPTION',
      hintKey: 'USER.DIALOG.TYPEUSERNAME',
      hintParam: 'USER.DIALOG.DELETE_DESCRIPTION',
      confirmationKey: 'USER.DIALOG.USERNAME',
      confirmation: user.preferredLoginName,
    };

    if (user?.userId) {
      const authUser = this.user();
      const isMe = authUser?.id === user.userId;

      let dialogRef;

      if (isMe) {
        dialogRef = this.dialog.open(WarnDialogComponent, {
          data: authUserData,
          width: '400px',
        });
      } else {
        dialogRef = this.dialog.open(WarnDialogComponent, {
          data: mgmtUserData,
          width: '400px',
        });
      }

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.userService
            .deleteUser(user.userId)
            .then(() => {
              setTimeout(() => {
                this.refreshPage();
              }, 1000);
              this.selection.clear();
              this.toast.showInfo('USER.TOAST.DELETED', true);
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
    return selected ? selected.findIndex((user) => user.state !== UserState.ACTIVE) > -1 : false;
  }

  public get multipleDeactivatePossible(): boolean {
    const selected = this.selection.selected;
    return selected ? selected.findIndex((user) => user.state !== UserState.INACTIVE) > -1 : false;
  }
}
