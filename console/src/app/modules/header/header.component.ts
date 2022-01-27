import { Component, ElementRef, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Router } from '@angular/router';
import { BehaviorSubject, catchError, debounceTime, finalize, from, map, Observable, of, Subject } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Org, OrgNameQuery, OrgQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';

@Component({
  selector: 'cnsl-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent implements OnDestroy {
  @ViewChild('input', { static: false }) input!: ElementRef;

  @Input() public isDarkTheme: boolean = true;
  @Input() public user!: User.AsObject;
  @Input() public labelpolicy!: LabelPolicy.AsObject;
  public showOrgContext: boolean = false;

  public orgs$: Observable<Org.AsObject[]> = of([]);
  @Input() public org!: Org.AsObject;
  @Output() public changedActiveOrg: EventEmitter<Org.AsObject> = new EventEmitter();
  public filterControl: FormControl = new FormControl('');
  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public showAccount: boolean = false;
  public hideAdminWarn: boolean = true;
  private destroy$: Subject<void> = new Subject();

  public BreadcrumbType: any = BreadcrumbType;
  constructor(
    public authenticationService: AuthenticationService,
    private authService: GrpcAuthService,
    public mgmtService: ManagementService,
    public breadcrumbService: BreadcrumbService,
    public router: Router,
  ) {
    this.filterControl.valueChanges.pipe(debounceTime(300)).subscribe((value) => {
      this.loadOrgs(value.trim().toLowerCase());
    });

    this.hideAdminWarn = localStorage.getItem('hideAdministratorWarning') === 'true' ? true : false;
  }

  public ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public loadOrgs(filter?: string): void {
    let query;
    if (filter) {
      query = new OrgQuery();
      const orgNameQuery = new OrgNameQuery();
      orgNameQuery.setName(filter);
      orgNameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
      query.setNameQuery(orgNameQuery);
    }

    this.orgLoading$.next(true);
    this.orgs$ = from(this.authService.listMyProjectOrgs(10, 0, query ? [query] : undefined)).pipe(
      map((resp) => {
        return resp.resultList;
      }),
      catchError(() => of([])),
      finalize(() => {
        this.orgLoading$.next(false);
        this.focusFilter();
      }),
    );
  }

  public closeAccountCard(): void {
    if (this.showAccount) {
      this.showAccount = false;
    }
  }

  public focusFilter(): void {
    setTimeout(() => {
      this.input.nativeElement.focus();
    }, 0);
  }

  public setActiveOrg(org: Org.AsObject): void {
    this.org = org;
    this.authService.setActiveOrg(org);
    this.changedActiveOrg.emit(org);
  }

  public get isOnSystem(): boolean {
    return ['/system', '/views', '/failed-events'].includes(this.router.url);
  }
}
