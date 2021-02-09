import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { debounceTime } from 'rxjs/operators';
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

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';

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
        {
            titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.TITLE',
            descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.DESCRIPTION',
            type: OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB,
            prefix: 'WEB',
            background: 'rgb(80, 110, 110)',
        },
        {
            titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.TITLE',
            descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.DESCRIPTION',
            type: OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE,
            prefix: 'N',
            background: '#595d80',
        },
        {
            titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.TITLE',
            descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.DESCRIPTION',
            type: OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT,
            prefix: 'UA',
            background: '#6a506e',
        },
    ];

    public authMethod: any[] =
        [
            {
                key: 'CODE',
                titleI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.CODE.TITLE',
                descI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.CODE.DESCRIPTION',
                checked: false,
                disabled: false,
                prefix: 'CODE',
                background: 'rgb(80, 110, 110)',
                responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
                grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
                authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC,
                recommended: false,
            },
            {
                key: 'PKCE',
                titleI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.PKCE.TITLE',
                descI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.PKCE.DESCRIPTION',
                checked: false,
                disabled: false,
                prefix: 'PKCE',
                background: '#595d80',
                responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
                grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
                authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
                recommended: true,
            },
            {
                key: 'ALTERNATIVE',
                titleI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.ALTERNATIVE.TITLE',
                descI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.ALTERNATIVE.DESCRIPTION',
                checked: false,
                disabled: false,
                prefix: 'ALT',
                background: '#6a506e',
                responseType: OIDCResponseType.OIDCRESPONSETYPE_CODE,
                grantType: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE,
                authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST,
            },
            {
                key: 'IMPLICIT',
                titleI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.IMPLICIT.TITLE',
                descI18nKey: 'APP.OIDC.SELECTION.AUTHMETHOD.IMPLICIT.DESCRIPTION',
                checked: false,
                disabled: false,
                prefix: 'IMP',
                background: 'rgb(110, 80, 80)',
                responseType: OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN,
                grantType: OIDCGrantType.OIDCGRANTTYPE_IMPLICIT,
                authMethod: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
                notRecommended: true,
            },
        ];

    public oidcAuthMethodType: { type: OIDCAuthMethodType, checked: boolean, disabled: boolean; }[] = [
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC, checked: false, disabled: false },
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE, checked: false, disabled: false },
        { type: OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST, checked: false, disabled: false },
    ];

    // stepper
    firstFormGroup!: FormGroup;
    secondFormGroup!: FormGroup;
    // thirdFormGroup!: FormGroup;

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

    public redirectControl: FormControl = new FormControl('');
    public postRedirectControl: FormControl = new FormControl('');

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
            applicationType: ['', [Validators.required]],
        });

        this.firstFormGroup.valueChanges.subscribe(value => {
            if (this.firstFormGroup.valid) {
                switch (value.applicationType) {
                    case OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE:
                        this.oidcResponseTypes[0].checked = true;
                        this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];

                        this.oidcApp.grantTypesList =
                            [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE];
                        this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;

                        // this.redirectControl = new FormControl('', [nativeValidator as ValidatorFn]);
                        // this.postRedirectControl = new FormControl('', [nativeValidator as ValidatorFn]);
                        break;
                    case OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB:
                        this.oidcAuthMethodType[0].disabled = false;
                        this.oidcAuthMethodType[1].disabled = false;
                        this.oidcAuthMethodType[2].disabled = false;
                        this.authMethodType?.setValue(OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC);
                        this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC;

                        this.oidcResponseTypes[0].checked = true;
                        this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];
                        this.changeResponseType();

                        this.oidcApp.grantTypesList =
                            [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE];

                        // this.redirectControl = new FormControl('', [webValidator as ValidatorFn]);
                        // this.postRedirectControl = new FormControl('', [webValidator as ValidatorFn]);
                        break;
                    case OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT:
                        this.oidcResponseTypes[0].checked = true;
                        this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];

                        this.oidcApp.grantTypesList =
                            [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE, OIDCGrantType.OIDCGRANTTYPE_IMPLICIT];

                        this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;

                        // this.redirectControl = new FormControl('', [webValidator as ValidatorFn]);
                        // this.postRedirectControl = new FormControl('', [webValidator as ValidatorFn]);
                        break;
                }

                this.oidcApp.name = this.name?.value;
                this.oidcApp.applicationType = this.applicationType?.value;
            }
        });

        this.secondFormGroup = this.fb.group({
            authMethod: ['', [Validators.required]],
        });
        this.secondFormGroup.valueChanges.subscribe(key => {
            switch (key) {
                case 'CODE':
                    this.oidcResponseTypes[0].checked = true;
                    this.oidcResponseTypes[1].checked = false;
                    this.oidcResponseTypes[2].checked = false;
                    this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];

                    this.oidcApp.grantTypesList = [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE];
                    this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;
                    this.changeResponseType();
                    break;
                case 'PKCE':
                    this.oidcResponseTypes[0].checked = true;
                    this.oidcResponseTypes[1].checked = false;
                    this.oidcResponseTypes[2].checked = false;
                    this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];

                    this.oidcApp.grantTypesList = [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE];
                    this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;
                    this.changeResponseType();
                    break;
                case 'ALTERNATIVE':
                    this.oidcResponseTypes[0].checked = true;
                    this.oidcResponseTypes[1].checked = false;
                    this.oidcResponseTypes[2].checked = false;
                    this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];

                    this.oidcApp.grantTypesList = [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE];
                    this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;
                    this.changeResponseType();
                    break;
                case 'IMPLICIT':
                    this.oidcResponseTypes[0].checked = true;
                    this.oidcResponseTypes[1].checked = false;
                    this.oidcResponseTypes[2].checked = false;
                    this.oidcApp.responseTypesList = [OIDCResponseType.OIDCRESPONSETYPE_CODE];

                    this.oidcApp.grantTypesList = [OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE];
                    this.oidcApp.authMethodType = OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE;
                    this.changeResponseType();
                    break;
            }
        });
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
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

    public addUri(input: any, target: string): void {
        const value = input.value.trim();

        if (value !== '') {
            if (target === 'REDIRECT' && this.redirectControl.valid) {
                this.oidcApp.redirectUrisList.push(value);
                if (input) {
                    input.value = '';
                }
            } else if (target === 'POSTREDIRECT' && this.redirectControl.valid) {
                this.oidcApp.postLogoutRedirectUrisList.push(value);
                if (input) {
                    input.value = '';
                }
            }
        }
    }

    public removeUri(uri: string, target: string): void {
        if (target === 'REDIRECT') {
            const index = this.oidcApp.redirectUrisList.indexOf(uri);

            if (index !== undefined && index >= 0) {
                this.oidcApp.redirectUrisList.splice(index, 1);
            }
        } else if (target === 'POSTREDIRECT') {
            const index = this.oidcApp.postLogoutRedirectUrisList.indexOf(uri);

            if (index !== undefined && index >= 0) {
                this.oidcApp.postLogoutRedirectUrisList.splice(index, 1);
            }
        }
    }

    changeResponseType(): void {
        this.oidcApp.responseTypesList = this.oidcResponseTypes.filter(gt => gt.checked).map(gt => gt.type);
    };

    moreThanOneOption(options: Array<{ type: OIDCGrantType, checked: boolean, disabled: boolean; }>): boolean {
        return options.filter(option => option.disabled === false).length > 1;
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

    get authMethodType(): AbstractControl | null {
        return this.secondFormGroup.get('authMethodType');
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

