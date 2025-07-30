import { Component } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.scss'],
})
export class UserListComponent {
  constructor(
    protected readonly translate: TranslateService,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }
}
