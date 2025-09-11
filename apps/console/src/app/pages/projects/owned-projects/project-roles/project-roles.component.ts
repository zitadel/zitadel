import { Component, Input } from '@angular/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

const ROUTEPARAM = 'projectid';

@Component({
  selector: 'cnsl-project-roles',
  templateUrl: './project-roles.component.html',
  styleUrls: ['./project-roles.component.scss'],
})
export class ProjectRolesComponent {
  @Input() public projectId: string = '';

  constructor(private breadcrumbService: BreadcrumbService) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
      new Breadcrumb({
        type: BreadcrumbType.PROJECT,
        name: '',
        param: { key: ROUTEPARAM, value: this.projectId },
        routerLink: ['/projects', this.projectId],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);
  }
}
