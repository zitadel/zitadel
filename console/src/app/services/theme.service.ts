import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';

@Injectable()
export class ThemeService {
    private _darkTheme: Subject<boolean> = new Subject<boolean>();
    public isDarkTheme: Observable<boolean> = this._darkTheme.asObservable();

    setDarkTheme(isDarkTheme: boolean): void {
        this._darkTheme.next(isDarkTheme);
    }
}
