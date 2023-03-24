import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Buffer } from 'buffer';
import { Subject, Subscription } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
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
  AddSAMLAppRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import {
  BASIC_AUTH_METHOD,
  CODE_METHOD,
  getPartialConfigFromAuthMethod,
  IMPLICIT_METHOD,
  PKCE_METHOD,
  PK_JWT_METHOD,
  POST_METHOD,
} from '../authmethods';
import { API_TYPE, AppCreateType, NATIVE_TYPE, RadioItemAppType, SAML_TYPE, USER_AGENT_TYPE, WEB_TYPE } from '../authtypes';

const MAX_ALLOWED_SIZE = 1 * 1024 * 1024;

@Component({
  selector: 'cnsl-app-create',
  templateUrl: './app-create.component.html',
  styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent implements OnInit, OnDestroy {
  private subscription: Subscription = new Subscription();
  private destroyed$: Subject<void> = new Subject();
  public devmode: boolean = false;
  public projectId: string = '';
  public loading: boolean = false;

  public currentCreateStep: number = 1;

  public oidcAppRequest: AddOIDCAppRequest = new AddOIDCAppRequest();
  public apiAppRequest: AddAPIAppRequest = new AddAPIAppRequest();
  public samlAppRequest: AddSAMLAppRequest = new AddSAMLAppRequest();

  public oidcResponseTypes: { type: OIDCResponseType; checked: boolean; disabled: boolean }[] = [
    { type: OIDCResponseType.OIDC_RESPONSE_TYPE_CODE, checked: false, disabled: false },
    { type: OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN, checked: false, disabled: false },
    { type: OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN, checked: false, disabled: false },
  ];

  public oidcAppTypes: OIDCAppType[] = [
    OIDCAppType.OIDC_APP_TYPE_WEB,
    OIDCAppType.OIDC_APP_TYPE_NATIVE,
    OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
  ];
  public appTypes: any = [WEB_TYPE, NATIVE_TYPE, USER_AGENT_TYPE, API_TYPE, SAML_TYPE];

  public authMethods: RadioItemAuthType[] = [PKCE_METHOD, CODE_METHOD, PK_JWT_METHOD, POST_METHOD];

  // set to oidc first
  public authMethodTypes: {
    type: OIDCAuthMethodType | APIAuthMethodType;
    checked: boolean;
    disabled: boolean;
    api?: boolean;
    oidc?: boolean;
  }[] = [
    { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC, checked: false, disabled: false, oidc: true },
    { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE, checked: false, disabled: false, oidc: true },
    { type: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST, checked: false, disabled: false, oidc: true },
  ];

  // stepper
  public firstFormGroup!: UntypedFormGroup;
  public secondFormGroup!: UntypedFormGroup;
  public samlConfigForm!: UntypedFormGroup;

  public redirectUrisList: string[] = [];
  public postLogoutRedirectUrisList: string[] = [];

  // devmode
  public form!: UntypedFormGroup;

  public AppCreateType: any = AppCreateType;
  public OIDCAppType: any = OIDCAppType;
  public OIDCGrantType: any = OIDCGrantType;
  public OIDCAuthMethodType: any = OIDCAuthMethodType;

  public oidcGrantTypes: {
    type: OIDCGrantType;
    checked: boolean;
    disabled: boolean;
  }[] = [
    { type: OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE, checked: true, disabled: false },
    { type: OIDCGrantType.OIDC_GRANT_TYPE_IMPLICIT, checked: false, disabled: true },
    { type: OIDCGrantType.OIDC_GRANT_TYPE_REFRESH_TOKEN, checked: false, disabled: true },
  ];

  public readonly separatorKeysCodes: number[] = [ENTER, COMMA, SPACE];
  public requestRedirectValuesSubject$: Subject<void> = new Subject();

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private toast: ToastService,
    private dialog: MatDialog,
    private mgmtService: ManagementService,
    private fb: UntypedFormBuilder,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
  ) {
    this.form = this.fb.group({
      name: ['', [requiredValidator]],
      appType: ['', [requiredValidator]],
      // apptype OIDC
      responseTypesList: ['', []],
      grantTypesList: ['', []],
      authMethodType: ['', []],
      // apptype SAML
      metadataUrl: ['', []],
    });

    this.initForm();

    this.firstFormGroup = this.fb.group({
      name: ['', [requiredValidator]],
      appType: [WEB_TYPE, [requiredValidator]],
    });

    this.samlConfigForm = this.fb.group({
      metadataUrl: ['', []],
    });

    this.firstFormGroup.valueChanges.subscribe((value) => {
      if (this.firstFormGroup.valid) {
        this.oidcAppRequest.setName(this.name?.value);
        this.apiAppRequest.setName(this.name?.value);
        this.samlAppRequest.setName(this.name?.value);

        if (this.isStepperOIDC) {
          const oidcAppType = (this.appType?.value as RadioItemAppType).oidcAppType;
          if (oidcAppType !== undefined) {
            this.oidcAppRequest.setAppType(oidcAppType);
          }

          switch (this.appType?.value.oidcAppType) {
            case OIDCAppType.OIDC_APP_TYPE_NATIVE:
              this.authMethods = [PKCE_METHOD];

              // automatically set to PKCE and skip step
              this.oidcAppRequest.setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE]);
              this.oidcAppRequest.setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE]);
              this.oidcAppRequest.setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE);

              break;
            case OIDCAppType.OIDC_APP_TYPE_WEB:
              // PK_JWT_METHOD.recommended = false;
              this.authMethods = [PKCE_METHOD, CODE_METHOD, PK_JWT_METHOD, POST_METHOD];

              this.authMethod?.setValue(PKCE_METHOD.key);
              break;
            case OIDCAppType.OIDC_APP_TYPE_USER_AGENT:
              this.authMethods = [PKCE_METHOD, IMPLICIT_METHOD];

              this.authMethod?.setValue(PKCE_METHOD.key);
              break;
          }
        } else if (this.isStepperAPI) {
          // PK_JWT_METHOD.recommended = true;
          this.authMethods = [PK_JWT_METHOD, BASIC_AUTH_METHOD];

          this.authMethod?.setValue(PK_JWT_METHOD.key);
        }
      }
    });

    this.secondFormGroup = this.fb.group({
      authMethod: [this.authMethods[0].key, [requiredValidator]],
    });

    this.secondFormGroup.valueChanges.subscribe((form) => {
      const partialConfig = getPartialConfigFromAuthMethod(form.authMethod);

      if (this.isStepperOIDC && partialConfig && partialConfig.oidc) {
        this.oidcAppRequest.setResponseTypesList(partialConfig.oidc?.responseTypesList ?? []);

        this.oidcAppRequest.setGrantTypesList(partialConfig.oidc?.grantTypesList ?? []);

        this.oidcAppRequest.setAuthMethodType(
          partialConfig.oidc?.authMethodType ?? OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
        );
      } else if (this.isStepperAPI && partialConfig && partialConfig.api) {
        this.apiAppRequest.setAuthMethodType(
          partialConfig.api?.authMethodType ?? APIAuthMethodType.API_AUTH_METHOD_TYPE_BASIC,
        );
      }
    });

    this.samlConfigForm.valueChanges.subscribe((form) => {
      if (form.metadataUrl && form.metadataUrl.length > 0) {
        this.samlAppRequest.setMetadataUrl(form.metadataUrl);
      }
    });
  }

  public ngOnInit(): void {
    this.subscription = this.route.params.subscribe((params) => this.getData(params));

    const projectId = this.route.snapshot.paramMap.get('projectid');
    if (projectId) {
      const breadcrumbs = [
        new Breadcrumb({
          type: BreadcrumbType.ORG,
          routerLink: ['/org'],
        }),
        new Breadcrumb({
          type: BreadcrumbType.PROJECT,
          name: '',
          param: { key: 'projectid', value: projectId },
          routerLink: ['/projects', projectId],
          isZitadel: false,
        }),
      ];
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    }
  }

  public ngOnDestroy(): void {
    this.subscription?.unsubscribe();
    this.destroyed$.next();
  }

  public initForm(): void {
    this.form.valueChanges.pipe(takeUntil(this.destroyed$), debounceTime(150)).subscribe(() => {
      this.oidcAppRequest.setName(this.formname?.value);
      this.apiAppRequest.setName(this.formname?.value);
      this.samlAppRequest.setName(this.formname?.value);

      this.oidcAppRequest.setResponseTypesList(this.formresponseTypesList?.value);
      this.oidcAppRequest.setGrantTypesList(this.grantTypesList?.value);

      this.oidcAppRequest.setAuthMethodType(this.authMethodType?.value);
      this.apiAppRequest.setAuthMethodType(this.authMethodType?.value);

      if (this.formMetadataUrl?.value) {
        this.samlAppRequest.setMetadataUrl(this.formMetadataUrl?.value);
      }

      const oidcAppType = (this.formappType?.value as RadioItemAppType).oidcAppType;
      if (oidcAppType !== undefined) {
        this.oidcAppRequest.setAppType(oidcAppType);
      }
    });

    this.formappType?.valueChanges.pipe(takeUntil(this.destroyed$)).subscribe(() => {
      this.setDevFormValidators();
    });
  }

  public setDevFormValidators(): void {
    if (this.isDevOIDC) {
      const grantTypesControl = new UntypedFormControl('', [requiredValidator]);
      const responseTypesControl = new UntypedFormControl('', [requiredValidator]);

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

  public changeStep(event: StepperSelectionEvent): void {
    this.currentCreateStep = event.selectedIndex + 1;

    if (event.selectedIndex >= 2) {
      this.requestRedirectValuesSubject$.next();
    }
  }

  private async getData({ projectid }: Params): Promise<void> {
    this.projectId = projectid;
    this.oidcAppRequest.setProjectId(projectid);
    this.apiAppRequest.setProjectId(projectid);
    this.samlAppRequest.setProjectId(projectid);
  }

  public close(): void {
    this._location.back();
  }

  public onDropXML(filelist: FileList): void {
    const file = filelist.item(0);
    this.metadataUrl?.setValue('');
    if (file) {
      if (file.size > MAX_ALLOWED_SIZE) {
        this.toast.showInfo('POLICY.PRIVATELABELING.MAXSIZEEXCEEDED', true);
      } else {
        const reader = new FileReader();
        reader.onload = ((aXML) => {
          return (e) => {
            const xmlBase64 = e.target?.result;
            if (xmlBase64 && typeof xmlBase64 === 'string') {
              const cropped = xmlBase64.replace('data:text/xml;base64,', '');
              this.samlAppRequest.setMetadataXml(cropped);
            }
          };
        })(file);
        reader.readAsDataURL(file);
      }
    }
  }

  public createApp(): void {
    const appOIDCCheck = this.devmode ? this.isDevOIDC : this.isStepperOIDC;
    const appAPICheck = this.devmode ? this.isDevAPI : this.isStepperAPI;
    const appSAMLCheck = this.devmode ? this.isDevSAML : this.isStepperSAML;

    if (appOIDCCheck) {
      this.requestRedirectValuesSubject$.next();

      this.loading = true;
      this.mgmtService
        .addOIDCApp(this.oidcAppRequest)
        .then((resp) => {
          this.loading = false;
          this.toast.showInfo('APP.TOAST.CREATED', true);
          if (resp.clientId || resp.clientSecret) {
            this.showSavedDialog(resp);
          } else {
            this.router.navigate(['projects', this.projectId, 'apps', resp.appId]);
          }
        })
        .catch((error) => {
          this.loading = false;
          this.toast.showError(error);
        });
    } else if (appAPICheck) {
      this.loading = true;
      this.toast.showInfo('APP.TOAST.CREATED', true);
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
        .catch((error) => {
          this.loading = false;
          this.toast.showError(error);
        });
    } else if (appSAMLCheck) {
      this.loading = true;
      this.toast.showInfo('APP.TOAST.CREATED', true);
      this.mgmtService
        .addSAMLApp(this.samlAppRequest)
        .then((resp) => {
          this.loading = false;
          this.router.navigate(['projects', this.projectId, 'apps', resp.appId]);
        })
        .catch((error) => {
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
        clientId: clientId,
      },
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
    return (
      this.oidcGrantTypes
        .filter((gt) => gt.checked)
        .map((gt) => gt.type)
        .findIndex((t) => t === type) > -1
    );
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

  get grantTypesList(): AbstractControl | null {
    return this.form.get('grantTypesList');
  }

  get formappType(): AbstractControl | null {
    return this.form.get('appType');
  }

  get formMetadataUrl(): AbstractControl | null {
    return this.form.get('metadataUrl');
  }
  // get formapplicationType(): AbstractControl | null {
  //     return this.form.get('applicationType');
  // }

  get authMethodType(): AbstractControl | null {
    return this.form.get('authMethodType');
  }

  get isDevOIDC(): boolean {
    return (this.formappType?.value as RadioItemAppType).createType === AppCreateType.OIDC;
  }

  get isStepperOIDC(): boolean {
    return (this.appType?.value as RadioItemAppType).createType === AppCreateType.OIDC;
  }

  get isDevAPI(): boolean {
    return (this.formappType?.value as RadioItemAppType).createType === AppCreateType.API;
  }

  get isDevSAML(): boolean {
    return (this.formappType?.value as RadioItemAppType).createType === AppCreateType.SAML;
  }

  get isStepperAPI(): boolean {
    return (this.appType?.value as RadioItemAppType).createType === AppCreateType.API;
  }

  get isStepperSAML(): boolean {
    return (this.appType?.value as RadioItemAppType).createType === AppCreateType.SAML;
  }

  get decodedBase64(): string {
    const samlReq = this.samlAppRequest.toObject();
    if (samlReq && samlReq.metadataXml && typeof samlReq.metadataXml === 'string') {
      return Buffer.from(samlReq.metadataXml, 'base64').toString('ascii');
    } else {
      return '';
    }
  }

  set decodedBase64(xmlString) {
    if (this.samlAppRequest) {
      const base64 = Buffer.from(xmlString, 'ascii').toString('base64');
      this.samlAppRequest.setMetadataXml(base64);
    }
  }

  public get metadataUrl(): AbstractControl | null {
    return this.samlConfigForm.get('metadataUrl');
  }
}
