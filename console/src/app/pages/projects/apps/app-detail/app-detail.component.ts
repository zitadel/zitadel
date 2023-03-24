import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatLegacyCheckboxChange as MatCheckboxChange } from '@angular/material/legacy-checkbox';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Buffer } from 'buffer';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { Subject, Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
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
  SAMLConfig,
} from 'src/app/proto/generated/zitadel/app_pb';
import {
  GetOIDCInformationResponse,
  UpdateAPIAppConfigRequest,
  UpdateOIDCAppConfigRequest,
  UpdateSAMLAppConfigRequest,
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
  PKCE_METHOD,
  PK_JWT_METHOD,
  POST_METHOD,
} from '../authmethods';
import { AuthMethodDialogComponent } from './auth-method-dialog/auth-method-dialog.component';

const MAX_ALLOWED_SIZE = 1 * 1024 * 1024;

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
  public app?: App.AsObject;

  public environmentMap: { [key: string]: string } = {};
  public wellKnownMap: { [key: string]: string } = {};

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
  public oidcForm!: UntypedFormGroup;
  public oidcTokenForm!: UntypedFormGroup;
  public apiForm!: UntypedFormGroup;
  public samlForm!: UntypedFormGroup;

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
  public InfoSectionType: any = InfoSectionType;
  public copied: string = '';

  public settingsList: SidenavSetting[] = [{ id: 'configuration', i18nKey: 'APP.CONFIGURATION' }];
  public currentSetting: string | undefined = this.settingsList[0].id;

  constructor(
    public translate: TranslateService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private fb: UntypedFormBuilder,
    private _location: Location,
    private dialog: MatDialog,
    private mgmtService: ManagementService,
    private authService: GrpcAuthService,
    private router: Router,
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
    });

    this.oidcTokenForm = this.fb.group({
      accessTokenType: [{ value: '', disabled: true }],
      accessTokenRoleAssertion: [{ value: false, disabled: true }],
      idTokenRoleAssertion: [{ value: false, disabled: true }],
      idTokenUserinfoAssertion: [{ value: false, disabled: true }],
      clockSkewSeconds: [{ value: 0, disabled: true }],
    });

    this.apiForm = this.fb.group({
      authMethodType: [{ value: '', disabled: true }],
    });

    this.samlForm = this.fb.group({
      metadataUrl: [{ value: '', disabled: true }],
      metadataXml: [{ value: '', disabled: true }],
    });

    this.http.get('./assets/environment.json').subscribe((env: any) => {
      this.environmentMap = {
        issuer: env.issuer,
        adminServiceUrl: `${env.api}/admin/v1`,
        mgmtServiceUrl: `${env.api}/management/v1`,
        authServiceUrl: `${env.api}/auth/v1`,
      };

      this.http.get(`${env.issuer}/.well-known/openid-configuration`).subscribe((wellKnown: any) => {
        this.wellKnownMap = {
          authorization_endpoint: wellKnown.authorization_endpoint,
          end_session_endpoint: wellKnown.end_session_endpoint,
          introspection_endpoint: wellKnown.introspection_endpoint,
          token_endpoint: wellKnown.token_endpoint,
          userinfo_endpoint: wellKnown.userinfo_endpoint,
        };
      });
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
        name: this.app?.name,
        titleKey: 'APP.NAMEDIALOG.TITLE',
        descKey: 'APP.NAMEDIALOG.DESCRIPTION',
        labelKey: 'APP.NAMEDIALOG.NAME',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((name) => {
      if (name) {
        this.app!.name = name;
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

  private async getData(projectId: string, appId: string): Promise<void> {
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
                  { id: 'token', i18nKey: 'APP.TOKEN' },
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
              } else if (this.app.samlConfig) {
                this.settingsList = [{ id: 'configuration', i18nKey: 'APP.CONFIGURATION' }];
              }

              if (allowed) {
                this.oidcForm.enable();
                this.oidcTokenForm.enable();
                this.apiForm.enable();
                this.samlForm.enable();
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
                this.oidcTokenForm.controls['clockSkewSeconds'].setValue(inSecs);
              }
              if (this.app.oidcConfig) {
                this.oidcForm.patchValue(this.app.oidcConfig);
                this.oidcTokenForm.patchValue(this.app.oidcConfig);
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
              });
            }
          })
          .catch((error) => {
            this.toast.showError(error);
            this.errorMessage = error.message;
          });
      });
    this.docs = await this.mgmtService.getOIDCInformation();
  }

  private getAuthMethodOptions(type: string): void {
    if (type === 'OIDC') {
      switch (this.app?.oidcConfig?.appType) {
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

  public onDropXML(filelist: FileList): void {
    const file = filelist.item(0);
    if (file) {
      if (file.size > MAX_ALLOWED_SIZE) {
        this.toast.showInfo('POLICY.PRIVATELABELING.MAXSIZEEXCEEDED', true);
      } else {
        this.metadataUrl?.setValue('');
        const reader = new FileReader();
        reader.onload = ((aXML) => {
          return (e) => {
            const xmlBase64 = e.target?.result;
            if (xmlBase64 && typeof xmlBase64 === 'string' && this.app?.samlConfig) {
              const samlConfig = new SAMLConfig();
              const cropped = xmlBase64.replace('data:text/xml;base64,', '');
              samlConfig.setMetadataXml(cropped);
              this.app.samlConfig.metadataXml = cropped;
            }
          };
        })(file);
        reader.readAsDataURL(file);
      }
    }
  }

  public authMethodFromPartialConfig(config: { oidc?: OIDCConfig.AsObject; api?: APIConfig.AsObject }): string {
    const key = getAuthMethodFromPartialConfig(config);
    return key;
  }

  public setPartialConfigFromAuthMethod(authMethod: string): void {
    const partialConfig = getPartialConfigFromAuthMethod(authMethod);
    if (partialConfig && partialConfig.oidc && this.app?.oidcConfig) {
      this.app!.oidcConfig.responseTypesList = (partialConfig.oidc as Partial<OIDCConfig.AsObject>).responseTypesList ?? [];

      this.app!.oidcConfig.grantTypesList = (partialConfig.oidc as Partial<OIDCConfig.AsObject>).grantTypesList ?? [];

      this.app!.oidcConfig.authMethodType =
        (partialConfig.oidc as Partial<OIDCConfig.AsObject>).authMethodType ?? OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;

      this.oidcForm.patchValue(this.app.oidcConfig);
      this.oidcTokenForm.patchValue(this.app.oidcConfig);
    } else if (partialConfig && partialConfig.api && this.app?.apiConfig) {
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
      if (resp && this.projectId && this.app?.id) {
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
    if (state === AppState.APP_STATE_ACTIVE && this.app) {
      this.mgmtService
        .reactivateApp(this.projectId, this.app.id)
        .then(() => {
          this.app!.state = state;
          this.toast.showInfo('APP.TOAST.REACTIVATED', true);
        })
        .catch((error: any) => {
          this.toast.showError(error);
        });
    } else if (state === AppState.APP_STATE_INACTIVE && this.app) {
      this.mgmtService
        .deactivateApp(this.projectId, this.app.id)
        .then(() => {
          this.app!.state = state;
          this.toast.showInfo('APP.TOAST.DEACTIVATED', true);
        })
        .catch((error: any) => {
          this.toast.showError(error);
        });
    }
  }

  public saveApp(): void {
    if (this.app) {
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
      if (this.app?.oidcConfig) {
        //   configuration
        this.app.oidcConfig.responseTypesList = this.responseTypesList?.value;
        this.app.oidcConfig.grantTypesList = this.grantTypesList?.value;
        this.app.oidcConfig.appType = this.appType?.value;
        this.app.oidcConfig.authMethodType = this.authMethodType?.value;

        // token
        this.app.oidcConfig.accessTokenType = this.accessTokenType?.value;
        this.app.oidcConfig.accessTokenRoleAssertion = this.accessTokenRoleAssertion?.value;
        this.app.oidcConfig.idTokenRoleAssertion = this.idTokenRoleAssertion?.value;
        this.app.oidcConfig.idTokenUserinfoAssertion = this.idTokenUserinfoAssertion?.value;

        // redirects
        this.app.oidcConfig.redirectUrisList = this.redirectUrisList;
        this.app.oidcConfig.postLogoutRedirectUrisList = this.postLogoutRedirectUrisList;
        this.app.oidcConfig.additionalOriginsList = this.additionalOriginsList;
        this.app.oidcConfig.devMode = this.devMode?.value;

        const req = new UpdateOIDCAppConfigRequest();
        req.setProjectId(this.projectId);
        req.setAppId(this.app.id);

        // configuration
        req.setResponseTypesList(this.app.oidcConfig.responseTypesList);
        req.setAuthMethodType(this.app.oidcConfig.authMethodType);
        req.setGrantTypesList(this.app.oidcConfig.grantTypesList);
        req.setAppType(this.app.oidcConfig.appType);

        // token
        req.setAccessTokenType(this.app.oidcConfig.accessTokenType);
        req.setAccessTokenRoleAssertion(this.app.oidcConfig.accessTokenRoleAssertion);
        req.setIdTokenRoleAssertion(this.app.oidcConfig.idTokenRoleAssertion);
        req.setIdTokenUserinfoAssertion(this.app.oidcConfig.idTokenUserinfoAssertion);

        // redirects
        req.setRedirectUrisList(this.app.oidcConfig.redirectUrisList);
        req.setAdditionalOriginsList(this.app.oidcConfig.additionalOriginsList);
        req.setPostLogoutRedirectUrisList(this.app.oidcConfig.postLogoutRedirectUrisList);
        req.setDevMode(this.app.oidcConfig.devMode);

        if (this.clockSkewSeconds?.value) {
          const dur = new Duration();
          dur.setSeconds(Math.floor(this.clockSkewSeconds?.value));
          dur.setNanos(Math.floor(this.clockSkewSeconds?.value % 1) * 10000);
          req.setClockSkew(dur);
        }

        this.mgmtService
          .updateOIDCAppConfig(req)
          .then(() => {
            if (this.app?.oidcConfig) {
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
    if (this.apiForm.valid && this.app?.apiConfig) {
      this.app.apiConfig.authMethodType = this.apiAuthMethodType?.value;

      const req = new UpdateAPIAppConfigRequest();
      req.setProjectId(this.projectId);
      req.setAppId(this.app.id);
      req.setAuthMethodType(this.app.apiConfig.authMethodType);

      this.mgmtService
        .updateAPIAppConfig(req)
        .then(() => {
          if (this.app?.apiConfig) {
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

  public saveSAMLApp(): void {
    if (this.samlForm.valid && this.app?.samlConfig) {
      const req = new UpdateSAMLAppConfigRequest();
      req.setProjectId(this.projectId);
      req.setAppId(this.app.id);

      if (this.app.samlConfig) {
        req.setMetadataUrl(this.app.samlConfig?.metadataUrl);
        req.setMetadataXml(this.app.samlConfig?.metadataXml);
      }

      this.mgmtService
        .updateSAMLAppConfig(req)
        .then(() => {
          this.toast.showInfo('APP.TOAST.APIUPDATED', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public regenerateOIDCClientSecret(): void {
    if (this.app) {
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
  }

  public changeAuthMethod(): void {
    const ref = this.dialog.open(AuthMethodDialogComponent, {
      data: {
        radioItemAuthType: this.currentRadioItemAuthType,
        initialAuthMethod: this.initialAuthMethod,
        currentAuthMethod: this.currentAuthMethod,
        authMethods: this.authMethods,
        isOIDC: this.app?.oidcConfig !== undefined,
      },
    });

    ref.afterClosed().subscribe((authMethod: string) => {
      if (authMethod) {
        this.setPartialConfigFromAuthMethod(authMethod);
      }
    });
  }

  public regenerateAPIClientSecret(): void {
    if (this.app) {
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
  }

  public get currentRadioItemAuthType(): RadioItemAuthType | undefined {
    return this.authMethods.find((i) => i.key === this.initialAuthMethod);
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

  public get devMode(): UntypedFormControl | null {
    return this.oidcForm.get('devMode') as UntypedFormControl;
  }

  public get accessTokenType(): AbstractControl | null {
    return this.oidcTokenForm.get('accessTokenType');
  }

  public get idTokenRoleAssertion(): AbstractControl | null {
    return this.oidcTokenForm.get('idTokenRoleAssertion');
  }

  public get accessTokenRoleAssertion(): AbstractControl | null {
    return this.oidcTokenForm.get('accessTokenRoleAssertion');
  }

  public get idTokenUserinfoAssertion(): AbstractControl | null {
    return this.oidcTokenForm.get('idTokenUserinfoAssertion');
  }

  public get clockSkewSeconds(): AbstractControl | null {
    return this.oidcTokenForm.get('clockSkewSeconds');
  }

  public get metadataUrl(): AbstractControl | null {
    return this.samlForm.get('metadataUrl');
  }

  get decodedBase64(): string {
    if (
      this.app &&
      this.app.samlConfig &&
      this.app.samlConfig.metadataXml &&
      typeof this.app.samlConfig.metadataXml === 'string'
    ) {
      return Buffer.from(this.app?.samlConfig.metadataXml, 'base64').toString('ascii');
    } else {
      return '';
    }
  }

  set decodedBase64(xmlString: string) {
    if (this.app && this.app.samlConfig && this.app.samlConfig.metadataXml) {
      const base64 = Buffer.from(xmlString, 'ascii').toString('base64');

      if (this.app.samlConfig) {
        this.app.samlConfig.metadataXml = base64;
      }
    }
  }
}
