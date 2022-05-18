import { Component, Input, OnInit } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

export interface SidenavSetting {
  id: string;
  i18nKey: string;
  groupI18nKey?: string;
  requiredRoles?: { [serviceType in PolicyComponentServiceType]: string[] };
  showWarn?: boolean;
}

@Component({
  selector: 'cnsl-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.scss'],
  providers: [{ provide: NG_VALUE_ACCESSOR, useExisting: SidenavComponent, multi: true }],
})
export class SidenavComponent implements ControlValueAccessor, OnInit {
  @Input() public title: string = '';
  @Input() public description: string = '';
  @Input() public indented: boolean = false;
  @Input() public currentSetting?: string | undefined = undefined;
  @Input() public settingsList: SidenavSetting[] = [];
  @Input() public queryParam: string = '';

  constructor(private router: Router, private route: ActivatedRoute) {}

  ngOnInit(): void {
    if (!this.value) {
      this.value = this.settingsList[0].id;
    }
  }

  private onChange = (current: string | undefined) => {};
  private onTouch = (current: string | undefined) => {};

  @Input() get value(): string | undefined {
    return this.currentSetting;
  }

  set value(setting: string | undefined) {
    this.currentSetting = setting;

    if (setting || setting === undefined) {
      this.onChange(setting);
      this.onTouch(setting);
    }

    if (this.queryParam && setting) {
      this.router.navigate([], {
        relativeTo: this.route,
        queryParams: {
          [this.queryParam]: setting,
        },
        replaceUrl: true,
        queryParamsHandling: 'merge',
        skipLocationChange: false,
      });
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
