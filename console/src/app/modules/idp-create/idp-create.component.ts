import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { OidcIdpConfigCreate } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-idp-create',
    templateUrl: './idp-create.component.html',
    styleUrls: ['./idp-create.component.scss'],
})
export class IdpCreateComponent implements OnInit, OnDestroy {
    private subscription?: Subscription;
    public projectId: string = '';

    public formGroup!: FormGroup;
    public createSteps: number = 1;
    public currentCreateStep: number = 1;

    constructor(
        private router: Router,
        private route: ActivatedRoute,
        private toast: ToastService,
        private adminService: AdminService,
        private _location: Location,
    ) {
        this.formGroup = new FormGroup({
            name: new FormControl('', [Validators.required]),
            logoSrc: new FormControl('', []),
            clientId: new FormControl('', [Validators.required]),
            clientSecret: new FormControl('', [Validators.required]),
            issuer: new FormControl('', [Validators.required]),
            scopesList: new FormControl([], []),
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
        const req: OidcIdpConfigCreate = new OidcIdpConfigCreate();

        req.setName(this.name?.value);
        req.setClientId(this.clientId?.value);
        req.setClientSecret(this.clientSecret?.value);
        req.setIssuer(this.issuer?.value);
        req.setLogoSrc(this.logoSrc?.value);

        this.adminService.CreateOidcIdp(req).then((idp) => {
            this.router.navigate(['idp', idp.getId()]);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public close(): void {
        this._location.back();
    }

    public get name(): AbstractControl | null {
        return this.formGroup.get('name');
    }

    public get logoSrc(): AbstractControl | null {
        return this.formGroup.get('logoSrc');
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
}
