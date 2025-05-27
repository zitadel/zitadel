import { ChangeDetectionStrategy, Component, effect, EventEmitter, Input, Output, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

export interface SidenavSetting {
  id: string;
  i18nKey: string;
  groupI18nKey?: string;
  requiredRoles?: {
    [PolicyComponentServiceType.MGMT]?: string[];
    [PolicyComponentServiceType.ADMIN]?: string[];
  };
  showWarn?: boolean;
  beta?: boolean;
}

@Component({
  selector: 'cnsl-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidenavComponent {
  @Input() public navigate: boolean = true;
  @Input() public indented: boolean = false;
  @Input({ required: true }) public settingsList: SidenavSetting[] = [];
  @Input({ required: true })
  public set setting(setting: SidenavSetting | null) {
    if (!setting) {
      return;
    }
    this.setting$.set(setting);
  }

  @Output()
  public settingChange = new EventEmitter<SidenavSetting>();

  protected readonly setting$ = signal<SidenavSetting | null>(null);

  protected PolicyComponentServiceType = PolicyComponentServiceType;

  constructor(
    private readonly router: Router,
    private readonly route: ActivatedRoute,
  ) {
    effect(
      () => {
        const setting = this.setting$();
        if (setting === null) {
          return;
        }

        this.settingChange.emit(setting);

        if (!this.navigate) {
          return;
        }

        this.router
          .navigate([], {
            relativeTo: this.route,
            queryParams: {
              id: setting ? setting.id : undefined,
            },
            replaceUrl: true,
            queryParamsHandling: 'merge',
            skipLocationChange: false,
          })
          .then();
      },
      { allowSignalWrites: true },
    );
  }

  protected trackSettings(_: number, setting: SidenavSetting): string {
    return setting.id;
  }
}
