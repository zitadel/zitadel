import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subject, Subscription } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import {
    APIApplicationCreate,
    APIAuthMethodType,
    Application,
    OIDCApplicationCreate,
    OIDCApplicationType,
    OIDCAuthMethodType,
    OIDCConfig,
    OIDCGrantType,
    OIDCResponseType,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    WEB_TYPE,
    NATIVE_TYPE,
    USER_AGENT_TYPE,
    API_TYPE,
    RadioItemAppType,
    AppCreateType
} from '../authtypes';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import { CODE_METHOD, getPartialConfigFromAuthMethod, IMPLICIT_METHOD, BASIC_AUTH_METHOD, PKCE_METHOD, PK_JWT_METHOD, POST_METHOD } from '../authmethods';
import { StepperSelectionEvent } from '@angular/cdk/stepper';


@Component({
    selector: 'app-app-create',
    templateUrl: './app-create.component.html',
    styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent implements OnInit, OnDestroy {
    private subscription?: Subscription;
    private destroyed$: Subject<void> = new Subject();
    public devmode: boolean = false;
    public projectId: string = '';
    public loading: boolean = false;

    public oidcApp: OIDCApplicationCreate.AsObject = new OIDCApplicationCreate().toObject();
    public apiApp: APIApplicationCreate.AsObject = new APIApplicationCreate().toObject();

    public oidcResponseTypes: { type: OIDCResponseType, checked: boolean; disabled: boolean; }[] = [
        { type: OIDCResponseType.OIDCRESPONSETYPE_CODE, checked: false, disabled: false },
        { type: OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN, checked: false, disabled: false },
        { type: OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN_TOKEN, checked: false, disabled: false },
    ];

    public oidcAppTypes: OIDCApplicationType[] = [
        OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB,
        OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE,
        OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT,
    ];
    public appTypes: any = [
        WEB_TYPE,
        NATIVE_TYPE,
        USER_AGENT_TYPE,
        API_TYPE,
    ];

    public authMethods: RadioItemAuthType[] = [
        PKCE_METHOD,
        CODE_METHOD,
        PK_JWT_METHOD,
        POST_METHOD,
    ];

    // set to oidc first
    public authMethodTypes: { type: OIDCAuthMethodType | APIAuthMethodType, checked: boolean, disabled: boolean; api?: boolean; oidc?: boolean; }[] = [
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC, checked: false, disabled: false, oidc: true },
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE, checked: false, disabled: false, oidc: true },
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST, checked: false, disabled: false, oidc: true },
    ];

    // stepper
    firstFormGroup!: FormGroup;
    secondFormGroup!: FormGroup;

    // devmode
    public form!: FormGroup;

    public AppCreateType: any = AppCreateType;
    public OIDCApplicationType: any = OIDCApplicationType;
    public OIDCGrantType: any = OIDCGrantType;
    public OIDCAuthMethodType: any = OIDCAuthMethodType;

    public oidcGrantTypes: {
        type: OIDCGrantType,
        checked: boolean,
        disabled: boolean,
    }[] = [
            { type: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE, checked: true, disabled: false },
            { type: OIDCGrantType.OIDCGRANTTYPE_IMPLICIT, checked: false, disabled: true },
            // { type: OIDCGrantType.OIDCGRANTTYPE_REFRESH_TOKEN, checked: false, disabled: true },
            // TODO show when implemented
        ];

    public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
    public requestRedirectValuesSubject$: Subject<void> = new Subject();

    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private toast: ToastService,
        private dialog: MatDialog,
        private mgmtService: ManagementService,
        private fb: FormBuilder,
        private _location: Location,
    ) {
        this.form = this.fb.group({
            name: ['', [Validators.required]],
            responseTypesList: ['', [Validators.required]],
            grantTypesList: ['', [Validators.required]],
            appType: ['', [Validators.required]],
            authMethodType: ['', [Validators.required]],
        });

        this.initForm();

        this.firstFormGroup = this.fb.group({
            name: ['', [Validators.required]],
            appType: [WEB_TYPE, [Validators.required]],
        });

        this.firstFormGroup.valueChanges.subscribe(value => {
            if (this.firstFormGroup.valid) {
                this.oidcApp.name = this.name?.value;
                this.apiApp.name = this.name?.value;

                if (this.isStepperOIDC) {
                    const oidcAppType = (this.appType?.value as RadioItemAppType).oidcApplicationType;
                    if (oidcAppType !== undefined) {
                        this.oidcApp.applicationType = oidcAppType;
                    }

                    switch (this.oidcApp.applicationType) {
                        case OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE:
                            this.authMethods = [
                                PKCE_METHOD,
                            ];

                            // automatically set to PKCE and skip step
                            this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];
                            this.oidcApp.grantTypesList = [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE];
                            this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;

                            break;
                        case OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB:
                            // PK_JWT_METHOD.recommended = false;
                            this.authMethods = [
                                PKCE_METHOD,
                                CODE_METHOD,
                                PK_JWT_METHOD,
                                POST_METHOD,
                            ];

                            this.authMethod?.setValue(PKCE_METHOD.key);
                            break;
                        case OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT:
                            this.authMethods = [
                                PKCE_METHOD,
                                IMPLICIT_METHOD,
                            ];

                            this.authMethod?.setValue(PKCE_METHOD.key);
                            break;
                    }
                } else if (this.isStepperAPI) {
                    // PK_JWT_METHOD.recommended = true;
                    this.authMethods = [
                        PK_JWT_METHOD,
                        BASIC_AUTH_METHOD,
                    ];

                    this.authMethod?.setValue(PK_JWT_METHOD.key);
                }
            }
        });

        this.secondFormGroup = this.fb.group({
            authMethod: [this.authMethods[0].key, [Validators.required]],
        });
        this.secondFormGroup.valueChanges.subscribe(form => {
            const partialConfig = getPartialConfigFromAuthMethod(form.authMethod);

            if (this.isStepperOIDC && partialConfig && partialConfig.oidc) {
                this.oidcApp.responseTypesList = partialConfig.oidc?.responseTypesList ?? [];
                this.oidcApp.grantTypesList = partialConfig.oidc?.grantTypesList ?? [];
                this.oidcApp.authMethodType = partialConfig.oidc?.authMethodType ?? OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;
            } else if (this.isStepperAPI && partialConfig && partialConfig.api) {
                this.apiApp.authMethodType = partialConfig.api?.authMethodType ?? APIAuthMethodType.APIAUTHMETHODTYPE_BASIC;
            }
        });
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
        this.destroyed$.next();
    }

    public initForm(): void {
        this.form.valueChanges.pipe(
            takeUntil(this.destroyed$),
            debounceTime(150)).subscribe(() => {
                this.oidcApp.name = this.formname?.value;
                this.apiApp.name = this.formname?.value;

                this.oidcApp.responseTypesList = this.formresponseTypesList?.value;
                this.oidcApp.grantTypesList = this.formgrantTypesList?.value;

                this.oidcApp.authMethodType = this.formauthMethodType?.value;
                this.apiApp.authMethodType = this.formauthMethodType?.value;

                const oidcAppType = (this.formappType?.value as RadioItemAppType).oidcApplicationType;
                if (oidcAppType !== undefined) {
                    this.oidcApp.applicationType = oidcAppType;
                }
            });

        this.formappType?.valueChanges.pipe(takeUntil(this.destroyed$)).subscribe(() => {
            this.setDevFormValidators();
        });
    }

    public setDevFormValidators(): void {
        if (this.isDevOIDC) {
            const grantTypesControl = new FormControl('', [Validators.required]);
            const responseTypesControl = new FormControl('', [Validators.required]);

            this.form.addControl('grantTypesList', grantTypesControl);
            this.form.addControl('responseTypesList', responseTypesControl);

            this.authMethodTypes = [
                { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC, checked: false, disabled: false, oidc: true },
                { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE, checked: false, disabled: false, oidc: true },
                { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST, checked: false, disabled: false, oidc: true },
            ];
            this.authMethod?.setValue(OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC);
        } else if (this.isDevAPI) {
            this.form.removeControl('grantTypesList');
            this.form.removeControl('responseTypesList');

            this.authMethodTypes = [
                { type: APIAuthMethodType.APIAUTHMETHODTYPE_PRIVATE_KEY_JWT, checked: false, disabled: false, api: true },
                { type: APIAuthMethodType.APIAUTHMETHODTYPE_BASIC, checked: false, disabled: false, api: true },
            ];
            this.authMethod?.setValue(APIAuthMethodType.APIAUTHMETHODTYPE_PRIVATE_KEY_JWT);
        }
        this.form.updateValueAndValidity();
    }

    public changeStep(event: StepperSelectionEvent) {
        if (event.selectedIndex >= 2) {
            this.requestRedirectValuesSubject$.next();
        };
    }

    private async getData({ projectid }: Params): Promise<void> {
        this.projectId = projectid;
        this.oidcApp.projectId = projectid;
        this.apiApp.projectId = projectid;
    }

    public close(): void {
        this._location.back();
    }

    public createApp(): void {
        const appOIDCCheck = this.devmode ? this.isDevOIDC : this.isStepperOIDC;
        const appAPICheck = this.devmode ? this.isDevAPI : this.isStepperAPI;

        if (appOIDCCheck) {
            this.requestRedirectValuesSubject$.next();

            this.loading = true;
            this.mgmtService
                .CreateOIDCApp(this.oidcApp)
                .then((data: Application) => {
                    this.loading = false;
                    const response = data.toObject();
                    if (response.oidcConfig?.authMethodType !== OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE) {
                        this.showSavedDialog(response);
                    } else {
                        this.router.navigate(['projects', this.projectId, 'apps', response.id]);
                    }
                })
                .catch(error => {
                    this.loading = false;
                    this.toast.showError(error);
                });
        } else if (appAPICheck) {
            this.loading = true;
            this.mgmtService
                .CreateAPIApplication(this.apiApp)
                .then((data: Application) => {
                    this.loading = false;
                    const response = data.toObject();
                    if (response.apiConfig?.authMethodType == APIAuthMethodType.APIAUTHMETHODTYPE_BASIC) {
                        this.showSavedDialog(response);
                    } else {
                        this.router.navigate(['projects', this.projectId, 'apps', response.id]);
                    }
                })
                .catch(error => {
                    this.loading = false;
                    this.toast.showError(error);
                });
        }
    }

    public showSavedDialog(app: Application.AsObject): void {
        if (app.oidcConfig?.clientSecret !== undefined) {
            const dialogRef = this.dialog.open(AppSecretDialogComponent, {
                data: app.oidcConfig,
            });

            dialogRef.afterClosed().subscribe(() => {
                this.router.navigate(['projects', this.projectId, 'apps', app.id]);
            });
        }
        else if (app.apiConfig?.clientSecret !== undefined) {
            const dialogRef = this.dialog.open(AppSecretDialogComponent, {
                data: app.apiConfig,
            });

            dialogRef.afterClosed().subscribe(() => {
                this.router.navigate(['projects', this.projectId, 'apps', app.id]);
            });
        } else {
            this.router.navigate(['projects', this.projectId, 'apps', app.id]);
        }
    }

    get name(): AbstractControl | null {
        return this.firstFormGroup.get('name');
    }
    get appType(): AbstractControl | null {
        return this.firstFormGroup.get('appType');
    }
    public grantTypeChecked(type: OIDCGrantType): boolean {
        return this.oidcGrantTypes.filter(gt => gt.checked).map(gt => gt.type).findIndex(t => t === type) > -1;
    }
    get responseTypesList(): AbstractControl | null {
        return this.secondFormGroup.get('responseTypesList');
    }
    get authMethod(): AbstractControl | null {
        return this.secondFormGroup.get('authMethod');
    }

    // devmode

    get formname(): AbstractControl | null {
        return this.form.get('name');
    }
    get formresponseTypesList(): AbstractControl | null {
        return this.form.get('responseTypesList');
    }
    get formgrantTypesList(): AbstractControl | null {
        return this.form.get('grantTypesList');
    }
    get formappType(): AbstractControl | null {
        return this.form.get('appType');
    }
    // get formapplicationType(): AbstractControl | null {
    //     return this.form.get('applicationType');
    // }
    get formauthMethodType(): AbstractControl | null {
        return this.form.get('authMethodType');
    }

    get isDevOIDC(): boolean {
        return (this.formappType?.value as RadioItemAppType).createType == AppCreateType.OIDC;
    }

    get isStepperOIDC(): boolean {
        return (this.appType?.value as RadioItemAppType).createType == AppCreateType.OIDC;
    }

    get isDevAPI(): boolean {
        return (this.formappType?.value as RadioItemAppType).createType == AppCreateType.API;
    }

    get isStepperAPI(): boolean {
        return (this.appType?.value as RadioItemAppType).createType == AppCreateType.API;
    }
};

