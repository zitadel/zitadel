import { ChangeDetectionStrategy, Component, DestroyRef, inject, Signal, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { firstValueFrom, map } from 'rxjs';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { Location } from '@angular/common';
import { ToastService } from 'src/app/services/toast.service';
import { frameworksWithOidcConfiguration, OIDC_CONFIGURATIONS } from 'src/app/utils/framework';
import { generateRandomProjectName, generateFrameworkAppName } from 'src/app/utils/name-generator';
import { ThemeService } from 'src/app/services/theme.service';
import { MatDialog } from '@angular/material/dialog';
import { AppSecretDialogComponent } from 'src/app/pages/projects/apps/app-secret-dialog/app-secret-dialog.component';
import { takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop';
import { NavigationService } from 'src/app/services/navigation.service';
import { ApplicationService } from 'src/app/services/application.service';
import { CreateMutationResult, injectMutation } from '@tanstack/angular-query-experimental';
import { ProjectService } from 'src/app/services/project.service';
import { OIDCAppType, OIDCAuthMethodType } from '@zitadel/proto/zitadel/app/v2beta/oidc_pb';
import { UserService } from 'src/app/services/user.service';
import { AuthorizationService } from 'src/app/services/authorization.service';
import { TranslatePipe } from '@ngx-translate/core';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { MatCheckbox } from '@angular/material/checkbox';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { SearchProjectAutocompleteModule } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';
import { ProjectAutocompleteType } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.component';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { MatButton, MatIconButton } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';

type Framework = (typeof frameworksWithOidcConfiguration)[number];

@Component({
  selector: 'cnsl-app-quick-create',
  templateUrl: './app-quick-create.component.html',
  styleUrls: ['./app-quick-create.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    TranslatePipe,
    CreateLayoutModule,
    MatCheckbox,
    ReactiveFormsModule,
    SearchProjectAutocompleteModule,
    InfoSectionModule,
    MatIconModule,
    MatIconButton,
    MatButton,
    MatProgressSpinnerModule,
    FormsModule,
  ],
})
export class AppQuickCreateComponent {
  public InfoSectionType = InfoSectionType;

  public selectedFramework = signal<Framework | undefined>(undefined);
  public readonly initialParam: Signal<string>;

  protected frameworks = frameworksWithOidcConfiguration
    .filter((f) => !('client' in f) || !f.client) // Filter out client libraries/SDKs
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

      const aIndex = popularityOrder.indexOf(a.id.toLowerCase());
      const bIndex = popularityOrder.indexOf(b.id.toLowerCase());

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

  private readonly destroyRef = inject(DestroyRef);
  private readonly router = inject(Router);
  private readonly toast = inject(ToastService);
  private readonly breadcrumbService = inject(BreadcrumbService);
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly location = inject(Location);
  protected readonly themeService = inject(ThemeService);
  private readonly dialog = inject(MatDialog);
  private readonly navigation = inject(NavigationService);
  protected readonly createOidcApplicationMutation = injectMutation(
    inject(ApplicationService).createApplicationMutationOptions<'oidcRequest'>,
  );
  protected readonly addProjectRoleMutation = injectMutation(inject(ProjectService).addProjectRoleMutationsOptions);
  protected readonly createProjectMutation = injectMutation(inject(ProjectService).createProjectMutationOptions);
  protected readonly createAuthorizationMutation = injectMutation(
    inject(AuthorizationService).createAuthorizationMutationOptions,
  );

  private readonly userId = inject(UserService).userId;
  protected readonly createProject = signal(true);
  protected readonly project = signal<
    | {
        project: Project.AsObject | GrantedProject.AsObject;
        type: ProjectType;
        name: string;
      }
    | undefined
  >(undefined);
  protected readonly projectName = signal('');

  constructor() {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    this.breadcrumbService.setBreadcrumb([bread]);

    const framework$ = this.activatedRoute.queryParamMap.pipe(map((params) => params.get('framework') ?? ''));
    this.initialParam = toSignal(framework$, { initialValue: '' });
  }

  public getAppTypeLabel(framework: Framework): string {
    const { appType } = OIDC_CONFIGURATIONS[framework.id];

    if (appType === OIDCAppType.OIDC_APP_TYPE_WEB) {
      return 'Web Application';
    }
    if (appType === OIDCAppType.OIDC_APP_TYPE_USER_AGENT) {
      return 'User Agent (SPA)';
    }
    if (appType === OIDCAppType.OIDC_APP_TYPE_NATIVE) {
      return 'Native Application';
    }

    throw new Error('Unknown application type');
  }

  public getAuthMethod(framework: Framework): string {
    const { authMethodType } = OIDC_CONFIGURATIONS[framework.id];

    if (authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE) {
      return 'OIDC with PKCE';
    }
    if (authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC) {
      return 'OIDC with Client Secret';
    }

    return 'OIDC';
  }

  public async close() {
    if (this.navigation.isBackPossible) {
      this.location.back();
      return;
    }

    const project = this.createProjectMutation.data();
    if (project) {
      await this.router.navigate(['/projects', project.id]);
      return;
    }
    await this.router.navigate(['/projects']);
  }

  private async setupAdminRoleAndGrant(projectId: string): Promise<void> {
    try {
      // Create universal admin role
      await this.addProjectRoleMutation.mutateAsync({
        projectId,
        roleKey: 'admin',
        displayName: 'Administrator',
        group: 'Management',
      });
    } catch (error) {
      console.warn('Failed to setup admin role and authorization:', error);
      // Don't fail the entire flow if role creation fails
      // The user can still manually set up roles later
      return;
    }

    // Grant admin role to the current user
    await this.grantAdminRoleToCurrentUser(projectId);
  }

  private async grantAdminRoleToCurrentUser(projectId: string): Promise<void> {
    const userId = this.userId();
    if (!userId) {
      return;
    }

    try {
      await this.createAuthorizationMutation.mutateAsync({
        userId,
        projectId,
        roleKeys: ['admin'],
      });
    } catch (error) {
      console.warn('Failed to grant admin role to current user:', error);
      // Don't fail if the grant creation fails
    }
  }

  async createProjectAndContinue(framework: Framework) {
    try {
      const { project } = this.project() ?? {};
      if (!this.createProject() && project) {
        // They selected an existing project
        const projectId = 'id' in project ? project.id : project.projectId;

        await this.createApp(framework, projectId);
        return;
      }
      const projectName = this.projectName().trim();
      if (!this.createProject() && !projectName) {
        throw new Error('Please select an existing project or enter a new project name');
      }

      const response = await this.createProjectMutation.mutateAsync({ name: projectName || generateRandomProjectName() });
      // Create admin role and grant it to the current user
      await this.setupAdminRoleAndGrant(response.id);

      await this.createApp(framework, response.id);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async createApp(framework: Framework, projectId: string): Promise<void> {
    const oidcConfiguration = OIDC_CONFIGURATIONS[framework.id];

    const resp = await this.createOidcApplicationMutation.mutateAsync({
      projectId,
      name: generateFrameworkAppName(framework.title),
      creationRequestType: {
        case: 'oidcRequest',
        value: {
          ...oidcConfiguration,
          devMode: true,
        },
      },
    });

    this.toast.showInfo('APP.TOAST.CREATED', true);

    const usesClientSecret = oidcConfiguration.authMethodType === OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC;
    if (usesClientSecret && resp.creationResponseType.value.clientSecret) {
      const closed$ = this.dialog
        .open(AppSecretDialogComponent, {
          data: {
            clientSecret: resp.creationResponseType.value.clientSecret,
          },
          width: '800px',
        })
        .afterClosed()
        .pipe(takeUntilDestroyed(this.destroyRef));

      await firstValueFrom(closed$);
    }

    await this.router.navigate(['projects', projectId, 'apps', resp.appId], {
      queryParams: { new: true, framework: framework.id },
    });
  }

  public selectProject(project: {
    project: Project.AsObject | GrantedProject.AsObject;
    type: ProjectType;
    name: string;
  }): void {
    if (!project) {
      return;
    }
    this.project.set(project);
  }

  protected readonly ProjectAutocompleteType = ProjectAutocompleteType;
}
