import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { Location } from '@angular/common';
import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatChipInputEvent } from '@angular/material/chips';
import { ActivatedRoute } from '@angular/router';
import { Observable, Subject } from 'rxjs';
import { switchMap, take, takeUntil } from 'rxjs/operators';
import { UpdateIDPOIDCConfigRequest, UpdateIDPRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { IDPStylingType, OIDCMappingField } from 'src/app/proto/generated/zitadel/idp_pb';
import { UpdateOrgIDPOIDCConfigRequest, UpdateOrgIDPRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

@Component({
  selector: 'app-idp',
  templateUrl: './idp.component.html',
  styleUrls: ['./idp.component.scss'],
})
export class IdpComponent implements OnDestroy {
  public mappingFields: OIDCMappingField[] = [];
  public styleFields: IDPStylingType[] = [];

  public showIdSecretSection: boolean = false;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  private service!: ManagementService | AdminService;
  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];

  private destroy$: Subject<void> = new Subject();
  public projectId: string = '';

  public idpForm!: FormGroup;
  public oidcConfigForm!: FormGroup;

  public canWrite: Observable<boolean> = this.authService.isAllowed([this.serviceType === PolicyComponentServiceType.ADMIN ?
    'iam.idp.write' : this.serviceType === PolicyComponentServiceType.MGMT ?
      'org.idp.write' : '']);

  constructor(
    private toast: ToastService,
    private injector: Injector,
    private route: ActivatedRoute,
    private _location: Location,
    private authService: GrpcAuthService,
  ) {
    this.idpForm = new FormGroup({
      id: new FormControl({ disabled: true, value: '' }, [Validators.required]),
      name: new FormControl('', [Validators.required]),
      stylingType: new FormControl('', [Validators.required]),
      autoRegister: new FormControl(false, [Validators.required]),
    });

    this.oidcConfigForm = new FormGroup({
      clientId: new FormControl('', [Validators.required]),
      clientSecret: new FormControl(''),
      issuer: new FormControl('', [Validators.required]),
      scopesList: new FormControl([], []),
      idpDisplayNameMapping: new FormControl(0),
      usernameMapping: new FormControl(0),
    });

    this.route.data.pipe(
      takeUntil(this.destroy$),
      switchMap(data => {
        this.serviceType = data.serviceType;
        switch (this.serviceType) {
          case PolicyComponentServiceType.MGMT:
            this.service = this.injector.get(ManagementService as Type<ManagementService>);

            break;
          case PolicyComponentServiceType.ADMIN:
            this.service = this.injector.get(AdminService as Type<AdminService>);

            break;
        }

        this.mappingFields = [
          OIDCMappingField.OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
          OIDCMappingField.OIDC_MAPPING_FIELD_EMAIL];
        this.styleFields = [
          IDPStylingType.STYLING_TYPE_UNSPECIFIED,
          IDPStylingType.STYLING_TYPE_GOOGLE];

        return this.route.params.pipe(take(1));
      })).subscribe((params) => {
        const { id } = params;
        if (id) {
          this.checkWrite();

          if (this.serviceType === PolicyComponentServiceType.MGMT) {

            (this.service as ManagementService).getOrgIDPByID(id).then(resp => {
              if (resp.idp) {
                const idpObject = resp.idp;
                this.idpForm.patchValue(idpObject);
                if (idpObject.oidcConfig) {
                  this.oidcConfigForm.patchValue(idpObject.oidcConfig);
                }
              }
            });
          } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
            (this.service as AdminService).getIDPByID(id).then(resp => {
              if (resp.idp) {
                const idpObject = resp.idp;
                this.idpForm.patchValue(idpObject);
                if (idpObject.oidcConfig) {
                  this.oidcConfigForm.patchValue(idpObject.oidcConfig);
                }
              }
            });
          }
        }
      });
  }

  public checkWrite(): void {
    this.canWrite.pipe(take(1)).subscribe(canWrite => {
      if (canWrite) {
        this.idpForm.enable();
        this.oidcConfigForm.enable();
      } else {
        this.idpForm.disable();
        this.oidcConfigForm.disable();
      }
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public updateIdp(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new UpdateOrgIDPRequest();

      req.setIdpId(this.id?.value);
      req.setName(this.name?.value);
      req.setStylingType(this.stylingType?.value);
      req.setAutoRegister(this.autoRegister?.value);

      (this.service as ManagementService).updateOrgIDP(req).then(() => {
        this.toast.showInfo('IDP.TOAST.SAVED', true);
        // this.router.navigate(['idp', ]);
      }).catch(error => {
        this.toast.showError(error);
      });
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      const req = new UpdateIDPRequest();

      req.setIdpId(this.id?.value);
      req.setName(this.name?.value);
      req.setStylingType(this.stylingType?.value);
      req.setAutoRegister(this.autoRegister?.value);

      (this.service as AdminService).updateIDP(req).then(() => {
        this.toast.showInfo('IDP.TOAST.SAVED', true);
        // this.router.navigate(['idp', ]);
      }).catch(error => {
        this.toast.showError(error);
      });
    }
  }

  public updateOidcConfig(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const req = new UpdateOrgIDPOIDCConfigRequest();

      req.setIdpId(this.id?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);
      req.setUsernameMapping(this.usernameMapping?.value);
      req.setDisplayNameMapping(this.idpDisplayNameMapping?.value);

      (this.service as ManagementService).updateOrgIDPOIDCConfig(req).then((oidcConfig) => {
        this.toast.showInfo('IDP.TOAST.SAVED', true);
        // this.router.navigate(['idp', ]);
      }).catch(error => {
        this.toast.showError(error);
      });
    } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      const req = new UpdateIDPOIDCConfigRequest();

      req.setIdpId(this.id?.value);
      req.setClientId(this.clientId?.value);
      req.setClientSecret(this.clientSecret?.value);
      req.setIssuer(this.issuer?.value);
      req.setScopesList(this.scopesList?.value);
      req.setUsernameMapping(this.usernameMapping?.value);
      req.setDisplayNameMapping(this.idpDisplayNameMapping?.value);

      (this.service as AdminService).updateIDPOIDCConfig(req).then((oidcConfig) => {
        this.toast.showInfo('IDP.TOAST.SAVED', true);
        // this.router.navigate(['idp', ]);
      }).catch(error => {
        this.toast.showError(error);
      });
    }
  }

  public close(): void {
    this._location.back();
  }

  public addScope(event: MatChipInputEvent): void {
    const input = event.chipInput?.inputElement;
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

  public get autoRegister(): AbstractControl | null {
    return this.idpForm.get('autoRegister');
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
