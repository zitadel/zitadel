import { SelectionModel } from '@angular/cdk/collections';
import { Component, DestroyRef, EventEmitter, Input, OnInit, Output, signal, Signal, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort, SortDirection } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import {
  combineLatestWith,
  defer,
  delay,
  distinctUntilChanged,
  EMPTY,
  from,
  Observable,
  of,
  ReplaySubject,
  shareReplay,
  switchMap,
  toArray,
} from 'rxjs';
import { catchError, filter, finalize, map, startWith, take } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { ToastService } from 'src/app/services/toast.service';
import { UserService } from 'src/app/services/user.service';
import { takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop';
import { SearchQuery as UserSearchQuery } from 'src/app/proto/generated/zitadel/user_pb';
import { Type, UserFieldName } from '@zitadel/proto/zitadel/user/v2/query_pb';
import { UserState, User } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { ListUsersRequestSchema, ListUsersResponse } from '@zitadel/proto/zitadel/user/v2/user_service_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { UserState as UserStateV1 } from 'src/app/proto/generated/zitadel/user_pb';

type ListUsersRequest = MessageInitShape<typeof ListUsersRequestSchema>;
type QueriesArray = NonNullable<ListUsersRequest['queries']>;
type QueryWrapper = QueriesArray extends readonly (infer T)[] ? T : never;
type Query = NonNullable<QueryWrapper extends { query?: infer Q } ? Q : never>;

@Component({
  selector: 'cnsl-user-table',
  templateUrl: './user-table.component.html',
  styleUrls: ['./user-table.component.scss'],
  animations: [enterAnimations],
})
export class UserTableComponent implements OnInit {
  protected readonly Type = Type;
  protected readonly refresh$ = new ReplaySubject<true>(1);

  @Input() public canWrite$: Observable<boolean> = of(false);
  @Input() public canDelete$: Observable<boolean> = of(false);

  protected readonly dataSize: Signal<number>;
  protected readonly loading = signal(false);

  private readonly paginator$ = new ReplaySubject<PaginatorComponent>(1);
  @ViewChild(PaginatorComponent) public set paginator(paginator: PaginatorComponent) {
    this.paginator$.next(paginator);
  }
  private readonly sort$ = new ReplaySubject<MatSort>(1);
  @ViewChild(MatSort) public set sort(sort: MatSort) {
    this.sort$.next(sort);
  }

  protected readonly INITIAL_PAGE_SIZE = 20;

  protected readonly dataSource: MatTableDataSource<User> = new MatTableDataSource<User>();
  protected readonly selection: SelectionModel<User> = new SelectionModel<User>(true, []);
  protected readonly users$: Observable<ListUsersResponse>;
  protected readonly type$: Observable<Type>;
  protected readonly searchQueries$ = new ReplaySubject<UserSearchQuery[]>(1);
  protected readonly myUser: Signal<User | undefined>;

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

  @Output() public changedSelection: EventEmitter<Array<User>> = new EventEmitter();

  protected readonly UserState = UserState;

  protected ActionKeysType = ActionKeysType;
  protected filterOpen: boolean = false;

  constructor(
    protected readonly router: Router,
    public readonly translate: TranslateService,
    private readonly userService: UserService,
    private readonly toast: ToastService,
    private readonly dialog: MatDialog,
    private readonly route: ActivatedRoute,
    private readonly destroyRef: DestroyRef,
    private readonly authenticationService: AuthenticationService,
    private readonly authService: GrpcAuthService,
  ) {
    this.type$ = this.getType$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.users$ = this.getUsers(this.type$).pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.myUser = toSignal(this.getMyUser());

    this.dataSize = toSignal(
      this.users$.pipe(
        map((users) => Number(users.details?.totalResult ?? users.result.length)),
        distinctUntilChanged(),
      ),
      { initialValue: 0 },
    );
  }

  ngOnInit(): void {
    this.selection.changed.pipe(takeUntilDestroyed(this.destroyRef)).subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });

    this.users$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((users) => (this.dataSource.data = users.result));

    this.route.queryParamMap
      .pipe(
        map((params) => params.get('deferredReload')),
        filter(Boolean),
        take(1),
        delay(2000),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe(() => this.refresh$.next(true));
  }

  setType(type: Type) {
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
  }

  private getMyUser() {
    return this.userService.user$.pipe(
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }

  private getType$(): Observable<Type> {
    return this.route.queryParamMap.pipe(
      map((params) => params.get('type')),
      filter(Boolean),
      map((type) => (type === 'machine' ? Type.MACHINE : Type.HUMAN)),
      startWith(Type.HUMAN),
      distinctUntilChanged(),
    );
  }

  private getDirection$() {
    return this.sort$.pipe(
      switchMap((sort) =>
        sort.sortChange.pipe(
          map(({ direction }) => direction),
          startWith(sort.direction),
        ),
      ),
      distinctUntilChanged(),
    );
  }

  private getSortingColumn$() {
    return this.sort$.pipe(
      switchMap((sort) =>
        sort.sortChange.pipe(
          map(({ active }) => active),
          startWith(sort.active),
        ),
      ),
      map((active) => {
        switch (active) {
          case 'displayName':
            return UserFieldName.DISPLAY_NAME;
          case 'username':
            return UserFieldName.USER_NAME;
          case 'preferredLoginName':
            // TODO: replace with preferred username sorting once implemented
            return UserFieldName.USER_NAME;
          case 'email':
            return UserFieldName.EMAIL;
          case 'state':
            return UserFieldName.STATE;
          case 'creationDate':
            return UserFieldName.CREATION_DATE;
          default:
            return undefined;
        }
      }),
      distinctUntilChanged(),
    );
  }

  private getQueries(type$: Observable<Type>): Observable<Query[]> {
    const activeOrgId$ = this.getActiveOrgId();

    return this.searchQueries$.pipe(
      startWith([]),
      combineLatestWith(type$, activeOrgId$),
      switchMap(([queries, type, organizationId]) =>
        from(queries).pipe(
          map((query) => this.searchQueryToV2(query.toObject())),
          startWith({ case: 'typeQuery' as const, value: { type } }),
          startWith(organizationId ? { case: 'organizationIdQuery' as const, value: { organizationId } } : undefined),
          filter(Boolean),
          toArray(),
        ),
      ),
    );
  }

  private searchQueryToV2(query: UserSearchQuery.AsObject): Query | undefined {
    if (query.userNameQuery) {
      return {
        case: 'userNameQuery' as const,
        value: {
          userName: query.userNameQuery.userName,
          method: query.userNameQuery.method as unknown as any,
        },
      };
    } else if (query.displayNameQuery) {
      return {
        case: 'displayNameQuery' as const,
        value: {
          displayName: query.displayNameQuery.displayName,
          method: query.displayNameQuery.method as unknown as any,
        },
      };
    } else if (query.emailQuery) {
      return {
        case: 'emailQuery' as const,
        value: {
          emailAddress: query.emailQuery.emailAddress,
          method: query.emailQuery.method as unknown as any,
        },
      };
    } else if (query.stateQuery) {
      return {
        case: 'stateQuery' as const,
        value: {
          state: this.toV2State(query.stateQuery.state),
        },
      };
    } else {
      return undefined;
    }
  }

  private toV2State(state: UserStateV1) {
    switch (state) {
      case UserStateV1.USER_STATE_ACTIVE:
        return UserState.ACTIVE;
      case UserStateV1.USER_STATE_INACTIVE:
        return UserState.INACTIVE;
      case UserStateV1.USER_STATE_DELETED:
        return UserState.DELETED;
      case UserStateV1.USER_STATE_LOCKED:
        return UserState.LOCKED;
      case UserStateV1.USER_STATE_INITIAL:
        return UserState.INITIAL;
      default:
        throw new Error(`Invalid UserState ${state}`);
    }
  }

  private getUsers(type$: Observable<Type>) {
    const queries$ = this.getQueries(type$);
    const direction$ = this.getDirection$();
    const sortingColumn$ = this.getSortingColumn$();

    const page$ = this.paginator$.pipe(switchMap((paginator) => paginator.page));
    const pageSize$ = page$.pipe(
      map(({ pageSize }) => pageSize),
      startWith(this.INITIAL_PAGE_SIZE),
      distinctUntilChanged(),
    );
    const pageIndex$ = page$.pipe(
      map(({ pageIndex }) => pageIndex),
      startWith(0),
      distinctUntilChanged(),
    );

    return this.refresh$.pipe(
      startWith(true),
      combineLatestWith(queries$, direction$, sortingColumn$, pageSize$, pageIndex$),
      switchMap(([_, queries, direction, sortingColumn, pageSize, pageIndex]) => {
        return this.fetchUsers(queries, direction, sortingColumn, pageSize, pageIndex);
      }),
    );
  }

  private fetchUsers(
    queries: Query[],
    direction: SortDirection,
    sortingColumn: UserFieldName | undefined,
    pageSize: number,
    pageIndex: number,
  ) {
    return defer(() => {
      const req = {
        query: {
          limit: pageSize,
          offset: BigInt(pageIndex * pageSize),
          asc: direction === 'asc',
        },
        sortingColumn,
        queries: queries.map((query) => ({ query })),
      };

      this.loading.set(true);
      return this.userService.listUsers(req);
    }).pipe(
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
      finalize(() => this.loading.set(false)),
    );
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected() ? this.selection.clear() : this.dataSource.data.forEach((row) => this.selection.select(row));
  }

  public async deactivateSelectedUsers(): Promise<void> {
    const usersToDeactivate = this.selection.selected
      .filter((u) => u.state === UserState.ACTIVE)
      .map((value) => {
        return this.userService.deactivateUser(value.userId);
      });

    try {
      await Promise.all(usersToDeactivate);
    } catch (error) {
      this.toast.showError(error);
      return;
    }

    this.toast.showInfo('USER.TOAST.SELECTEDDEACTIVATED', true);
    this.selection.clear();
    setTimeout(() => {
      this.refresh$.next(true);
    }, 1000);
  }

  public async reactivateSelectedUsers(): Promise<void> {
    const usersToReactivate = this.selection.selected
      .filter((u) => u.state === UserState.INACTIVE)
      .map((value) => {
        return this.userService.reactivateUser(value.userId);
      });

    try {
      await Promise.all(usersToReactivate);
    } catch (error) {
      this.toast.showError(error);
      return;
    }

    this.toast.showInfo('USER.TOAST.SELECTEDREACTIVATED', true);
    this.selection.clear();
    setTimeout(() => {
      this.refresh$.next(true);
    }, 1000);
  }

  public deleteUser(user: User): void {
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
      const authUser = this.myUser();
      console.log('my user', authUser);
      const isMe = authUser?.userId === user.userId;

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

      dialogRef
        .afterClosed()
        .pipe(
          filter(Boolean),
          switchMap(() => this.userService.deleteUser(user.userId)),
        )
        .subscribe({
          next: () => {
            setTimeout(() => {
              this.refresh$.next(true);
            }, 1000);
            this.selection.clear();
            this.toast.showInfo('USER.TOAST.DELETED', true);
          },
          error: (err) => this.toast.showError(err),
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

  private getActiveOrgId() {
    return this.authenticationService.authenticationChanged.pipe(
      startWith(true),
      filter(Boolean),
      switchMap(() =>
        from(this.authService.getActiveOrg()).pipe(
          catchError((err) => {
            this.toast.showError(err);
            return of(undefined);
          }),
        ),
      ),
      map((org) => org?.id),
    );
  }
}
