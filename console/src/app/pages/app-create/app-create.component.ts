import { Component, OnDestroy, signal } from '@angular/core';
import { ActivatedRoute, Navigation, Params, Router } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { ProjectAutocompleteType } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.component';
import { AddProjectRequest, AddProjectResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { Framework } from 'src/app/components/quickstart/quickstart.component';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { NavigationService } from 'src/app/services/navigation.service';
import { Location } from '@angular/common';

@Component({
  selector: 'cnsl-app-create',
  templateUrl: './app-create.component.html',
  styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent implements OnDestroy {
  public InfoSectionType: any = InfoSectionType;
  public project?: {
    project: Project.AsObject | GrantedProject.AsObject;
    type: ProjectType;
    name: string;
  } = undefined;
  public ProjectAutocompleteType: any = ProjectAutocompleteType;
  public projectName: string = '';

  public error = signal('');
  public framework = signal<Framework | undefined>(undefined);
  public customFramework = signal<boolean>(false);
  public initialParam = signal<string>('');
  public destroy$: Subject<void> = new Subject();

  public frameworks: Framework[] = frameworkDefinition.map((f) => {
    return {
      ...f,
      fragment: '',
      imgSrcDark: `assets${f.imgSrcDark}`,
      imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
    };
  });
  constructor(
    private router: Router,
    private mgmtService: ManagementService,
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

  public createProjectAndContinue() {
    const project = new AddProjectRequest();
    project.setName(this.projectName);

    return this.mgmtService
      .addProject(project.toObject())
      .then((resp: AddProjectResponse.AsObject) => {
        this.error.set('');
        if (this.customFramework()) {
          this.router.navigate(['/projects', resp.id, 'apps', 'create']);
        } else if (this.framework()) {
          this.router.navigate(['/projects', resp.id, 'apps', 'integrate'], {
            queryParams: { framework: this.framework()?.id },
          });
        }
      })
      .catch((error) => {
        const { message } = error;
        this.error.set(message);
      });
  }
}
