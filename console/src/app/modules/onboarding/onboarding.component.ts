import { Component } from '@angular/core';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-onboarding',
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent {
  constructor(public authService: GrpcAuthService, breadcrumbService: BreadcrumbService) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };

    breadcrumbService.setBreadcrumb([bread]);
  }
}
