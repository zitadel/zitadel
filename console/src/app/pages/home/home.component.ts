import { Component } from '@angular/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  public dark: boolean = true;

  constructor(public authService: GrpcAuthService, breadcrumbService: BreadcrumbService) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;
  }
}
