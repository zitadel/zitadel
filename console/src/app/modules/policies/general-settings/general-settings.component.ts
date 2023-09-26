import { Component, OnInit } from '@angular/core';
import { SetDefaultLanguageResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-general-settings',
  templateUrl: './general-settings.component.html',
  styleUrls: ['./general-settings.component.scss'],
})
export class GeneralSettingsComponent implements OnInit {
  public defaultLanguage: string = '';
  public defaultLanguageOptions: string[] = [];

  public loading: boolean = false;
  constructor(
    private service: AdminService,
    private toast: ToastService,
  ) {}

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.service.getDefaultLanguage().then((langResp) => {
      this.defaultLanguage = langResp.language;
    });
    this.service.getSupportedLanguages().then((supportedResp) => {
      this.defaultLanguageOptions = supportedResp.languagesList;
    });
  }

  private updateData(): Promise<SetDefaultLanguageResponse.AsObject> {
    return (this.service as AdminService).setDefaultLanguage(this.defaultLanguage);
  }

  public savePolicy(): void {
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
