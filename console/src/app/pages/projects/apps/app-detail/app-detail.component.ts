import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { ChangeDetectorRef, Component, OnDestroy, OnInit, signal } from '@angular/core';
import { AbstractControl, FormControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../../../docs/frameworks.json';
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';
import { Buffer } from 'buffer';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { mergeMap, Subject, Subscription } from 'rxjs';
import { map, startWith, switchMap, take } from 'rxjs/operators';
import { EnvVar } from 'src/app/components/env-vars-block/env-vars-block.component';
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
  LoginV1,
  LoginV2,
  LoginVersion,
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
import { AdminService } from 'src/app/services/admin.service';

import { Environment, EnvironmentService } from 'src/app/services/environment.service';
import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import { ListMilestonesRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { ListQuery } from 'src/app/proto/generated/zitadel/object_pb';
import { MilestoneQuery, IsReachedQuery, MilestoneType } from 'src/app/proto/generated/zitadel/milestone/v1/milestone_pb';
import {
  BASIC_AUTH_METHOD,
  CODE_METHOD,
  CUSTOM_METHOD,
  DEVICE_CODE_METHOD,
  getAuthMethodFromPartialConfig,
  getPartialConfigFromAuthMethod,
  IMPLICIT_METHOD,
  PK_JWT_METHOD,
  PKCE_METHOD,
  POST_METHOD,
} from '../authmethods';
import { AuthMethodDialogComponent } from './auth-method-dialog/auth-method-dialog.component';

const MAX_ALLOWED_SIZE = 1 * 1024 * 1024;

@Component({
  selector: 'cnsl-app-detail',
  templateUrl: './app-detail.component.html',
  styleUrls: ['./app-detail.component.scss'],
  standalone: false,
})
export class AppDetailComponent implements OnInit, OnDestroy {
  public editState: boolean = false;
  public currentAuthMethod: string = CUSTOM_METHOD.key;
  public initialAuthMethod: string = CUSTOM_METHOD.key;
  public canWrite: boolean = false;
  public errorMessage: string = '';
  public removable: boolean = true;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public authMethods: RadioItemAuthType[] = [];
  private subscription?: Subscription;
  public projectId: string = '';
  public appId: string = '';
  public app?: App.AsObject;

  public apiURLs$ = this.envSvc.env.pipe(
    mergeMap((env) =>
      this.wellknownURLs$.pipe(
        map((wellknown) => {
          return [
            ['Issuer', env.issuer],
            ['Admin Service URL', `${env.api}/admin/v1`],
            ['Management Service URL', `${env.api}/management/v1`],
            ['Auth Service URL', `${env.api}/auth/v1`],
            ...wellknown.filter(
              ([k, v]) => k === 'Revocation Endpoint' || k === 'JKWS URI' || k === 'Introspection Endpoint',
            ),
          ];
        }),
      ),
    ),
  );

  public issuer$ = this.apiURLs$.pipe(map((urls) => urls.find(([k, v]) => k === 'Issuer')?.[1]));

  public samlURLs$ = this.envSvc.env.pipe(
    map((env) => {
      return {
        samlCertificateURL: `${env.issuer}/saml/v2/certificate`,
        samlSSO: `${env.issuer}/saml/v2/SSO`,
        samlSLO: `${env.issuer}/saml/v2/SLO`,
      };
    }),
  );

  public wellknownURLs$ = this.envSvc.wellknown.pipe(
    map((wellknown) => {
      return [
        ['Authorization Endpoint', wellknown.authorization_endpoint],
        ['Device Authorization Endpoint', wellknown.device_authorization_endpoint],
        ['End Session Endpoint', wellknown.end_session_endpoint],
        ['Introspection Endpoint', wellknown.introspection_endpoint],
        ['JKWS URI', wellknown.jwks_uri],
        ['Revocation Endpoint', wellknown.revocation_endpoint],
        ['Token Endpoint', wellknown.token_endpoint],
        ['Userinfo Endpoint', wellknown.userinfo_endpoint],
      ];
    }),
  );

  public oidcResponseTypes: OIDCResponseType[] = [
    OIDCResponseType.OIDC_RESPONSE_TYPE_CODE,
    OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN,
    OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN,
  ];
  public oidcGrantTypes: OIDCGrantType[] = [
    OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE,
    OIDCGrantType.OIDC_GRANT_TYPE_IMPLICIT,
    OIDCGrantType.OIDC_GRANT_TYPE_DEVICE_CODE,
    OIDCGrantType.OIDC_GRANT_TYPE_REFRESH_TOKEN,
    OIDCGrantType.OIDC_GRANT_TYPE_TOKEN_EXCHANGE,
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

  public OIDCAppType = OIDCAppType;
  public OIDCAuthMethodType = OIDCAuthMethodType;
  public APIAuthMethodType = APIAuthMethodType;
  public OIDCTokenType = OIDCTokenType;
  public OIDCGrantType = OIDCGrantType;

  public ChangeType = ChangeType;

  public requestRedirectValuesSubject$: Subject<void> = new Subject();
  public InfoSectionType = InfoSectionType;
  public copied: string = '';

  public settingsList: SidenavSetting[] = [
    { id: 'quickstart', i18nKey: 'APP.QUICKSTART' },
    { id: 'configuration', i18nKey: 'APP.CONFIGURATION' },
  ];
  public currentSetting = this.settingsList[0];
  public framework: string | null = null;
  public frameworks = frameworkDefinition
    .filter((f) => OIDC_CONFIGURATIONS[f.id as unknown as keyof typeof OIDC_CONFIGURATIONS])
    .map((f) => {
      return {
        ...f,
        fragment: f.id || '',
        imgSrcDark: `assets${f.imgSrcDark}`,
        imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
      };
    });

  public isNew = signal<boolean>(false);
  public selectedScenario: 'new' | 'existing' = 'new';
  public isAuthenticated = signal<boolean | null>(null); // null = loading

  private appDataSubject = new Subject<void>();
  private cachedIssuer: string = '';
  private cachedDomain: string = '';

  public environmentVariables$ = this.appDataSubject.pipe(
    startWith(null),
    switchMap(() => this.envSvc.env),
    map((env: Environment) => {
      if (!this.app) return '';

      const issuer = env?.issuer || 'https://your-domain.zitadel.cloud';
      const clientId = this.clientId?.value || 'your-client-id';
      let envVars = `ZITADEL_DOMAIN=${issuer}\n`;
      envVars += `ZITADEL_CLIENT_ID=${clientId}\n`;
      envVars += `ZITADEL_PROJECT_ID=${this.projectId}\n`;

      if (this.app.oidcConfig) {
        if (
          this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC ||
          this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST
        ) {
          envVars += `ZITADEL_CLIENT_SECRET=your-client-secret\n`;
        }

        const redirectUris = this.app.oidcConfig.redirectUrisList || [];
        if (redirectUris.length > 0) {
          envVars += `ZITADEL_CALLBACK_URL=${redirectUris[0]}\n`;
        }

        const postLogoutUris = this.app.oidcConfig.postLogoutRedirectUrisList || [];
        if (postLogoutUris.length > 0) {
          envVars += `ZITADEL_POST_LOGOUT_URL=${postLogoutUris[0]}\n`;
        }
      }

      return envVars;
    }),
  );

  public envVarsArray: EnvVar[] = [];
  public isExampleEnvMode: boolean = false;

  private getEnvVarPrefix(): string {
    if (this.selectedScenario !== 'new') {
      return '';
    }

    const framework = this.framework?.toLowerCase();

    const prefixMap: { [key: string]: string } = {
      angular: 'NG_APP_',
      vue: 'VITE_',
      qwik: 'VITE_',
      react: 'VITE_',
    };

    return prefixMap[framework || ''] || '';
  }

  private async updateEnvVars(): Promise<void> {
    try {
      const env = await this.envSvc.env.pipe(take(1)).toPromise();

      if (!this.app || !env) {
        this.envVarsArray = [];
        return;
      }

      const issuer = env.issuer || 'https://your-domain.zitadel.cloud';
      const domain = env.issuer?.replace('https://', '').replace('http://', '') || 'your-domain.zitadel.cloud';
      const clientId = this.app.oidcConfig?.clientId || 'your-client-id';

      // Cache the values for use in replacePlaceholders
      this.cachedIssuer = issuer;
      this.cachedDomain = domain;

      // Get framework configuration
      const frameworkInfo = this.getFrameworkInfo();
      const hasExampleEnv = frameworkInfo?.hasExampleEnv;

      // For frameworks with .env.example files (hasExampleEnv=true),
      // only show minimal ZITADEL-specific values that need to be updated
      if (hasExampleEnv && this.selectedScenario === 'new') {
        this.isExampleEnvMode = true;
        const prefix = this.getEnvVarPrefix();
        const envVars: EnvVar[] = [
          {
            key: `${prefix}ZITADEL_DOMAIN`,
            value: issuer,
          },
          {
            key: `${prefix}ZITADEL_CLIENT_ID`,
            value: clientId,
          },
        ];

        // Add client secret if applicable
        if (
          this.app.oidcConfig &&
          (this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC ||
            this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST)
        ) {
          envVars.push({
            key: `${prefix}ZITADEL_CLIENT_SECRET`,
            value: 'your-client-secret',
          });
        }

        // Add callback URL if configured
        if (this.app.oidcConfig) {
          const redirectUris = this.app.oidcConfig.redirectUrisList || [];
          if (redirectUris.length > 0) {
            envVars.push({
              key: `${prefix}ZITADEL_CALLBACK_URL`,
              value: redirectUris[0],
            });
          }

          // Add post-logout redirect URL if configured
          const postLogoutUris = this.app.oidcConfig.postLogoutRedirectUrisList || [];
          if (postLogoutUris.length > 0) {
            envVars.push({
              key: `${prefix}ZITADEL_POST_LOGOUT_URL`,
              value: postLogoutUris[0],
            });
          }
        }

        this.envVarsArray = envVars;
        this.cdr.detectChanges();
        return;
      }

      // Fallback to framework-specific envSetup configuration (for non-example frameworks)
      const envSetup = frameworkInfo?.envSetup;

      // Only use custom envSetup if:
      // 1. Framework has envSetup defined AND
      // 2. User selected "new" scenario (wants to use ZITADEL's example app)
      // For "existing" scenario, always use standard ZITADEL variable names
      if (envSetup && envSetup.variables && this.selectedScenario === 'new') {
        this.isExampleEnvMode = false;
        const envVars: EnvVar[] = envSetup.variables.map((variable: any) => {
          let value = variable.placeholder;

          // Replace template placeholders with actual values
          value = value
            .replace('{{ZITADEL_DOMAIN}}', domain)
            .replace('{{ZITADEL_ISSUER}}', issuer)
            .replace('{{CLIENT_ID}}', clientId)
            .replace('{{PROJECT_ID}}', this.projectId)
            .replace('{{PROJECT_NAME}}', this.app?.name || 'project');

          // Handle special cases
          if (variable.name.includes('SECRET') && value === variable.placeholder) {
            value = 'your-client-secret';
          }

          // Handle redirect URIs
          if (variable.name.includes('REDIRECT') || variable.name.includes('CALLBACK')) {
            const redirectUris = this.app?.oidcConfig?.redirectUrisList || [];
            if (redirectUris.length > 0) {
              value = redirectUris[0];
            }
          }

          // Handle post logout URIs
          if (variable.name.includes('POST_LOGOUT') || variable.name.includes('LOGOUT')) {
            const postLogoutUris = this.app?.oidcConfig?.postLogoutRedirectUrisList || [];
            if (postLogoutUris.length > 0) {
              value = postLogoutUris[0];
            }
          }

          return {
            key: variable.name,
            value: value,
          };
        });

        this.envVarsArray = envVars;
        this.cdr.detectChanges();
        return;
      }

      // Fallback to default behavior if no envSetup or hasExampleEnv
      this.isExampleEnvMode = false;
      const prefix = this.getEnvVarPrefix();

      const envVars: EnvVar[] = [
        {
          key: `${prefix}ZITADEL_DOMAIN`,
          value: issuer,
        },
        {
          key: `${prefix}ZITADEL_CLIENT_ID`,
          value: clientId,
        },
        {
          key: `${prefix}ZITADEL_PROJECT_ID`,
          value: this.projectId,
        },
      ];

      if (this.app.oidcConfig) {
        if (
          this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC ||
          this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST
        ) {
          envVars.push({
            key: `${prefix}ZITADEL_CLIENT_SECRET`,
            value: 'your-client-secret',
          });
        }

        const redirectUris = this.app.oidcConfig.redirectUrisList || [];
        if (redirectUris.length > 0) {
          envVars.push({
            key: `${prefix}ZITADEL_CALLBACK_URL`,
            value: redirectUris[0],
          });
        }

        const postLogoutUris = this.app.oidcConfig.postLogoutRedirectUrisList || [];
        if (postLogoutUris.length > 0) {
          envVars.push({
            key: `${prefix}ZITADEL_POST_LOGOUT_URL`,
            value: postLogoutUris[0],
          });
        }
      }

      this.envVarsArray = envVars;
      this.cdr.detectChanges();
    } catch (error) {
      console.error('Error updating env vars:', error);
      this.envVarsArray = [];
      this.cdr.detectChanges();
    }
  }

  constructor(
    private envSvc: EnvironmentService,
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
    private adminService: AdminService,
    private cdr: ChangeDetectorRef,
  ) {
    this.oidcForm = this.fb.group({
      devMode: [{ value: false, disabled: true }],
      skipNativeAppSuccessPage: [{ value: false, disabled: true }],
      clientId: [{ value: '', disabled: true }],
      responseTypesList: [{ value: [], disabled: true }],
      grantTypesList: [{ value: [], disabled: true }],
      appType: [{ value: '', disabled: true }],
      authMethodType: [{ value: '', disabled: true }],
      loginV2: [{ value: false, disabled: true }],
      loginV2BaseURL: [{ value: '', disabled: true }],
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
      entityId: ['', []],
      acsURL: ['', []],
      metadataXml: [{ value: '', disabled: true }],
      loginV2: [{ value: false, disabled: true }],
      loginV2BaseURL: [{ value: '', disabled: true }],
    });

    this.samlForm.valueChanges.subscribe(() => {
      if (!this.app) {
        this.app = new App().toObject();
        this.updateEnvVars();
      }

      let minimalMetadata =
        this.entityId?.value && this.acsURL?.value
          ? `<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="${this.entityId?.value}">
    <md:SPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol urn:oasis:names:tc:SAML:1.1:protocol">
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="${this.acsURL?.value}" index="0"/>
    </md:SPSSODescriptor>
</md:EntityDescriptor>`
          : '';

      if (!minimalMetadata && !this.metadataUrl?.value) {
        return;
      }

      if (!this.app.samlConfig) {
        this.app.samlConfig = new SAMLConfig().toObject();
      }

      if (minimalMetadata) {
        const base64 = Buffer.from(minimalMetadata, 'utf-8').toString('base64');
        this.app.samlConfig.metadataXml = base64;
        this.app.samlConfig.metadataUrl = '';
      }

      if (this.metadataUrl?.value) {
        this.app.samlConfig.metadataXml = '';
        this.app.samlConfig.metadataUrl = this.metadataUrl?.value;
      }
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
    const isNew = this.route.snapshot.queryParamMap.get('new');
    const framework = this.route.snapshot.queryParamMap.get('framework');

    this.isNew.set(isNew === 'true');
    this.framework = framework;

    // Reorder sidebar based on new query param
    // If new=true: Quickstart first, Config second (users want to get started)
    // If new!=true: Config first, Quickstart last (users want to manage settings)
    if (isNew === 'true') {
      this.settingsList = [
        { id: 'quickstart', i18nKey: 'APP.QUICKSTART' },
        { id: 'configuration', i18nKey: 'APP.CONFIGURATION' },
      ];
      this.currentSetting = this.settingsList[0]; // Start with quickstart
    } else {
      this.settingsList = [
        { id: 'configuration', i18nKey: 'APP.CONFIGURATION' },
        { id: 'quickstart', i18nKey: 'APP.QUICKSTART' },
      ];
      this.currentSetting = this.settingsList[0]; // Start with configuration
    }

    if (projectId && appId) {
      this.projectId = projectId;
      this.appId = appId;
      this.getData(projectId, appId).then();
      this.checkAuthenticationStatus();
    } else if (projectId && isNew === 'true') {
      // Handle new app creation
      this.projectId = projectId;
      this.checkAuthenticationStatus();
      // Enable forms for new apps so users can see framework defaults
      this.authService
        .isAllowed(['project.app.write$', 'project.app.write:' + projectId])
        .pipe(take(1))
        .subscribe((allowed) => {
          if (allowed) {
            this.oidcForm.enable();
            this.oidcForm.controls['clientId'].disable(); // Keep clientId disabled as it's generated
            this.oidcTokenForm.enable();
            this.apiForm.enable();
            this.samlForm.enable();
            // Initialize form with framework defaults if framework is selected
            this.initializeFormWithFrameworkDefaults();
          }
        });
    }
  }

  public ngOnDestroy(): void {
    this.subscription?.unsubscribe();
    this.appDataSubject.complete();
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
              this.updateEnvVars();

              // TODO: duplicates should be handled in the API
              if (this.app.oidcConfig?.complianceProblemsList && this.app.oidcConfig?.complianceProblemsList.length) {
                this.app.oidcConfig.complianceProblemsList = this.app.oidcConfig?.complianceProblemsList.filter(
                  (element, index) => {
                    return this.app?.oidcConfig?.complianceProblemsList.findIndex((e) => e.key === element.key) === index;
                  },
                );
              }

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

                if (
                  this.app.oidcConfig.grantTypesList.length === 1 &&
                  this.app.oidcConfig.grantTypesList[0] === OIDCGrantType.OIDC_GRANT_TYPE_DEVICE_CODE
                ) {
                  this.settingsList = this.buildSettingsList([
                    { id: 'token', i18nKey: 'APP.TOKEN' },
                    { id: 'urls', i18nKey: 'APP.URLS' },
                  ]);
                } else {
                  this.settingsList = this.buildSettingsList([
                    { id: 'token', i18nKey: 'APP.TOKEN' },
                    { id: 'redirect-uris', i18nKey: 'APP.OIDC.REDIRECTSECTIONTITLE' },
                    { id: 'additional-origins', i18nKey: 'APP.ADDITIONALORIGINS' },
                    { id: 'urls', i18nKey: 'APP.URLS' },
                  ]);
                }

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
                  this.settingsList = this.buildSettingsList([{ id: 'urls', i18nKey: 'APP.URLS' }]);
                  this.currentSetting = this.settingsList[0];
                } else {
                  this.settingsList = this.buildSettingsList([{ id: 'urls', i18nKey: 'APP.URLS' }]);
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
                this.settingsList = this.buildSettingsList([{ id: 'urls', i18nKey: 'APP.URLS' }]);
                if (this.app.samlConfig?.loginVersion?.loginV1) {
                  this.samlForm.controls['loginV2'].setValue(false);
                } else if (this.app.samlConfig?.loginVersion?.loginV2) {
                  this.samlForm.controls['loginV2'].setValue(true);
                  this.samlForm.controls['loginV2BaseURL'].setValue(this.app.samlConfig.loginVersion.loginV2.baseUri);
                }
              }

              if (allowed) {
                this.oidcForm.enable();
                this.oidcForm.controls['clientId'].disable();
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
              if (this.app.oidcConfig?.loginVersion?.loginV1) {
                this.oidcForm.controls['loginV2'].setValue(false);
              } else if (this.app.oidcConfig?.loginVersion?.loginV2) {
                this.oidcForm.controls['loginV2'].setValue(true);
                this.oidcForm.controls['loginV2BaseURL'].setValue(this.app.oidcConfig.loginVersion.loginV2.baseUri);
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

              this.appDataSubject.next();
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
          this.authMethods = [PKCE_METHOD, DEVICE_CODE_METHOD, CUSTOM_METHOD];
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
        this.entityId?.setValue('');
        this.acsURL?.setValue('');
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
        this.app.oidcConfig.devMode = !!this.devMode?.value;
        this.app.oidcConfig.skipNativeAppSuccessPage = !!this.skipNativeAppSuccessPage?.value;

        const req = new UpdateOIDCAppConfigRequest();
        req.setProjectId(this.projectId);
        req.setAppId(this.app.id);

        // configuration
        req.setResponseTypesList(this.app.oidcConfig.responseTypesList);
        req.setAuthMethodType(this.app.oidcConfig.authMethodType);
        req.setGrantTypesList(this.app.oidcConfig.grantTypesList);
        req.setAppType(this.app.oidcConfig.appType);
        const login = new LoginVersion();
        if (this.oidcLoginV2?.value) {
          const loginV2 = new LoginV2();
          loginV2.setBaseUri(this.oidcLoginV2BaseURL?.value);
          login.setLoginV2(loginV2);
        } else {
          login.setLoginV1(new LoginV1());
        }
        req.setLoginVersion(login);

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
        req.setSkipNativeAppSuccessPage(this.app.oidcConfig.skipNativeAppSuccessPage);

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
            setTimeout(() => {
              this.getData(this.projectId, this.appId);
            }, 1000);
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
              this.settingsList = this.buildSettingsList([{ id: 'urls', i18nKey: 'APP.URLS' }]);
              this.currentSetting = this.settingsList[0];
            } else {
              this.settingsList = this.buildSettingsList([{ id: 'urls', i18nKey: 'APP.URLS' }]);
              this.currentSetting = this.settingsList[0];
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

      if (this.app.samlConfig?.metadataUrl?.length > 0) {
        req.setMetadataUrl(this.app.samlConfig?.metadataUrl);
      }
      if (this.app.samlConfig?.metadataXml?.length > 0) {
        req.setMetadataXml(this.app.samlConfig?.metadataXml);
      }

      const login = new LoginVersion();
      if (this.samlLoginV2?.value) {
        const loginV2 = new LoginV2();
        loginV2.setBaseUri(this.samlLoginV2BaseURL?.value);
        login.setLoginV2(loginV2);
      } else {
        login.setLoginV1(new LoginV1());
      }
      req.setLoginVersion(login);

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

  public get oidcLoginV2(): FormControl<boolean> | null {
    return this.oidcForm.get('loginV2') as FormControl<boolean>;
  }

  public get oidcLoginV2BaseURL(): AbstractControl | null {
    return this.oidcForm.get('loginV2BaseURL');
  }

  public get apiAuthMethodType(): AbstractControl | null {
    return this.apiForm.get('authMethodType') as UntypedFormControl;
  }

  public get devMode(): FormControl<boolean> | null {
    return this.oidcForm.get('devMode') as FormControl<boolean>;
  }

  public get skipNativeAppSuccessPage(): FormControl<boolean> | null {
    return this.oidcForm.get('skipNativeAppSuccessPage') as FormControl<boolean>;
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

  public get entityId(): AbstractControl | null {
    return this.samlForm.get('entityId');
  }

  public get acsURL(): AbstractControl | null {
    return this.samlForm.get('acsURL');
  }

  public get samlLoginV2(): FormControl<boolean> | null {
    return this.samlForm.get('loginV2') as FormControl<boolean>;
  }

  public get samlLoginV2BaseURL(): AbstractControl | null {
    return this.samlForm.get('loginV2BaseURL');
  }

  get decodedBase64(): string {
    if (
      this.app &&
      this.app.samlConfig &&
      this.app.samlConfig.metadataXml &&
      typeof this.app.samlConfig.metadataXml === 'string'
    ) {
      return Buffer.from(this.app?.samlConfig.metadataXml, 'base64').toString('utf-8');
    } else {
      return '';
    }
  }

  set decodedBase64(xmlString: string) {
    if (this.app && this.app.samlConfig && this.app.samlConfig.metadataXml) {
      const base64 = Buffer.from(xmlString, 'utf-8').toString('base64');

      if (this.app.samlConfig) {
        this.app.samlConfig.metadataXml = base64;
      }
    }
  }

  public getEnvironmentVariables(): string {
    if (!this.app) return '';

    let issuer = 'https://your-domain.zitadel.cloud';

    this.envSvc.env.pipe(take(1)).subscribe((env) => {
      if (env && env.issuer) {
        issuer = env.issuer;
      }
    });

    let envVars = `ZITADEL_DOMAIN=${issuer}\n`;
    const clientId = this.clientId?.value || 'your-client-id';
    envVars += `ZITADEL_CLIENT_ID=${clientId}\n`;
    envVars += `ZITADEL_PROJECT_ID=${this.projectId}\n`;

    if (this.app.oidcConfig) {
      if (
        this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC ||
        this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST
      ) {
        envVars += `ZITADEL_CLIENT_SECRET=your-client-secret\n`;
      }

      const redirectUris = this.app.oidcConfig.redirectUrisList || [];
      if (redirectUris.length > 0) {
        envVars += `ZITADEL_REDIRECT_URI=${redirectUris[0]}\n`;
      }

      const postLogoutUris = this.app.oidcConfig.postLogoutRedirectUrisList || [];
      if (postLogoutUris.length > 0) {
        envVars += `ZITADEL_POST_LOGOUT_REDIRECT_URI=${postLogoutUris[0]}\n`;
      }
    }

    return envVars;
  }

  public getFrameworkInfo(): {
    id?: string;
    title: string;
    example?: string;
    imgSrcDark?: string;
    imgSrcLight?: string;
    docsLink?: string;
    sdk?: boolean;
    sdkLink?: string;
    sdkCommand?: string;
    exampleLink?: string;
    client?: boolean;
    startCommand?: string;
    installCommand?: string;
    sdkName?: string;
    sdkDocLink?: string;
    hasExampleEnv?: boolean;
    envSetup?: {
      type: string;
      filename?: string;
      description: string;
      variables: Array<{
        name: string;
        description: string;
        placeholder: string;
        required: boolean;
      }>;
    };
  } | null {
    if (!this.framework) return null;

    const frameworkInfo = this.frameworks.find((f) => f.id === this.framework);

    return frameworkInfo || null;
  }

  public getCloneCommand(existing?: boolean): string {
    const frameworkInfo = this.getFrameworkInfo();
    if (!frameworkInfo) {
      return 'git clone https://github.com/zitadel/zitadel-examples';
    }

    if (existing && frameworkInfo.sdk && frameworkInfo.sdkCommand) {
      return frameworkInfo.sdkCommand;
    }

    return `git clone ${frameworkInfo.example}`;
  }

  private replacePlaceholders(command: string): string {
    if (!this.app) return command;

    const issuer = this.cachedIssuer || 'https://your-domain.zitadel.cloud';
    const domain = this.cachedDomain || 'your-domain.zitadel.cloud';
    const clientId = this.app.oidcConfig?.clientId || 'your-client-id';

    return command
      .replace(/\{\{ZITADEL_DOMAIN\}\}/g, domain)
      .replace(/\{\{ZITADEL_ISSUER\}\}/g, issuer)
      .replace(/\{\{ZITADEL_CLIENT_ID\}\}/g, clientId)
      .replace(/\{\{CLIENT_ID\}\}/g, clientId)
      .replace(/\{\{PROJECT_ID\}\}/g, this.projectId)
      .replace(/\{\{PROJECT_NAME\}\}/g, this.app?.name || 'project');
  }

  public getBuildCommand(): string {
    const frameworkInfo = this.getFrameworkInfo();
    if (!frameworkInfo) {
      return 'cd <your-repo> && npm install';
    }

    const example = frameworkInfo.example?.split('/').pop()?.replace('.git', '') || 'your-repo';
    const installCmd = frameworkInfo.installCommand || 'npm install';
    const replacedCmd = this.replacePlaceholders(installCmd);

    return `cd ${example}\n${replacedCmd}`;
  }

  public getRunCommand(): string {
    const frameworkInfo = this.getFrameworkInfo();
    if (!frameworkInfo) {
      return 'npm start';
    }

    const startCmd = frameworkInfo.startCommand || 'npm start';
    return this.replacePlaceholders(startCmd);
  }

  public getSdkCommand(): string {
    const frameworkInfo = this.getFrameworkInfo();
    if (!frameworkInfo || !frameworkInfo.sdkCommand) {
      return 'npm install @zitadel/your-sdk';
    }

    return frameworkInfo.sdkCommand || 'npm install @zitadel/your-sdk';
  }

  // Helper methods to determine available content for framework
  public hasExample(): boolean {
    const frameworkInfo = this.getFrameworkInfo();
    return !!frameworkInfo?.example;
  }

  public hasExampleLink(): boolean {
    const frameworkInfo = this.getFrameworkInfo();
    return !!frameworkInfo?.exampleLink;
  }

  // Helper to build settings list with quickstart in correct position
  private buildSettingsList(additionalSettings: SidenavSetting[] = []): SidenavSetting[] {
    const quickstart = { id: 'quickstart', i18nKey: 'APP.QUICKSTART' };
    const configuration = { id: 'configuration', i18nKey: 'APP.CONFIGURATION' };

    // If new=true, quickstart goes first
    // Otherwise, quickstart goes last
    if (this.isNew()) {
      return [quickstart, configuration, ...additionalSettings];
    } else {
      return [configuration, ...additionalSettings, quickstart];
    }
  }

  public hasSdk(): boolean {
    const frameworkInfo = this.getFrameworkInfo();
    return !!(frameworkInfo?.sdk || frameworkInfo?.sdkCommand);
  }

  public hasBuildCommands(): boolean {
    const frameworkInfo = this.getFrameworkInfo();
    return !!(frameworkInfo?.installCommand && frameworkInfo?.startCommand);
  }

  public hasAnyFrameworkContent(): boolean {
    return this.hasExample() || this.hasSdk() || this.hasBuildCommands();
  }

  public getExistingAppFinalStepNumber(): number {
    const hasSdk = this.hasSdk();
    const hasFrameworkContent = this.hasAnyFrameworkContent();
    const hasDocsLink = !!this.getFrameworkInfo()?.docsLink;

    // If has SDK and framework content with docs link: 1 (SDK) + 2 (env vars) + 3 (docs) = 4
    if (hasSdk && hasFrameworkContent && hasDocsLink) {
      return 4;
    }

    // If has SDK OR (framework content with docs link): step 3
    if (hasSdk || (hasFrameworkContent && hasDocsLink)) {
      return 3;
    }

    // If no SDK and no framework content: step 3 (generic guidance + env vars + integration)
    if (!hasSdk && !hasFrameworkContent) {
      return 3;
    }

    // Default case: step 2
    return 2;
  }

  public requiresClientSecret(): boolean {
    return !!(
      this.app?.oidcConfig &&
      (this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC ||
        this.app.oidcConfig.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST)
    );
  }

  public copyToClipboard(text: string, type: string): void {
    navigator.clipboard.writeText(text).then(() => {
      this.copied = type;
      setTimeout(() => {
        this.copied = '';
      }, 2000);
    });
  }

  public copyEnvironmentVariables(): void {
    this.environmentVariables$.pipe(take(1)).subscribe((envVars: string) => {
      this.copyToClipboard(envVars, 'env');
    });
  }

  public downloadEnvironmentVariables(): void {
    // Convert envVarsArray to env file format
    const envVars = this.envVarsArray.map((v) => `${v.key}=${v.value}`).join('\n');

    const blob = new Blob([envVars], { type: 'text/plain;charset=utf-8' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = '.env';
    link.setAttribute('download', '.env'); // Ensure download attribute is set
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);

    this.copied = 'download';
    setTimeout(() => {
      this.copied = '';
    }, 2000);
  }

  public switchToTab(tabId: string): void {
    const setting = this.settingsList.find((s) => s.id === tabId);
    if (setting) {
      this.currentSetting = setting;
    }
  }

  public switchScenario(scenario: 'new' | 'existing'): void {
    this.selectedScenario = scenario;
    // Regenerate environment variables with appropriate prefixes for the selected scenario
    this.updateEnvVars();
  }

  public selectFramework(framework: any): void {
    this.framework = framework.id;
    this.updateEnvVars();
    this.initializeFormWithFrameworkDefaults();
  }

  public selectFrameworkById(frameworkId: string): void {
    this.framework = frameworkId;
    this.updateEnvVars();
    this.initializeFormWithFrameworkDefaults();

    // Add framework as query param to URL
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { framework: frameworkId },
      queryParamsHandling: 'merge',
    });
  }

  private async checkAuthenticationStatus(): Promise<void> {
    try {
      const milestonesListQuery = new ListQuery();
      milestonesListQuery.setAsc(true);
      milestonesListQuery.setLimit(20);

      const milestoneIsReachedQuery = new IsReachedQuery().setReached(true);
      const milestonesQuery = new MilestoneQuery().setIsReachedQuery(milestoneIsReachedQuery);
      const milestonesReq = new ListMilestonesRequest().setQuery(milestonesListQuery).setQueriesList([milestonesQuery]);

      const response = await this.adminService.listMilestones(milestonesReq);

      const isAuthMilestoneReached = response.resultList.some(
        (milestone) => milestone.type === MilestoneType.MILESTONE_TYPE_AUTHENTICATION_SUCCEEDED_ON_APPLICATION,
      );

      this.isAuthenticated.set(isAuthMilestoneReached);
    } catch (error) {
      console.error('Failed to check authentication status:', error);
      this.isAuthenticated.set(false);
    }
  }

  private initializeFormWithFrameworkDefaults(): void {
    // Only initialize form defaults for new apps
    if (!this.isNew() || !this.framework) {
      return;
    }

    const frameworkConfig = OIDC_CONFIGURATIONS[this.framework as keyof typeof OIDC_CONFIGURATIONS];
    if (!frameworkConfig) {
      return;
    }

    // Initialize OIDC form with framework defaults
    this.oidcForm.patchValue({
      appType: frameworkConfig.appType,
      authMethodType: frameworkConfig.authMethodType,
      responseTypesList: frameworkConfig.responseTypes,
      grantTypesList: frameworkConfig.grantTypes,
    });

    // Initialize OIDC token form with framework defaults
    this.oidcTokenForm.patchValue({
      idTokenRoleAssertion: (frameworkConfig as any).idTokenRoleAssertion,
      idTokenUserinfoAssertion: (frameworkConfig as any).idTokenUserinfoAssertion,
    });

    // Set redirect URIs
    this.redirectUrisList = frameworkConfig.redirectUris;
    this.postLogoutRedirectUrisList = frameworkConfig.postLogoutRedirectUris;
  }
}
