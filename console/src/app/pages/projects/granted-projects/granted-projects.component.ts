import { Component } from '@angular/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-granted-projects',
  templateUrl: './granted-projects.component.html',
  styleUrls: ['./granted-projects.component.scss'],
})
export class GrantedProjectsComponent {
  constructor(breadcrumbService: BreadcrumbService) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }
}
