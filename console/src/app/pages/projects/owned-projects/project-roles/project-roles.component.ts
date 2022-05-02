import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

const ROUTEPARAM = 'projectid';

@Component({
  selector: 'cnsl-project-roles',
  templateUrl: './project-roles.component.html',
  styleUrls: ['./project-roles.component.scss'],
})
export class ProjectRolesComponent {
  public projectId: string = '';

  constructor(private route: ActivatedRoute, private breadcrumbService: BreadcrumbService) {
    const projectId = this.route.snapshot.paramMap.get(ROUTEPARAM);
    if (projectId) {
      this.projectId = projectId;

      const breadcrumbs = [
        new Breadcrumb({
          type: BreadcrumbType.IAM,
          name: 'Instance',
          routerLink: ['/instance'],
        }),
        new Breadcrumb({
          type: BreadcrumbType.ORG,
          routerLink: ['/org'],
        }),
        new Breadcrumb({
          type: BreadcrumbType.PROJECT,
          name: '',
          param: { key: ROUTEPARAM, value: projectId },
          routerLink: ['/projects', projectId],
        }),
      ];
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    }
  }
}
