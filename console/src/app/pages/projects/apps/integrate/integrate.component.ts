import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Buffer } from 'buffer';
import { Subject, Subscription } from 'rxjs';
import { debounceTime, map, takeUntil } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import {
  APIAuthMethodType,
  OIDCAppType,
  OIDCAuthMethodType,
  OIDCGrantType,
  OIDCResponseType,
} from 'src/app/proto/generated/zitadel/app_pb';
import {
  AddAPIAppRequest,
  AddAPIAppResponse,
  AddOIDCAppRequest,
  AddOIDCAppResponse,
  AddSAMLAppRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import {
  BASIC_AUTH_METHOD,
  CODE_METHOD,
  DEVICE_CODE_METHOD,
  getPartialConfigFromAuthMethod,
  IMPLICIT_METHOD,
  PKCE_METHOD,
  PK_JWT_METHOD,
  POST_METHOD,
} from '../authmethods';
import { API_TYPE, AppCreateType, NATIVE_TYPE, RadioItemAppType, SAML_TYPE, USER_AGENT_TYPE, WEB_TYPE } from '../authtypes';
import { EnvironmentService } from 'src/app/services/environment.service';

@Component({
  selector: 'cnsl-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
})
export class IntegrateAppComponent implements OnInit, OnDestroy {
  private destroyed$: Subject<void> = new Subject();
  public projectId: string = '';
  public loading: boolean = false;

  public currentCreateStep: number = 1;

  public oidcAppRequest: AddOIDCAppRequest = new AddOIDCAppRequest();
  public apiAppRequest: AddAPIAppRequest = new AddAPIAppRequest();

  public oidcResponseTypes: { type: OIDCResponseType; checked: boolean; disabled: boolean }[] = [
    { type: OIDCResponseType.OIDC_RESPONSE_TYPE_CODE, checked: false, disabled: false },
    { type: OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN, checked: false, disabled: false },
    { type: OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN, checked: false, disabled: false },
  ];

  public oidcAppTypes: OIDCAppType[] = [
    OIDCAppType.OIDC_APP_TYPE_WEB,
    OIDCAppType.OIDC_APP_TYPE_NATIVE,
    OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
  ];

  public appTypes: any = [WEB_TYPE, NATIVE_TYPE, USER_AGENT_TYPE, API_TYPE, SAML_TYPE];

  public authMethods: RadioItemAuthType[] = [PKCE_METHOD, CODE_METHOD, PK_JWT_METHOD, POST_METHOD];

  // set to oidc first
  public authMethodTypes: {
    type: OIDCAuthMethodType | APIAuthMethodType;
    checked: boolean;
    disabled: boolean;
    api?: boolean;
    oidc?: boolean;
  }[] = [
    { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC, checked: false, disabled: false, oidc: true },
    { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE, checked: false, disabled: false, oidc: true },
    { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST, checked: false, disabled: false, oidc: true },
  ];

  // stepper
  public firstFormGroup!: UntypedFormGroup;
  public redirectUrisList: string[] = [];
  public postLogoutRedirectUrisList: string[] = [];

  public oidcGrantTypes: {
    type: OIDCGrantType;
    checked: boolean;
    disabled: boolean;
  }[] = [
    { type: OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE, checked: true, disabled: false },
    { type: OIDCGrantType.OIDC_GRANT_TYPE_IMPLICIT, checked: false, disabled: true },
    { type: OIDCGrantType.OIDC_GRANT_TYPE_REFRESH_TOKEN, checked: false, disabled: true },
    { type: OIDCGrantType.OIDC_GRANT_TYPE_DEVICE_CODE, checked: false, disabled: true },
  ];

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
  public requestRedirectValuesSubject$: Subject<void> = new Subject();
  public samlCertificateURL$ = this.envSvc.env.pipe(map((env) => `${env.issuer}/saml/v2/certificate`));

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private toast: ToastService,
    private dialog: MatDialog,
    private mgmtService: ManagementService,
    private fb: UntypedFormBuilder,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
    private envSvc: EnvironmentService,
  ) {
    this.firstFormGroup = this.fb.group({
      name: ['', [requiredValidator]],
      appType: [WEB_TYPE, [requiredValidator]],
    });

    this.firstFormGroup.valueChanges.subscribe((value) => {
      if (this.firstFormGroup.valid) {
        this.oidcAppRequest.setName(this.name?.value);
      }
    });
  }

  public get redirectUris() {
    return this.oidcAppRequest.toObject().redirectUrisList;
  }

  public set redirectUris(value: string[]) {
    this.oidcAppRequest.setRedirectUrisList(value);
  }

  public get postLogoutUrisList() {
    return this.oidcAppRequest.toObject().postLogoutRedirectUrisList;
  }

  public set postLogoutUrisList(value: string[]) {
    this.oidcAppRequest.setPostLogoutRedirectUrisList(value);
  }

  public ngOnInit(): void {
    this.route.params.pipe(takeUntil(this.destroyed$)).subscribe((params) => this.getData(params));

    const projectId = this.route.snapshot.paramMap.get('projectid');
    if (projectId) {
      const breadcrumbs = [
        new Breadcrumb({
          type: BreadcrumbType.ORG,
          routerLink: ['/org'],
        }),
        new Breadcrumb({
          type: BreadcrumbType.PROJECT,
          name: '',
          param: { key: 'projectid', value: projectId },
          routerLink: ['/projects', projectId],
          isZitadel: false,
        }),
      ];
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    }
  }

  public ngOnDestroy(): void {
    this.destroyed$.next();
  }

  public changeStep(event: StepperSelectionEvent): void {
    this.currentCreateStep = event.selectedIndex + 1;

    if (event.selectedIndex >= 2) {
      this.requestRedirectValuesSubject$.next();
    }
  }

  private async getData({ projectid }: Params): Promise<void> {
    this.projectId = projectid;
    this.oidcAppRequest.setProjectId(projectid);
  }

  public close(): void {
    this._location.back();
  }

  public createApp(): void {
    this.requestRedirectValuesSubject$.next();

    this.loading = true;
    this.mgmtService
      .addOIDCApp(this.oidcAppRequest)
      .then((resp) => {
        this.loading = false;
        this.toast.showInfo('APP.TOAST.CREATED', true);
        if (resp.clientId || resp.clientSecret) {
          this.showSavedDialog(resp);
        } else {
          this.router.navigate(['projects', this.projectId, 'apps', resp.appId], { queryParams: { new: true } });
        }
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
  }

  public showSavedDialog(added: AddOIDCAppResponse.AsObject | AddAPIAppResponse.AsObject): void {
    let clientSecret = '';
    if (added.clientSecret) {
      clientSecret = added.clientSecret;
    }
    let clientId = '';
    if (added.clientId) {
      clientId = added.clientId;
    }
    const dialogRef = this.dialog.open(AppSecretDialogComponent, {
      data: {
        clientSecret: clientSecret,
        clientId: clientId,
      },
    });

    dialogRef.afterClosed().subscribe(() => {
      this.router.navigate(['projects', this.projectId, 'apps', added.appId], { queryParams: { new: true } });
    });
  }

  get name(): AbstractControl | null {
    return this.firstFormGroup.get('name');
  }

  get appType(): AbstractControl | null {
    return this.firstFormGroup.get('appType');
  }

  get isStepperOIDC(): boolean {
    return (this.appType?.value as RadioItemAppType).createType === AppCreateType.OIDC;
  }
}
