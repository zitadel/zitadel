import { Component, forwardRef, Input, OnInit } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
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
}

@Component({
  selector: 'cnsl-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.scss'],
  providers: [{ provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SidenavComponent), multi: true }],
})
export class SidenavComponent implements ControlValueAccessor {
  @Input() public title: string = '';
  @Input() public description: string = '';
  @Input() public indented: boolean = false;
  @Input() public currentSetting?: string | undefined = undefined;
  @Input() public settingsList: SidenavSetting[] = [];
  @Input() public queryParam: string = '';

  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  constructor(
    private router: Router,
    private route: ActivatedRoute,
  ) {}

  private onChange = (current: string | undefined) => {};
  private onTouch = (current: string | undefined) => {};

  @Input() get value(): string | undefined {
    return this.currentSetting;
  }

  set value(setting: string | undefined) {
    this.currentSetting = setting;

    if (setting || setting === undefined || setting === '') {
      this.onChange(setting);
      this.onTouch(setting);
    }

    if (this.queryParam && setting) {
      this.router
        .navigate([], {
          relativeTo: this.route,
          queryParams: {
            [this.queryParam]: setting,
          },
          replaceUrl: true,
          queryParamsHandling: 'merge',
          skipLocationChange: false,
        })
        .then();
    }
  }

  public writeValue(value: any) {
    this.value = value;
  }

  public registerOnChange(fn: any) {
    this.onChange = fn;
  }

  public registerOnTouched(fn: any) {
    this.onTouch = fn;
  }
}
