import { Component, OnDestroy } from '@angular/core';
import { Subject, takeUntil } from 'rxjs';
import { ThemeService } from 'src/app/services/theme.service';

@Component({
  selector: 'cnsl-theme-setting',
  templateUrl: './theme-setting.component.html',
  styleUrls: ['./theme-setting.component.scss'],
})
export class ThemeSettingComponent implements OnDestroy {
  public darkTheme: boolean = true;
  private destroy$: Subject<void> = new Subject();
  constructor(public themeService: ThemeService) {
    themeService.isDarkTheme.pipe(takeUntil(this.destroy$)).subscribe((isDark) => (this.darkTheme = isDark));
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public change(event: any): void {
    if (event.target) {
      const checked = event.target.checked;
      this.themeService.setDarkTheme(checked);
    }
  }

  public setTheme(): void {
    this.themeService.setDarkTheme(!this.darkTheme);
  }
}
