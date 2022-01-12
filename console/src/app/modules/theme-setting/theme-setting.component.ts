import { Component } from '@angular/core';
import { ThemeService } from 'src/app/services/theme.service';

@Component({
  selector: 'cnsl-theme-setting',
  templateUrl: './theme-setting.component.html',
  styleUrls: ['./theme-setting.component.scss'],
})
export class ThemeSettingComponent {
  constructor(public themeService: ThemeService) {
    // const theme = localStorage.getItem('theme');
    // this.isDarkTheme = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;
  }

  public change(event: any): void {
    if (event.target) {
      const checked = event.target.checked;
      this.themeService.setDarkTheme(checked);
    }
  }
}
