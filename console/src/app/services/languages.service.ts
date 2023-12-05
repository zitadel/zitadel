import { from, Observable, share, switchMap, take } from 'rxjs';
import { map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { AdminService } from './admin.service';

interface languagesResponse {
  languagesList: Array<string>;
}

@Injectable({
  providedIn: 'root',
})
export class LanguagesService {
  private supportedLanguages$!: Observable<string[]>;

  constructor(private adminSvc: AdminService) {}

  public allowedLanguages(): Observable<string[]> {
    return this.toObservable(this.adminSvc.getAllowedLanguages());
  }
  // By accepting an observable, we can reuse the results of that API call.
  public notAllowedLanguages(allowedLanguages: Observable<string[]>): Observable<string[]> {
    return allowedLanguages.pipe(
      switchMap((allowed) =>
        // this.supportedLanguages always returns the same observable, so we don't have to worry about API calls.
        this.supportedLanguages().pipe(
          // cut the allowed languages from the supported list
          map((supported) => supported.filter((s) => !allowed.includes(s))),
        ),
      ),
    );
  }
  public supportedLanguages(): Observable<string[]> {
    if (!this.supportedLanguages$) {
      this.supportedLanguages$ = this.toObservable(this.adminSvc.getSupportedLanguages());
    }
    return this.supportedLanguages$;
  }
  private toObservable(resp: Promise<languagesResponse>): Observable<Array<string>> {
    return from(resp).pipe(map(({ languagesList }) => languagesList));
  }
}
