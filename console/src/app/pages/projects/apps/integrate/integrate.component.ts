import { Location } from '@angular/common';
import { Component, OnInit, effect, signal } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, Subject, combineLatest } from 'rxjs';
import { map } from 'rxjs/operators';
import { OIDCAppType } from 'src/app/proto/generated/zitadel/app_pb';
import { AddAPIAppResponse, AddOIDCAppRequest, AddOIDCAppResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AppSecretDialogComponent } from '../app-secret-dialog/app-secret-dialog.component';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { frameworks, OIDC_CONFIGURATIONS } from 'src/app/utils/framework';
import { NavigationService } from 'src/app/services/navigation.service';
import { NameDialogComponent } from 'src/app/modules/name-dialog/name-dialog.component';

@Component({
  selector: 'cnsl-integrate',
  templateUrl: './integrate.component.html',
  styleUrls: ['./integrate.component.scss'],
  standalone: false,
})
export class IntegrateAppComponent implements OnInit {
  private destroy$: Subject<void> = new Subject();
  public projectId: string = '';
  public loading: boolean = false;
  public InfoSectionType = InfoSectionType;
  protected readonly framework = signal<(typeof frameworks)[number] | undefined>(undefined);
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
      const fw = this.framework();
      if (fw) {
        const request = OIDC_CONFIGURATIONS[fw.id as unknown as keyof typeof OIDC_CONFIGURATIONS];
        // request.setProjectId(this.projectId);
        // request.setName(fw.title);
        // request.setDevMode(true);
        this.requestRedirectValuesSubject$.next();
        this.showRenameWarning.next(false);

        // this.oidcAppRequest.next(request);
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
          this.router.navigate(['projects', this.projectId, 'apps', resp.appId], { queryParams: { new: true } }).then();
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
