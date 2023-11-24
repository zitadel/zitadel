import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import { SetDefaultLanguageResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import {AbstractControl, FormControl, UntypedFormBuilder, UntypedFormGroup, Validators} from "@angular/forms";
import {LanguagesService} from "../../../services/languages.service";
import {AsyncSubject, forkJoin, from, Observable, Subject} from "rxjs";
import {GrpcAuthService} from "../../../services/grpc-auth.service";
import {i18nValidator} from "../../form-field/validators/validators";

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

  public form!: UntypedFormGroup;
  public formLoaded$ = new AsyncSubject<boolean>();

  public canWriteRestrictions$: Observable<boolean> = this.authService.isAllowed(["iam.restrictions.write"]);
  public canWriteDefaultLanguage$: Observable<boolean> = this.authService.isAllowed(["iam.write"]);

  private remoteState: RemoteState | null = null;

  public loading: boolean = false;
  constructor(
    private service: AdminService,
    private toast: ToastService,
    private fb: UntypedFormBuilder,
    private languagesSvc: LanguagesService,
    private authService: GrpcAuthService,
    private cdr: ChangeDetectorRef,
  ) {
    this.form = this.fb.group({
      defaultLang: ['', [i18nValidator(
        "SETTING.LANGUAGES.DEFAULT_MUST_BE_ALLOWED",
        (control: AbstractControl) => {
          const alloweds = (this.form?.getRawValue() as LocalState)?.allowed
          return alloweds && alloweds.find(allowed => allowed == control.value) !== undefined
        }
      )]],
      allowed: [[''], []],
      notAllowed: [[''], []],
    });
    const allowed$ = this.languagesSvc.allowedLanguages(this.service)
    const notAllowed$: Observable<string[]> = this.languagesSvc.notAllowedLanguages(this.service, allowed$);
    const defaultLang$ = from(this.service.getDefaultLanguage());
    const sub = forkJoin([allowed$, notAllowed$, defaultLang$]).subscribe({
      next: ([allowed, notAllowed, {language: defaultLang}]) => {
        this.remoteState = {defaultLang, allowed,};
        this.form.setValue(<LocalState>{notAllowed, ...this.remoteState});
        this.formLoaded$.next(true);
        this.formLoaded$.complete();
        console.log("da")
        this.formLoaded$.asObservable().subscribe((x) => {console.log("hier", x)})
        this.formLoaded$.subscribe((x) => {console.log("hier", x)})
        this.formLoaded$.subscribe((x) => {console.log("hier", x)})
        this.formLoaded$.subscribe((x) => {console.log("hier", x)})
        this.formLoaded$.subscribe((x) => {console.log("hier", x)})
        this.cdr.detectChanges()
      },
      error: this.toast.showError,
      complete: () => {
        sub.unsubscribe()
      },
    })
  }

  private discardChanges(): void {
    this.form.reset()
  }

  get selectedAllowedLanguages() { return this.form.get('allowed')?.value as string[] }

  public save(): void {
    const newState: RemoteState = this.form.getRawValue();
    if (newState.defaultLang !== this.remoteState?.defaultLang) {
      this.service.setDefaultLanguage(newState.defaultLang).then(() => {
        this.toast.showInfo("SETTING.LANGUAGE.SAVED", true);
      }).catch(error => {
        this.toast.showError(error);
      });
    }
    if (newState.allowed.length != this.remoteState?.allowed.length ||
      newState.allowed.every((item, i) => this.remoteState?.allowed[i] === item)) {
      this.service.setDefaultLanguage(newState.defaultLang).then(() => {
        this.toast.showInfo("SETTING.LANGUAGES.SAVED", true);
      }).catch(error => {
        this.toast.showError(error);
      });
    }
  }
}
