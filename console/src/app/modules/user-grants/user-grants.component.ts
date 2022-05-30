import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatInput } from '@angular/material/input';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { Router } from '@angular/router';
import { tap } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { Type, UserGrant, UserGrantQuery } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { ActionKeysType } from '../action-keys/action-keys.component';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { UserGrantRoleDialogComponent } from '../user-grant-role-dialog/user-grant-role-dialog.component';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { UserGrantContext, UserGrantsDataSource } from './user-grants-datasource';

export enum UserGrantListSearchKey {
  DISPLAY_NAME,
  ORG_NAME,
  PROJECT_NAME,
  ROLE_KEY,
}
@Component({
  selector: 'cnsl-user-grants',
  templateUrl: './user-grants.component.html',
  styleUrls: ['./user-grants.component.scss'],
  animations: [enterAnimations],
})
export class UserGrantsComponent implements OnInit, AfterViewInit {
  public userGrantListSearchKey: UserGrantListSearchKey | undefined = undefined;
  public UserGrantListSearchKey: any = UserGrantListSearchKey;

  public INITIAL_PAGE_SIZE: number = 50;
  @Input() context: UserGrantContext = UserGrantContext.NONE;
  @Input() refreshOnPreviousRoutes: string[] = [];

  public dataSource!: UserGrantsDataSource;
  public selection: SelectionModel<UserGrant.AsObject> = new SelectionModel<UserGrant.AsObject>(true, []);
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<UserGrant.AsObject>;

  @Input() disableWrite: boolean = false;
  @Input() disableDelete: boolean = false;

  @Input() userId: string = '';
  @Input() projectId: string = '';
  @Input() grantId: string = '';
  @ViewChild('input') public filter!: MatInput;

  public projectRoleOptions: Role.AsObject[] = [];
  public routerLink: any = [''];

  public loadedId: string = '';
  public loadedProjectId: string = '';
  public grantToEdit: string = '';

  public UserGrantContext: any = UserGrantContext;
  public Type: any = Type;
  public ActionKeysType: any = ActionKeysType;
  @Input() public type: Type | undefined = undefined;

  public filterOpen: boolean = false;

  constructor(
    private userService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    private router: Router,
  ) {}

  @Input() public displayedColumns: string[] = [
    'select',
    'user',
    'org',
    'projectId',
    'type',
    'creationDate',
    'changeDate',
    'roleNamesList',
    'actions',
  ];

  ngOnInit(): void {
    this.dataSource = new UserGrantsDataSource(this.userService);

    switch (this.context) {
      case UserGrantContext.OWNED_PROJECT:
        if (this.projectId) {
          // this.getProjectRoleOptions(this.projectId);
          this.routerLink = ['/grant-create', 'project', this.projectId];
        }
        break;
      case UserGrantContext.GRANTED_PROJECT:
        if (this.grantId) {
          this.routerLink = ['/grant-create', 'project', this.projectId, 'grant', this.grantId];
          // this.getGrantRoleOptions(this.grantId, this.projectId);
        }
        break;
      case UserGrantContext.USER:
        if (this.userId) {
          this.routerLink = ['/grant-create', 'user', this.userId];
        }
        break;
      case UserGrantContext.NONE:
        this.routerLink = ['/grant-create'];
    }

    this.loadGrantsPage(this.type);
  }

  public ngAfterViewInit(): void {
    this.paginator.page.pipe(tap(() => this.loadGrantsPage(this.type))).subscribe();
  }

  public setType(type: Type | undefined): void {
    this.type = type;
    this.loadGrantsPage(type);
  }

  public getType(grant: UserGrant.AsObject): string {
    if (grant.projectGrantId) {
      return 'Project Grant';
    } else if (grant.projectId) {
      return 'Project';
    } else {
      return '';
    }
  }

  public gotoCreateLink(rL: any): void {
    this.router.navigate(rL);
  }

  private loadGrantsPage(type: Type | undefined, searchQueries?: UserGrantQuery[]): void {
    let queries: UserGrantQuery[] = [];

    this.dataSource.loadGrants(
      this.context,
      this.paginator?.pageIndex ?? 0,
      this.paginator?.pageSize ?? this.INITIAL_PAGE_SIZE,
      {
        projectId: this.projectId,
        grantId: this.grantId,
        userId: this.userId,
      },
      searchQueries ? [...searchQueries, ...queries] : queries,
    );
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.grantsSubject.value.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource.grantsSubject.value.forEach((row) => this.selection.select(row));
  }

  public openEditDialog(grant: UserGrant.AsObject): void {
    const dialogRef = this.dialog.open(UserGrantRoleDialogComponent, {
      data: {
        projectId: grant.projectId,
        grantId: grant.projectGrantId,
        selectedRoleKeysList: grant.roleKeysList,
        i18nTitle: 'GRANTS.EDIT.TITLE',
      },
      width: '600px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && resp.roles) {
        this.userService
          .updateUserGrant(grant.id, grant.userId, resp.roles)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
            grant.roleKeysList = resp.roles;
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  updateRoles(grant: UserGrant.AsObject, selectionChange: MatSelectChange): void {
    this.userService
      .updateUserGrant(grant.id, grant.userId, selectionChange.value)
      .then(() => {
        this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deleteGrant(event: any, grant: UserGrant.AsObject): void {
    event.stopPropagation();

    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'GRANTS.DIALOG.DELETE_TITLE',
        descriptionKey: 'GRANTS.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.userService
          .removeUserGrant(grant.id, grant.userId)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.REMOVED', true);
            const data = this.dataSource.grantsSubject.getValue();

            const index = data.findIndex((i) => i.id === grant.id);
            if (index > -1) {
              data.splice(index, 1);
              this.dataSource.grantsSubject.next(data);
            }
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public deleteGrantSelection(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'GRANTS.DIALOG.BULK_DELETE_TITLE',
        descriptionKey: 'GRANTS.DIALOG.BULK_DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.userService
          .bulkRemoveUserGrant(this.selection.selected.map((grant) => grant.id))
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.BULKREMOVED', true);
            const data = this.dataSource.grantsSubject.getValue();
            this.selection.selected.forEach((item) => {
              const index = data.findIndex((i) => i.id === item.id);
              if (index > -1) {
                data.splice(index, 1);
                this.dataSource.grantsSubject.next(data);
              }
            });
            this.selection.clear();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public changePage(event?: PageEvent): void {
    this.dataSource.loadGrants(
      this.context,
      event?.pageIndex ?? this.paginator.pageIndex,
      event?.pageSize ?? this.paginator.pageSize,
      {
        projectId: this.projectId,
        grantId: this.grantId,
        userId: this.userId,
      },
    );
  }

  public applySearchQuery(searchQueries?: UserGrantQuery[]): void {
    this.selection.clear();
    this.loadGrantsPage(this.type, searchQueries);
  }

  public setFilter(key: UserGrantListSearchKey): void {
    setTimeout(() => {
      if (this.filter) {
        (this.filter as any).nativeElement.focus();
      }
    }, 100);

    if (this.userGrantListSearchKey !== key) {
      this.userGrantListSearchKey = key;
    } else {
      this.userGrantListSearchKey = undefined;
      this.loadGrantsPage(this.type);
    }
  }
}
