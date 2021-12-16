import { Component } from '@angular/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-owned-projects',
  templateUrl: './owned-projects.component.html',
  styleUrls: ['./owned-projects.component.scss'],
})
export class OwnedProjectsComponent {
  constructor(breadcrumbService: BreadcrumbService) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }
}
