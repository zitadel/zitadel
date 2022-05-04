import { Component } from '@angular/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-instance-settings',
  templateUrl: './instance-settings.component.html',
  styleUrls: ['./instance-settings.component.scss'],
})
export class InstanceSettingsComponent {
  constructor(breadcrumbService: BreadcrumbService) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    breadcrumbService.setBreadcrumb(breadcrumbs);
  }
}
