import {from, Observable, share, take} from "rxjs";
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
    return this.toSharedObservable(svc.getAllowedLanguages());
  }
  public supportedLanguages(svc: LanguagesProvider): Observable<string[]> {
    if (!this.supportedLanguages$) {
      this.supportedLanguages$ = this.toSharedObservable(svc.getSupportedLanguages());
    }
    return this.supportedLanguages$;
  }
  private toSharedObservable(resp: Promise<languagesResponse>): Observable<Array<string>> {
    return from(resp).pipe(
      take(1),
      map(({ languagesList }) => languagesList ),
      share(),
    );
  }
}
