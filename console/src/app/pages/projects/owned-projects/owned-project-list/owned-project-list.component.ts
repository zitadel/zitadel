import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { ProjectView } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-owned-project-list',
    templateUrl: './owned-project-list.component.html',
    styleUrls: ['./owned-project-list.component.scss'],
    animations: [
        trigger('list', [
            transition(':enter', [
                query('@animate',
                    stagger(80, animateChild()),
                ),
            ]),
        ]),
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
export class OwnedProjectListComponent implements OnInit, OnDestroy {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;

    public dataSource: MatTableDataSource<ProjectView.AsObject> =
        new MatTableDataSource<ProjectView.AsObject>();

    @ViewChild(MatPaginator) public paginator!: MatPaginator;

    public ownedProjectList: ProjectView.AsObject[] = [];
    public displayedColumns: string[] = ['select', 'name', 'state', 'creationDate', 'changeDate', 'actions'];
    public selection: SelectionModel<ProjectView.AsObject> = new SelectionModel<ProjectView.AsObject>(true, []);

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public grid: boolean = true;
    private subscription?: Subscription;

    public zitadelProjectId: string = '';

    constructor(private router: Router,
        private route: ActivatedRoute,
        public translate: TranslateService,
        private mgmtService: ManagementService,
        private toast: ToastService,
        private dialog: MatDialog,
    ) {
        this.mgmtService.GetIam().then(iam => {
            this.zitadelProjectId = iam.toObject().iamProjectId;
        });
    }

    public ngOnInit(): void {
        this.route.queryParams.pipe(take(1)).subscribe(params => {
            this.getData();
            if (params.deferredReload) {
                setTimeout(() => {
                    this.getData();
                }, 2000);
            }
        });
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.data.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.data.forEach(row => this.selection.select(row));
    }

    public changePage(event: PageEvent): void {
        this.getData(event.pageSize, event.pageIndex);
    }

    public addProject(): void {
        this.router.navigate(['/projects', 'create']);
    }

    private async getData(limit?: number, offset?: number): Promise<void> {
        this.loadingSubject.next(true);
        this.mgmtService.SearchProjects(limit, offset).then(res => {
            const response = res.toObject();
            this.ownedProjectList = response.resultList;
            this.totalResult = response.totalResult;
            if (this.totalResult > 10) {
                this.grid = false;
            }
            if (response.viewTimestamp) {
                this.viewTimestamp = response.viewTimestamp;
            }
            this.dataSource.data = this.ownedProjectList;
            this.loadingSubject.next(false);
        }).catch(error => {
            console.error(error);
            this.toast.showError(error);
            this.loadingSubject.next(false);
        });

        this.ownedProjectList = [];
    }

    public reactivateSelectedProjects(): void {
        const promises = this.selection.selected.map(project => {
            this.mgmtService.ReactivateProject(project.projectId);
        });

        Promise.all(promises).then(() => {
            this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }


    public deactivateSelectedProjects(): void {
        const promises = this.selection.selected.map(project => {
            this.mgmtService.DeactivateProject(project.projectId);
        });

        Promise.all(promises).then(() => {
            this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public refreshPage(): void {
        this.selection.clear();
        this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
    }

    public deleteProject(item: ProjectView.AsObject): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'PROJECT.PAGES.DIALOG.DELETE.TITLE',
                descriptionKey: 'PROJECT.PAGES.DIALOG.DELETE.DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (this.zitadelProjectId && resp && item.projectId !== this.zitadelProjectId) {
                this.mgmtService.RemoveProject(item.projectId).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.DELETED', true);
                    setTimeout(() => {
                        this.refreshPage();
                    }, 1000);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }
}
