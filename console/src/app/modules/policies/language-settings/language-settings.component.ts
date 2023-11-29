import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {SetDefaultLanguageResponse, SetRestrictionsRequest} from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import {AbstractControl, FormControl, UntypedFormBuilder, UntypedFormGroup, Validators} from "@angular/forms";
import {LanguagesService} from "../../../services/languages.service";
import {AsyncSubject, BehaviorSubject, concat, forkJoin, from, Observable, Subject} from "rxjs";
import {GrpcAuthService} from "../../../services/grpc-auth.service";
import {i18nValidator} from "../../form-field/validators/validators";
import {CdkDrag, CdkDragDrop, moveItemInArray, transferArrayItem} from "@angular/cdk/drag-drop";

interface RemoteState {
  defaultLang: string,
  allowed: string[]
}

interface LocalState extends RemoteState {
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

  public localState$ = new BehaviorSubject<LocalState>({allowed: [], notAllowed: [], defaultLang: ""});
  public remoteState$ = new BehaviorSubject<RemoteState>({allowed: [], defaultLang: ""});

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
        this.remoteState$.next({defaultLang, allowed});
        this.localState$.next({notAllowed, ...{allowed: [...allowed], defaultLang}});
        this.cdr.detectChanges()
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

  public save(): void {
    const newState = this.localState$.value;
    const remoteState = this.remoteState$.value
    const defaultLangChanged = newState.defaultLang !== remoteState.defaultLang
    const allowedLangsChanged = newState.allowed.length != remoteState.allowed.length ||
      !newState.allowed.every((item, i) => remoteState.allowed[i] === item)
    if (!defaultLangChanged && !allowedLangsChanged) {
      this.toast.showInfo("SETTING.LANGUAGES.NO_CHANGES", true);
      return
    }
    let saveCalls = []
    if (defaultLangChanged && allowedLangsChanged) {
      // Temporarily enable all languages to avoid conflicts
      saveCalls.push(this.service.setRestrictions(undefined, []))
    }
    if (defaultLangChanged) {
      saveCalls.push(from(this.service.setDefaultLanguage(newState.defaultLang)))
    }
    if (allowedLangsChanged) {
      saveCalls.push(from(this.service.setRestrictions(undefined, newState.allowed)))
    }
    const sub = concat(...saveCalls).subscribe({
      next: () => {
        this.remoteState$.next({
          defaultLang: newState.defaultLang,
          allowed: [...newState.allowed],
        });
        this.toast.showInfo("SETTING.LANGUAGES.SAVED", true);
      },
      error: this.toast.showError,
      complete: () => {
        sub.unsubscribe()
      },
    })
  }
}
