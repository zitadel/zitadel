import { Component, ElementRef, Input, OnDestroy, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { BehaviorSubject, Subject } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
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

  @Input() public org!: Org.AsObject;
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
    this.hideAdminWarn = localStorage.getItem('hideAdministratorWarning') === 'true' ? true : false;
  }

  public toggleAdminHide(): void {
    this.hideAdminWarn = !this.hideAdminWarn;
    localStorage.setItem('hideAdministratorWarning', this.hideAdminWarn.toString());
  }

  public ngOnInit(): void {
    //Called after the constructor, initializing input properties, and the first call to ngOnChanges.
    //Add 'implements OnInit' to the class.
    console.log(this.org);
  }

  public ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
