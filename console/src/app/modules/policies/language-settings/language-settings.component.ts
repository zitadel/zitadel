import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {SetDefaultLanguageResponse, SetRestrictionsRequest} from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import {AbstractControl, FormControl, UntypedFormBuilder, UntypedFormGroup, Validators} from "@angular/forms";
import {LanguagesService} from "../../../services/languages.service";
import {AsyncSubject, BehaviorSubject, concat, forkJoin, from, Observable, of, Subject} from "rxjs";
import {GrpcAuthService} from "../../../services/grpc-auth.service";
import {i18nValidator} from "../../form-field/validators/validators";
import {CdkDrag, CdkDragDrop, moveItemInArray, transferArrayItem} from "@angular/cdk/drag-drop";
import {catchError, map} from "rxjs/operators";

interface State {
  defaultLang: string,
  allowed: string[]
  notAllowed: string[]
}

@Component({
  selector: 'cnsl-language-settings',
  templateUrl: './language-settings.component.html',
  styleUrls: ['./language-settings.component.scss'],
})
export class LanguageSettingsComponent {

  public canWriteRestrictions$: Observable<boolean> = this.authService.isAllowed(["iam.restrictions.write"]);
  public canWriteDefaultLanguage$: Observable<boolean> = this.authService.isAllowed(["iam.write"]);

  public localState$ = new BehaviorSubject<State>({allowed: [], notAllowed: [], defaultLang: ""});
  public remoteState$ = new BehaviorSubject<State>({allowed: [], notAllowed: [], defaultLang: ""});

  public loading: boolean = false;
  constructor(
    private service: AdminService,
    private toast: ToastService,
    private fb: UntypedFormBuilder,
    private languagesSvc: LanguagesService,
    private authService: GrpcAuthService,
    private cdr: ChangeDetectorRef,
  ) {
    const allowedInit$ = this.languagesSvc.allowedLanguages(this.service)
    const notAllowedInit$ = this.languagesSvc.notAllowedLanguages(this.service, allowedInit$);
    const defaultLang$ = from(this.service.getDefaultLanguage());
    const sub = forkJoin([allowedInit$, notAllowedInit$, defaultLang$]).subscribe({
      next: ([allowed, notAllowed, {language: defaultLang}]) => {
        this.remoteState$.next({notAllowed: [...notAllowed], ...{allowed: [...allowed], defaultLang}});
        this.localState$.next({notAllowed: [...notAllowed], ...{allowed: [...allowed], defaultLang}});
      },
      error: this.toast.showError,
      complete: () => {
        sub.unsubscribe()
      },
    })
  }

  drop(event: CdkDragDrop<string[]>) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
    } else {
      transferArrayItem(
        event.previousContainer.data,
        event.container.data,
        event.previousIndex,
        event.currentIndex,
      );
    }
  }

  setLocalDefaultLang(lang: any): void {
    this.localState$.next({...this.localState$.value, defaultLang: lang});
  }

  defaultLangPredicate = (lang: CdkDrag<string>) => {
    return !!lang?.data && lang.data !== this.localState$.value.defaultLang;
  }

  public isRemotelyDisallowed$(lang: string): Observable<boolean>{
    return this.remoteState$.pipe(
      map(({allowed}) => !allowed.includes(lang))
    )
  }

  public save(): void {
    const newState = this.localState$.value;
    const remoteState = this.remoteState$.value
    const sub = concat(
      from(this.service.setDefaultLanguage(newState.defaultLang)).pipe(
        // We just ignore if the instance is unchanged
        catchError((err, caught) => (err as {message: string}).message.includes('INST-DS3rq') ? of(true) : caught)
      ),
      from(this.service.setRestrictions(undefined, newState.allowed))
    ).subscribe({
      next: () => {
        this.remoteState$.next({
          defaultLang: newState.defaultLang,
          allowed: [...newState.allowed],
          notAllowed: [...newState.notAllowed],
        });
        this.toast.showInfo("SETTING.LANGUAGES.SAVED", true);
      },
      error: this.toast.showError,
      complete: () => {
        sub.unsubscribe()
      },
    })
  }

  public discard(): void {
    const remoteState = this.remoteState$.value;
    this.localState$.next({
      defaultLang: remoteState.defaultLang,
      allowed: [...remoteState.allowed],
      notAllowed: [...remoteState.notAllowed],
    });
  }
}
