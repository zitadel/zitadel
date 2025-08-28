import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { UntypedFormBuilder } from '@angular/forms';
import { LanguagesService } from '../../../services/languages.service';
import { BehaviorSubject, concat, forkJoin, from, Observable, of, Subject, switchMap, take, takeUntil } from 'rxjs';
import { GrpcAuthService } from '../../../services/grpc-auth.service';
import { CdkDrag, CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { catchError, map } from 'rxjs/operators';

interface State {
  allowed: string[];
  notAllowed: string[];
}

@Component({
  selector: 'cnsl-language-settings',
  templateUrl: './language-settings.component.html',
  styleUrls: ['./language-settings.component.scss'],
})
export class LanguageSettingsComponent {
  public canWriteRestrictions$: Observable<boolean> = this.authService.isAllowed(['iam.restrictions.write']);
  public canWriteDefaultLanguage$: Observable<boolean> = this.authService.isAllowed(['iam.write']);

  public localState$ = new BehaviorSubject<State>({ allowed: [], notAllowed: [] });
  public remoteState$ = new BehaviorSubject<State>({ allowed: [], notAllowed: [] });
  public defaultLang$ = new BehaviorSubject<string>('');

  public loading: boolean = false;
  constructor(
    private service: AdminService,
    private toast: ToastService,
    private langSvc: LanguagesService,
    private authService: GrpcAuthService,
  ) {
    const sub = forkJoin([
      langSvc.allowed$.pipe(take(1)),
      langSvc.notAllowed$.pipe(take(1)),
      from(this.service.getDefaultLanguage()).pipe(take(1)),
    ]).subscribe({
      next: ([allowed, notAllowed, { language: defaultLang }]) => {
        this.defaultLang$.next(defaultLang);
        this.remoteState$.next({ notAllowed: [...notAllowed], ...{ allowed: [...allowed] } });
        this.localState$.next({ notAllowed: [...notAllowed], ...{ allowed: [...allowed] } });
      },
      error: this.toast.showError,
      complete: () => {
        sub.unsubscribe();
      },
    });
  }

  drop(event: CdkDragDrop<string[]>) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
    } else {
      transferArrayItem(event.previousContainer.data, event.container.data, event.previousIndex, event.currentIndex);
    }
  }

  public defaultLangPredicate = (lang: CdkDrag<string>) => {
    return !!lang?.data && lang.data !== this.defaultLang$.value;
  };

  public isRemotelyAllowed$(lang: string): Observable<boolean> {
    return this.remoteState$.pipe(map(({ allowed }) => allowed.includes(lang)));
  }

  public allowAll(): void {
    this.localState$.next({ allowed: [...this.allLocalLangs()], notAllowed: [] });
  }

  public disallowAll(): void {
    const disallowed = this.allLocalLangs().filter((lang) => lang !== this.defaultLang$.value);
    this.localState$.next({ allowed: [this.defaultLang$.value], notAllowed: disallowed });
  }

  public submit(): void {
    const { allowed, notAllowed } = this.localState$.value;
    const sub = from(this.service.setRestrictions(undefined, allowed)).subscribe({
      next: () => {
        this.remoteState$.next({
          allowed: [...allowed],
          notAllowed: [...notAllowed],
        });
        this.langSvc.newAllowed(allowed);
        this.toast.showInfo('SETTING.LANGUAGES.ALLOWED_SAVED', true);
      },
      error: this.toast.showError,
      complete: () => {
        sub.unsubscribe();
      },
    });
  }

  public discard(): void {
    this.localState$.next(this.remoteState$.value);
  }

  public setDefaultLang(lang: string): void {
    const sub = from(this.service.setDefaultLanguage(lang)).subscribe({
      next: () => {
        this.defaultLang$.next(lang);
        this.toast.showInfo('SETTING.LANGUAGES.DEFAULT_SAVED', true);
      },
      error: this.toast.showError,
      complete: () => {
        sub.unsubscribe();
      },
    });
  }

  private allLocalLangs(): string[] {
    return [...this.localState$.value.allowed, ...this.localState$.value.notAllowed];
  }
}
