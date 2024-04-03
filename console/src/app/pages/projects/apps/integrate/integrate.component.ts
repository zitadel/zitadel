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
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { Framework } from 'src/app/components/quickstart/quickstart.component';
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';
import { NavigationService } from 'src/app/services/navigation.service';
import { NameDialogComponent } from 'src/app/modules/name-dialog/name-dialog.component';

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
  public showRenameWarning: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
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
        this.showRenameWarning.next(false);

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
        this.showRenameWarning.next(false);
        this.toast.showInfo('APP.TOAST.CREATED', true);
        if (resp.clientSecret) {
          this.showSavedDialog(resp);
        } else {
          this.router.navigate(['projects', this.projectId, 'apps', resp.appId], { queryParams: { new: true } });
        }
      })
      .catch((error) => {
        if (error.code === 6) {
          this.showRenameWarning.next(true);
        }
        this.loading = false;
        this.toast.showError(error);
      });
  }

  public editName() {
    const dialogRef = this.dialog.open(NameDialogComponent, {
      data: {
        name: this.oidcAppRequest.getValue()?.getName() ?? '',
        titleKey: 'APP.NAMEDIALOG.TITLE',
        descKey: 'APP.NAMEDIALOG.DESCRIPTION',
        labelKey: 'APP.NAMEDIALOG.NAME',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((name) => {
      if (name && name !== this.framework()?.title) {
        const request = this.oidcAppRequest.getValue();
        request.setName(name);
        this.showRenameWarning.next(false);
        this.oidcAppRequest.next(request);
      }
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
