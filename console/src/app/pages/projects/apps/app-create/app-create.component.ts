import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import {
    Application,
    OIDCApplicationCreate,
    OIDCApplicationType,
    OIDCAuthMethodType,
    OIDCGrantType,
    OIDCResponseType,
} from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';

@Component({
    selector: 'app-app-create',
    templateUrl: './app-create.component.html',
    styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent implements OnInit, OnDestroy {
    private subscription?: Subscription;
    public projectId: string = '';
    public oidcApp: OIDCApplicationCreate.AsObject = new OIDCApplicationCreate().toObject();
    public oidcResponseTypes: OIDCResponseType[] = [
        OIDCResponseType.OIDCRESPONSETYPE_CODE,
        OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN,
        OIDCResponseType.OIDCRESPONSETYPE_ID_TOKEN_TOKEN,
    ];
    public oidcGrantTypes: {
        type: OIDCGrantType,
        checked: boolean,
    }[] = [
            { type: OIDCGrantType.OIDCGRANTTYPE_AUTHORIZATION_CODE, checked: false },
            { type: OIDCGrantType.OIDCGRANTTYPE_IMPLICIT, checked: false },
            { type: OIDCGrantType.OIDCGRANTTYPE_REFRESH_TOKEN, checked: false },
        ];
    public oidcAppTypes: OIDCApplicationType[] = [
        OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB,
        OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT,
        OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE,
    ];
    public oidcAuthMethodType: OIDCAuthMethodType[] = [
        OIDCAuthMethodType.OIDCAUTHMETHODTYPE_BASIC,
        OIDCAuthMethodType.OIDCAUTHMETHODTYPE_NONE,
        OIDCAuthMethodType.OIDCAUTHMETHODTYPE_POST,
    ];

    firstFormGroup!: FormGroup;
    secondFormGroup!: FormGroup;
    thirdFormGroup!: FormGroup;

    public postLogoutRedirectUrisList: string[] = [];

    public addOnBlur: boolean = true;
    public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private toast: ToastService,
        private dialog: MatDialog,
        private projectService: ProjectService,
        private fb: FormBuilder,
        private _location: Location,
    ) {
        this.firstFormGroup = this.fb.group({
            name: ['', [Validators.required]],
            applicationType: ['', [Validators.required]],
        });
        // this.secondFormGroup = this.fb.group({
        //     responseTypesList: ['', []],
        // });
        this.secondFormGroup = this.fb.group({
            authMethodType: ['', [Validators.required]],
        });

        this.secondFormGroup.valueChanges.subscribe(value => {
            console.log(value);
        });
        // this.form = this.fb.group({
        //     name: ['', [Validators.required]],
        //     responseTypesList: ['', []],
        //     grantTypesList: ['', []],
        //     applicationType: ['', []],
        //     authMethodType: [],
        // });
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

    // public saveOIDCApp(): void {
    //     this.oidcApp.name = this.name?.value;
    //     this.oidcApp.applicationType = this.applicationType?.value;
    //     this.oidcApp.grantTypesList = this.grantTypesList?.value;
    //     this.oidcApp.responseTypesList = this.responseTypesList?.value;
    //     // this.oidcApp.authMethodType = this.authMethodType?.value;

    //     this.projectService
    //         .CreateOIDCApp(this.oidcApp)
    //         .then((data: Application) => {
    //             this.showSavedDialog(data.toObject());
    //         })
    //         .catch(error => {
    //             this.toast.showError(error);
    //         });
    // }

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

    public addUri(event: MatChipInputEvent, target: string): void {
        const input = event.input;
        const value = event.value.trim();

        if (value !== '') {
            if (target === 'REDIRECT') {
                this.oidcApp.redirectUrisList.push(value);
            } else if (target === 'POSTREDIRECT') {
                this.oidcApp.postLogoutRedirectUrisList.push(value);
            }
        }

        if (input) {
            input.value = '';
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

    get name(): AbstractControl | null {
        return this.firstFormGroup.get('name');
    }

    get applicationType(): AbstractControl | null {
        return this.firstFormGroup.get('applicationType');
    }

    // get grantTypesList(): AbstractControl | null {
    //     return this.secondFormGroup.get('grantTypesList');
    // }

    getCheckedOidcGrantTypes(): OIDCGrantType[] {
        return this.oidcGrantTypes.filter(gt => gt.checked).map(gt => gt.type);
    }

    get responseTypesList(): AbstractControl | null {
        return this.secondFormGroup.get('responseTypesList');
    }

    // get applicationType(): AbstractControl | null {
    //     return this.form.get('applicationType');
    // }

    get authMethodType(): AbstractControl | null {
        return this.secondFormGroup.get('authMethodType');
    }
}

