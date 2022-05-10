import { Component, Input, ViewChild } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, catchError, finalize, from, map, Observable, of } from 'rxjs';
import { Org, OrgQuery, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';

enum OrgListSearchKey {
  NAME = 'NAME',
}

@Component({
  selector: 'cnsl-org-table',
  templateUrl: './org-table.component.html',
  styleUrls: ['./org-table.component.scss'],
})
export class OrgTableComponent {
  public orgSearchKey: OrgListSearchKey | undefined = undefined;

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild('input') public filter!: Input;

  public dataSource!: MatTableDataSource<Org.AsObject>;
  public displayedColumns: string[] = ['name', 'state', 'primaryDomain', 'creationDate', 'changeDate'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public activeOrg!: Org.AsObject;
  public OrgListSearchKey: any = OrgListSearchKey;
  public initialLimit: number = 20;
  public timestamp: Timestamp.AsObject | undefined = undefined;
  public totalResult: number = 0;
  public filterOpen: boolean = false;
  public OrgState: any = OrgState;
  public copied: string = '';
  constructor(private authService: GrpcAuthService, private router: Router) {
    this.loadOrgs(this.initialLimit, 0);
    this.authService.getActiveOrg().then((org) => (this.activeOrg = org));
  }

  public loadOrgs(limit: number, offset: number, queries?: OrgQuery[]): void {
    this.loadingSubject.next(true);

    from(this.authService.listMyProjectOrgs(limit, offset, queries))
      .pipe(
        map((resp) => {
          this.timestamp = resp.details?.viewTimestamp;
          this.totalResult = resp.details?.totalResult ?? 0;
          return resp.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((views) => {
        this.dataSource = new MatTableDataSource(views);
      });
  }

  public selectOrg(item: Org.AsObject, event?: any): void {
    this.authService.setActiveOrg(item);
  }

  public refresh(): void {
    this.loadOrgs(this.paginator.length, this.paginator.pageSize * this.paginator.pageIndex);
  }

  public applySearchQuery(searchQueries: OrgQuery[]): void {
    this.loadOrgs(this.paginator.pageSize, this.paginator.pageSize * this.paginator.pageIndex, searchQueries);
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
    this.authService.setActiveOrg(org);
    this.router.navigate(['/org']);
  }

  public changePage(event: PageEvent): void {
    this.loadOrgs(event.pageSize, event.pageIndex * event.pageSize);
  }

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }
}
