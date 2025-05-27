import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatInput } from '@angular/material/input';
import { MatTable } from '@angular/material/table';
import { Router } from '@angular/router';
import { tap } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import {
  Type,
  UserGrant as MgmtUserGrant,
  UserGrant,
} from 'src/app/proto/generated/zitadel/user_pb';
import { GroupGrantQuery, GroupGrantState, GroupGrant } from 'src/app/proto/generated/zitadel/group_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { ActionKeysType } from '../action-keys/action-keys.component';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { UserGrantRoleDialogComponent } from '../user-grant-role-dialog/user-grant-role-dialog.component';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { GroupGrantContext, GroupGrantsDataSource } from './group-grants-datasource';
import { Org, OrgIDQuery, OrgQuery, OrgState } from 'src/app/proto/generated/zitadel/org_pb';

export enum GroupGrantListSearchKey {
  DISPLAY_NAME,
  ORG_NAME,
  PROJECT_NAME,
  ROLE_KEY,
}

type GroupGrantAsObject = GroupGrant.AsObject;

@Component({
  selector: 'cnsl-group-grants',
  templateUrl: './group-grants.component.html',
  styleUrls: ['./group-grants.component.scss'],
  animations: [enterAnimations],
})
export class GroupGrantsComponent implements OnInit, AfterViewInit {
  public GroupGrantListSearchKey: any = GroupGrantListSearchKey;

  public INITIAL_PAGE_SIZE: number = 50;
  @Input() context: GroupGrantContext = GroupGrantContext.NONE;
  @Input() refreshOnPreviousRoutes: string[] = [];

  public dataSource: GroupGrantsDataSource = new GroupGrantsDataSource(this.groupService);
  public selection: SelectionModel<GroupGrantAsObject> = new SelectionModel<GroupGrantAsObject>(true, []);
  @ViewChild(PaginatorComponent) public paginator?: PaginatorComponent;
  @ViewChild(MatTable) public table?: MatTable<GroupGrantAsObject>;

  @Input() disableWrite: boolean = false;
  @Input() disableDelete: boolean = false;

  @Input() groupId: string = '';
  @Input() projectId: string = '';
  @Input() grantId: string = '';
  @ViewChild('input') public filter!: MatInput;

  public projectRoleOptions: Role.AsObject[] = [];
  public routerLink: any = undefined;

  public loadedId: string = '';
  public loadedProjectId: string = '';
  public grantToEdit: string = '';

  public Type: any = Type;
  public ActionKeysType: any = ActionKeysType;
  public GroupGrantState: any = GroupGrantState;
  @Input() public type: Type | undefined = undefined;

  public filterOpen: boolean = false;
  public myOrgs: Array<Org.AsObject> = [];
  constructor(
    private authService: GrpcAuthService,
    private groupService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    private router: Router,
  ) {}

  @Input() public displayedColumns: string[] = [
    'select',
    'groupName',
    'org',
    'projectId',
    'creationDate',
    'changeDate',
    'state',
    'roleNamesList',
    'actions',
  ];

  ngOnInit(): void {
    switch (this.context) {
      case GroupGrantContext.OWNED_PROJECT:
        if (this.projectId) {
          this.routerLink = ['/grant-create/groups/', 'project', this.projectId];
        }
        break;
      case GroupGrantContext.GRANTED_PROJECT:
        if (this.grantId) {
          this.routerLink = ['/grant-create/groups/', 'project', this.projectId, 'grant', this.grantId];
        }
        break;
      case GroupGrantContext.GROUP:
        if (this.groupId) {
          this.routerLink = ['/grant-create/groups/', 'group', this.groupId];
        }
        break;
      case GroupGrantContext.NONE:
        this.routerLink = ['/grant-create/groups/'];
    }
    this.loadGrantsPage();
  }

  public ngAfterViewInit(): void {
    this.paginator?.page.pipe(tap(() => this.loadGrantsPage())).subscribe();
  }

  public gotoCreateLink(rL: any): void {
    this.router.navigate(rL);
  }

  private loadGrantsPage(searchQueries?: GroupGrantQuery[]): void {
    let queries: GroupGrantQuery[] = [];

    this.dataSource.loadGrants(
      this.context,
      this.paginator?.pageIndex ?? 0,
      this.paginator?.pageSize ?? this.INITIAL_PAGE_SIZE,
      {
        projectId: this.projectId,
        grantId: this.grantId,
        groupId: this.groupId,
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

  public openEditDialog(grant: GroupGrantAsObject): void {
    const dialogRef = this.dialog.open(UserGrantRoleDialogComponent, {
      data: {
        projectId: grant.projectId,
        grantId: grant?.projectGrantId,
        selectedRoleKeysList: grant.roleKeysList,
        i18nTitle: 'GRANTS.EDIT.TITLE',
      },
      width: '600px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && resp.roles) {
        this.groupService
          .updateGroupGrant(
            (grant as GroupGrant.AsObject).id,
            grant.groupId,
            resp.roles,
          )
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

  public deleteGrant(event: any, grant: GroupGrant.AsObject): void {
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
        this.groupService
          .removeGroupGrant(grant.id, grant.groupId)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.REMOVED', true);
            const data = this.dataSource.grantsSubject.getValue();

            const index = data.findIndex(
              (i) => (i as GroupGrant.AsObject).id && (i as GroupGrant.AsObject).id === grant.id,
            );
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
        this.groupService
          .bulkRemoveGroupGrant(this.selection.selected.map((grant) => (grant as GroupGrant.AsObject).id))
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.BULKREMOVED', true);
            const data = this.dataSource.grantsSubject.getValue();
            this.selection.selected.forEach((item) => {
              const index = data.findIndex((i) => (i as GroupGrant.AsObject).id === (item as GroupGrant.AsObject).id);
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
      event?.pageIndex ?? this.paginator?.pageIndex ?? 0,
      event?.pageSize ?? this.paginator?.pageSize ?? this.INITIAL_PAGE_SIZE,
      {
        projectId: this.projectId,
        grantId: this.grantId,
        groupId: this.groupId,
      },
    );
  }

  public applySearchQuery(searchQueries?: GroupGrantQuery[]): void {
    this.selection.clear();
    this.loadGrantsPage(searchQueries);
  }

  public setFilter(key: GroupGrantListSearchKey): void {
    setTimeout(() => {
      if (this.filter) {
        (this.filter as any).nativeElement.focus();
      }
    }, 100);

    if (this.GroupGrantListSearchKey !== key) {
      this.GroupGrantListSearchKey = key;
    } else {
      this.GroupGrantListSearchKey = undefined;
      this.loadGrantsPage();
    }
  }
}
