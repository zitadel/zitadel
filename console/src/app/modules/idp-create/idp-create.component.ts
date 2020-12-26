import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, OnDestroy, OnInit, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { take } from 'rxjs/operators';
import {
    OidcIdpConfigCreate as AdminOidcIdpConfigCreate,
    OIDCMappingField as authMappingFields,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    OidcIdpConfigCreate as MgmtOidcIdpConfigCreate,
    OIDCMappingField as mgmtMappingFields,
} from '../../proto/generated/zitadel/management_pb';
import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

@Component({
    selector: 'app-idp-create',
    templateUrl: './idp-create.component.html',
    styleUrls: ['./idp-create.component.scss'],
})
export class IdpCreateComponent implements OnInit, OnDestroy {
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    private service!: ManagementService | AdminService;
    public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
    public mappingFields: mgmtMappingFields[] | authMappingFields[] = [];

    private subscription?: Subscription;
    public projectId: string = '';

    public formGroup!: FormGroup;
    public createSteps: number = 1;
    public currentCreateStep: number = 1;
    public loading: boolean = false;
    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private toast: ToastService,
        private injector: Injector,
        private _location: Location,
    ) {
        this.formGroup = new FormGroup({
            name: new FormControl('', [Validators.required]),
            clientId: new FormControl('', [Validators.required]),
            clientSecret: new FormControl('', [Validators.required]),
            issuer: new FormControl('', [Validators.required]),
            scopesList: new FormControl(['openid', 'profile', 'email'], []),
            idpDisplayNameMapping: new FormControl(0),
            usernameMapping: new FormControl(0),
        });

        this.route.data.pipe(take(1)).subscribe(data => {
            this.serviceType = data.serviceType;
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    this.service = this.injector.get(ManagementService as Type<ManagementService>);
                    this.mappingFields = [
                        mgmtMappingFields.OIDCMAPPINGFIELD_PREFERRED_USERNAME,
                        mgmtMappingFields.OIDCMAPPINGFIELD_EMAIL];
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.service = this.injector.get(AdminService as Type<AdminService>);
                    this.mappingFields = [
                        authMappingFields.OIDCMAPPINGFIELD_PREFERRED_USERNAME,
                        authMappingFields.OIDCMAPPINGFIELD_EMAIL];
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

    private getData({ projectid }: Params): void {
        this.projectId = projectid;
    }

    public addIdp(): void {
        let req: AdminOidcIdpConfigCreate | MgmtOidcIdpConfigCreate;

        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                req = new MgmtOidcIdpConfigCreate();
                break;
            case PolicyComponentServiceType.ADMIN:
                req = new AdminOidcIdpConfigCreate();
                break;
        }

        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setIssuer(this.issuer?.value);
        req.setScopesList(this.scopesList?.value);
        req.setIdpDisplayNameMapping(this.idpDisplayNameMapping?.value);
        req.setUsernameMapping(this.usernameMapping?.value);
        this.loading = true;
        this.service.CreateOidcIdp(req).then((idp) => {
            setTimeout(() => {
                this.loading = false;
                this.router.navigate([
                    this.serviceType === PolicyComponentServiceType.MGMT ? 'org' :
                        this.serviceType === PolicyComponentServiceType.ADMIN ? 'iam' : '',
                    'idp', idp.getId()]);
            }, 2000);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public close(): void {
        this._location.back();
    }

    public addScope(event: MatChipInputEvent): void {
        const input = event.input;
        const value = event.value.trim();

        if (value !== '') {
            if (this.scopesList?.value) {
                this.scopesList.value.push(value);
                if (input) {
                    input.value = '';
                }
            }
        }
    }

    public removeScope(uri: string): void {
        if (this.scopesList?.value) {
            const index = this.scopesList.value.indexOf(uri);

            if (index !== undefined && index >= 0) {
                this.scopesList.value.splice(index, 1);
            }
        }
    }

    public get name(): AbstractControl | null {
        return this.formGroup.get('name');
    }

    public get clientId(): AbstractControl | null {
        return this.formGroup.get('clientId');
    }

    public get clientSecret(): AbstractControl | null {
        return this.formGroup.get('clientSecret');
    }

    public get issuer(): AbstractControl | null {
        return this.formGroup.get('issuer');
    }
    public get scopesList(): AbstractControl | null {
        return this.formGroup.get('scopesList');
    }

    public get idpDisplayNameMapping(): AbstractControl | null {
        return this.formGroup.get('idpDisplayNameMapping');
    }

    public get usernameMapping(): AbstractControl | null {
        return this.formGroup.get('usernameMapping');
    }

}
