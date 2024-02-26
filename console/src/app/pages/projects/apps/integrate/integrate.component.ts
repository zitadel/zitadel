import { C, COMMA, ENTER, SPACE } from '@angular/cdk/keycodes';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit, Signal, computed, effect, signal } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Buffer } from 'buffer';
import { BehaviorSubject, Subject, Subscription, combineLatest } from 'rxjs';
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
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';
import { NavigationService } from 'src/app/services/navigation.service';

@Component({
  selector: 'cnsl-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
})
export class IntegrateAppComponent implements OnInit, OnDestroy {
  private destroy$: Subject<void> = new Subject();
  public projectId: string = '';
  public loading: boolean = false;
  public InfoSectionType: any = InfoSectionType;
  public framework = signal<Framework | undefined>(undefined);
  public oidcAppRequest: BehaviorSubject<AddOIDCAppRequest> = new BehaviorSubject(new AddOIDCAppRequest());

  public OIDCAppType: any = OIDCAppType;
  public requestRedirectValuesSubject$: Subject<void> = new Subject();

  constructor(
    private activatedRoute: ActivatedRoute,
    private router: Router,
    private toast: ToastService,
    private dialog: MatDialog,
    private mgmtService: ManagementService,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
    public navigation: NavigationService,
  ) {
    effect(() => {
      const fwId = this.framework()?.id;
      const fw = this.framework();
      if (fw && fwId) {
        const request = OIDC_CONFIGURATIONS[fwId];
        request.setProjectId(this.projectId);
        request.setName(fw.title);
        request.setDevMode(true);
        this.requestRedirectValuesSubject$.next();

        console.log(request.toObject());
        this.oidcAppRequest.next(request);
        return request;
      } else {
        const request = new AddOIDCAppRequest();
        this.oidcAppRequest.next(request);
        return request;
      }
    });
  }

  public projectName$ = combineLatest([this.mgmtService.ownedProjects, this.mgmtService.grantedProjects]).pipe(
    map(([projects, grantedProjects]) => {
      const project = projects.find((project) => project.id === this.activatedRoute.snapshot.paramMap.get('projectid'));

      const grantedproject = grantedProjects.find(
        (grantedproject) => grantedproject.projectId === this.activatedRoute.snapshot.paramMap.get('projectid'),
      );

      return project?.name ?? grantedproject?.projectName ?? '';
    }),
  );

  public setFramework(framework: Framework | undefined) {
    this.framework.set(framework);
  }

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
      this.projectId = projectId;
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    }
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public close(): void {
    if (this.navigation.isBackPossible) {
      this._location.back();
    } else {
      this.router.navigate(['/projects', this.projectId]);
    }
  }

  public createApp(): void {
    this.loading = true;
    this.mgmtService
      .addOIDCApp(this.oidcAppRequest.getValue())
      .then((resp) => {
        this.loading = false;
        this.toast.showInfo('APP.TOAST.CREATED', true);
        if (resp.clientSecret) {
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

  public get redirectUris() {
    return this.oidcAppRequest.getValue().toObject().redirectUrisList;
  }

  public set redirectUris(value: string[]) {
    const request = this.oidcAppRequest.getValue();
    request.setRedirectUrisList(value);
    this.oidcAppRequest.next(request);
  }

  public get postLogoutUrisList() {
    return this.oidcAppRequest.getValue().toObject().postLogoutRedirectUrisList;
  }

  public set postLogoutUrisList(value: string[]) {
    const request = this.oidcAppRequest.getValue();
    request.setPostLogoutRedirectUrisList(value);
    this.oidcAppRequest.next(request);
  }
}
