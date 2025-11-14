import { Component, OnDestroy, signal } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
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
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { OIDC_CONFIGURATIONS } from 'src/app/utils/framework';
import { generateRandomProjectName, generateFrameworkAppName } from 'src/app/utils/name-generator';
import { ThemeService } from 'src/app/services/theme.service';
import { OIDCAppType, OIDCAuthMethodType } from 'src/app/proto/generated/zitadel/app_pb';
import { MatDialog } from '@angular/material/dialog';
import { AppSecretDialogComponent } from 'src/app/pages/projects/apps/app-secret-dialog/app-secret-dialog.component';

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
  public selectedFramework = signal<Framework | undefined>(undefined);
  public customFramework = signal<boolean>(false);
  public initialParam = signal<string>('');
  public createProject: boolean = true; // Auto-create project by default (regular property for ngModel)
  public destroy$: Subject<void> = new Subject();

  public loading: boolean = false;
  public showRenameWarning: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  public frameworks: Framework[] = frameworkDefinition
    .filter((f) => !f.client) // Filter out client libraries/SDKs
    .map((f) => {
      return {
        ...f,
        fragment: '',
        imgSrcDark: `assets${f.imgSrcDark}`,
        imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
      };
    })
    .sort((a, b) => {
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
    private authService: GrpcAuthService,
    public themeService: ThemeService,
    private dialog: MatDialog,
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
        // Auto-select the framework if provided via query param
        const matchingFramework = this.frameworks.find((f) => f.id === framework);
        if (matchingFramework) {
          this.selectedFramework.set(matchingFramework);
          this.framework.set(matchingFramework);
        }
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

  public selectFrameworkForPreview(framework: Framework) {
    this.selectedFramework.set(framework);
  }

  public getAppTypeLabel(frameworkId?: string): string {
    if (!frameworkId) return '';
    const config = OIDC_CONFIGURATIONS[frameworkId];
    if (!config) return '';

    const appType = config.getAppType();
    switch (appType) {
      case OIDCAppType.OIDC_APP_TYPE_WEB:
        return 'Web Application';
      case OIDCAppType.OIDC_APP_TYPE_USER_AGENT:
        return 'User Agent (SPA)';
      case OIDCAppType.OIDC_APP_TYPE_NATIVE:
        return 'Native Application';
      default:
        return '';
    }
  }

  public getAuthMethod(frameworkId?: string): string {
    if (!frameworkId) return '';
    const config = OIDC_CONFIGURATIONS[frameworkId];
    if (!config) return '';

    const authMethod = config.getAuthMethodType();

    if (authMethod === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE) {
      return 'OIDC with PKCE';
    } else if (authMethod === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC) {
      return 'OIDC with Client Secret';
    }
    return 'OIDC';
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
    request.setName(generateFrameworkAppName(this.framework()?.title));
    request.setDevMode(true);

    const authMethodType = OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getAuthMethodType();
    request.setAppType(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getAppType());
    request.setAuthMethodType(authMethodType);
    request.setResponseTypesList(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getResponseTypesList());
    request.setGrantTypesList(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getGrantTypesList());
    request.setRedirectUrisList(OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getRedirectUrisList());
    request.setPostLogoutRedirectUrisList(
      OIDC_CONFIGURATIONS[this.framework()?.id || 'other']?.getPostLogoutRedirectUrisList(),
    );

    const usesClientSecret = authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC;

    this.mgmtService
      .addOIDCApp(request)
      .then((resp) => {
        this.loading = false;
        this.showRenameWarning.next(false);
        this.toast.showInfo('APP.TOAST.CREATED', true);

        if (usesClientSecret && resp.clientSecret) {
          this.dialog.open(AppSecretDialogComponent, {
            data: {
              clientSecret: resp.clientSecret,
            },
            width: '800px',
          });
        }

        this.router.navigate(['projects', projectId, 'apps', resp.appId], {
          queryParams: { new: true, framework: this.framework()?.id },
        });
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
    this.framework.set(this.frameworks.find((f) => f.id === frameworkId));

    // Determine if we're creating a new project or using existing
    let shouldCreateProject = this.createProject;
    let projectNameForCreation = '';

    // If checkbox is unchecked, check if they selected an existing project or typed a new name
    if (!this.createProject) {
      if (this.project) {
        // They selected an existing project
        const projectId = (this.project.project as Project.AsObject).id
          ? (this.project.project as Project.AsObject).id
          : (this.project.project as GrantedProject.AsObject).projectId;

        if (frameworkId === 'other' || !frameworkId) {
          this.router.navigate(['/projects', projectId, 'apps', 'create']);
        } else {
          this.createApp(projectId);
        }
        return;
      } else if (this.projectName && this.projectName.trim().length > 0) {
        // They typed a new project name - create project with that name
        shouldCreateProject = true;
        projectNameForCreation = this.projectName.trim();
      } else {
        // No project selected and no name typed
        this.error.set('Please select an existing project or enter a new project name');
        return;
      }
    }

    // Create new project (either because checkbox is checked or they typed a new name)
    const project = new AddProjectRequest();
    project.setName(projectNameForCreation || generateRandomProjectName());

    return this.mgmtService
      .addProject(project.toObject())
      .then((resp: AddProjectResponse.AsObject) => {
        this.error.set('');
        this.projectId = resp.id;

        // Create admin role and grant it to the current user
        return this.setupAdminRoleAndGrant(resp.id);
      })
      .then(() => {
        if (frameworkId === 'other' || !frameworkId) {
          this.router.navigate(['/projects', this.projectId, 'apps', 'create']);
        } else {
          this.createApp(this.projectId);
        }
      })
      .catch((error) => {
        const { message } = error;
        this.error.set(message);
      });
  }

  private async setupAdminRoleAndGrant(projectId: string): Promise<void> {
    try {
      // Create universal admin role
      await this.mgmtService.addProjectRole(projectId, 'admin', 'Administrator', 'Management');

      // Grant admin role to the current user
      await this.grantAdminRoleToCurrentUser(projectId);
    } catch (error) {
      console.warn('Failed to setup admin role and authorization:', error);
      // Don't fail the entire flow if role creation fails
      // The user can still manually set up roles later
    }
  }

  private async grantAdminRoleToCurrentUser(projectId: string): Promise<void> {
    try {
      // Get current user info from the auth service
      const userInfo = await this.authService.getMyUser();

      if (userInfo && userInfo.user?.id) {
        await this.mgmtService.addUserGrant(userInfo.user.id, ['admin'], projectId);
      }
    } catch (error) {
      console.warn('Failed to grant admin role to current user:', error);
      // Don't fail if the grant creation fails
    }
  }
}
