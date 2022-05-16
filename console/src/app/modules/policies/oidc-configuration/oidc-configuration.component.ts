import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { SetDefaultLanguageResponse, UpdateOIDCSettingsRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { OIDCSettings } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-oidc-configuration',
  templateUrl: './oidc-configuration.component.html',
  styleUrls: ['./oidc-configuration.component.scss'],
})
export class OIDCConfigurationComponent implements OnInit {
  public oidcSettings!: OIDCSettings.AsObject;

  public loading: boolean = false;
  public form!: FormGroup;
  constructor(private service: AdminService, private fb: FormBuilder, private toast: ToastService) {
    this.form = this.fb.group({
      accessTokenLifetime: [12, [Validators.required]],
      idTokenLifetime: [12, [Validators.required]],
      refreshTokenExpiration: [30, [Validators.required]],
      refreshTokenIdleExpiration: [90, [Validators.required]],
    });
  }

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.service
      .getOIDCSettings()
      .then((oidcConfiguration) => {
        if (oidcConfiguration.settings) {
          this.oidcSettings = oidcConfiguration.settings;

          this.accessTokenLifetime?.setValue(
            oidcConfiguration.settings.accessTokenLifetime?.seconds
              ? oidcConfiguration.settings.accessTokenLifetime?.seconds / 60 / 60
              : 12,
          );
          this.idTokenLifetime?.setValue(
            oidcConfiguration.settings.idTokenLifetime?.seconds
              ? oidcConfiguration.settings.idTokenLifetime?.seconds / 60 / 60
              : 12,
          );
          this.refreshTokenExpiration?.setValue(
            oidcConfiguration.settings.refreshTokenExpiration?.seconds
              ? oidcConfiguration.settings.refreshTokenExpiration?.seconds / 60 / 60 / 24
              : 30,
          );
          this.refreshTokenIdleExpiration?.setValue(
            oidcConfiguration.settings.refreshTokenIdleExpiration?.seconds
              ? oidcConfiguration.settings.refreshTokenIdleExpiration?.seconds / 60 / 60 / 24
              : 90,
          );
        }
      })
      .catch((error) => {
        if (error.code === 5) {
        } else {
          this.toast.showError(error);
        }
      });
  }

  private updateData(): Promise<SetDefaultLanguageResponse.AsObject> {
    const req = new UpdateOIDCSettingsRequest();

    const accessToken = new Duration().setSeconds((this.accessTokenLifetime?.value ?? 12) * 60 * 60);
    req.setAccessTokenLifetime(accessToken);

    const idToken = new Duration().setSeconds((this.idTokenLifetime?.value ?? 12) * 60 * 60);
    req.setIdTokenLifetime(idToken);

    const refreshToken = new Duration().setSeconds((this.refreshTokenExpiration?.value ?? 30) * 60 * 60 * 24);
    req.setRefreshTokenExpiration(refreshToken);

    const refreshIdleToken = new Duration().setSeconds((this.refreshTokenIdleExpiration?.value ?? 90) * 60 * 60 * 24);
    req.setRefreshTokenIdleExpiration(refreshIdleToken);

    return (this.service as AdminService).updateOIDCSettings(req);
  }

  public savePolicy(): void {
    const prom = this.updateData();
    if (prom) {
      prom
        .then(() => {
          this.toast.showInfo('SETTING.SMTP.SAVED', true);
          this.loading = true;
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public get accessTokenLifetime(): AbstractControl | null {
    return this.form.get('accessTokenLifetime');
  }

  public get idTokenLifetime(): AbstractControl | null {
    return this.form.get('idTokenLifetime');
  }

  public get refreshTokenExpiration(): AbstractControl | null {
    return this.form.get('refreshTokenExpiration');
  }

  public get refreshTokenIdleExpiration(): AbstractControl | null {
    return this.form.get('refreshTokenIdleExpiration');
  }
}
