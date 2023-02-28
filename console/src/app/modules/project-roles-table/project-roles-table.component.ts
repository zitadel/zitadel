import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { MatLegacyTable as MatTable } from '@angular/material/legacy-table';
import { Router } from '@angular/router';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PaginatorComponent } from '../paginator/paginator.component';
import { ProjectRoleDetailDialogComponent } from '../project-role-detail-dialog/project-role-detail-dialog.component';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { ProjectRolesDataSource } from './project-roles-table-datasource';

@Component({
  selector: 'cnsl-project-roles-table',
  templateUrl: './project-roles-table.component.html',
  styleUrls: ['./project-roles-table.component.scss'],
})
export class ProjectRolesTableComponent implements OnInit {
  @Input() public projectId: string = '';
  @Input() public grantId: string = '';
  @Input() public disabled: boolean = false;
  @Input() public actionsVisible: boolean = false;
  @Input() public selectedKeys: string[] = [];
  @Input() public showSelectionActionButton: boolean = true;
  @ViewChild(PaginatorComponent) public paginator?: PaginatorComponent;
  @ViewChild(MatTable) public table?: MatTable<Role.AsObject>;
  public dataSource: ProjectRolesDataSource = new ProjectRolesDataSource(this.mgmtService);
  public selection: SelectionModel<string> = new SelectionModel<string>(true, []);
  @Output() public changedSelection: EventEmitter<Array<string>> = new EventEmitter();
  @Input() public displayedColumns: string[] = ['key', 'displayname', 'group', 'creationDate', 'changeDate', 'actions'];

  constructor(
    private mgmtService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    private router: Router,
  ) {}

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }

  public ngOnInit(): void {
    this.dataSource.loadRoles(this.projectId, this.grantId, 0, 25, 'asc');

    this.dataSource.rolesSubject.subscribe((roles) => {
      const selectedRoles: Role.AsObject[] = roles.filter((role) => this.selectedKeys.includes(role.key));
      this.selection.select(...selectedRoles.map((r) => r.key));
    });

    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  public selectAllOfGroup(group: string): void {
    const groupRoles: Role.AsObject[] = this.dataSource.rolesSubject.getValue().filter((role) => role.group === group);
    this.selection.select(...groupRoles.map((r) => r.key));
  }

  private loadRolesPage(): void {
    this.dataSource.loadRoles(this.projectId, this.grantId, this.paginator?.pageIndex ?? 0, this.paginator?.pageSize ?? 25);
  }

  public changePage(): void {
    this.loadRolesPage();
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.totalResult;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource.rolesSubject.value.forEach((row: Role.AsObject) => this.selection.select(row.key));
  }

  public deleteRole(role: Role.AsObject): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'PROJECT.ROLE.DIALOG.DELETE_TITLE',
        descriptionKey: 'PROJECT.ROLE.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const index = this.dataSource.rolesSubject.value.findIndex((iter) => iter.key === role.key);

        this.mgmtService.removeProjectRole(this.projectId, role.key).then(() => {
          this.toast.showInfo('PROJECT.TOAST.ROLEREMOVED', true);

          if (index > -1) {
            this.dataSource.rolesSubject.value.splice(index, 1);
            this.dataSource.rolesSubject.next(this.dataSource.rolesSubject.value);
          }
        });
      }
    });
  }

  public removeRole(role: Role.AsObject, index: number): void {
    this.mgmtService
      .removeProjectRole(this.projectId, role.key)
      .then(() => {
        this.toast.showInfo('PROJECT.TOAST.ROLEREMOVED', true);
        this.dataSource.rolesSubject.value.splice(index, 1);
        this.dataSource.rolesSubject.next(this.dataSource.rolesSubject.value);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public openDetailDialog(role: Role.AsObject): void {
    this.dialog.open(ProjectRoleDetailDialogComponent, {
      data: {
        role,
        projectId: this.projectId,
      },
      width: '400px',
    });
  }

  public refreshPage(): void {
    this.dataSource.loadRoles(this.projectId, this.grantId, this.paginator?.pageIndex ?? 0, this.paginator?.pageSize ?? 25);
  }

  public get selectionAllowed(): boolean {
    return this.displayedColumns.includes('select');
  }
}
