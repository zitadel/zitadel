import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { ProjectAutocompleteType } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.component';
import { Project } from 'src/app/proto/generated/zitadel/project_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-app-create',
  templateUrl: './app-create.component.html',
  styleUrls: ['./app-create.component.scss'],
})
export class AppCreateComponent {
  public projectId: string = '';
  public ProjectAutocompleteType: any = ProjectAutocompleteType;

  constructor(private router: Router, breadcrumbService: BreadcrumbService) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }

  public goToAppCreatePage(): void {
    this.router.navigate(['/projects', this.projectId, 'apps', 'create']);
  }

  public close(): void {
    window.history.back();
  }

  public selectProject(project: Project.AsObject): void {
    if (project.id) {
      this.projectId = project.id;
    }
  }
}
