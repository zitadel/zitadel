import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup } from '@angular/forms';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { Subject, Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { CnslLinks } from 'src/app/modules/links/links.component';
import { NameDialogComponent } from 'src/app/modules/name-dialog/name-dialog.component';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import {
    APIAuthMethodType,
    APIConfig,
    App,
    AppState,
    OIDCAppType,
    OIDCAuthMethodType,
    OIDCConfig,
    OIDCGrantType,
    OIDCResponseType,
    OIDCTokenType,
} from 'src/app/proto/generated/zitadel/app_pb';
import {
    GetOIDCInformationResponse,
    UpdateAPIAppConfigRequest,
    UpdateOIDCAppConfigRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import {
    BASIC_AUTH_METHOD,
    CODE_METHOD,
    CUSTOM_METHOD,
    getAuthMethodFromPartialConfig,
    getPartialConfigFromAuthMethod,
    IMPLICIT_METHOD,
    PK_JWT_METHOD,
    PKCE_METHOD,
    POST_METHOD,
} from '../authmethods';

@Component({
  selector: 'cnsl-app-detail',
  templateUrl: './app-detail.component.html',
  styleUrls: ['./app-detail.component.scss'],
})
export class AppDetailComponent implements OnInit, OnDestroy {
  public editState: boolean = false;
  public currentAuthMethod: string = CUSTOM_METHOD.key;
  public initialAuthMethod: string = CUSTOM_METHOD.key;
  public canWrite: boolean = false;
  public errorMessage: string = '';
  public removable: boolean = true;
  public addOnBlur: boolean = true;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public authMethods: RadioItemAuthType[] = [];
  private subscription?: Subscription;
  public projectId: string = '';
  public app!: App.AsObject;

  public environmentMap: { [key: string]: string } = {};

  public oidcResponseTypes: OIDCResponseType[] = [
    OIDCResponseType.OIDC_RESPONSE_TYPE_CODE,
    OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN,
    OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN,
  ];
  public oidcGrantTypes: OIDCGrantType[] = [
    OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE,
    OIDCGrantType.OIDC_GRANT_TYPE_IMPLICIT,
    OIDCGrantType.OIDC_GRANT_TYPE_REFRESH_TOKEN,
  ];
  public oidcAppTypes: OIDCAppType[] = [
    OIDCAppType.OIDC_APP_TYPE_WEB,
    OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
    OIDCAppType.OIDC_APP_TYPE_NATIVE,
  ];

  public oidcAuthMethodType: OIDCAuthMethodType[] = [
    OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC,
    OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST,
    OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
  ];

  public oidcTokenTypes: OIDCTokenType[] = [OIDCTokenType.OIDC_TOKEN_TYPE_BEARER, OIDCTokenType.OIDC_TOKEN_TYPE_JWT];

  public AppState: any = AppState;
  public oidcForm!: FormGroup;
  public apiForm!: FormGroup;

  public redirectUrisList: string[] = [];
  public postLogoutRedirectUrisList: string[] = [];
  public additionalOriginsList: string[] = [];

  public isZitadel: boolean = false;
  public docs!: GetOIDCInformationResponse.AsObject;

  public OIDCAppType: any = OIDCAppType;
  public OIDCAuthMethodType: any = OIDCAuthMethodType;
  public APIAuthMethodType: any = APIAuthMethodType;
  public OIDCTokenType: any = OIDCTokenType;
  public OIDCGrantType: any = OIDCGrantType;

  public ChangeType: any = ChangeType;

  public requestRedirectValuesSubject$: Subject<void> = new Subject();
  public copiedKey: any = '';
  public nextLinks: Array<CnslLinks> = [];
  public InfoSectionType: any = InfoSectionType;
  public copied: string = '';

  public settingsList: SidenavSetting[] = [{ id: 'configuration', i18nKey: 'APP.CONFIGURATION' }];
  public currentSetting: string | undefined = this.settingsList[0].id;

  constructor(
    public translate: TranslateService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private fb: FormBuilder,
    private _location: Location,
    private dialog: MatDialog,
    private mgmtService: ManagementService,
    private authService: GrpcAuthService,
    private router: Router,
    private snackbar: MatSnackBar,
    private breadcrumbService: BreadcrumbService,
    private http: HttpClient,
  ) {
    this.oidcForm = this.fb.group({
      devMode: [{ value: false, disabled: true }, []],
      clientId: [{ value: '', disabled: true }],
      responseTypesList: [{ value: [], disabled: true }],
      grantTypesList: [{ value: [], disabled: true }],
      appType: [{ value: '', disabled: true }],
      authMethodType: [{ value: '', disabled: true }],
      accessTokenType: [{ value: '', disabled: true }],
      accessTokenRoleAssertion: [{ value: false, disabled: true }],
      idTokenRoleAssertion: [{ value: false, disabled: true }],
      idTokenUserinfoAssertion: [{ value: false, disabled: true }],
      clockSkewSeconds: [{ value: 0, disabled: true }],
    });

    this.apiForm = this.fb.group({
      authMethodType: [{ value: '', disabled: true }],
    });

    this.http.get('./assets/environment.json').subscribe((env: any) => {
      this.environmentMap = {
        issuer: env.issuer,
        adminServiceUrl: env.api,
        mgmtServiceUrl: env.api,
        authServiceUrl: env.api,
      };
    });
  }

  public formatClockSkewLabel(seconds: number): string {
    return seconds + 's';
  }

  public additionalOriginsListChanged(origins: string[]): void {
    this.additionalOriginsList = origins;
  }

  public openNameDialog(): void {
    const dialogRef = this.dialog.open(NameDialogComponent, {
      data: {
        name: this.app.name,
        titleKey: 'APP.NAMEDIALOG.TITLE',
        descKey: 'APP.NAMEDIALOG.DESCRIPTION',
        labelKey: 'APP.NAMEDIALOG.NAME',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((name) => {
      if (name) {
        this.app.name = name;
        this.saveApp();
      }
    });
  }

  public ngOnInit(): void {
    const projectId = this.route.snapshot.paramMap.get('projectid');
    const appId = this.route.snapshot.paramMap.get('appid');

    if (projectId && appId) {
      this.projectId = projectId;
      this.getData(projectId, appId);
    }
  }

  public ngOnDestroy(): void {
    this.subscription?.unsubscribe();
  }

  private initLinks(): void {
    this.nextLinks = [
      {
        i18nTitle: 'APP.PAGES.NEXTSTEPS.0.TITLE',
        i18nDesc: 'APP.PAGES.NEXTSTEPS.0.DESC',
        routerLink: ['/projects', this.projectId, 'roles'],
        iconClasses: 'las la-user-tag',
      },
      {
        i18nTitle: 'APP.PAGES.NEXTSTEPS.1.TITLE',
        i18nDesc: 'APP.PAGES.NEXTSTEPS.1.DESC',
        routerLink: ['/users', 'create'],
        iconClasses: 'las la-user-plus',
      },
      {
        i18nTitle: 'APP.PAGES.NEXTSTEPS.2.TITLE',
        i18nDesc: 'APP.PAGES.NEXTSTEPS.2.DESC',
        href: 'https://docs.zitadel.ch',
        iconClasses: 'las la-people-carry',
      },
    ];
  }

  private async getData(projectId: string, appId: string): Promise<void> {
    this.initLinks();

    this.mgmtService.getIAM().then((iam) => {
      this.isZitadel = iam.iamProjectId === this.projectId;
    });
    this.authService
      .isAllowed(['project.app.write$', 'project.app.write:' + projectId])
      .pipe(take(1))
      .subscribe((allowed) => {
        this.canWrite = allowed;
        this.mgmtService
          .getAppByID(projectId, appId)
          .then((app) => {
            if (app.app) {
              this.app = app.app;

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
                }),
                new Breadcrumb({
                  type: BreadcrumbType.APP,
                  name: app.app.name,
                  param: { key: 'appid', value: appId },
                  routerLink: ['/projects', projectId, 'apps', appId],
                }),
              ];
              this.breadcrumbService.setBreadcrumb(breadcrumbs);

              if (this.app.oidcConfig) {
                this.getAuthMethodOptions('OIDC');

                this.settingsList = [
                  { id: 'configuration', i18nKey: 'APP.CONFIGURATION' },
                  { id: 'redirect-uris', i18nKey: 'APP.OIDC.REDIRECTSECTIONTITLE' },
                  { id: 'additional-origins', i18nKey: 'APP.ADDITIONALORIGINS' },
                  { id: 'urls', i18nKey: 'APP.URLS' },
                ];

                this.initialAuthMethod = this.authMethodFromPartialConfig({ oidc: this.app.oidcConfig });
                this.currentAuthMethod = this.initialAuthMethod;
                if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                  if (!this.authMethods.includes(CUSTOM_METHOD)) {
                    this.authMethods.push(CUSTOM_METHOD);
                  }
                } else {
                  this.authMethods = this.authMethods.filter((element) => element !== CUSTOM_METHOD);
                }
              } else if (this.app.apiConfig) {
                this.getAuthMethodOptions('API');

                this.initialAuthMethod = this.authMethodFromPartialConfig({ api: this.app.apiConfig });

                if (this.initialAuthMethod === 'BASIC') {
                  this.settingsList = [{ id: 'urls', i18nKey: 'APP.URLS' }];
                  this.currentSetting = 'urls';
                } else {
                  this.settingsList = [
                    { id: 'configuration', i18nKey: 'APP.CONFIGURATION' },
                    { id: 'urls', i18nKey: 'APP.URLS' },
                  ];
                }
                this.currentAuthMethod = this.initialAuthMethod;
                if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                  if (!this.authMethods.includes(CUSTOM_METHOD)) {
                    this.authMethods.push(CUSTOM_METHOD);
                  }
                } else {
                  this.authMethods = this.authMethods.filter((element) => element !== CUSTOM_METHOD);
                }
              }

              if (allowed) {
                this.oidcForm.enable();
                this.apiForm.enable();
              }

              if (this.app.oidcConfig?.redirectUrisList) {
                this.redirectUrisList = this.app.oidcConfig.redirectUrisList;
              }
              if (this.app.oidcConfig?.postLogoutRedirectUrisList) {
                this.postLogoutRedirectUrisList = this.app.oidcConfig.postLogoutRedirectUrisList;
              }
              if (this.app.oidcConfig?.additionalOriginsList) {
                this.additionalOriginsList = this.app.oidcConfig.additionalOriginsList;
              }

              if (this.app.oidcConfig?.clockSkew) {
                const inSecs = this.app.oidcConfig?.clockSkew.seconds + this.app.oidcConfig?.clockSkew.nanos / 100000;
                this.oidcForm.controls['clockSkewSeconds'].setValue(inSecs);
              }
              if (this.app.oidcConfig) {
                this.oidcForm.patchValue(this.app.oidcConfig);
              }
              if (this.app.apiConfig) {
                this.apiForm.patchValue(this.app.apiConfig);
              }

              this.oidcForm.valueChanges.subscribe((oidcConfig) => {
                this.initialAuthMethod = this.authMethodFromPartialConfig({ oidc: oidcConfig });
                if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                  if (!this.authMethods.includes(CUSTOM_METHOD)) {
                    this.authMethods.push(CUSTOM_METHOD);
                  }
                } else {
                  this.authMethods = this.authMethods.filter((element) => element !== CUSTOM_METHOD);
                }

                this.showSaveSnack();
              });

              this.apiForm.valueChanges.subscribe((apiConfig) => {
                this.initialAuthMethod = this.authMethodFromPartialConfig({ api: apiConfig });
                if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                  if (!this.authMethods.includes(CUSTOM_METHOD)) {
                    this.authMethods.push(CUSTOM_METHOD);
                  }
                } else {
                  this.authMethods = this.authMethods.filter((element) => element !== CUSTOM_METHOD);
                }

                this.showSaveSnack();
              });
            }
          })
          .catch((error) => {
            console.error(error);
            this.toast.showError(error);
            this.errorMessage = error.message;
          });
      });
    this.docs = await this.mgmtService.getOIDCInformation();
  }

  private async showSaveSnack(): Promise<void> {
    const message = await this.translate.get('APP.TOAST.CONFIGCHANGED').toPromise();
    const action = await this.translate.get('ACTIONS.SAVENOW').toPromise();

    const snackRef = this.snackbar.open(message, action, { duration: 5000, verticalPosition: 'top' });
    snackRef.onAction().subscribe(() => {
      if (this.app.oidcConfig) {
        this.saveOIDCApp();
      } else if (this.app.apiConfig) {
        this.saveAPIApp();
      }
    });
  }

  private getAuthMethodOptions(type: string): void {
    if (type === 'OIDC') {
      switch (this.app.oidcConfig?.appType) {
        case OIDCAppType.OIDC_APP_TYPE_NATIVE:
          this.authMethods = [PKCE_METHOD, CUSTOM_METHOD];
          break;
        case OIDCAppType.OIDC_APP_TYPE_WEB:
          this.authMethods = [PKCE_METHOD, CODE_METHOD, PK_JWT_METHOD, POST_METHOD];
          break;
        case OIDCAppType.OIDC_APP_TYPE_USER_AGENT:
          this.authMethods = [PKCE_METHOD, IMPLICIT_METHOD];
          break;
      }
    }
    if (type === 'API') {
      this.authMethods = [PK_JWT_METHOD, BASIC_AUTH_METHOD];
    }
  }

  public authMethodFromPartialConfig(config: { oidc?: OIDCConfig.AsObject; api?: APIConfig.AsObject }): string {
    const key = getAuthMethodFromPartialConfig(config);
    return key;
  }

  public setPartialConfigFromAuthMethod(authMethod: string): void {
    const partialConfig = getPartialConfigFromAuthMethod(authMethod);
    if (partialConfig && partialConfig.oidc && this.app.oidcConfig) {
      this.app.oidcConfig.responseTypesList = (partialConfig.oidc as Partial<OIDCConfig.AsObject>).responseTypesList ?? [];

      this.app.oidcConfig.grantTypesList = (partialConfig.oidc as Partial<OIDCConfig.AsObject>).grantTypesList ?? [];

      this.app.oidcConfig.authMethodType =
        (partialConfig.oidc as Partial<OIDCConfig.AsObject>).authMethodType ?? OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;

      this.oidcForm.patchValue(this.app.oidcConfig);
    } else if (partialConfig && partialConfig.api && this.app.apiConfig) {
      this.app.apiConfig.authMethodType =
        (partialConfig.api as Partial<APIConfig.AsObject>).authMethodType ?? APIAuthMethodType.API_AUTH_METHOD_TYPE_BASIC;

      this.apiForm.patchValue(this.app.apiConfig);
    }
  }

  public deleteApp(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'APP.PAGES.DIALOG.DELETE.TITLE',
        descriptionKey: 'APP.PAGES.DIALOG.DELETE.DESCRIPTION',
      },
      width: '400px',
    });
    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && this.projectId && this.app.id) {
        this.mgmtService
          .removeApp(this.projectId, this.app.id)
          .then(() => {
            this.toast.showInfo('APP.TOAST.DELETED', true);

            this.router.navigate(['/projects', this.projectId]);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public changeState(state: AppState): void {
    if (state === AppState.APP_STATE_ACTIVE) {
      this.mgmtService
        .reactivateApp(this.projectId, this.app.id)
        .then(() => {
          this.app.state = state;
          this.toast.showInfo('APP.TOAST.REACTIVATED', true);
        })
        .catch((error: any) => {
          this.toast.showError(error);
        });
    } else if (state === AppState.APP_STATE_INACTIVE) {
      this.mgmtService
        .deactivateApp(this.projectId, this.app.id)
        .then(() => {
          this.app.state = state;
          this.toast.showInfo('APP.TOAST.DEACTIVATED', true);
        })
        .catch((error: any) => {
          this.toast.showError(error);
        });
    }
  }

  public saveApp(): void {
    this.mgmtService
      .updateApp(this.projectId, this.app.id, this.app.name)
      .then(() => {
        this.toast.showInfo('APP.TOAST.UPDATED', true);
        this.editState = false;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public toggleRefreshToken(event: MatCheckboxChange): void {
    const c = this.grantTypesList?.value;

    if (event.checked) {
      if (!c.includes(OIDCGrantType.OIDC_GRANT_TYPE_REFRESH_TOKEN)) {
        this.grantTypesList?.setValue([OIDCGrantType.OIDC_GRANT_TYPE_REFRESH_TOKEN, ...c]);
      }
    } else {
      const index = (this.grantTypesList?.value as OIDCGrantType[]).findIndex(
        (gt) => gt === OIDCGrantType.OIDC_GRANT_TYPE_REFRESH_TOKEN,
      );
      if (index > -1) {
        const copy = Object.assign([], this.grantTypesList?.value);
        copy.splice(index, 1);
        this.grantTypesList?.setValue(copy);
      }
    }
  }

  public saveOIDCApp(): void {
    this.requestRedirectValuesSubject$.next();
    if (this.oidcForm.valid) {
      if (this.app.oidcConfig) {
        this.app.oidcConfig.responseTypesList = this.responseTypesList?.value;
        this.app.oidcConfig.grantTypesList = this.grantTypesList?.value;
        this.app.oidcConfig.appType = this.appType?.value;
        this.app.oidcConfig.authMethodType = this.authMethodType?.value;
        this.app.oidcConfig.redirectUrisList = this.redirectUrisList;
        this.app.oidcConfig.postLogoutRedirectUrisList = this.postLogoutRedirectUrisList;
        this.app.oidcConfig.additionalOriginsList = this.additionalOriginsList;
        this.app.oidcConfig.devMode = this.devMode?.value;
        this.app.oidcConfig.accessTokenType = this.accessTokenType?.value;
        this.app.oidcConfig.accessTokenRoleAssertion = this.accessTokenRoleAssertion?.value;
        this.app.oidcConfig.idTokenRoleAssertion = this.idTokenRoleAssertion?.value;
        this.app.oidcConfig.idTokenUserinfoAssertion = this.idTokenUserinfoAssertion?.value;

        const req = new UpdateOIDCAppConfigRequest();
        req.setProjectId(this.projectId);
        req.setAppId(this.app.id);
        req.setRedirectUrisList(this.app.oidcConfig.redirectUrisList);
        req.setResponseTypesList(this.app.oidcConfig.responseTypesList);
        req.setAdditionalOriginsList(this.app.oidcConfig.additionalOriginsList);
        req.setAuthMethodType(this.app.oidcConfig.authMethodType);
        req.setPostLogoutRedirectUrisList(this.app.oidcConfig.postLogoutRedirectUrisList);
        req.setGrantTypesList(this.app.oidcConfig.grantTypesList);
        req.setAppType(this.app.oidcConfig.appType);
        req.setDevMode(this.app.oidcConfig.devMode);
        req.setAccessTokenType(this.app.oidcConfig.accessTokenType);
        req.setAccessTokenRoleAssertion(this.app.oidcConfig.accessTokenRoleAssertion);
        req.setIdTokenRoleAssertion(this.app.oidcConfig.idTokenRoleAssertion);
        req.setIdTokenUserinfoAssertion(this.app.oidcConfig.idTokenUserinfoAssertion);
        if (this.clockSkewSeconds?.value) {
          const dur = new Duration();
          dur.setSeconds(Math.floor(this.clockSkewSeconds?.value));
          dur.setNanos(Math.floor(this.clockSkewSeconds?.value % 1) * 10000);
          req.setClockSkew(dur);
        }
        this.mgmtService
          .updateOIDCAppConfig(req)
          .then(() => {
            if (this.app.oidcConfig) {
              const config = { oidc: this.app.oidcConfig };
              this.currentAuthMethod = this.authMethodFromPartialConfig(config);
            }
            this.toast.showInfo('APP.TOAST.OIDCUPDATED', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    }
  }

  public saveAPIApp(): void {
    if (this.apiForm.valid && this.app.apiConfig) {
      this.app.apiConfig.authMethodType = this.apiAuthMethodType?.value;

      const req = new UpdateAPIAppConfigRequest();
      req.setProjectId(this.projectId);
      req.setAppId(this.app.id);
      req.setAuthMethodType(this.app.apiConfig.authMethodType);

      this.mgmtService
        .updateAPIAppConfig(req)
        .then(() => {
          if (this.app.apiConfig) {
            const config = { api: this.app.apiConfig };
            this.currentAuthMethod = this.authMethodFromPartialConfig(config);

            if (this.currentAuthMethod === 'BASIC') {
              this.settingsList = [{ id: 'urls', i18nKey: 'APP.URLS' }];
              this.currentSetting = 'urls';
            } else {
              this.settingsList = [
                { id: 'configuration', i18nKey: 'APP.CONFIGURATION' },
                { id: 'urls', i18nKey: 'APP.URLS' },
              ];
              this.currentSetting = 'configuration';
            }
          }
          this.toast.showInfo('APP.TOAST.APIUPDATED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public regenerateOIDCClientSecret(): void {
    this.mgmtService
      .regenerateOIDCClientSecret(this.app.id, this.projectId)
      .then((resp) => {
        this.toast.showInfo('APP.TOAST.CLIENTSECRETREGENERATED', true);
        this.dialog.open(AppSecretDialogComponent, {
          data: {
            // clientId: data.toObject() as ClientSecret.AsObject.clientId,
            clientSecret: resp.clientSecret,
          },
          width: '400px',
        });
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public regenerateAPIClientSecret(): void {
    this.mgmtService
      .regenerateAPIClientSecret(this.app.id, this.projectId)
      .then((resp) => {
        this.toast.showInfo('APP.TOAST.CLIENTSECRETREGENERATED', true);
        this.dialog.open(AppSecretDialogComponent, {
          data: {
            clientSecret: resp.clientSecret,
          },
          width: '400px',
        });
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public navigateBack(): void {
    this._location.back();
  }

  public get clientId(): AbstractControl | null {
    return this.oidcForm.get('clientId');
  }

  public get responseTypesList(): AbstractControl | null {
    return this.oidcForm.get('responseTypesList');
  }

  public get grantTypesList(): AbstractControl | null {
    return this.oidcForm.get('grantTypesList');
  }

  public get appType(): AbstractControl | null {
    return this.oidcForm.get('appType');
  }

  public get authMethodType(): AbstractControl | null {
    return this.oidcForm.get('authMethodType');
  }

  public get apiAuthMethodType(): AbstractControl | null {
    return this.apiForm.get('authMethodType');
  }

  public get devMode(): FormControl | null {
    return this.oidcForm.get('devMode') as FormControl;
  }

  public get accessTokenType(): AbstractControl | null {
    return this.oidcForm.get('accessTokenType');
  }

  public get idTokenRoleAssertion(): AbstractControl | null {
    return this.oidcForm.get('idTokenRoleAssertion');
  }

  public get accessTokenRoleAssertion(): AbstractControl | null {
    return this.oidcForm.get('accessTokenRoleAssertion');
  }

  public get idTokenUserinfoAssertion(): AbstractControl | null {
    return this.oidcForm.get('idTokenUserinfoAssertion');
  }

  public get clockSkewSeconds(): AbstractControl | null {
    return this.oidcForm.get('clockSkewSeconds');
  }
}
