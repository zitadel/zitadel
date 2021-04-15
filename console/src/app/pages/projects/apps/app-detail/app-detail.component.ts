import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { Subject, Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { CnslLinks } from 'src/app/modules/links/links.component';
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
    selector: 'app-app-detail',
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

    public oidcTokenTypes: OIDCTokenType[] = [
        OIDCTokenType.OIDC_TOKEN_TYPE_BEARER,
        OIDCTokenType.OIDC_TOKEN_TYPE_JWT,
    ];

    public AppState: any = AppState;
    public appNameForm!: FormGroup;
    public oidcForm!: FormGroup;
    public apiForm!: FormGroup;

    public redirectUrisList: string[] = [];
    public postLogoutRedirectUrisList: string[] = [];

    public isZitadel: boolean = false;
    public docs!: GetOIDCInformationResponse.AsObject;

    public OIDCAppType: any = OIDCAppType;
    public OIDCAuthMethodType: any = OIDCAuthMethodType;
    public APIAuthMethodType: any = APIAuthMethodType;
    public OIDCTokenType: any = OIDCTokenType;

    public ChangeType: any = ChangeType;

    public requestRedirectValuesSubject$: Subject<void> = new Subject();
    public copiedKey: any = '';
    public environmentMap: { [key: string]: string; } = {};
    public nextLinks: Array<CnslLinks> = [];

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
        private http: HttpClient,
        private snackbar: MatSnackBar,
    ) {
        this.http.get('./assets/environment.json')
            .toPromise().then((env: any) => {

                this.environmentMap = {
                    issuer: env.issuer,
                    adminServiceUrl: env.adminServiceUrl,
                    mgmtServiceUrl: env.mgmtServiceUrl,
                    authServiceUrl: env.adminServiceUrl,
                };
            });

        this.appNameForm = this.fb.group({
            state: [{ value: '', disabled: true }, []],
            name: [{ value: '', disabled: true }, [Validators.required]],
        });

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
    }

    public formatClockSkewLabel(seconds: number): string {
        return seconds + 's';
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    private initLinks(): void {
        this.nextLinks = [
            {
                i18nTitle: 'APP.PAGES.NEXTSTEPS.0.TITLE',
                i18nDesc: 'APP.PAGES.NEXTSTEPS.0.DESC',
                routerLink: ['/projects', this.projectId],
            },
            {
                i18nTitle: 'APP.PAGES.NEXTSTEPS.1.TITLE',
                i18nDesc: 'APP.PAGES.NEXTSTEPS.1.DESC',
                routerLink: ['/users', 'create'],
            }, {
                i18nTitle: 'APP.PAGES.NEXTSTEPS.2.TITLE',
                i18nDesc: 'APP.PAGES.NEXTSTEPS.2.DESC',
                href: 'https://docs.zitadel.ch',
            },
        ];
    }

    private async getData({ projectid, id }: Params): Promise<void> {
        this.projectId = projectid;

        this.initLinks();

        this.mgmtService.getIAM().then(iam => {
            this.isZitadel = iam.iamProjectId === this.projectId;
        });
        this.authService.isAllowed(['project.app.write$', 'project.app.write:' + projectid])
            .pipe(take(1))
            .subscribe((allowed) => {
                this.canWrite = allowed;
                this.mgmtService.getAppByID(projectid, id).then(app => {
                    if (app.app) {
                        this.app = app.app;
                        this.appNameForm.patchValue(this.app);

                        if (this.app.oidcConfig) {
                            this.getAuthMethodOptions('OIDC');

                            this.initialAuthMethod = this.authMethodFromPartialConfig({ oidc: this.app.oidcConfig });
                            this.currentAuthMethod = this.initialAuthMethod;
                            if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                                if (!this.authMethods.includes(CUSTOM_METHOD)) {
                                    this.authMethods.push(CUSTOM_METHOD);
                                }
                            } else {
                                this.authMethods = this.authMethods.filter(element => element !== CUSTOM_METHOD);
                            }
                        } else if (this.app.apiConfig) {
                            this.getAuthMethodOptions('API');

                            this.initialAuthMethod = this.authMethodFromPartialConfig({ api: this.app.apiConfig });
                            this.currentAuthMethod = this.initialAuthMethod;
                            if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                                if (!this.authMethods.includes(CUSTOM_METHOD)) {
                                    this.authMethods.push(CUSTOM_METHOD);
                                }
                            } else {
                                this.authMethods = this.authMethods.filter(element => element !== CUSTOM_METHOD);
                            }
                        }

                        if (allowed) {
                            this.appNameForm.enable();
                            this.oidcForm.enable();
                            this.apiForm.enable();
                        }

                        if (this.app.oidcConfig?.redirectUrisList) {
                            this.redirectUrisList = this.app.oidcConfig.redirectUrisList;
                        }
                        if (this.app.oidcConfig?.postLogoutRedirectUrisList) {
                            this.postLogoutRedirectUrisList = this.app.oidcConfig.postLogoutRedirectUrisList;
                        }
                        if (this.app.oidcConfig?.clockSkew) {
                            const inSecs = this.app.oidcConfig?.clockSkew.seconds +
                                this.app.oidcConfig?.clockSkew.nanos / 100000;
                            this.oidcForm.controls['clockSkewSeconds'].setValue(inSecs);
                        }
                        if (this.app.oidcConfig) {
                            this.oidcForm.patchValue(this.app.oidcConfig);
                        }

                        this.oidcForm.valueChanges.subscribe((oidcConfig) => {
                            this.initialAuthMethod = this.authMethodFromPartialConfig({ oidc: oidcConfig });
                            if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                                if (!this.authMethods.includes(CUSTOM_METHOD)) {
                                    this.authMethods.push(CUSTOM_METHOD);
                                }
                            } else {
                                this.authMethods = this.authMethods.filter(element => element !== CUSTOM_METHOD);
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
                                this.authMethods = this.authMethods.filter(element => element !== CUSTOM_METHOD);
                            }

                            this.showSaveSnack();
                        });
                    }
                }).catch(error => {
                    console.error(error);
                    this.toast.showError(error);
                    this.errorMessage = error.message;
                });
            });
        this.docs = (await this.mgmtService.getOIDCInformation());
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
                    this.authMethods = [
                        PKCE_METHOD,
                        CUSTOM_METHOD,
                    ];
                    break;
                case OIDCAppType.OIDC_APP_TYPE_WEB:
                    this.authMethods = [
                        PKCE_METHOD,
                        CODE_METHOD,
                        PK_JWT_METHOD,
                        POST_METHOD,
                    ];
                    break;
                case OIDCAppType.OIDC_APP_TYPE_USER_AGENT:
                    this.authMethods = [
                        PKCE_METHOD,
                        IMPLICIT_METHOD,
                    ];
                    break;
            }
        }
        if (type === 'API') {
            this.authMethods = [
                PK_JWT_METHOD,
                BASIC_AUTH_METHOD,
            ];
        }
    }

    public authMethodFromPartialConfig(config: { oidc?: OIDCConfig.AsObject, api?: APIConfig.AsObject; }): string {
        const key = getAuthMethodFromPartialConfig(config);
        return key;
    }

    public setPartialConfigFromAuthMethod(authMethod: string): void {
        const partialConfig = getPartialConfigFromAuthMethod(authMethod);
        if (partialConfig && partialConfig.oidc && this.app.oidcConfig) {
            this.app.oidcConfig.responseTypesList =
                (partialConfig.oidc as Partial<OIDCConfig.AsObject>).responseTypesList
                ?? [];

            this.app.oidcConfig.grantTypesList =
                (partialConfig.oidc as Partial<OIDCConfig.AsObject>).grantTypesList
                ?? [];

            this.app.oidcConfig.authMethodType =
                (partialConfig.oidc as Partial<OIDCConfig.AsObject>).authMethodType
                ?? OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;

            this.oidcForm.patchValue(this.app.oidcConfig);
        } else if (partialConfig && partialConfig.api && this.app.apiConfig) {
            this.app.apiConfig.authMethodType =
                (partialConfig.api as Partial<APIConfig.AsObject>).authMethodType
                ?? APIAuthMethodType.API_AUTH_METHOD_TYPE_BASIC;

            this.apiAuthMethodType?.setValue(this.app.apiConfig.authMethodType);
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
        dialogRef.afterClosed().subscribe(resp => {
            if (resp && this.projectId && this.app.id) {
                this.mgmtService.removeApp(this.projectId, this.app.id).then(() => {
                    this.toast.showInfo('APP.TOAST.DELETED', true);

                    this.router.navigate(['/projects', this.projectId]);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    public changeState(state: AppState): void {
        if (state === AppState.APP_STATE_ACTIVE) {
            this.mgmtService.reactivateApp(this.projectId, this.app.id).then(() => {
                this.app.state = state;
                this.toast.showInfo('APP.TOAST.REACTIVATED', true);
            }).catch((error: any) => {
                this.toast.showError(error);
            });
        } else if (state === AppState.APP_STATE_INACTIVE) {
            this.mgmtService.deactivateApp(this.projectId, this.app.id).then(() => {
                this.app.state = state;
                this.toast.showInfo('APP.TOAST.DEACTIVATED', true);
            }).catch((error: any) => {
                this.toast.showError(error);
            });
        }
    }

    public saveApp(): void {
        if (this.appNameForm.valid) {
            this.app.name = this.name?.value;

            this.mgmtService
                .updateApp(this.projectId, this.app.id, this.name?.value)
                .then(() => {
                    this.toast.showInfo('APP.TOAST.UPDATED', true);
                    this.editState = false;
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }


    public saveOIDCApp(): void {
        this.requestRedirectValuesSubject$.next();
        if (this.appNameForm.valid) {
            this.app.name = this.name?.value;
        }

        if (this.oidcForm.valid) {
            if (this.app.oidcConfig) {
                this.app.oidcConfig.responseTypesList = this.responseTypesList?.value;
                this.app.oidcConfig.grantTypesList = this.grantTypesList?.value;
                this.app.oidcConfig.appType = this.appType?.value;
                this.app.oidcConfig.authMethodType = this.authMethodType?.value;
                this.app.oidcConfig.redirectUrisList = this.redirectUrisList;
                this.app.oidcConfig.postLogoutRedirectUrisList = this.postLogoutRedirectUrisList;
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
                    dur.setNanos((Math.floor(this.clockSkewSeconds?.value % 1) * 10000));
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
                    .catch(error => {
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
                    }
                    this.toast.showInfo('APP.TOAST.OIDCUPDATED', true);
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public regenerateOIDCClientSecret(): void {
        this.mgmtService.regenerateOIDCClientSecret(this.app.id, this.projectId).then(resp => {
            this.toast.showInfo('APP.TOAST.CLIENTSECRETREGENERATED', true);
            this.dialog.open(AppSecretDialogComponent, {
                data: {
                    // clientId: data.toObject() as ClientSecret.AsObject.clientId,
                    clientSecret: resp.clientSecret,
                },
                width: '400px',
            });

        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public regenerateAPIClientSecret(): void {
        this.mgmtService.regenerateAPIClientSecret(this.app.id, this.projectId).then(resp => {
            this.toast.showInfo('APP.TOAST.CLIENTSECRETREGENERATED', true);
            this.dialog.open(AppSecretDialogComponent, {
                data: {
                    // clientId: data.toObject().clientId ?? '',
                    clientSecret: resp.clientSecret,
                },
                width: '400px',
            });

        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public navigateBack(): void {
        this._location.back();
    }

    public get name(): AbstractControl | null {
        return this.appNameForm.get('name');
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

    public get devMode(): AbstractControl | null {
        return this.oidcForm.get('devMode');
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
