import {from, Observable, share, switchMap, take} from "rxjs";
import {map} from "rxjs/operators";
import {Injectable} from "@angular/core";

interface LanguagesProvider {
  getAllowedLanguages(): Promise<languagesResponse>;
  getSupportedLanguages(): Promise<languagesResponse>;
}

interface languagesResponse {
  languagesList: Array<string>
}

@Injectable({
  providedIn: 'root',
})
export class LanguagesService {

  private supportedLanguages$!: Observable<string[]>;

  public allowedLanguages(svc: LanguagesProvider): Observable<string[]> {
    return this.toObservable(svc.getAllowedLanguages());
  }
  // By accepting an observable, we can reuse the results of that API call.
  public notAllowedLanguages(svc: LanguagesProvider, allowedLanguages: Observable<string[]>): Observable<string[]> {
    return allowedLanguages.pipe(
      switchMap((allowed) =>
        // this.supportedLanguages always returns the same observable, so we don't have to worry about API calls.
        this.supportedLanguages(svc).pipe(
          // cut the allowed languages from the supported list
          map((supported) => supported.filter((s) => !allowed.includes(s))),
        ),
      ),
    );
  }
  public supportedLanguages(svc: LanguagesProvider): Observable<string[]> {
    if (!this.supportedLanguages$) {
      this.supportedLanguages$ = this.toObservable(svc.getSupportedLanguages());
    }
    return this.supportedLanguages$;
  }
  private toObservable(resp: Promise<languagesResponse>): Observable<Array<string>> {
    return from(resp).pipe(map(({ languagesList }) => languagesList ));
  }
}
