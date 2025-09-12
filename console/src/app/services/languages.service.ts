import { forkJoin, Observable, ReplaySubject, Subscription } from 'rxjs';
import { map, withLatestFrom } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { AdminService } from './admin.service';

@Injectable({
  providedIn: 'root',
})
export class LanguagesService {
  private supportedSubject$ = new ReplaySubject<string[]>(1);
  public supported$: Observable<string[]> = this.supportedSubject$.asObservable();
  private allowedSubject$ = new ReplaySubject<string[]>(1);
  public allowed$: Observable<string[]> = this.allowedSubject$.asObservable();
  public notAllowed$: Observable<string[]> = this.allowed$.pipe(
    withLatestFrom(this.supported$),
    map(([allowed, supported]) => {
      return supported.filter((s) => !allowed.includes(s));
    }),
  );
  public restricted$: Observable<boolean> = this.notAllowed$.pipe(
    map((notallowed) => {
      return notallowed.length > 0;
    }),
  );

  constructor(private adminSvc: AdminService) {
    const sub: Subscription = forkJoin([
      this.adminSvc.getSupportedLanguages(),
      this.adminSvc.getAllowedLanguages(),
    ]).subscribe({
      next: ([{ languagesList: supported }, { languagesList: allowed }]) => {
        this.supportedSubject$.next(supported);
        this.allowedSubject$.next(allowed);
      },
      complete: () => sub.unsubscribe(),
    });
  }

  public newAllowed(languages: string[]) {
    this.allowedSubject$.next(languages);
  }

  public isNotAllowed(language: string): Observable<boolean> {
    return this.notAllowed$.pipe(map((notAllowed) => notAllowed.includes(language)));
  }
}
