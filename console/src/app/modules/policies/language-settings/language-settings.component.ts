import { Component, OnInit } from '@angular/core';
import { SetDefaultLanguageResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import {AbstractControl, UntypedFormBuilder, UntypedFormGroup, Validators} from "@angular/forms";
import {LanguagesService} from "../../../services/languages.service";
import {forkJoin, from, Observable} from "rxjs";
import {GrpcAuthService} from "../../../services/grpc-auth.service";
import {i18nValidator} from "../../form-field/validators/validators";

@Component({
  selector: 'cnsl-language-settings',
  templateUrl: './language-settings.component.html',
  styleUrls: ['./language-settings.component.scss'],
})
export class LanguageSettingsComponent implements OnInit {

  public form!: UntypedFormGroup;

  public canWriteRestrictions$: Observable<boolean> = this.authService.isAllowed(["iam.restrictions.write"]);
  public canWriteDefaultLanguage$: Observable<boolean> = this.authService.isAllowed(["iam.write"]);

  public allowed$: Observable<string[]>;
  public notAllowed$: Observable<string[]>;
  private remoteState: { defaultLang: string, allowed: string[] } | null = null;

  public loading: boolean = false;
  constructor(
    private service: AdminService,
    private toast: ToastService,
    private fb: UntypedFormBuilder,
    private languagesSvc: LanguagesService,
    private authService: GrpcAuthService,
  ) {
    this.form = this.fb.group({
      defaultLang: ['', [i18nValidator(
        "SETTING.LANGUAGES.DEFAULT_MUST_BE_ALLOWED",
        (control: AbstractControl) => control.parent?.get('allowed')?.value.contains(control.value))
      ]],
      allowed: [[''], []],
    });
    const allowed$ = this.languagesSvc.allowedLanguages(this.service);
    const notAllowed$: Observable<string[]> = this.languagesSvc.notAllowedLanguages(this.service, allowed$);
    const defaultLang$ = from(this.service.getDefaultLanguage());
    forkJoin([allowed$, notAllowed$, defaultLang$]).subscribe(([allowed, notAllowed, { language: defaultLang }]) => {
      this.remoteState = { defaultLang, allowed, };
      this.form.setValue(this.remoteState);
    }).unsubscribe();
  }

  private discardChanges(): void {
    this.form.reset()
  }

  public save(): void {
    const newValue = this.form.getRawValue();
    if (newValue.defaultLang !== this.remoteState?.defaultLang) {
      this.service.setDefaultLanguage(newValue.defaultLang).then(() => {
        this.toast.showInfo("POLICY.LANGUAGE.SAVED", true);
        this.fetchData();
      }).catch(error => {
        this.toast.showError(error);
      });
    }

    const prom = this.updateData();
    this.loading = true;
    if (prom) {
      prom
        .then(() => {
          this.toast.showInfo('POLICY.LOGIN_POLICY.SAVED', true);
          this.loading = false;
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.loading = false;
          this.toast.showError(error);
        });
    }
  }
}
