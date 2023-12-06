import { combineLatest, forkJoin, from, merge, Observable, of, share, switchMap, take, zip } from 'rxjs';
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
  public supportedLanguages$: Observable<string[]> = from(this.adminSvc.getSupportedLanguages()).pipe(
    map((list) => list.languagesList),
  );
  public allowedLanguages$: Observable<string[]> = from(this.adminSvc.getAllowedLanguages()).pipe(
    map((list) => list.languagesList),
  );
  public notAllowedLanguages$: Observable<string[]> = combineLatest([this.supportedLanguages$, this.allowedLanguages$]).pipe(
    switchMap(([supported, allowed]) => [supported.filter((s) => !allowed.includes(s))]), // TODO return valid array
  );
  constructor(private adminSvc: AdminService) {}
}
