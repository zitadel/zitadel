import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

declare const tinycolor: any;

export interface Color {
  name: string;
  hex: string;
  rgb: string;
  contrastColor: string;
}

@Injectable()
export class ThemeService {
  private _darkTheme: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
  public isDarkTheme: Observable<boolean> = this._darkTheme.asObservable();
  public loading: boolean = false;

  constructor() {
    const theme = localStorage.getItem('theme');
    if (theme) {
      if (theme === 'light-theme') {
        this.setDarkTheme(false);
      } else {
        this.setDarkTheme(true);
      }
    }
  }

  setDarkTheme(isDarkTheme: boolean): void {
    this._darkTheme.next(isDarkTheme);
  }
}
