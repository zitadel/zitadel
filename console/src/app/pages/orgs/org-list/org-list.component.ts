import { AfterViewInit, Component, Input, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Org, OrgNameQuery, OrgQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

enum OrgListSearchKey {
  NAME = 'NAME',
}

@Component({
  selector: 'app-org-list',
  templateUrl: './org-list.component.html',
  styleUrls: ['./org-list.component.scss'],
  animations: [
    enterAnimations,
  ],
})
export class OrgListComponent implements AfterViewInit {
  public orgSearchKey: OrgListSearchKey | undefined = undefined;

  @ViewChild(MatPaginator) public paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;
  @ViewChild('input') public filter!: Input;

  public dataSource!: MatTableDataSource<Org.AsObject>;
  public displayedColumns: string[] = ['select', 'id', 'name', 'creationDate'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public activeOrg!: Org.AsObject;
  public OrgListSearchKey: any = OrgListSearchKey;

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
      query = new OrgQuery();
      const orgNameQuery = new OrgNameQuery();
      orgNameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
      orgNameQuery.setName(filter);
      query.setNameQuery(orgNameQuery);
    }

    from(this.authService.listMyProjectOrgs(limit, offset, query ? [query] : undefined)).pipe(
      map(resp => {
        return resp.resultList;
      }),
      catchError(() => of([])),
      finalize(() => this.loadingSubject.next(false)),
    ).subscribe(views => {
      this.dataSource = new MatTableDataSource(views);
      this.dataSource.paginator = this.paginator;
      console.log(views);
      this.dataSource.sort = this.sort;
    });
  }

  public selectOrg(item: Org.AsObject, event?: any): void {
    this.authService.setActiveOrg(item);
  }

  public refresh(): void {
    this.loadOrgs(this.paginator.length, this.paginator.pageSize * this.paginator.pageIndex);
  }

  public setFilter(key: OrgListSearchKey): void {
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
