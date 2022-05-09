import { Component, forwardRef, Input, OnInit } from '@angular/core';
import { NG_VALUE_ACCESSOR } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';

export interface SidenavSetting {
  id: string;
  i18nKey: string;
  groupI18nKey?: string;
  requiredRoles?: { [serviceType in PolicyComponentServiceType]: string[] };
}

@Component({
  selector: 'cnsl-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.scss'],
  providers: [{ provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SidenavComponent), multi: true }],
})
export class SidenavComponent implements OnInit {
  @Input() public title: string = '';
  @Input() public description: string = '';
  @Input() public indented: boolean = false;
  @Input() public currentSetting: string | undefined = undefined;
  @Input() public settingsList: SidenavSetting[] = [];
  @Input() public queryParam: string = '';

  constructor(private router: Router, private route: ActivatedRoute) {}

  ngOnInit(): void {
    if (!this.currentSetting) {
      this.value = this.settingsList[0].id;
    }
  }

  private onChange: any = () => {};
  private onTouch: any = () => {};

  set value(setting: string | undefined) {
    this.currentSetting = setting;
    this.onChange(setting);
    this.onTouch(setting);

    if (this.queryParam && setting) {
      this.router.navigate([], {
        relativeTo: this.route,
        queryParams: {
          [this.queryParam]: setting,
        },
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
