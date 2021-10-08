import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { ThemeService } from 'src/app/services/theme.service';

@Component({
  selector: 'cnsl-theme-setting',
  templateUrl: './theme-setting.component.html',
  styleUrls: ['./theme-setting.component.scss'],
})
export class ThemeSettingComponent implements OnInit, OnDestroy {
  public isDarkTheme: boolean = true;
  private sub: Subscription = new Subscription();
  constructor(private themeService: ThemeService) {
    const theme = localStorage.getItem('theme');
    this.isDarkTheme = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;
  }

  ngOnInit(): void {
    this.sub = this.themeService.isDarkTheme.subscribe(d => this.isDarkTheme = d);
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public change(checked: boolean): void {
    this.themeService.setDarkTheme(checked);
  }
}
