import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { debounceTime } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import {
    Application,
    OIDCApplicationCreate,
    OIDCApplicationType,
    OIDCAuthMethodType,
    OIDCGrantType,
    OIDCResponseType,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    WEB_TYPE,
    NATIVE_TYPE,
    USER_AGENT_TYPE
} from '../authtypes';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import { CODE_METHOD, getPartialConfigFromAuthMethod, IMPLICIT_METHOD, PKCE_METHOD, PK_JWT_METHOD, POST_METHOD } from '../authmethods';


@Component({
    selector: 'app-app-create',
    templateUrl: './app-create.component.html',
    styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent implements OnInit, OnDestroy {
    private subscription?: Subscription;
    public devmode: boolean = false;
    public projectId: string = '';
    public loading: boolean = false;
    public oidcApp: OIDCApplicationCreate.AsObject = new OIDCApplicationCreate().toObject();

    public oidcResponseTypes: { type: OIDCResponseType, checked: boolean; disabled: boolean; }[] = [
        { type: OIDCResponseType.OIDCRESPONSETYPE_CODE, checked: false, disabled: false },
        { type: OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN, checked: false, disabled: false },
        { type: OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN_TOKEN, checked: false, disabled: false },
    ];

    public oidcAppTypes: any = [
        WEB_TYPE,
        NATIVE_TYPE,
        USER_AGENT_TYPE,
    ];

    public authMethods: RadioItemAuthType[] = [
        PKCE_METHOD,
        CODE_METHOD,
        POST_METHOD,
    ];

    public oidcAuthMethodType: { type: OIDCAuthMethodType, checked: boolean, disabled: boolean; }[] = [
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC, checked: false, disabled: false },
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE, checked: false, disabled: false },
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST, checked: false, disabled: false },
    ];

    // stepper
    firstFormGroup!: FormGroup;
    secondFormGroup!: FormGroup;

    // devmode
    public form!: FormGroup;

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
            applicationType: ['', [Validators.required]],
            authMethodType: ['', [Validators.required]],
        });

        this.form.valueChanges.pipe(debounceTime(300)).subscribe((value) => {
            this.oidcApp.name = this.formname?.value;
            this.oidcApp.applicationType = this.formapplicationType?.value;
            this.oidcApp.responseTypesList = this.formresponseTypesList?.value;
            this.oidcApp.grantTypesList = this.formgrantTypesList?.value;
            this.oidcApp.authMethodType = this.formauthMethodType?.value;
        });

        this.firstFormGroup = this.fb.group({
            name: ['', [Validators.required]],
            applicationType: [OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB, [Validators.required]],
        });

        this.firstFormGroup.valueChanges.subscribe(value => {
            if (this.firstFormGroup.valid) {
                this.oidcApp.name = this.name?.value;
                this.oidcApp.applicationType = this.applicationType?.value;

                switch (this.applicationType?.value) {
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
                        this.authMethods = [
                            PKCE_METHOD,
                            CODE_METHOD,
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
            }
        });

        this.secondFormGroup = this.fb.group({
            authMethod: [this.authMethods[0].key, [Validators.required]],
        });
        this.secondFormGroup.valueChanges.subscribe(form => {
            const partialConfig = getPartialConfigFromAuthMethod(form.authMethod);

            if (partialConfig) {
                this.oidcApp.responseTypesList = partialConfig.responseTypesList ?? [];
                this.oidcApp.grantTypesList = partialConfig.grantTypesList ?? [];
                this.oidcApp.authMethodType = partialConfig.authMethodType ?? OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;
            }
        });
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    public changedAppType(type: OIDCApplicationType) {
        this.firstFormGroup.controls['applicationType'].setValue(type);
    }

    public changedAppAuthMethod(methodKey: string) {
        console.log(methodKey);
        this.secondFormGroup.controls['authMethod'].setValue(methodKey);
    }

    private async getData({ projectid }: Params): Promise<void> {
        this.projectId = projectid;
        this.oidcApp.projectId = projectid;
    }

    public close(): void {
        this._location.back();
    }

    public saveOIDCApp(): void {
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
    }

    public showSavedDialog(app: Application.AsObject): void {
        if (app.oidcConfig !== undefined) {
            const dialogRef = this.dialog.open(AppSecretDialogComponent, {
                data: app.oidcConfig,
            });

            dialogRef.afterClosed().subscribe(result => {
                this.router.navigate(['projects', this.projectId, 'apps', app.id]);
            });
        } else {
            this.router.navigate(['projects', this.projectId, 'apps', app.id]);
        }
    }

    get name(): AbstractControl | null {
        return this.firstFormGroup.get('name');
    }
    get applicationType(): AbstractControl | null {
        return this.firstFormGroup.get('applicationType');
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
    get formapplicationType(): AbstractControl | null {
        return this.form.get('applicationType');
    }
    get formauthMethodType(): AbstractControl | null {
        return this.form.get('authMethodType');
    }
};

