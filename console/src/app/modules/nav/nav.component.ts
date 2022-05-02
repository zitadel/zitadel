import { animate, keyframes, style, transition, trigger } from '@angular/animations';
import { BreakpointObserver } from '@angular/cdk/layout';
import { Component, ElementRef, Input, OnDestroy, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Router } from '@angular/router';
import { BehaviorSubject, map, Observable, Subject } from 'rxjs';
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
  animations: [
    trigger('navrow', [
      transition(':enter', [
        animate('.2s ease-in', keyframes([style({ opacity: 0, height: '0' }), style({ opacity: 1, height: '*' })])),
      ]),
      transition(':leave', [
        animate('.2s ease-out', keyframes([style({ opacity: 1, height: '*' }), style({ opacity: 0, height: '0' })])),
      ]),
    ]),
    trigger('navroworg', [
      transition(':enter', [
        animate(
          '.2s ease-in',
          keyframes([
            style({ opacity: 0, transform: 'translateY(-50%)' }),
            style({ opacity: 1, transform: 'translateY(0%)' }),
          ]),
        ),
      ]),
      transition(':leave', [
        animate(
          '.2s ease-out',
          keyframes([
            style({ opacity: 1, transform: 'translateY(0%)' }),
            style({ opacity: 0, transform: 'translateY(-50%)' }),
          ]),
        ),
      ]),
    ]),
    trigger('navrowproject', [
      transition(':enter', [
        animate(
          '.2s ease-in',
          keyframes([
            style({ opacity: 0, transform: 'translateY(+50%)' }),
            style({ opacity: 1, transform: 'translateY(0%)' }),
          ]),
        ),
      ]),
      transition(':leave', [
        animate(
          '.2s ease-out',
          keyframes([
            style({ opacity: 1, transform: 'translateY(0%)' }),
            style({ opacity: 0, transform: 'translateY(+50%)' }),
          ]),
        ),
      ]),
    ]),
  ],
})
export class NavComponent implements OnDestroy {
  @ViewChild('input', { static: false }) input!: ElementRef;

  @Input() public isDarkTheme: boolean = true;
  @Input() public user!: User.AsObject;
  @Input() public labelpolicy!: LabelPolicy.AsObject;
  public isHandset$: Observable<boolean> = this.breakpointObserver.observe('(max-width: 599px)').pipe(
    map((result) => {
      return result.matches;
    }),
  );

  @Input() public org!: Org.AsObject;
  public filterControl: FormControl = new FormControl('');
  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public showAccount: boolean = false;
  public hideAdminWarn: boolean = true;
  private destroy$: Subject<void> = new Subject();

  public BreadcrumbType: any = BreadcrumbType;

  constructor(
    public authenticationService: AuthenticationService,
    public breadcrumbService: BreadcrumbService,
    public mgmtService: ManagementService,
    private router: Router,
    private breakpointObserver: BreakpointObserver,
  ) {
    this.hideAdminWarn = localStorage.getItem('hideAdministratorWarning') === 'true' ? true : false;
  }

  public toggleAdminHide(): void {
    this.hideAdminWarn = !this.hideAdminWarn;
    localStorage.setItem('hideAdministratorWarning', this.hideAdminWarn.toString());
  }

  public ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public get isUserLinkActive(): boolean {
    const url = this.router.url;
    return url.substring(0, 6) === '/users';
  }
}
