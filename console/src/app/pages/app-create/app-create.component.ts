import { Component, signal } from '@angular/core';
import { Router } from '@angular/router';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { ProjectAutocompleteType } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.component';
import { AddProjectRequest, AddProjectResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';

@Component({
  selector: 'cnsl-app-create',
  templateUrl: './app-create.component.html',
  styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent {
  public InfoSectionType: any = InfoSectionType;
  public project?: {
    project: Project.AsObject | GrantedProject.AsObject;
    type: ProjectType;
    name: string;
  } = undefined;
  public ProjectAutocompleteType: any = ProjectAutocompleteType;
  public projectName: string = '';

  public error = signal('');

  constructor(
    private router: Router,
    private mgmtService: ManagementService,
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

  public createProjectAndContinue() {
    const project = new AddProjectRequest();
    project.setName(this.projectName);

    return this.mgmtService
      .addProject(project.toObject())
      .then((resp: AddProjectResponse.AsObject) => {
        this.error.set('');
        this.router.navigate(['/projects', resp.id, 'apps', 'integrate']);
      })
      .catch((error) => {
        const { message } = error;
        this.error.set(message);
      });
  }
}
