import { ChangeDetectionStrategy, Component, DestroyRef, effect, inject, Signal, signal } from '@angular/core';
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
import { OIDCAppType, OIDCAuthMethodType, OIDCGrantType, OIDCResponseType } from '@zitadel/proto/zitadel/app/v2beta/oidc_pb';
import { APIAuthMethodType } from '@zitadel/proto/zitadel/app/v2beta/api_pb';
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
import {
  WEB_TYPE,
  NATIVE_TYPE,
  USER_AGENT_TYPE,
  API_TYPE,
  SAML_TYPE,
  RadioItemAppType,
  AppCreateType,
} from 'src/app/pages/projects/apps/authtypes';
import { RadioItemAuthType } from 'src/app/modules/app-radio/app-auth-method-radio/app-auth-method-radio.component';
import {
  PKCE_METHOD,
  CODE_METHOD,
  PK_JWT_METHOD,
  POST_METHOD,
  BASIC_AUTH_METHOD,
  DEVICE_CODE_METHOD,
  IMPLICIT_METHOD,
  getPartialConfigFromAuthMethod,
} from 'src/app/pages/projects/apps/authmethods';
import { AppRadioModule } from 'src/app/modules/app-radio/app-radio.module';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';

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
    AppRadioModule,
    FormFieldModule,
  ],
})
export class AppQuickCreateComponent {
  public InfoSectionType = InfoSectionType;

  public selectedFramework = signal<Framework | undefined>(undefined);

  // Custom app creation mode
  public customAppMode = signal(false);
  public customAppType = signal<RadioItemAppType | undefined>(undefined);
  public customAuthMethod = signal<string>('');
  public customProjectName = generateRandomProjectName();
  public customAppName = '';

  // App types and auth methods for custom mode
  public appTypes: RadioItemAppType[] = [WEB_TYPE, NATIVE_TYPE, USER_AGENT_TYPE, API_TYPE];
  public authMethods: RadioItemAuthType[] = [
    PKCE_METHOD,
    CODE_METHOD,
    PK_JWT_METHOD,
    POST_METHOD,
    BASIC_AUTH_METHOD,
    DEVICE_CODE_METHOD,
    IMPLICIT_METHOD,
  ];
  public readonly initialParam: Signal<string>;

