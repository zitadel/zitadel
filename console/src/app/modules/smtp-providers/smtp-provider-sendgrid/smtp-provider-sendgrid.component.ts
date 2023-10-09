import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import { ActivatedRoute } from '@angular/router';
import { Subject, take } from 'rxjs';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import {
  AddGoogleProviderRequest as AdminAddGoogleProviderRequest,
  GetProviderByIDRequest as AdminGetProviderByIDRequest,
  UpdateGoogleProviderRequest as AdminUpdateGoogleProviderRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { Options, Provider } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  AddGoogleProviderRequest as MgmtAddGoogleProviderRequest,
  GetProviderByIDRequest as MgmtGetProviderByIDRequest,
  UpdateGoogleProviderRequest as MgmtUpdateGoogleProviderRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { MatLegacyCheckboxChange } from '@angular/material/legacy-checkbox';

@Component({
  selector: 'cnsl-provider-sendgrid',
  templateUrl: './smtp-provider-sendgrid.component.html',
  styleUrls: ['./smtp-provider-sendgrid.component.scss'],
})
export class SMTPProviderSendgridComponent {
  public showOptional: boolean = false;
  public options: Options = new Options().setIsCreationAllowed(true).setIsLinkingAllowed(true);
  public id: string | null = '';
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  // private service!: ManagementService | AdminService;

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  public loading: boolean = false;

  public provider?: Provider.AsObject;
  public updateClientSecret: boolean = false;

  // stepper
  public currentCreateStep: number = 1;
  public requestRedirectValuesSubject$: Subject<void> = new Subject();
  public firstFormGroup!: UntypedFormGroup;
  public secondFormGroup!: UntypedFormGroup;

  private host: string = 'smtp.sendgrid.net';
  private unencryptedPort: number = 587;
  private encryptedPort: number = 465;
  private tls: boolean = false;

  constructor(
    private authService: GrpcAuthService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
    private fb: UntypedFormBuilder,
  ) {
    this.firstFormGroup = this.fb.group({
      tls: [false],
      hostAndPort: [`${this.host}:${this.unencryptedPort}`],
      user: ['apiKey'],
      password: [''],
    });

    this.secondFormGroup = this.fb.group({
      senderAddress: ['', [requiredValidator]],
      senderName: ['', [requiredValidator]],
      replyToAddress: [''],
    });

    // this.authService
    //   .isAllowed(
    //     this.serviceType === PolicyComponentServiceType.ADMIN
    //       ? ['iam.idp.write']
    //       : this.serviceType === PolicyComponentServiceType.MGMT
    //       ? ['org.idp.write']
    //       : [],
    //   )
    //   .pipe(take(1))
    //   .subscribe((allowed) => {
    //     if (allowed) {
    //       this.form.enable();
    //     } else {
    //       this.form.disable();
    //     }
    //   });

    // this.route.data.pipe(take(1)).subscribe((data) => {
    //   this.serviceType = data['serviceType'];

    //   switch (this.serviceType) {
    //     case PolicyComponentServiceType.MGMT:
    //       this.service = this.injector.get(ManagementService as Type<ManagementService>);

    //       const bread: Breadcrumb = {
    //         type: BreadcrumbType.ORG,
    //         routerLink: ['/org'],
    //       };

    //       this.breadcrumbService.setBreadcrumb([bread]);
    //       break;
    //     case PolicyComponentServiceType.ADMIN:
    //       this.service = this.injector.get(AdminService as Type<AdminService>);

    //       const iamBread = new Breadcrumb({
    //         type: BreadcrumbType.ORG,
    //         name: 'Instance',
    //         routerLink: ['/instance'],
    //       });
    //       this.breadcrumbService.setBreadcrumb([iamBread]);
    //       break;
    //   }

    //   this.id = this.route.snapshot.paramMap.get('id');
    //   if (this.id) {
    //     this.clientSecret?.setValidators([]);
    //     this.getData(this.id);
    //   }
    // });
  }

  public changeStep(event: StepperSelectionEvent): void {
    this.currentCreateStep = event.selectedIndex + 1;

    if (event.selectedIndex >= 2) {
      this.requestRedirectValuesSubject$.next();
    }
  }

  public toggleTLS(event: MatLegacyCheckboxChange) {
    this.hostAndPort?.setValue(`${this.host}:${event.checked ? this.encryptedPort : this.unencryptedPort}`);
  }

  public get hostAndPort(): AbstractControl | null {
    return this.firstFormGroup.get('hostAndPort');
  }

  // private getData(id: string): void {
  //   const req =
  //     this.serviceType === PolicyComponentServiceType.ADMIN
  //       ? new AdminGetProviderByIDRequest()
  //       : new MgmtGetProviderByIDRequest();
  //   req.setId(id);
  //   this.service
  //     .getProviderByID(req)
  //     .then((resp) => {
  //       this.provider = resp.idp;
  //       this.loading = false;
  //       if (this.provider?.config?.google) {
  //         this.form.patchValue(this.provider.config.google);
  //         this.name?.setValue(this.provider.name);
  //       }
  //     })
  //     .catch((error) => {
  //       this.toast.showError(error);
  //       this.loading = false;
  //     });
  // }

  // public submitForm(): void {
  //   this.provider ? this.updateGoogleProvider() : this.addGoogleProvider();
  // }

  // public addGoogleProvider(): void {
  //   const req =
  //     this.serviceType === PolicyComponentServiceType.MGMT
  //       ? new MgmtAddGoogleProviderRequest()
  //       : new AdminAddGoogleProviderRequest();

  //   req.setName(this.name?.value);
  //   req.setClientId(this.clientId?.value);
  //   req.setClientSecret(this.clientSecret?.value);
  //   req.setScopesList(this.scopesList?.value);
  //   req.setProviderOptions(this.options);

  //   this.loading = true;
  //   this.service
  //     .addGoogleProvider(req)
  //     .then((idp) => {
  //       setTimeout(() => {
  //         this.loading = false;
  //         this.close();
  //       }, 2000);
  //     })
  //     .catch((error) => {
  //       this.toast.showError(error);
  //       this.loading = false;
  //     });
  // }

  // public updateGoogleProvider(): void {
  //   if (this.provider) {
  //     if (this.serviceType === PolicyComponentServiceType.MGMT) {
  //       const req = new MgmtUpdateGoogleProviderRequest();
  //       req.setId(this.provider.id);
  //       req.setName(this.name?.value);
  //       req.setClientId(this.clientId?.value);
  //       req.setScopesList(this.scopesList?.value);
  //       req.setProviderOptions(this.options);

  //       if (this.updateClientSecret) {
  //         req.setClientSecret(this.clientSecret?.value);
  //       }

  //       this.loading = true;
  //       (this.service as ManagementService)
  //         .updateGoogleProvider(req)
  //         .then((idp) => {
  //           setTimeout(() => {
  //             this.loading = false;
  //             this.close();
  //           }, 2000);
  //         })
  //         .catch((error) => {
  //           this.toast.showError(error);
  //           this.loading = false;
  //         });
  //     } else if (PolicyComponentServiceType.ADMIN) {
  //       const req = new AdminUpdateGoogleProviderRequest();
  //       req.setId(this.provider.id);
  //       req.setName(this.name?.value);
  //       req.setClientId(this.clientId?.value);
  //       req.setScopesList(this.scopesList?.value);
  //       req.setProviderOptions(this.options);

  //       if (this.updateClientSecret) {
  //         req.setClientSecret(this.clientSecret?.value);
  //       }

  //       this.loading = true;
  //       (this.service as AdminService)
  //         .updateGoogleProvider(req)
  //         .then((idp) => {
  //           setTimeout(() => {
  //             this.loading = false;
  //             this.close();
  //           }, 2000);
  //         })
  //         .catch((error) => {
  //           this.loading = false;
  //           this.toast.showError(error);
  //         });
  //     }
  //   }
  // }

  public close(): void {
    this._location.back();
  }
}
