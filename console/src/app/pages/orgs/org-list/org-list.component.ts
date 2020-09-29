import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { MyProjectOrgSearchKey, MyProjectOrgSearchQuery, Org, SearchMethod } from 'src/app/proto/generated/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-org-list',
    templateUrl: './org-list.component.html',
    styleUrls: ['./org-list.component.scss'],
})
export class OrgListComponent implements AfterViewInit {
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatSort) sort!: MatSort;

    public dataSource!: MatTableDataSource<Org.AsObject>;
    public displayedColumns: string[] = ['select', 'id', 'name'];
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public activeOrg!: Org.AsObject;
    public orgList: Org.AsObject[] = [];

    public selection: SelectionModel<Org.AsObject> = new SelectionModel<Org.AsObject>(true, []);
    public selectedIndex: number = -1;
    public loading: boolean = false;

    public notPinned: Array<Org.AsObject> = [];

    constructor(
        private authService: GrpcAuthService,
        private toast: ToastService,
        private router: Router,
    ) {
        this.loading = true;
        this.loadOrgs(10, 0);

        this.authService.GetActiveOrg().then(org => this.activeOrg = org);
    }

    public ngAfterViewInit(): void {
        this.loadOrgs(10, 0);
    }

    public loadOrgs(limit: number, offset: number, filter?: string): void {
        this.loadingSubject.next(true);
        let query;
        if (filter) {
            query = new MyProjectOrgSearchQuery();
            query.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS);
            query.setKey(MyProjectOrgSearchKey.MYPROJECTORGSEARCHKEY_ORG_NAME);
            query.setValue(filter);
        }

        from(this.authService.SearchMyProjectOrgs(limit, offset, query ? [query] : undefined)).pipe(
            map(resp => {
                return resp.toObject().resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(views => {
            this.dataSource = new MatTableDataSource(views);
            this.dataSource.paginator = this.paginator;
            this.dataSource.sort = this.sort;
        });
    }

    public selectOrg(item: Org.AsObject, event?: any): void {
        this.authService.setActiveOrg(item);
    }

    public routeToOrg(item: Org.AsObject): void {
        this.router.navigate(['/orgs', item.id]);
    }

    public refresh(): void {
        this.loadOrgs(this.paginator.length, this.paginator.pageSize * this.paginator.pageIndex);
    }

    public applyFilter(event: Event): void {
        const filterValue = (event.target as HTMLInputElement).value;
        this.loadOrgs(
            this.paginator.pageSize,
            this.paginator.pageIndex * this.paginator.pageSize,
            filterValue.trim().toLowerCase(),
        );
    }
}
