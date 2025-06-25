import { ChangeDetectionStrategy, Component, effect, Signal, signal } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { NewOrganizationService } from '../../services/new-organization.service';
import { ToastService } from '../../services/toast.service';
import { AsyncPipe, NgIf, NgTemplateOutlet } from '@angular/common';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { OrganizationSelectorComponent } from './organization-selector/organization-selector.component';
import { CdkOverlayOrigin } from '@angular/cdk/overlay';
import { MatSelectModule } from '@angular/material/select';
import { InputModule } from '../input/input.module';
import { HeaderButtonComponent } from './header-button/header-button.component';
import { HeaderDropdownComponent } from './header-dropdown/header-dropdown.component';
import { InstanceSelectorComponent } from './instance-selector/instance-selector.component';
import { HasRolePipeModule } from '../../pipes/has-role-pipe/has-role-pipe.module';
import { map } from 'rxjs/operators';
import { toSignal } from '@angular/core/rxjs-interop';
import { BreakpointObserver } from '@angular/cdk/layout';
import { NewAdminService } from '../../services/new-admin.service';
import { NewAuthService } from '../../services/new-auth.service';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'cnsl-new-header',
  templateUrl: './new-header.component.html',
  styleUrls: ['./new-header.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    MatToolbarModule,
    OrganizationSelectorComponent,
    CdkOverlayOrigin,
    MatSelectModule,
    NgIf,
    InputModule,
    HeaderButtonComponent,
    HeaderDropdownComponent,
    InstanceSelectorComponent,
    NgTemplateOutlet,
    AsyncPipe,
    HasRolePipeModule,
    RouterLink,
  ],
})
export class NewHeaderComponent {
  protected readonly listMyZitadelPermissionsQuery = this.newAuthService.listMyZitadelPermissionsQuery();
  protected readonly myInstanceQuery = this.adminService.getMyInstanceQuery();
  protected readonly organizationsQuery = injectQuery(() => ({
    ...this.newOrganizationService.listOrganizationsQueryOptions(),
    enabled: (this.listMyZitadelPermissionsQuery.data() ?? []).includes('org.read'),
  }));
  protected readonly isInstanceDropdownOpen = signal(false);
  protected readonly isOrgDropdownOpen = signal(false);
  protected readonly instanceSelectorSecondStep = signal(false);
  protected readonly activeOrganizationQuery = this.newOrganizationService.activeOrganizationQuery();
  protected readonly isHandset: Signal<boolean>;

  constructor(
    private readonly newOrganizationService: NewOrganizationService,
    private readonly toastService: ToastService,
    private readonly breakpointObserver: BreakpointObserver,
    private readonly adminService: NewAdminService,
    private readonly newAuthService: NewAuthService,
  ) {
    this.isHandset = this.getIsHandset();

    effect(() => {
      if (this.listMyZitadelPermissionsQuery.isError()) {
        this.toastService.showError(this.listMyZitadelPermissionsQuery.error());
      }
    });

    effect(() => {
      if (this.organizationsQuery.isError()) {
        this.toastService.showError(this.organizationsQuery.error());
      }
    });

    effect(() => {
      if (this.myInstanceQuery.isError()) {
        this.toastService.showError(this.myInstanceQuery.error());
      }
    });
  }

  private getIsHandset() {
    const mediaQuery = '(max-width: 599px)';
    const isHandset$ = this.breakpointObserver.observe(mediaQuery).pipe(map(({ matches }) => matches));
    return toSignal(isHandset$, { initialValue: this.breakpointObserver.isMatched(mediaQuery) });
  }
}
