import { Component, ElementRef, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { BehaviorSubject, Subject } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';

@Component({
  selector: 'cnsl-nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss'],
})
export class NavComponent implements OnInit, OnDestroy {
  @ViewChild('input', { static: false }) input!: ElementRef;

  @Input() public isDarkTheme: boolean = true;
  @Input() public user!: User.AsObject;
  @Input() public labelpolicy!: LabelPolicy.AsObject;

  @Input() public org!: Org.AsObject;
  public filterControl: FormControl = new FormControl('');
  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public showAccount: boolean = false;
  public hideAdminWarn: boolean = true;
  private destroy$: Subject<void> = new Subject();

  public BreadcrumbType: any = BreadcrumbType;

  constructor(
    public authenticationService: AuthenticationService,
    // private router: Router,
    public breadcrumbService: BreadcrumbService,
    public mgmtService: ManagementService,
  ) {
    this.hideAdminWarn = localStorage.getItem('hideAdministratorWarning') === 'true' ? true : false;
  }

  ngOnInit(): void {
    // this.router.events
    //   .pipe(
    //     filter((e) => e instanceof ActivationEnd),
    //     map((e) => (e instanceof ActivationEnd ? e.snapshot : {})),
    //     takeUntil(this.destroy$),
    //   )
    //   .subscribe((params) => {
    //     console.log(params);
    //     // Do whatever you want here!!!!
    //   });
  }

  public toggleAdminHide(): void {
    this.hideAdminWarn = !this.hideAdminWarn;
    localStorage.setItem('hideAdministratorWarning', this.hideAdminWarn.toString());
  }

  public ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
