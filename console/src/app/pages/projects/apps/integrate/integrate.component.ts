import { COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Buffer } from 'buffer';
import { Subject, Subscription } from 'rxjs';
import { debounceTime, map, takeUntil } from 'rxjs/operators';
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
  DEVICE_CODE_METHOD,
  getPartialConfigFromAuthMethod,
  IMPLICIT_METHOD,
  PKCE_METHOD,
  PK_JWT_METHOD,
  POST_METHOD,
} from '../authmethods';
import { API_TYPE, AppCreateType, NATIVE_TYPE, RadioItemAppType, SAML_TYPE, USER_AGENT_TYPE, WEB_TYPE } from '../authtypes';
import { EnvironmentService } from 'src/app/services/environment.service';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { Framework } from 'src/app/components/quickstart/quickstart.component';

@Component({
  selector: 'cnsl-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
})
export class IntegrateAppComponent implements OnInit, OnDestroy {
  private destroy$: Subject<void> = new Subject();
  public projectId: string = '';
  public loading: boolean = false;
  public oidcAppRequest: AddOIDCAppRequest = new AddOIDCAppRequest();
  public InfoSectionType: any = InfoSectionType;
  public framework = signal<Framework | undefined>(undefined);

  constructor(
    private activatedRoute: ActivatedRoute,
    private router: Router,
    private toast: ToastService,
    private dialog: MatDialog,
    private mgmtService: ManagementService,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
  ) {}

  public ngOnInit(): void {
    const projectId = this.activatedRoute.snapshot.paramMap.get('projectid');
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
    this.destroy$.next();
    this.destroy$.complete();
  }

  public close(): void {
    this._location.back();
  }

  public createApp(): void {
    // this.requestRedirectValuesSubject$.next();

    this.loading = true;
    this.mgmtService
      .addOIDCApp(this.oidcAppRequest)
      .then((resp) => {
        this.loading = false;
        this.toast.showInfo('APP.TOAST.CREATED', true);
        if (resp.clientId || resp.clientSecret) {
          this.showSavedDialog(resp);
        } else {
          this.router.navigate(['projects', this.projectId, 'apps', resp.appId], { queryParams: { new: true } });
        }
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
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
      this.router.navigate(['projects', this.projectId, 'apps', added.appId], { queryParams: { new: true } });
    });
  }
}
