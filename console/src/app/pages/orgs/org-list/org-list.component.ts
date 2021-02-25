import { AfterViewInit, Component, Input, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { MyProjectOrgSearchKey, MyProjectOrgSearchQuery, Org, SearchMethod } from 'src/app/proto/generated/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
    selector: 'app-org-list',
    templateUrl: './org-list.component.html',
    styleUrls: ['./org-list.component.scss'],
    animations: [
        enterAnimations,
    ],
})
export class OrgListComponent implements AfterViewInit {
    public orgSearchKey: MyProjectOrgSearchKey | undefined = undefined;

    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatSort) sort!: MatSort;
    @ViewChild('input') public filter!: Input;

    public dataSource!: MatTableDataSource<Org.AsObject>;
    public displayedColumns: string[] = ['select', 'id', 'name'];
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public activeOrg!: Org.AsObject;
    public MyProjectOrgSearchKey: any = MyProjectOrgSearchKey;

    constructor(
        private authService: GrpcAuthService,
        private router: Router,
    ) {
        this.loadOrgs(10, 0);

        this.authService.getActiveOrg().then(org => this.activeOrg = org);
    }

    public ngAfterViewInit(): void {
        this.loadOrgs(10, 0);
    }

    public loadOrgs(limit: number, offset: number, filter?: string): void {
        this.loadingSubject.next(true);
        let query;
        if (filter) {
            query = new MyProjectOrgSearchQuery();
            query.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS_IGNORE_CASE);
            query.setKey(MyProjectOrgSearchKey.MYPROJECTORGSEARCHKEY_ORG_NAME);
            query.setValue(filter);
        }

        from(this.authService.listMyProjectOrgs(limit, offset, query ? [query] : undefined)).pipe(
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

    public refresh(): void {
        this.loadOrgs(this.paginator.length, this.paginator.pageSize * this.paginator.pageIndex);
    }

    public setFilter(key: MyProjectOrgSearchKey): void {
        setTimeout(() => {
            if (this.filter) {
                (this.filter as any).nativeElement.focus();
            }
        }, 100);

        if (this.orgSearchKey !== key) {
            this.orgSearchKey = key;
        } else {
            this.orgSearchKey = undefined;
            this.refresh();
        }
    }

    public applyFilter(event: Event): void {
        const filterValue = (event.target as HTMLInputElement).value;
        this.loadOrgs(
            this.paginator.pageSize,
            this.paginator.pageIndex * this.paginator.pageSize,
            filterValue.trim().toLowerCase(),
        );
    }

    public setAndNavigateToOrg(org: Org.AsObject): void {
        this.authService.setActiveOrg(org);
        this.router.navigate(['/org']);
    }
}
