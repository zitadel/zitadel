import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatButtonToggleChange } from '@angular/material/button-toggle';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router, RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { Subject, Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { CnslLinks } from 'src/app/modules/links/links.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import {
    Application,
    AppState,
    OIDCApplicationType,
    OIDCAuthMethodType,
    OIDCConfig,
    OIDCConfigUpdate,
    OIDCGrantType,
    OIDCResponseType,
    OIDCTokenType,
    ZitadelDocs,
} from 'src/app/proto/generated/management_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import { CODE_METHOD, getAuthMethodFromPartialConfig, getPartialConfigFromAuthMethod, IMPLICIT_METHOD, PKCE_METHOD, PK_JWT_METHOD, POST_METHOD, CUSTOM_METHOD } from '../authmethods';

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

    public authMethods: RadioItemAuthType[] = [
        PKCE_METHOD,
        CODE_METHOD,
        POST_METHOD,
    ];
    private subscription?: Subscription;
    public projectId: string = '';
    public app!: Application.AsObject;
    public oidcResponseTypes: OIDCResponseType[] = [
        OIDCResponseType.OIDCRESPONSETYPE_CODE,
        OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN,
        OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN_TOKEN,
    ];
    public oidcGrantTypes: OIDCGrantType[] = [
        OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
        OIDCGrantType.OIDCGRANTTYPE_IMPLICIT,
        OIDCGrantType.OIDCGRANTTYPE_REFRESH_TOKEN,
    ];
    public oidcAppTypes: OIDCApplicationType[] = [
        OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB,
        OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT,
        OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE,
    ];

    public oidcAuthMethodType: OIDCAuthMethodType[] = [
        OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC,
        OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST,
        OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
    ];

    public oidcTokenTypes: OIDCTokenType[] = [
        OIDCTokenType.OIDCTOKENTYPE_BEARER,
        OIDCTokenType.OIDCTOKENTYPE_JWT,
    ];

    public AppState: any = AppState;
    public appNameForm!: FormGroup;
    public appForm!: FormGroup;

    public redirectUrisList: string[] = [];
    public postLogoutRedirectUrisList: string[] = [];

    public isZitadel: boolean = false;
    public docs!: ZitadelDocs.AsObject;

    public OIDCApplicationType: any = OIDCApplicationType;
    public OIDCAuthMethodType: any = OIDCAuthMethodType;
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
        this.appForm = this.fb.group({
            devMode: [{ value: false, disabled: true }, []],
            clientId: [{ value: '', disabled: true }],
            responseTypesList: [{ value: [], disabled: true }],
            grantTypesList: [{ value: [], disabled: true }],
            applicationType: [{ value: '', disabled: true }],
            authMethodType: [{ value: '', disabled: true }],
            accessTokenType: [{ value: '', disabled: true }],
            accessTokenRoleAssertion: [{ value: false, disabled: true }],
            idTokenRoleAssertion: [{ value: false, disabled: true }],
            idTokenUserinfoAssertion: [{ value: false, disabled: true }],
            clockSkewSeconds: [{ value: 0, disabled: true }],
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
                href: 'https://docs.zitadel.ch'
            },
        ];
    }

    private async getData({ projectid, id }: Params): Promise<void> {
        this.projectId = projectid;

        this.initLinks();

        this.mgmtService.GetIam().then(iam => {
            this.isZitadel = iam.toObject().iamProjectId === this.projectId;
        });
        this.authService.isAllowed(['project.app.write$', 'project.app.write:' + projectid]).pipe(take(1)).subscribe((allowed) => {
            this.canWrite = allowed;
            this.mgmtService.GetApplicationById(projectid, id).then(app => {
                this.app = app.toObject();
                this.appNameForm.patchValue(this.app);

                this.getAuthMethodOptions();
                if (this.app.oidcConfig) {
                    this.initialAuthMethod = this.authMethodFromPartialConfig(this.app.oidcConfig);
                    this.currentAuthMethod = this.initialAuthMethod;
                    if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                        if (!this.authMethods.includes(CUSTOM_METHOD)) {
                            this.authMethods.push(CUSTOM_METHOD);
                        }
                    } else {
                        this.authMethods = this.authMethods.filter(element => element != CUSTOM_METHOD);
                    }
                }

                if (allowed) {
                    this.appNameForm.enable();
                    this.appForm.enable();
                }

                if (this.app.oidcConfig?.redirectUrisList) {
                    this.redirectUrisList = this.app.oidcConfig.redirectUrisList;
                }
                if (this.app.oidcConfig?.postLogoutRedirectUrisList) {
                    this.postLogoutRedirectUrisList = this.app.oidcConfig.postLogoutRedirectUrisList;
                }
                if (this.app.oidcConfig?.clockSkew) {
                    const inSecs = this.app.oidcConfig?.clockSkew.seconds + this.app.oidcConfig?.clockSkew.nanos / 100000;
                    this.appForm.controls['clockSkewSeconds'].setValue(inSecs);
                }
                if (this.app.oidcConfig) {
                    this.appForm.patchValue(this.app.oidcConfig);
                }

                this.appForm.valueChanges.subscribe(oidcConfig => {
                    this.initialAuthMethod = this.authMethodFromPartialConfig(oidcConfig);
                    if (this.initialAuthMethod === CUSTOM_METHOD.key) {
                        if (!this.authMethods.includes(CUSTOM_METHOD)) {
                            this.authMethods.push(CUSTOM_METHOD);
                        }
                    } else {
                        this.authMethods = this.authMethods.filter(element => element != CUSTOM_METHOD);
                    }
                });
            }).catch(error => {
                console.error(error);
                this.toast.showError(error);
                this.errorMessage = error.message;
            });
        });
        this.docs = (await this.mgmtService.GetZitadelDocs()).toObject();
    }

    private getAuthMethodOptions(): void {
        switch (this.app.oidcConfig?.applicationType) {
            case OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE:
                this.authMethods = [
                    PKCE_METHOD,
                    CUSTOM_METHOD,
                ];
                break;
            case OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB:
                this.authMethods = [
                    PKCE_METHOD,
                    CODE_METHOD,
                    POST_METHOD,
                ];
                break;
            case OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT:
                this.authMethods = [
                    PKCE_METHOD,
                    IMPLICIT_METHOD,
                ];
                break;
        }
    }

    public authMethodFromPartialConfig(config: OIDCConfig.AsObject): string {
        const key = getAuthMethodFromPartialConfig(config);
        return key;
    }

    public setPartialConfigFromAuthMethod(authMethod: string): void {
        const partialConfig = getPartialConfigFromAuthMethod(authMethod);

        if (partialConfig && this.app.oidcConfig) {
            this.app.oidcConfig.responseTypesList = partialConfig.responseTypesList ?? [];
            this.app.oidcConfig.grantTypesList = partialConfig.grantTypesList ?? [];
            this.app.oidcConfig.authMethodType = partialConfig.authMethodType ?? OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;
            this.appForm.patchValue(this.app.oidcConfig);
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
                this.mgmtService.RemoveApplication(this.projectId, this.app.id).then(() => {
                    this.toast.showInfo('APP.TOAST.DELETED', true);

                    this.router.navigate(['/projects', this.projectId]);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    public changeState(event: MatButtonToggleChange): void {
        if (event.value === AppState.APPSTATE_ACTIVE) {
            this.mgmtService.ReactivateApplication(this.projectId, this.app.id).then(() => {
                this.toast.showInfo('APP.TOAST.REACTIVATED', true);
            }).catch((error: any) => {
                this.toast.showError(error);
            });
        } else if (event.value === AppState.APPSTATE_INACTIVE) {
            this.mgmtService.DeactivateApplication(this.projectId, this.app.id).then(() => {
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
                .UpdateApplication(this.projectId, this.app.id, this.name?.value)
                .then(() => {
                    this.toast.showInfo('APP.TOAST.OIDCUPDATED', true);
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

        if (this.appForm.valid) {
            if (this.app.oidcConfig) {
                this.app.oidcConfig.responseTypesList = this.responseTypesList?.value;
                this.app.oidcConfig.grantTypesList = this.grantTypesList?.value;
                this.app.oidcConfig.applicationType = this.applicationType?.value;
                this.app.oidcConfig.authMethodType = this.authMethodType?.value;
                this.app.oidcConfig.redirectUrisList = this.redirectUrisList;
                this.app.oidcConfig.postLogoutRedirectUrisList = this.postLogoutRedirectUrisList;
                this.app.oidcConfig.devMode = this.devMode?.value;
                this.app.oidcConfig.accessTokenType = this.accessTokenType?.value;
                this.app.oidcConfig.accessTokenRoleAssertion = this.accessTokenRoleAssertion?.value;
                this.app.oidcConfig.idTokenRoleAssertion = this.idTokenRoleAssertion?.value;
                this.app.oidcConfig.idTokenUserinfoAssertion = this.idTokenUserinfoAssertion?.value;


                const req = new OIDCConfigUpdate();
                req.setProjectId(this.projectId);
                req.setApplicationId(this.app.id);
                req.setRedirectUrisList(this.app.oidcConfig.redirectUrisList);
                req.setResponseTypesList(this.app.oidcConfig.responseTypesList);
                req.setAuthMethodType(this.app.oidcConfig.authMethodType);
                req.setPostLogoutRedirectUrisList(this.app.oidcConfig.postLogoutRedirectUrisList);
                req.setGrantTypesList(this.app.oidcConfig.grantTypesList);
                req.setApplicationType(this.app.oidcConfig.applicationType);
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
                    .UpdateOIDCAppConfig(req)
                    .then(() => {
                        if (this.app.oidcConfig) {
                            this.currentAuthMethod = this.authMethodFromPartialConfig(this.app.oidcConfig);
                        }
                        this.toast.showInfo('APP.TOAST.OIDCUPDATED', true);
                    })
                    .catch(error => {
                        this.toast.showError(error);
                    });
            }
        }
    }

    public regenerateOIDCClientSecret(): void {
        this.mgmtService.RegenerateOIDCClientSecret(this.app.id, this.projectId).then((data: OIDCConfig) => {
            this.toast.showInfo('APP.TOAST.OIDCCLIENTSECRETREGENERATED', true);
            this.dialog.open(AppSecretDialogComponent, {
                data: {
                    clientId: data.toObject().clientId,
                    clientSecret: data.toObject().clientSecret,
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
        return this.appForm.get('clientId');
    }

    public get responseTypesList(): AbstractControl | null {
        return this.appForm.get('responseTypesList');
    }

    public get grantTypesList(): AbstractControl | null {
        return this.appForm.get('grantTypesList');
    }

    public get applicationType(): AbstractControl | null {
        return this.appForm.get('applicationType');
    }

    public get authMethodType(): AbstractControl | null {
        return this.appForm.get('authMethodType');
    }

    public get devMode(): AbstractControl | null {
        return this.appForm.get('devMode');
    }

    public get accessTokenType(): AbstractControl | null {
        return this.appForm.get('accessTokenType');
    }

    public get idTokenRoleAssertion(): AbstractControl | null {
        return this.appForm.get('idTokenRoleAssertion');
    }

    public get accessTokenRoleAssertion(): AbstractControl | null {
        return this.appForm.get('accessTokenRoleAssertion');
    }

    public get idTokenUserinfoAssertion(): AbstractControl | null {
        return this.appForm.get('idTokenUserinfoAssertion');
    }

    public get clockSkewSeconds(): AbstractControl | null {
        return this.appForm.get('clockSkewSeconds');
    }
}
