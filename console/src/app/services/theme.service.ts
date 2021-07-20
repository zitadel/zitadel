import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';

declare const tinycolor: any;

export interface Color {
  name: string;
  hex: string;
  darkContrast: boolean;
}

@Injectable()
export class ThemeService {
  private _darkTheme: Subject<boolean> = new Subject<boolean>();
  public isDarkTheme: Observable<boolean> = this._darkTheme.asObservable();

  private primaryColorPalette: Color[] = [];
  private backgroundColorPalette: Color[] = [];

  setDarkTheme(isDarkTheme: boolean): void {
    this._darkTheme.next(isDarkTheme);
  }

  public updateTheme(colors: Color[], type: string, theme: string) {
    colors.forEach(color => {
      document.documentElement.style.setProperty(
        `--theme-${theme}-${type}-${color.name}`,
        color.hex
      );
      document.documentElement.style.setProperty(
        `--theme-${theme}-${type}-contrast-${color.name}`,
        color.darkContrast ? 'hsla(0, 0%, 0%, 0.87)' : '#ffffff'
      );
    });
  }

  public savePrimaryColor(color: string, isDark: boolean) {
    this.primaryColorPalette = this.computeColors(color);
    this.updateTheme(this.primaryColorPalette, 'primary', isDark ? 'dark' : 'light');
  }

  public saveBackgroundColor(color: string, isDark: boolean) {
    this.backgroundColorPalette = this.computeColors(color);
    this.updateTheme(this.backgroundColorPalette, 'background', isDark ? 'dark' : 'light');
  }

  private computeColors(hex: string): Color[] {
    return [
      this.getColorObject(tinycolor(hex).lighten(52), '50'),
      this.getColorObject(tinycolor(hex).lighten(37), '100'),
      this.getColorObject(tinycolor(hex).lighten(26), '200'),
      this.getColorObject(tinycolor(hex).lighten(12), '300'),
      this.getColorObject(tinycolor(hex).lighten(6), '400'),
      this.getColorObject(tinycolor(hex), '500'),
      this.getColorObject(tinycolor(hex).darken(6), '600'),
      this.getColorObject(tinycolor(hex).darken(12), '700'),
      this.getColorObject(tinycolor(hex).darken(18), '800'),
      this.getColorObject(tinycolor(hex).darken(24), '900'),
      this.getColorObject(tinycolor(hex).lighten(50).saturate(30), 'A100'),
      this.getColorObject(tinycolor(hex).lighten(30).saturate(30), 'A200'),
      this.getColorObject(tinycolor(hex).lighten(10).saturate(15), 'A400'),
      this.getColorObject(tinycolor(hex).lighten(5).saturate(5), 'A700')
    ];
  }

  private getColorObject(value: any, name: string): Color {
    const c = tinycolor(value);
    return {
      name: name,
      hex: c.toHexString(),
      darkContrast: c.isLight()
    };
  }
}
