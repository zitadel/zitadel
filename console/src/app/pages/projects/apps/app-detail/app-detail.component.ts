import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatButtonToggleChange } from '@angular/material/button-toggle';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import { ChangeType } from 'src/app/modules/changes/changes.component';
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

enum RedirectType {
    REDIRECT = 'redirect',
    POSTREDIRECT = 'postredirect',
}

@Component({
    selector: 'app-app-detail',
    templateUrl: './app-detail.component.html',
    styleUrls: ['./app-detail.component.scss'],
})
export class AppDetailComponent implements OnInit, OnDestroy {
    public canWrite: boolean = false;
    public errorMessage: string = '';
    public removable: boolean = true;
    public addOnBlur: boolean = true;
    public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

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

    public formatClockSkewLabel(seconds: number) {
        return seconds + 's';
    }

    public AppState: any = AppState;
    public appNameForm!: FormGroup;
    public appForm!: FormGroup;
    public redirectUrisList: string[] = [];
    public postLogoutRedirectUrisList: string[] = [];

    public RedirectType: any = RedirectType;

    public isZitadel: boolean = false;
    public docs!: ZitadelDocs.AsObject;

    public OIDCApplicationType: any = OIDCApplicationType;
    public OIDCAuthMethodType: any = OIDCAuthMethodType;
    public OIDCTokenType: any = OIDCTokenType;

    public redirectControl: FormControl = new FormControl({ value: '', disabled: true });
    public postRedirectControl: FormControl = new FormControl({ value: '', disabled: true });

    public ChangeType: any = ChangeType;
    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private fb: FormBuilder,
        private _location: Location,
        private dialog: MatDialog,
        private mgmtService: ManagementService,
        private authService: GrpcAuthService,
    ) {
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

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    private async getData({ projectid, id }: Params): Promise<void> {
        this.projectId = projectid;
        this.mgmtService.GetIam().then(iam => {
            this.isZitadel = iam.toObject().iamProjectId === this.projectId;
        });
        this.authService.isAllowed(['project.app.write$', 'project.app.write:' + id]).pipe(take(1)).subscribe((allowed) => {
            this.canWrite = allowed;
            this.mgmtService.GetApplicationById(projectid, id).then(app => {
                this.app = app.toObject();
                this.appNameForm.patchValue(this.app);
                console.log(this.app);
                if (allowed) {
                    this.appNameForm.enable();
                    this.appForm.enable();
                    this.redirectControl.enable();
                    this.postRedirectControl.enable();
                }

                if (this.app.oidcConfig?.redirectUrisList) {
                    this.redirectUrisList = this.app.oidcConfig.redirectUrisList;
                }
                if (this.app.oidcConfig?.postLogoutRedirectUrisList) {
                    this.postLogoutRedirectUrisList = this.app.oidcConfig.postLogoutRedirectUrisList;
                }
                if (this.app.oidcConfig?.clockSkew) {
                    const inSecs = this.app.oidcConfig?.clockSkew.seconds + this.app.oidcConfig?.clockSkew.nanos / 100000;
                    console.log(inSecs);
                    this.appForm.controls['clockSkewSeconds'].setValue(inSecs);
                }
                if (this.app.oidcConfig) {
                    this.appForm.patchValue(this.app.oidcConfig);
                }
            }).catch(error => {
                console.error(error);
                this.toast.showError(error);
                this.errorMessage = error.message;
            });
        });


        this.docs = (await this.mgmtService.GetZitadelDocs()).toObject();
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

    public add(input: any, target: RedirectType): void {
        if (target === RedirectType.POSTREDIRECT && this.postRedirectControl.valid) {
            if (input.value !== '' && input.value !== ' ' && input.value !== '/') {
                this.postLogoutRedirectUrisList.push(input.value);
            }
            if (input) {
                input.value = '';
            }
        } else if (target === RedirectType.REDIRECT && this.redirectControl.valid) {
            if (input.value !== '' && input.value !== ' ' && input.value !== '/') {
                this.redirectUrisList.push(input.value);
            }
            if (input) {
                input.value = '';
            }
        }
    }

    public remove(redirect: any, target: RedirectType): void {
        if (target === RedirectType.POSTREDIRECT) {
            const index = this.postLogoutRedirectUrisList.indexOf(redirect);

            if (index >= 0) {
                this.postLogoutRedirectUrisList.splice(index, 1);
            }
        } else if (target === RedirectType.REDIRECT) {
            const index = this.redirectUrisList.indexOf(redirect);

            if (index >= 0) {
                this.redirectUrisList.splice(index, 1);
            }
        }
    }

    public saveApp(): void {
        if (this.appNameForm.valid) {
            this.app.name = this.name?.value;

            this.mgmtService
                .UpdateApplication(this.projectId, this.app.id, this.name?.value)
                .then(() => {
                    this.toast.showInfo('APP.TOAST.OIDCUPDATED', true);
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }


    public saveOIDCApp(): void {
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
                console.log(req.toObject());
                this.mgmtService
                    .UpdateOIDCAppConfig(req)
                    .then(() => {
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
