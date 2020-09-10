import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, OnDestroy, OnInit, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute, Params } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap, take } from 'rxjs/operators';
import { OidcIdpConfigUpdate as AdminOidcIdpConfigUpdate } from 'src/app/proto/generated/admin_pb';
import { OidcIdpConfigUpdate as MgmtOidcIdpConfigUpdate } from 'src/app/proto/generated/management_pb';
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
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    private service!: ManagementService | AdminService;
    public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

    private subscription?: Subscription;
    public projectId: string = '';

    public formGroup!: FormGroup;
    public createSteps: number = 1;
    public currentCreateStep: number = 1;

    constructor(
        // private router: Router,
        private toast: ToastService,
        private injector: Injector,
        private route: ActivatedRoute,
        private _location: Location,
    ) {
        this.formGroup = new FormGroup({
            id: new FormControl({ disabled: true, value: '' }, [Validators.required]),
            name: new FormControl({ disabled: true, value: '' }, [Validators.required]),
            logoSrc: new FormControl({ disabled: true, value: '' }, [Validators.required]),
            clientId: new FormControl('', [Validators.required]),
            clientSecret: new FormControl('', [Validators.required]),
            issuer: new FormControl('', [Validators.required]),
            scopesList: new FormControl([], []),
        });

        this.route.data.pipe(switchMap(data => {
            console.log(data.serviceType);
            this.serviceType = data.serviceType;
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    this.service = this.injector.get(ManagementService as Type<ManagementService>);
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.service = this.injector.get(AdminService as Type<AdminService>);
                    break;
            }

            return this.route.params.pipe(take(1));
        })).subscribe((params) => {
            const { id } = params;
            if (id) {
                this.service.IdpByID(id).then(idp => {
                    this.formGroup.patchValue(idp.toObject());
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
        let req: AdminOidcIdpConfigUpdate | MgmtOidcIdpConfigUpdate;

        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                req = new MgmtOidcIdpConfigUpdate();
                break;
            case PolicyComponentServiceType.ADMIN:
                req = new AdminOidcIdpConfigUpdate();
                break;
        }

        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setIssuer(this.issuer?.value);
        req.setScopesList(this.scopesList?.value);

        this.service.UpdateOidcIdpConfig(req).then((idp) => {
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
            const index = this.scopesList.value.indexOf(uri);

            if (index !== undefined && index >= 0) {
                this.scopesList.value.splice(index, 1);
            }
        }
    }

    public get id(): AbstractControl | null {
        return this.formGroup.get('id');
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

    public get logoSrc(): AbstractControl | null {
        return this.formGroup.get('logoSrc');
    }

    public get issuer(): AbstractControl | null {
        return this.formGroup.get('issuer');
    }

    public get scopesList(): AbstractControl | null {
        return this.formGroup.get('scopesList');
    }
}
