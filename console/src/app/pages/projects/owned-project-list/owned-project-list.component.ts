import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable, Subscription } from 'rxjs';
import { ProjectView } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
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
    public dataSource: MatTableDataSource<ProjectView.AsObject> =
        new MatTableDataSource<ProjectView.AsObject>();

    public ownedProjectList: ProjectView.AsObject[] = [];
    public displayedColumns: string[] = ['select', 'name', 'orgName', 'orgDomain', 'type', 'state', 'creationDate', 'changeDate'];
    public selection: SelectionModel<ProjectView.AsObject> = new SelectionModel<ProjectView.AsObject>(true, []);

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public grid: boolean = true;
    private subscription?: Subscription;

    constructor(private router: Router,
        public translate: TranslateService,
        private projectService: ProjectService,
        private toast: ToastService,
    ) { }

    public ngOnInit(): void {
        console.log('asdf');
        this.getData(10, 0);
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

    private async getData(limit: number, offset: number): Promise<void> {
        console.log('getprojects');
        this.loadingSubject.next(true);
        this.projectService.SearchProjects(limit, offset).then(res => {
            this.ownedProjectList = res.toObject().resultList;
            this.totalResult = res.toObject().totalResult;
            this.dataSource.data = this.ownedProjectList;
            this.loadingSubject.next(false);
            console.log(this.ownedProjectList);
        }).catch(error => {
            console.error(error);
            this.toast.showError(error.message);
            this.loadingSubject.next(false);
        });
    }

    public dateFromTimestamp(date: Timestamp.AsObject): any {
        const ts: Date = new Date(date.seconds * 1000 + date.nanos / 1000);
        return ts;
    }

    public reactivateSelectedProjects(): void {
        const promises = this.selection.selected.map(project => {
            this.projectService.ReactivateProject(project.projectId);
        });

        Promise.all(promises).then(() => {
            this.toast.showInfo('Reactivated selected projects successfully');
        }).catch(error => {
            this.toast.showInfo(error.message);
        });
    }


    public deactivateSelectedProjects(): void {
        const promises = this.selection.selected.map(project => {
            this.projectService.DeactivateProject(project.projectId);
        });

        Promise.all(promises).then(() => {
            this.toast.showInfo('Deactivated selected projects Successfully');
        }).catch(error => {
            this.toast.showInfo(error.message);
        });
    }
}
