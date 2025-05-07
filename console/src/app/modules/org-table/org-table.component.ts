import { LiveAnnouncer } from '@angular/cdk/a11y';
import { Component, Input, ViewChild } from '@angular/core';
import { MatSort, Sort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, catchError, finalize, from, map, Observable, of, Subject, switchMap, takeUntil } from 'rxjs';
import { Org, OrgFieldName, OrgQuery, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { PaginatorComponent } from '../paginator/paginator.component';

enum OrgListSearchKey {
  NAME = 'NAME',
}

type Request = { limit: number; offset: number; queries: OrgQuery[] };

@Component({
  selector: 'cnsl-org-table',
  templateUrl: './org-table.component.html',
  styleUrls: ['./org-table.component.scss'],
})
export class OrgTableComponent {
  public orgSearchKey: OrgListSearchKey | undefined = undefined;

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild('input') public filter!: Input;

  public dataSource: MatTableDataSource<Org.AsObject> = new MatTableDataSource<Org.AsObject>([]);
  public displayedColumns: string[] = ['name', 'state', 'primaryDomain', 'creationDate', 'changeDate', 'actions'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public activeOrg!: Org.AsObject;
  public initialLimit: number = 20;
  public timestamp: Timestamp.AsObject | undefined = undefined;
  public totalResult: number = 0;
  public filterOpen: boolean = false;
  public OrgState: any = OrgState;
  public copied: string = '';
  @ViewChild(MatSort) public sort!: MatSort;

  private searchQueries: OrgQuery[] = [];
  private destroy$: Subject<void> = new Subject();
  private requestOrgs$: BehaviorSubject<Request> = new BehaviorSubject<Request>({
    limit: this.initialLimit,
    offset: 0,
    queries: [],
  });
  public defaultOrgId: string = '';
  private requestOrgsObservable$ = this.requestOrgs$.pipe(takeUntil(this.destroy$));

  constructor(
    private authService: GrpcAuthService,
    private mgmtService: ManagementService,
    private adminService: AdminService,
    private router: Router,
    private toast: ToastService,
    private _liveAnnouncer: LiveAnnouncer,
    private translate: TranslateService,
  ) {
    this.requestOrgs$.next({ limit: this.initialLimit, offset: 0, queries: this.searchQueries });
    this.authService.getActiveOrg().then((org) => (this.activeOrg = org));

    this.requestOrgsObservable$.pipe(switchMap((req) => this.loadOrgs(req))).subscribe((orgs) => {
      this.dataSource = new MatTableDataSource<Org.AsObject>(orgs);
    });

    this.mgmtService.getIAM().then((iam) => {
      this.defaultOrgId = iam.defaultOrgId;
    });
  }

  public loadOrgs(request: Request): Observable<Org.AsObject[]> {
    this.loadingSubject.next(true);

    let sortingField: OrgFieldName | undefined = undefined;
    if (this.sort?.active && this.sort?.direction)
      switch (this.sort.active) {
        case 'name':
          sortingField = OrgFieldName.ORG_FIELD_NAME_NAME;
          break;
      }

    return from(
      this.adminService.listOrgs(request.limit, request.offset, request.queries, sortingField, this.sort?.direction),
    ).pipe(
      map((resp) => {
        this.timestamp = resp.details?.viewTimestamp;
        this.totalResult = resp.details?.totalResult ?? 0;
        return resp.resultList;
      }),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
      finalize(() => this.loadingSubject.next(false)),
    );
  }

  public refresh(): void {
    this.requestOrgs$.next({
      limit: this.paginator.pageSize,
      offset: this.paginator.pageSize * this.paginator.pageIndex,
      queries: this.searchQueries,
    });
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
      this.refresh();
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  public setDefaultOrg(org: Org.AsObject) {
    this.adminService
      .setDefaultOrg(org.id)
      .then(() => {
        this.toast.showInfo('ORG.PAGES.DEFAULTORGSET', true);
        this.defaultOrgId = org.id;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public applySearchQuery(searchQueries: OrgQuery[]): void {
    this.searchQueries = searchQueries;
    this.requestOrgs$.next({
      limit: this.paginator ? this.paginator.pageSize : this.initialLimit,
      offset: this.paginator ? this.paginator.pageSize * this.paginator.pageIndex : 0,
      queries: this.searchQueries,
    });
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

  public setAndNavigateToOrg(org: Org.AsObject): void {
    if (org.state !== OrgState.ORG_STATE_REMOVED) {
      this.authService.setActiveOrg(org);
      this.router.navigate(['/org']);
    } else {
      this.translate.get('ORG.TOAST.ORG_WAS_DELETED').subscribe((data) => {
        this.toast.showInfo(data);
      });
    }
  }

  public changePage(): void {
    this.refresh();
  }

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }
}
