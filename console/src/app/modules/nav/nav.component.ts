import { Component, ElementRef, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { BehaviorSubject, from, Observable, of, Subject } from 'rxjs';
import { catchError, debounceTime, finalize, map } from 'rxjs/operators';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Org, OrgNameQuery, OrgQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ThemeService } from 'src/app/services/theme.service';

@Component({
  selector: 'cnsl-nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss'],
})
export class NavComponent implements OnDestroy {
  @ViewChild('input', { static: false }) input!: ElementRef;

  @Input() public isDarkTheme: boolean = true;
  @Input() public user!: User.AsObject;
  @Input() public labelpolicy!: LabelPolicy.AsObject;

  public orgs$: Observable<Org.AsObject[]> = of([]);
  @Input() public org!: Org.AsObject;
  @Output() public changedActiveOrg: EventEmitter<Org.AsObject> = new EventEmitter();
  public filterControl: FormControl = new FormControl('');
  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public showAccount: boolean = false;
  public hideAdminWarn: boolean = true;
  private destroy$: Subject<void> = new Subject();

  constructor(
    public authenticationService: AuthenticationService,
    private authService: GrpcAuthService,
    public mgmtService: ManagementService,
    private themeService: ThemeService,
  ) {
    this.filterControl.valueChanges.pipe(debounceTime(300)).subscribe((value) => {
      this.loadOrgs(value.trim().toLowerCase());
    });

    // this.getProjectCount();

    // this.authService.activeOrgChanged.pipe(takeUntil(this.destroy$)).subscribe((org) => {
    //   this.org = org;
    //   console.log(org);
    //   this.getProjectCount();
    // });

    // this.authenticationService.authenticationChanged.pipe(takeUntil(this.destroy$)).subscribe((authenticated) => {
    //   if (authenticated) {
    //     this.authService.getActiveOrg().then((org) => {
    //       this.org = org;
    //     });
    //   }
    // });

    this.hideAdminWarn = localStorage.getItem('hideAdministratorWarning') === 'true' ? true : false;
  }

  public toggleAdminHide(): void {
    this.hideAdminWarn = !this.hideAdminWarn;
    localStorage.setItem('hideAdministratorWarning', this.hideAdminWarn.toString());
  }

  // private getProjectCount(): void {
  //   this.authService.isAllowed(['project.read']).subscribe((allowed) => {
  //     if (allowed) {
  //       this.mgmtService.listProjects(0, 0);
  //       this.mgmtService.listGrantedProjects(0, 0);
  //     }
  //   });
  // }

  public ngOnInit(): void {
    //Called after the constructor, initializing input properties, and the first call to ngOnChanges.
    //Add 'implements OnInit' to the class.
    console.log(this.org);
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
}
