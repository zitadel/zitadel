import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
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
  public INITIAL_PAGE_SIZE: number = 50;
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
    this.loadRolesPage();
    this.selection.select(...this.selectedKeys);

    this.selection.changed.subscribe(() => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  private loadRolesPage(): void {
    this.dataSource.loadRoles(
      this.projectId,
      this.grantId,
      this.paginator?.pageIndex ?? 0,
      this.paginator?.pageSize ?? this.INITIAL_PAGE_SIZE,
    );
  }

  public changePage(): void {
    this.loadRolesPage();
  }

  private listIsAllSelected(list: string[]): boolean {
    return list.findIndex((key) => !this.selection.isSelected(key)) == -1;
  }

  private listIsAnySelected(list: string[]): boolean {
    return list.findIndex((key) => this.selection.isSelected(key)) != -1;
  }

  private listMasterToggle(list: string[]): void {
    if (this.listIsAllSelected(list)) this.selection.deselect(...list);
    else this.selection.select(...list);
  }

  private compilePageKeys(): string[] {
    return this.dataSource.rolesSubject.value.map((role) => role.key);
  }

  public masterToggle(): void {
    this.listMasterToggle(this.compilePageKeys());
  }

  public isAllSelected(): boolean {
    return this.listIsAllSelected(this.compilePageKeys());
  }

  public isAnySelected(): boolean {
    return this.listIsAnySelected(this.compilePageKeys());
  }

  public groupMasterToggle(group: string): void {
    this.listMasterToggle(this.dataSource.rolesSubject.value.filter((role) => role.group == group).map((role) => role.key));
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
        this.mgmtService.removeProjectRole(this.projectId, role.key).then(() => {
          this.toast.showInfo('PROJECT.TOAST.ROLEREMOVED', true);
          this.loadRolesPage();
        });
      }
    });
  }

  public openDetailDialog(role: Role.AsObject): void {
    const dialogRef = this.dialog.open(ProjectRoleDetailDialogComponent, {
      data: {
        role,
        projectId: this.projectId,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(() => this.loadRolesPage());
  }

  public refreshPage(): void {
    this.loadRolesPage();
  }

  public get selectionAllowed(): boolean {
    return this.displayedColumns.includes('select');
  }
}
