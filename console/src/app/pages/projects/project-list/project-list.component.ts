import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { GrantedProject, Project, ProjectQuery, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-project-list',
  templateUrl: './project-list.component.html',
  styleUrls: ['./project-list.component.scss'],
  animations: [
    trigger('list', [transition(':enter', [query('@animate', stagger(80, animateChild()))])]),
    trigger('animate', [
      transition(':enter', [
        style({ opacity: 0, transform: 'translateY(-100%)' }),
        animate('100ms', style({ opacity: 1, transform: 'translateY(0)' })),
      ]),
      transition(':leave', [
        style({ opacity: 1, transform: 'translateY(0)' }),
        animate('100ms', style({ opacity: 0, transform: 'translateY(100%)' })),
      ]),
    ]),
  ],
})
export class ProjectListComponent implements OnInit, OnDestroy {
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;

  public dataSource: MatTableDataSource<Project.AsObject | GrantedProject.AsObject> = new MatTableDataSource<
    Project.AsObject | GrantedProject.AsObject
  >();

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @Output() public emitAddProject: EventEmitter<void> = new EventEmitter();
  @Input() public projectType$: BehaviorSubject<any> = new BehaviorSubject(ProjectType.PROJECTTYPE_OWNED);
  public projectList: Project.AsObject[] | GrantedProject.AsObject[] = [];
  public displayedColumns: string[] = ['name', 'state', 'creationDate', 'changeDate', 'actions'];
  public selection: SelectionModel<Project.AsObject | GrantedProject.AsObject> = new SelectionModel<
    Project.AsObject | GrantedProject.AsObject
  >(true, []);

  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  public grid: boolean = true;
  public filterOpen: boolean = false;

  @Input() public zitadelProjectId: string = '';
  public ProjectState: any = ProjectState;
  public ProjectType: any = ProjectType;
  private destroy$: Subject<void> = new Subject();
  public INITIAL_PAGE_SIZE: number = 20;

  constructor(
    public translate: TranslateService,
    private mgmtService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    private router: Router,
  ) {}

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }

  public ngOnInit(): void {
    this.projectType$.pipe(takeUntil(this.destroy$)).subscribe((type) => {
      switch (type) {
        case ProjectType.PROJECTTYPE_OWNED:
          this.displayedColumns = ['name', 'state', 'creationDate', 'changeDate', 'actions'];
          break;
        case ProjectType.PROJECTTYPE_GRANTED:
          this.displayedColumns = ['name', 'projectOwnerName', 'state', 'creationDate', 'changeDate'];
          break;
      }

      this.getData(type, this.INITIAL_PAGE_SIZE, 0);
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected() ? this.selection.clear() : this.dataSource.data.forEach((row) => this.selection.select(row));
  }

  public changePage(type: ProjectType): void {
    this.getData(type, this.paginator.pageSize, this.paginator.pageSize * this.paginator.pageIndex);
  }

  public addProject(): void {
    this.emitAddProject.emit();
  }

  public applySearchQuery(type: ProjectType, searchQueries: ProjectQuery[]): void {
    this.selection.clear();
    this.getData(type, this.paginator.pageSize, this.paginator.pageSize * this.paginator.pageIndex, searchQueries);
  }

  private async getData(type: ProjectType, limit?: number, offset?: number, searchQueries?: ProjectQuery[]): Promise<void> {
    this.loadingSubject.next(true);
    switch (type) {
      case ProjectType.PROJECTTYPE_OWNED:
        this.mgmtService
          .listProjects(limit, offset, searchQueries)
          .then((resp) => {
            this.projectList = resp.resultList;
            if (resp.details?.totalResult) {
              this.totalResult = resp.details.totalResult;
            } else {
              this.totalResult = 0;
            }
            if (resp.details?.viewTimestamp) {
              this.viewTimestamp = resp.details?.viewTimestamp;
            }
            this.dataSource.data = this.projectList;
            this.loadingSubject.next(false);
          })
          .catch((error) => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
          });
        break;
      case ProjectType.PROJECTTYPE_GRANTED:
        this.mgmtService
          .listGrantedProjects(limit, offset, searchQueries)
          .then((resp) => {
            this.projectList = resp.resultList;
            if (resp.details?.totalResult) {
              this.totalResult = resp.details.totalResult;
            } else {
              this.totalResult = 0;
            }
            if (resp.details?.viewTimestamp) {
              this.viewTimestamp = resp.details?.viewTimestamp;
            }
            this.dataSource.data = this.projectList;
            this.loadingSubject.next(false);
          })
          .catch((error) => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
          });
        break;
    }
  }

  public reactivateSelectedProjects(): void {
    const promises = this.selection.selected.map((project) => {
      if ((project as Project.AsObject).id) {
        this.mgmtService.reactivateProject((project as Project.AsObject).id);
      }
    });

    Promise.all(promises)
      .then(() => {
        this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deactivateSelectedProjects(): void {
    const promises = this.selection.selected.map((project) => {
      if ((project as Project.AsObject).id) {
        this.mgmtService.deactivateProject((project as Project.AsObject).id);
      }
    });

    Promise.all(promises)
      .then(() => {
        this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public refreshPage(type: ProjectType): void {
    this.selection.clear();
    this.getData(type, this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }

  public deleteProject(id: string, name: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'PROJECT.PAGES.DIALOG.DELETE.TITLE',
        descriptionKey: 'PROJECT.PAGES.DIALOG.DELETE.DESCRIPTION',
        confirmationKey: 'PROJECT.PAGES.DIALOG.DELETE.TYPENAME',
        confirmation: name,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (this.zitadelProjectId && resp && id !== this.zitadelProjectId) {
        this.mgmtService
          .removeProject(id)
          .then(() => {
            this.toast.showInfo('PROJECT.TOAST.DELETED', true);
            setTimeout(() => {
              this.refreshPage(ProjectType.PROJECTTYPE_OWNED);
            }, 1000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }
}
