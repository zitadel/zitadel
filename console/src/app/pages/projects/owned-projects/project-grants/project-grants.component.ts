import { animate, state, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { Router } from '@angular/router';
import { tap } from 'rxjs/operators';
import { PageEvent, PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { GrantedProject, ProjectGrantState, Role } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { ProjectGrantsDataSource } from './project-grants-datasource';

@Component({
  selector: 'cnsl-project-grants',
  templateUrl: './project-grants.component.html',
  styleUrls: ['./project-grants.component.scss'],
  animations: [
    trigger('detailExpand', [
      state('collapsed', style({ height: '0px', minHeight: '0' })),
      state('expanded', style({ height: '*' })),
      transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
    ]),
  ],
})
export class ProjectGrantsComponent implements OnInit {
  public INITIAL_PAGESIZE: number = 10;

  @Input() public projectId: string = '';
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<GrantedProject.AsObject>;
  public dataSource: ProjectGrantsDataSource = new ProjectGrantsDataSource(this.mgmtService, this.toast);
  public selection: SelectionModel<GrantedProject.AsObject> = new SelectionModel<GrantedProject.AsObject>(true, []);
  public memberRoleOptions: Role.AsObject[] = [];
  public displayedColumns: string[] = ['grantedOrgName', 'state', 'creationDate', 'changeDate', 'roleNamesList', 'actions'];

  public ProjectGrantState: any = ProjectGrantState;

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
    this.dataSource.loadGrants(this.projectId, 0, this.INITIAL_PAGESIZE);
    this.getRoleOptions(this.projectId);
  }

  public loadGrantsPage(event: PageEvent): void {
    this.dataSource.loadGrants(this.projectId, event.pageIndex, event.pageSize);
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

  public getRoleOptions(projectId: string): void {
    this.mgmtService.listProjectRoles(projectId, 100, 0).then((resp) => {
      this.memberRoleOptions = resp.resultList;
    });
  }

  public updateRoles(grant: GrantedProject.AsObject, selectionChange: MatSelectChange): void {
    this.mgmtService
      .updateProjectGrant(grant.grantId, grant.projectId, selectionChange.value)
      .then(() => {
        this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTCHANGED', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public refreshPage(): void {
    this.selection.clear();
    this.dataSource.loadGrants(this.projectId, this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }

  public deleteGrant(grant: GrantedProject.AsObject): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'PROJECT.GRANT.DIALOG.DELETE_TITLE',
        descriptionKey: 'PROJECT.GRANT.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .removeProjectGrant(grant.grantId, grant.projectId)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.REMOVED', true);
            const data = this.dataSource.grantsSubject.getValue();
            this.selection.selected.forEach((item) => {
              const index = data.findIndex((i) => i.grantId === item.grantId);
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
}
