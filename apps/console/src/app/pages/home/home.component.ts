import { Component, effect } from '@angular/core';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ThemeService } from 'src/app/services/theme.service';
import { COLORS } from 'src/app/utils/color';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'cnsl-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  public greendark: string = COLORS[6][700];
  public greenlight = COLORS[6][200];

  public cyandark: string = COLORS[7][700];
  public cyanlight = COLORS[7][200];

  public bluedark: string = COLORS[9][700];
  public bluelight = COLORS[9][200];

  public dark: boolean = true;

  protected readonly PolicyComponentServiceType = PolicyComponentServiceType;

  private readonly permissions = this.newAuthService.listMyZitadelPermissionsQuery();

  constructor(
    public authService: GrpcAuthService,
    private readonly newAuthService: NewAuthService,
    breadcrumbService: BreadcrumbService,
    public themeService: ThemeService,
    private readonly router: Router,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.INSTANCE,
      routerLink: ['/'],
    };

    breadcrumbService.setBreadcrumb([bread]);

    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;

    effect(() => {
      const permission = this.permissions.data();
      if (!permission) {
        return;
      }
      if (permission.includes('iam.read')) {
        return;
      }
      this.router.navigate(['/org']).then();
    });
  }
}
