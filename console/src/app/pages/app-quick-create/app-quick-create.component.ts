import { Component, OnDestroy, signal } from '@angular/core';
import { ActivatedRoute, Navigation, Params, Router } from '@angular/router';
import { BehaviorSubject, Subject, takeUntil } from 'rxjs';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { ProjectAutocompleteType } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.component';
import { AddProjectRequest, AddProjectResponse, AddOIDCAppRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { Framework } from 'src/app/components/quickstart/quickstart.component';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { NavigationService } from 'src/app/services/navigation.service';
import { Location } from '@angular/common';
import { ToastService } from 'src/app/services/toast.service';
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';

@Component({
  selector: 'cnsl-app-quick-create',
  templateUrl: './app-quick-create.component.html',
  styleUrls: ['./app-quick-create.component.scss'],
  standalone: false,
})
export class AppQuickCreateComponent implements OnDestroy {
  public InfoSectionType: any = InfoSectionType;
  public project?: {
    project: Project.AsObject | GrantedProject.AsObject;
    type: ProjectType;
    name: string;
  } = undefined;
  public ProjectAutocompleteType: any = ProjectAutocompleteType;
  public projectName: string = '';
  public projectId: string = '';

  public error = signal('');
  public framework = signal<Framework | undefined>(undefined);
  public customFramework = signal<boolean>(false);
  public initialParam = signal<string>('');
  public destroy$: Subject<void> = new Subject();

  public loading: boolean = false;
  public showRenameWarning: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  public frameworks: Framework[] = frameworkDefinition
    .map((f) => {
      return {
        ...f,
        fragment: '',
        imgSrcDark: `assets${f.imgSrcDark}`,
        imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
      };
    })
    .sort((a, b) => {
      // Define popularity order (most popular first)
      const popularityOrder = [
        'react',
        'next',
        'vue',
        'angular',
        'client-node',
        'flutter',
        'client-go',
        'spring',
        'django',
        'symfony',
        'client-python',
        'client-java',
        'client-php',
        'client-ruby',
      ];

      const aIndex = popularityOrder.indexOf(a.id?.toLowerCase() || '');
      const bIndex = popularityOrder.indexOf(b.id?.toLowerCase() || '');

      // If both are in the popularity list, sort by their position
      if (aIndex !== -1 && bIndex !== -1) {
        return aIndex - bIndex;
      }

      // If only one is in the list, prioritize it
      if (aIndex !== -1) return -1;
      if (bIndex !== -1) return 1;

      // If neither is in the list, sort alphabetically
      return a.title.localeCompare(b.title);
    });
  constructor(
    private router: Router,
    private mgmtService: ManagementService,
    private toast: ToastService,
    breadcrumbService: BreadcrumbService,
    activatedRoute: ActivatedRoute,
    private _location: Location,
    private navigation: NavigationService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);

    activatedRoute.queryParams.pipe(takeUntil(this.destroy$)).subscribe((params: Params) => {
      const { framework } = params;
      if (framework) {
        this.initialParam.set(framework);
      }
    });
  }

  public findFramework(id: string) {
    if (id !== 'other') {
      this.customFramework.set(false);
      const temp = this.frameworks.find((f) => f.id === id);
      this.framework.set(temp);
    } else {
      this.customFramework.set(true);
    }
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public goToAppIntegratePage(): void {
    if (this.project && this.customFramework()) {
      const id = (this.project.project as Project.AsObject).id
        ? (this.project.project as Project.AsObject).id
        : (this.project.project as GrantedProject.AsObject).projectId
          ? (this.project.project as GrantedProject.AsObject).projectId
          : '';
      this.router.navigate(['/projects', id, 'apps', 'create']);
    } else if (this.project && this.framework()) {
      const id = (this.project.project as Project.AsObject).id
        ? (this.project.project as Project.AsObject).id
        : (this.project.project as GrantedProject.AsObject).projectId
          ? (this.project.project as GrantedProject.AsObject).projectId
          : '';
      this.router.navigate(['/projects', id, 'apps', 'integrate'], { queryParams: { framework: this.framework()?.id } });
    }
  }

  public close(): void {
    if (this.navigation.isBackPossible) {
      this._location.back();
    } else {
      if (this.project && this.framework()) {
        const id = (this.project.project as Project.AsObject).id
          ? (this.project.project as Project.AsObject).id
          : (this.project.project as GrantedProject.AsObject).projectId
            ? (this.project.project as GrantedProject.AsObject).projectId
            : '';
        this.router.navigate(['/projects', id]);
      } else {
        this.router.navigate(['/projects']);
      }
    }
  }

  public selectProject(project: {
    project: Project.AsObject | GrantedProject.AsObject;
    type: ProjectType;
    name: string;
  }): void {
    if (project) {
      this.project = project;
    }
  }

  public createApp(projectId: string): void {
    this.loading = true;
    const request = new AddOIDCAppRequest();
    request.setProjectId(projectId);
    request.setName(`My ${this.framework()?.title} App` || 'My App');
    request.setDevMode(true);
    request.setAppType(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getAppType());
    request.setAuthMethodType(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getAuthMethodType());
    request.setResponseTypesList(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getResponseTypesList());
    request.setGrantTypesList(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getGrantTypesList());
    request.setRedirectUrisList(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getRedirectUrisList());
    request.setPostLogoutRedirectUrisList(
      OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getPostLogoutRedirectUrisList(),
    );

    this.mgmtService
      .addOIDCApp(request)
      .then((resp) => {
        this.loading = false;
        this.showRenameWarning.next(false);
        this.toast.showInfo('APP.TOAST.CREATED', true);

        this.router.navigate(['projects', projectId, 'apps', resp.appId], { queryParams: { new: true } });
      })
      .catch((error) => {
        if (error.code === 6) {
          this.showRenameWarning.next(true);
        }
        this.loading = false;
        this.toast.showError(error);
      });
  }

  public createProjectAndContinue(frameworkId?: string) {
    const project = new AddProjectRequest();
    project.setName('my-new-project-5');
    this.framework.set(this.frameworks.find((f) => f.id === frameworkId));

    return this.mgmtService
      .addProject(project.toObject())
      .then((resp: AddProjectResponse.AsObject) => {
        this.error.set('');
        if (frameworkId === 'other' || !frameworkId) {
          this.router.navigate(['/projects', resp.id, 'apps', 'create']);
        } else {
          this.projectId = resp.id;
          this.createApp(resp.id);
        }
      })
      .catch((error) => {
        const { message } = error;
        this.error.set(message);
      });
  }
}
