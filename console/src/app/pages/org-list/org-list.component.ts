import { Component } from '@angular/core';
import { enterAnimations } from 'src/app/animations';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-org-list',
  templateUrl: './org-list.component.html',
  styleUrls: ['./org-list.component.scss'],
  animations: [enterAnimations],
})
export class OrgListComponent {
  constructor(breadcrumbService: BreadcrumbService) {
    const iamBread = new Breadcrumb({
      type: BreadcrumbType.INSTANCE,
      name: 'Instance',
      routerLink: ['/instance'],
    });

    breadcrumbService.setBreadcrumb([iamBread]);
  }
}
