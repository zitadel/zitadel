import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subject, Subscription } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
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
} from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import {
    BASIC_AUTH_METHOD,
    CODE_METHOD,
    getPartialConfigFromAuthMethod,
    IMPLICIT_METHOD,
    PK_JWT_METHOD,
    PKCE_METHOD,
    POST_METHOD,
} from '../authmethods';
import { API_TYPE, AppCreateType, NATIVE_TYPE, RadioItemAppType, USER_AGENT_TYPE, WEB_TYPE } from '../authtypes';


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

    public oidcAppRequest: AddOIDCAppRequest.AsObject = new AddOIDCAppRequest().toObject();
    public apiAppRequest: AddAPIAppRequest.AsObject = new AddAPIAppRequest().toObject();

    public oidcResponseTypes: { type: OIDCResponseType, checked: boolean; disabled: boolean; }[] = [
        { type: OIDCResponseType.OIDC_RESPONSE_TYPE_CODE, checked: false, disabled: false },
        { type: OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN, checked: false, disabled: false },
        { type: OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN, checked: false, disabled: false },
    ];

    public oidcAppTypes: OIDCAppType[] = [
        OIDCAppType.OIDC_APP_TYPE_WEB,
        OIDCAppType.OIDC_APP_TYPE_NATIVE,
        OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
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
        { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC, checked: false, disabled: false, oidc: true },
        { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE, checked: false, disabled: false, oidc: true },
        { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST, checked: false, disabled: false, oidc: true },
    ];

    // stepper
    firstFormGroup!: FormGroup;
    secondFormGroup!: FormGroup;

    // devmode
    public form!: FormGroup;

    public AppCreateType: any = AppCreateType;
    public OIDCAppType: any = OIDCAppType;
    public OIDCGrantType: any = OIDCGrantType;
    public OIDCAuthMethodType: any = OIDCAuthMethodType;

    public oidcGrantTypes: {
        type: OIDCGrantType,
        checked: boolean,
        disabled: boolean,
    }[] = [
            { type: OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE, checked: true, disabled: false },
            { type: OIDCGrantType.OIDC_GRANT_TYPE_IMPLICIT, checked: false, disabled: true },
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
                this.oidcAppRequest.name = this.name?.value;
                this.apiAppRequest.name = this.name?.value;

                if (this.isStepperOIDC) {
                    const oidcAppType = (this.appType?.value as RadioItemAppType).oidcAppType;
                    if (oidcAppType !== undefined) {
                        this.oidcAppRequest.appType = oidcAppType;
                    }

                    switch (this.oidcAppRequest.appType) {
                        case OIDCAppType.OIDC_APP_TYPE_NATIVE:
                            this.authMethods = [
                                PKCE_METHOD,
                            ];

                            // automatically set to PKCE and skip step
                            this.oidcAppRequest.responseTypesList = [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE];
                            this.oidcAppRequest.grantTypesList = [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE];
                            this.oidcAppRequest.authMethodType = OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;

                            break;
                        case OIDCAppType.OIDC_APP_TYPE_WEB:
                            // PK_JWT_METHOD.recommended = false;
                            this.authMethods = [
                                PKCE_METHOD,
                                CODE_METHOD,
                                PK_JWT_METHOD,
                                POST_METHOD,
                            ];

                            this.authMethod?.setValue(PKCE_METHOD.key);
                            break;
                        case OIDCAppType.OIDC_APP_TYPE_USER_AGENT:
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
                this.oidcAppRequest.responseTypesList = partialConfig.oidc?.responseTypesList ?? [];
                this.oidcAppRequest.grantTypesList = partialConfig.oidc?.grantTypesList ?? [];
                this.oidcAppRequest.authMethodType = partialConfig.oidc?.authMethodType ?? OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;
            } else if (this.isStepperAPI && partialConfig && partialConfig.api) {
                this.apiAppRequest.authMethodType = partialConfig.api?.authMethodType ?? APIAuthMethodType.API_AUTH_METHOD_TYPE_BASIC;
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
                this.oidcAppRequest.name = this.formname?.value;
                this.apiAppRequest.name = this.formname?.value;

                this.oidcAppRequest.responseTypesList = this.formresponseTypesList?.value;
                this.oidcAppRequest.grantTypesList = this.formgrantTypesList?.value;

                this.oidcAppRequest.authMethodType = this.formauthMethodType?.value;
                this.apiAppRequest.authMethodType = this.formauthMethodType?.value;

                const oidcAppType = (this.formappType?.value as RadioItemAppType).oidcAppType;
                if (oidcAppType !== undefined) {
                    this.oidcAppRequest.appType = oidcAppType;
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
                { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC, checked: false, disabled: false, oidc: true },
                { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE, checked: false, disabled: false, oidc: true },
                { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST, checked: false, disabled: false, oidc: true },
            ];
            this.authMethod?.setValue(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC);
        } else if (this.isDevAPI) {
            this.form.removeControl('grantTypesList');
            this.form.removeControl('responseTypesList');

            this.authMethodTypes = [
                { type: APIAuthMethodType.API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT, checked: false, disabled: false, api: true },
                { type: APIAuthMethodType.API_AUTH_METHOD_TYPE_BASIC, checked: false, disabled: false, api: true },
            ];
            this.authMethod?.setValue(APIAuthMethodType.API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT);
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
        console.log(this.projectId);
        this.oidcAppRequest.projectId = projectid;
        this.apiAppRequest.projectId = projectid;
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
                .addOIDCApp(this.oidcAppRequest)
                .then((resp) => {
                    this.loading = false;
                    if (resp.clientId || resp.clientSecret) {
                        this.showSavedDialog(resp);
                    } else {
                        this.router.navigate(['projects', this.projectId, 'apps', resp.appId]);
                    }
                })
                .catch(error => {
                    this.loading = false;
                    this.toast.showError(error);
                });
        } else if (appAPICheck) {
            this.loading = true;
            this.mgmtService
                .addAPIApp(this.apiAppRequest)
                .then((resp) => {
                    this.loading = false;

                    if (resp.clientId || resp.clientSecret) {
                        this.showSavedDialog(resp);
                    } else {
                        this.router.navigate(['projects', this.projectId, 'apps', resp.appId]);
                    }
                })
                .catch(error => {
                    this.loading = false;
                    this.toast.showError(error);
                });
        }
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
                clientId: clientId
            }
        });

        dialogRef.afterClosed().subscribe(() => {
            this.router.navigate(['projects', this.projectId, 'apps', added.appId]);
        });
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

