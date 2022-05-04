import { Component, OnInit } from '@angular/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-org-settings',
  templateUrl: './org-settings.component.html',
  styleUrls: ['./org-settings.component.scss'],
})
export class OrgSettingsComponent implements OnInit {
  constructor(breadcrumbService: BreadcrumbService) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ];
    breadcrumbService.setBreadcrumb(breadcrumbs);
  }

  ngOnInit(): void {}
}