  protected frameworks = frameworksWithOidcConfiguration
    .filter((f) => !('client' in f) || !f.client) // Filter out client libraries/SDKs
    .filter((f) => !('excludeFromAppCreation' in f) || !f.excludeFromAppCreation) // Filter out manually excluded frameworks
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
  protected readonly createApiApplicationMutation = injectMutation(
    inject(ApplicationService).createApplicationMutationOptions<'apiRequest'>,
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

    // Auto-select framework from query param
    effect(() => {
      const frameworkId = this.initialParam();
      if (frameworkId && !this.selectedFramework()) {
        const framework = this.frameworks.find((f) => f.id === frameworkId);
        if (framework) {
          this.selectedFramework.set(framework);
        }
      }
    });
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

  public enableCustomAppMode() {
    // Clear any selected framework first to avoid showing both panels
    this.selectedFramework.set(undefined);
    this.customAppMode.set(true);
    this.customProjectName = generateRandomProjectName();
    this.customAppName = generateFrameworkAppName('Custom');

    // Clear framework query param to prevent auto-reselection
    this.router.navigate([], {
      queryParams: { framework: null },
      queryParamsHandling: 'merge',
      replaceUrl: true,
    });
  }

  public disableCustomAppMode() {
    this.customAppMode.set(false);
    this.customAppType.set(undefined);
    this.customAuthMethod.set('');
    this.customAppName = '';
  }

  public selectFramework(framework: Framework) {
    this.selectedFramework.set(framework);
    // Disable custom app mode when a framework is selected
    if (this.customAppMode()) {
      this.disableCustomAppMode();
    }

    // Update query param to reflect the selected framework
    this.router.navigate([], {
      queryParams: { framework: framework.id },
      queryParamsHandling: 'merge',
      replaceUrl: true,
    });
  }

  public deselectFramework() {
    this.selectedFramework.set(undefined);

    // Clear framework query param to prevent auto-reselection
    this.router.navigate([], {
      queryParams: { framework: null },
      queryParamsHandling: 'merge',
      replaceUrl: true,
    });
  }

  public getAuthMethodDescription(key: string): string {
    const method = this.authMethods.find((m) => m.key === key);
    return method?.descI18nKey || '';
  }

  public getSelectedAuthMethod(): RadioItemAuthType | undefined {
    const key = this.customAuthMethod();
    return key ? this.authMethods.find((m) => m.key === key) : undefined;
  }

  public getFilteredAuthMethods(): RadioItemAuthType[] {
    const appType = this.customAppType();

    if (!appType) {
      return [];
    }

    if (appType.createType === AppCreateType.API) {
      // API: JWT and BASIC
      return this.authMethods.filter((m) => m.apiAuthMethod);
    }

    // OIDC apps - filter by specific app type
    if (appType.oidcAppType === OIDCAppType.OIDC_APP_TYPE_WEB) {
      // Web: PKCE, CODE, JWT, POST
      return this.authMethods.filter((m) => m.key === 'PKCE' || m.key === 'CODE' || m.key === 'PK_JWT' || m.key === 'POST');
    }

    if (appType.oidcAppType === OIDCAppType.OIDC_APP_TYPE_NATIVE) {
      // Native: PKCE and DEVICE_CODE
      return this.authMethods.filter((m) => m.key === 'PKCE' || m.key === 'DEVICECODE');
    }

    if (appType.oidcAppType === OIDCAppType.OIDC_APP_TYPE_USER_AGENT) {
      // User Agent: PKCE and IMPLICIT
      return this.authMethods.filter((m) => m.key === 'PKCE' || m.key === 'IMPLICIT');
    }

    // Fallback - show all OIDC methods
    return this.authMethods;
  }

  public onAppTypeChange(): void {
    // Reset auth method when app type changes if current selection is invalid
    const filteredMethods = this.getFilteredAuthMethods();
    const currentMethod = this.customAuthMethod();

    // Check if current method is still valid for the new app type
    const isCurrentMethodValid = filteredMethods.some((m) => m.key === currentMethod);

    if (!isCurrentMethodValid) {
      this.customAuthMethod.set('');
    }
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

  async createCustomProjectAndApp() {
    try {
      const { project } = this.project() ?? {};
      let projectId: string;

      if (!this.createProject() && project) {
        // They selected an existing project
        projectId = 'id' in project ? project.id : project.projectId;
      } else {
        const projectName = this.projectName().trim();
        if (!this.createProject() && !projectName) {
          throw new Error('Please select an existing project or enter a new project name');
        }

        const response = await this.createProjectMutation.mutateAsync({
          name: projectName || this.customProjectName,
        });
        projectId = response.id;

        // Create admin role and grant it to the current user
        await this.setupAdminRoleAndGrant(projectId);
      }

      await this.createCustomApp(projectId);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  private async createCustomApp(projectId: string): Promise<void> {
    const appType = this.customAppType();
    const authMethodKey = this.customAuthMethod();

    if (!appType) {
      throw new Error('Please select an application type');
    }

    const appName = this.customAppName.trim() || `${appType.prefix}-App`;

    if (appType.createType === AppCreateType.OIDC) {
      if (!authMethodKey) {
        throw new Error('Please select an authentication method');
      }

      let authMethodType: OIDCAuthMethodType;
      let responseTypes: OIDCResponseType[];
      let grantTypes: OIDCGrantType[];

      // Map auth method to proper OIDC configuration
      switch (authMethodKey) {
        case 'PKCE':
          authMethodType = OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;
          responseTypes = [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE];
          grantTypes = [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE];
          break;
        case 'CODE':
          authMethodType = OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC;
          responseTypes = [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE];
          grantTypes = [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE];
          break;
        case 'PK_JWT':
          authMethodType = OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT;
          responseTypes = [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE];
          grantTypes = [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE];
          break;
        case 'POST':
          authMethodType = OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_POST;
          responseTypes = [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE];
          grantTypes = [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE];
          break;
        case 'DEVICECODE':
          authMethodType = OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;
          responseTypes = [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE];
          grantTypes = [OIDCGrantType.OIDC_GRANT_TYPE_DEVICE_CODE];
          break;
        case 'IMPLICIT':
          authMethodType = OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE;
          responseTypes = [OIDCResponseType.OIDC_RESPONSE_TYPE_ID_TOKEN];
          grantTypes = [OIDCGrantType.OIDC_GRANT_TYPE_IMPLICIT];
          break;
        default:
          throw new Error('Unknown authentication method');
      }

      const resp = await this.createOidcApplicationMutation.mutateAsync({
        projectId,
        name: appName,
        creationRequestType: {
          case: 'oidcRequest',
          value: {
            appType: appType.oidcAppType!,
            responseTypes,
            grantTypes,
            authMethodType,
            devMode: true,
            redirectUris: ['http://localhost:3000/callback', 'http://localhost:3000/auth/callback'],
            postLogoutRedirectUris: ['http://localhost:3000'],
          },
        },
      });

      this.toast.showInfo('APP.TOAST.CREATED', true);

      // Navigate to app detail page and pass client secret if present
      const clientSecret = resp.creationResponseType.value.clientSecret;
      await this.router.navigate(['projects', projectId, 'apps', resp.appId], {
        queryParams: { new: true, custom: true },
        state: clientSecret ? { clientSecret } : undefined,
      });
    } else if (appType.createType === AppCreateType.API) {
      if (!authMethodKey) {
        throw new Error('Please select an authentication method');
      }

      let authMethodType: APIAuthMethodType;

      // Map auth method to API auth type
      switch (authMethodKey) {
        case 'PK_JWT':
          authMethodType = APIAuthMethodType.API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT;
          break;
        case 'BASIC':
          authMethodType = APIAuthMethodType.API_AUTH_METHOD_TYPE_BASIC;
          break;
        default:
          throw new Error('Unsupported authentication method for API apps');
      }

      const resp = await this.createApiApplicationMutation.mutateAsync({
        projectId,
        name: appName,
        creationRequestType: {
          case: 'apiRequest',
          value: {
            authMethodType,
          },
        },
      });

      this.toast.showInfo('APP.TOAST.CREATED', true);

      // Navigate to app detail page and pass client secret if present
      const clientSecret = resp.creationResponseType.value.clientSecret;
      await this.router.navigate(['projects', projectId, 'apps', resp.appId], {
        queryParams: { new: true, custom: true },
        state: clientSecret ? { clientSecret } : undefined,
      });
    }
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
