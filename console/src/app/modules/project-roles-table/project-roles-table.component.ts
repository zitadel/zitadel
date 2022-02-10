import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { Router } from '@angular/router';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import {
  ProjectRoleDetailComponent,
} from '../../pages/projects/owned-projects/project-roles/project-role-detail/project-role-detail.component';
import { PaginatorComponent } from '../paginator/paginator.component';
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
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<Role.AsObject>;
  public dataSource!: ProjectRolesDataSource;
  public selection: SelectionModel<Role.AsObject> = new SelectionModel<Role.AsObject>(true, []);
  @Output() public changedSelection: EventEmitter<Array<Role.AsObject>> = new EventEmitter();
  @Input() public displayedColumns: string[] = ['key', 'displayname', 'group', 'creationDate', 'changeDate', 'actions'];

  constructor(
    private mgmtService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    private router: Router,
  ) {
    this.dataSource = new ProjectRolesDataSource(this.mgmtService);
  }

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }

  public ngOnInit(): void {
    this.dataSource.loadRoles(this.projectId, this.grantId, 0, 25, 'asc');

    this.dataSource.rolesSubject.subscribe((roles) => {
      const selectedRoles: Role.AsObject[] = roles.filter((role) => this.selectedKeys.includes(role.key));
      this.selection.select(...selectedRoles);
      console.log(this.selectedKeys, this.dataSource.rolesSubject.getValue(), selectedRoles);
    });

    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  public selectAllOfGroup(group: string): void {
    const groupRoles: Role.AsObject[] = this.dataSource.rolesSubject.getValue().filter((role) => role.group === group);
    this.selection.select(...groupRoles);
  }

  private loadRolesPage(): void {
    this.dataSource.loadRoles(this.projectId, this.grantId, this.paginator.pageIndex, this.paginator.pageSize);
  }

  public changePage(): void {
    this.selection.clear();
    this.loadRolesPage();
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.rolesSubject.value.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource.rolesSubject.value.forEach((row: Role.AsObject) => this.selection.select(row));
  }

  public deleteRole(role: Role.AsObject): Promise<any> {
    const index = this.dataSource.rolesSubject.value.findIndex((iter) => iter.key === role.key);

    return this.mgmtService.removeProjectRole(this.projectId, role.key).then(() => {
      this.toast.showInfo('PROJECT.TOAST.ROLEREMOVED', true);

      if (index > -1) {
        this.dataSource.rolesSubject.value.splice(index, 1);
        this.dataSource.rolesSubject.next(this.dataSource.rolesSubject.value);
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
    this.dialog.open(ProjectRoleDetailComponent, {
      data: {
        role,
        projectId: this.projectId,
      },
      width: '400px',
    });
  }

  public refreshPage(): void {
    this.dataSource.loadRoles(this.projectId, this.grantId, this.paginator.pageIndex, this.paginator.pageSize);
  }

  public get selectionAllowed(): boolean {
    return this.displayedColumns.includes('select');
  }
}
