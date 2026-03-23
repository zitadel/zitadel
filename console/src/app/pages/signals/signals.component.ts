import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-signals-page',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: '<router-outlet></router-outlet>',
})
export class SignalsComponent {
  private readonly breadcrumbService = inject(BreadcrumbService);
  private readonly translate = inject(TranslateService);

  constructor() {
    this.breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.SIGNALS,
        name: this.translate.instant('MENU.SIGNALS'),
        routerLink: ['/signals'],
      }),
    ]);
  }
}
