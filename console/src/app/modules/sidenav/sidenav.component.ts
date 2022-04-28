import { Component, forwardRef, Input } from '@angular/core';
import { NG_VALUE_ACCESSOR } from '@angular/forms';

export interface SidenavSetting {
  id: string;
  i18nKey: string;
  featureRequired: string[] | false;
}

@Component({
  selector: 'cnsl-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.scss'],
  providers: [{ provide: NG_VALUE_ACCESSOR, useExisting: forwardRef(() => SidenavComponent), multi: true }],
})
export class SidenavComponent {
  @Input() public currentSetting: string | undefined = 'general';
  @Input() public settingsList: SidenavSetting[] = [];

  constructor() {}

  private onChange: any = () => {};
  private onTouch: any = () => {};

  set value(setting: string | undefined) {
    this.currentSetting = setting;
    this.onChange(setting);
    this.onTouch(setting);
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
