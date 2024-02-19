import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { ProjectAutocompleteType } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.component';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-app-create',
  templateUrl: './app-create.component.html',
  styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent {
  public project?: {
    project: Project.AsObject | GrantedProject.AsObject;
    type: ProjectType;
    name: string;
  } = undefined;
  public ProjectAutocompleteType: any = ProjectAutocompleteType;
  public projectName: string = '';

  constructor(
    private router: Router,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }

  public goToAppIntegratePage(): void {
    if (this.project) {
      const id = (this.project.project as Project.AsObject).id
        ? (this.project.project as Project.AsObject).id
        : (this.project.project as GrantedProject.AsObject).projectId
          ? (this.project.project as GrantedProject.AsObject).projectId
          : '';

      this.router.navigate(['/projects', id, 'apps', 'integrate']);
    }
  }

  public close(): void {
    window.history.back();
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

  public createProjectAndContinue(): void {}
}
