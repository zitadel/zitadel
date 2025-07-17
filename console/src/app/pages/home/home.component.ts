import { Component } from '@angular/core';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ThemeService } from 'src/app/services/theme.service';
import { COLORS } from 'src/app/utils/color';

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

  constructor(
    public authService: GrpcAuthService,
    breadcrumbService: BreadcrumbService,
    public themeService: ThemeService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.INSTANCE,
      routerLink: ['/'],
    };

    breadcrumbService.setBreadcrumb([bread]);

    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;
  }
}
