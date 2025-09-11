import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatTable } from '@angular/material/table';
import { merge } from 'rxjs';
import { tap } from 'rxjs/operators';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { App, AppState } from 'src/app/proto/generated/zitadel/app_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

import { ProjectApplicationsDataSource } from './applications-datasource';

@Component({
  selector: 'cnsl-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.scss'],
})
export class ApplicationsComponent implements AfterViewInit, OnInit {
  @Input() public projectId: string = '';
  @Input() public disabled: boolean = false;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<App.AsObject>;
  public dataSource: ProjectApplicationsDataSource = new ProjectApplicationsDataSource(this.mgmtService);
  public selection: SelectionModel<App.AsObject> = new SelectionModel<App.AsObject>(true, []);

  public displayedColumns: string[] = ['name', 'type', 'state', 'creationDate', 'changeDate'];
  public AppState: any = AppState;
  constructor(private mgmtService: ManagementService) {}

  ngOnInit(): void {
    this.dataSource.loadApps(this.projectId, 0, 25);
  }

  public ngAfterViewInit(): void {
    merge(this.paginator.page)
      .pipe(tap(() => this.loadRolesPage()))
      .subscribe();
  }

  private loadRolesPage(): void {
    this.dataSource.loadApps(this.projectId, this.paginator.pageIndex, this.paginator.pageSize);
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.appsSubject.value.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource.appsSubject.value.forEach((row: App.AsObject) => this.selection.select(row));
  }

  public refreshPage(): void {
    this.dataSource.loadApps(this.projectId, this.paginator.pageIndex, this.paginator.pageSize);
  }
}
