import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, OnDestroy, OnInit, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute, Params } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap, take } from 'rxjs/operators';
import {
    IdpStylingType as adminIdpStylingType,
    IdpUpdate as AdminIdpConfigUpdate,
    OidcIdpConfigUpdate as AdminOidcIdpConfigUpdate,
    OIDCMappingField as adminMappingFields,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
    IdpStylingType as mgmtIdpStylingType,
    IdpUpdate as MgmtIdpConfigUpdate,
    OidcIdpConfigUpdate as MgmtOidcIdpConfigUpdate,
    OIDCMappingField as mgmtMappingFields,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

@Component({
    selector: 'app-idp',
    templateUrl: './idp.component.html',
    styleUrls: ['./idp.component.scss'],
})
export class IdpComponent implements OnInit, OnDestroy {
    public mappingFields: mgmtMappingFields[] | adminMappingFields[] = [];
    public styleFields: mgmtIdpStylingType[] | adminIdpStylingType[] = [];

    public showIdSecretSection: boolean = false;
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    private service!: ManagementService | AdminService;
    public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

    private subscription?: Subscription;
    public projectId: string = '';

    public idpForm!: FormGroup;
    public oidcConfigForm!: FormGroup;

    constructor(
        private toast: ToastService,
        private injector: Injector,
        private route: ActivatedRoute,
        private _location: Location,
    ) {
        this.idpForm = new FormGroup({
            id: new FormControl({ disabled: true, value: '' }, [Validators.required]),
            name: new FormControl('', [Validators.required]),
            stylingType: new FormControl('', [Validators.required]),
        });

        this.oidcConfigForm = new FormGroup({
            clientId: new FormControl('', [Validators.required]),
            clientSecret: new FormControl(''),
            issuer: new FormControl('', [Validators.required]),
            scopesList: new FormControl([], []),
            idpDisplayNameMapping: new FormControl(0),
            usernameMapping: new FormControl(0),
        });

        this.route.data.pipe(switchMap(data => {
            this.serviceType = data.serviceType;
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    this.service = this.injector.get(ManagementService as Type<ManagementService>);
                    this.mappingFields = [
                        mgmtMappingFields.OIDCMAPPINGFIELD_PREFERRED_USERNAME,
                        mgmtMappingFields.OIDCMAPPINGFIELD_EMAIL];
                    this.styleFields = [
                        mgmtIdpStylingType.IDPSTYLINGTYPE_UNSPECIFIED,
                        mgmtIdpStylingType.IDPSTYLINGTYPE_GOOGLE];
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.service = this.injector.get(AdminService as Type<AdminService>);
                    this.mappingFields = [
                        adminMappingFields.OIDCMAPPINGFIELD_PREFERRED_USERNAME,
                        adminMappingFields.OIDCMAPPINGFIELD_EMAIL];
                    this.styleFields = [
                        adminIdpStylingType.IDPSTYLINGTYPE_UNSPECIFIED,
                        adminIdpStylingType.IDPSTYLINGTYPE_GOOGLE];
                    break;
            }

            return this.route.params.pipe(take(1));
        })).subscribe((params) => {
            const { id } = params;
            if (id) {
                this.service.IdpByID(id).then(idp => {
                    const idpObject = idp.toObject();
                    this.idpForm.patchValue(idpObject);
                    if (idpObject.oidcConfig) {
                        this.oidcConfigForm.patchValue(idpObject.oidcConfig);
                    }
                });
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

    public updateIdp(): void {
        let req: AdminIdpConfigUpdate | MgmtIdpConfigUpdate;

        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                req = new MgmtIdpConfigUpdate();
                break;
            case PolicyComponentServiceType.ADMIN:
                req = new AdminIdpConfigUpdate();
                break;
        }

        req.setId(this.id?.value);
        req.setName(this.name?.value);
        req.setStylingType(this.stylingType?.value);

        this.service.UpdateIdp(req).then((idp) => {
            this.toast.showInfo('IDP.TOAST.SAVED', true);
            // this.router.navigate(['idp', ]);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public updateOidcConfig(): void {
        let req: AdminOidcIdpConfigUpdate | MgmtOidcIdpConfigUpdate;

        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                req = new MgmtOidcIdpConfigUpdate();
                break;
            case PolicyComponentServiceType.ADMIN:
                req = new AdminOidcIdpConfigUpdate();
                break;
        }

        req.setIdpId(this.id?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setIssuer(this.issuer?.value);
        req.setScopesList(this.scopesList?.value);
        req.setUsernameMapping(this.usernameMapping?.value);
        req.setIdpDisplayNameMapping(this.idpDisplayNameMapping?.value);

        this.service.UpdateOidcIdpConfig(req).then((oidcConfig) => {
            this.toast.showInfo('IDP.TOAST.SAVED', true);
            // this.router.navigate(['idp', ]);
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
            const index = this.scopesList?.value.indexOf(uri);

            if (index !== undefined && index >= 0) {
                this.scopesList?.value.splice(index, 1);
            }
        }
    }

    public get backroutes(): string[] {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return ['/org', 'policy', 'login'];
            case PolicyComponentServiceType.ADMIN:
                return ['/iam', 'policy', 'login'];
        }
    }

    public get id(): AbstractControl | null {
        return this.idpForm.get('id');
    }

    public get name(): AbstractControl | null {
        return this.idpForm.get('name');
    }

    public get stylingType(): AbstractControl | null {
        return this.idpForm.get('stylingType');
    }

    public get clientId(): AbstractControl | null {
        return this.oidcConfigForm.get('clientId');
    }

    public get clientSecret(): AbstractControl | null {
        return this.oidcConfigForm.get('clientSecret');
    }

    public get issuer(): AbstractControl | null {
        return this.oidcConfigForm.get('issuer');
    }

    public get scopesList(): AbstractControl | null {
        return this.oidcConfigForm.get('scopesList');
    }

    public get idpDisplayNameMapping(): AbstractControl | null {
        return this.oidcConfigForm.get('idpDisplayNameMapping');
    }

    public get usernameMapping(): AbstractControl | null {
        return this.oidcConfigForm.get('usernameMapping');
    }
}
