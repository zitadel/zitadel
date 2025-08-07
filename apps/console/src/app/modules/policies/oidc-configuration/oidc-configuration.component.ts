import { Component, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { take } from 'rxjs';
import {
  AddOIDCSettingsRequest,
  AddOIDCSettingsResponse,
  UpdateOIDCSettingsRequest,
  UpdateOIDCSettingsResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { OIDCSettings } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

@Component({
  selector: 'cnsl-oidc-configuration',
  templateUrl: './oidc-configuration.component.html',
  styleUrls: ['./oidc-configuration.component.scss'],
})
export class OIDCConfigurationComponent implements OnInit {
  public oidcSettings!: OIDCSettings.AsObject;
  private settingsSet: boolean = false;

  public loading: boolean = false;
  public form!: UntypedFormGroup;
  constructor(
    private service: AdminService,
    private fb: UntypedFormBuilder,
    private toast: ToastService,
    private authService: GrpcAuthService,
  ) {
    this.form = this.fb.group({
      accessTokenLifetime: [{ disabled: true }, [requiredValidator]],
      idTokenLifetime: [{ disabled: true }, [requiredValidator]],
      refreshTokenExpiration: [{ disabled: true }, [requiredValidator]],
      refreshTokenIdleExpiration: [{ disabled: true }, [requiredValidator]],
    });
  }

  ngOnInit(): void {
    this.fetchData();
    this.authService
      .isAllowed(['iam.write'])
      .pipe(take(1))
      .subscribe((allowed) => {
        if (allowed) {
          this.form.enable();
        }
      });
  }

  private fetchData(): void {
    this.service
      .getOIDCSettings()
      .then((oidcConfiguration) => {
        if (oidcConfiguration.settings) {
          this.oidcSettings = oidcConfiguration.settings;
          this.settingsSet = true;

          this.accessTokenLifetime?.setValue(
            oidcConfiguration.settings.accessTokenLifetime?.seconds
              ? oidcConfiguration.settings.accessTokenLifetime?.seconds / 60 / 60
              : 0,
          );
          this.idTokenLifetime?.setValue(
            oidcConfiguration.settings.idTokenLifetime?.seconds
              ? oidcConfiguration.settings.idTokenLifetime?.seconds / 60 / 60
              : 0,
          );
          this.refreshTokenExpiration?.setValue(
            oidcConfiguration.settings.refreshTokenExpiration?.seconds
              ? oidcConfiguration.settings.refreshTokenExpiration?.seconds / 60 / 60 / 24
              : 0,
          );
          this.refreshTokenIdleExpiration?.setValue(
            oidcConfiguration.settings.refreshTokenIdleExpiration?.seconds
              ? oidcConfiguration.settings.refreshTokenIdleExpiration?.seconds / 60 / 60 / 24
              : 0,
          );
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  private updateData(): Promise<UpdateOIDCSettingsResponse.AsObject> {
    const req = new UpdateOIDCSettingsRequest();

    const accessToken = new Duration().setSeconds((this.accessTokenLifetime?.value ?? 0) * 60 * 60);
    req.setAccessTokenLifetime(accessToken);

    const idToken = new Duration().setSeconds((this.idTokenLifetime?.value ?? 0) * 60 * 60);
    req.setIdTokenLifetime(idToken);

    const refreshToken = new Duration().setSeconds((this.refreshTokenExpiration?.value ?? 0) * 60 * 60 * 24);
    req.setRefreshTokenExpiration(refreshToken);

    const refreshIdleToken = new Duration().setSeconds((this.refreshTokenIdleExpiration?.value ?? 0) * 60 * 60 * 24);
    req.setRefreshTokenIdleExpiration(refreshIdleToken);

    return (this.service as AdminService).updateOIDCSettings(req);
  }

  private addData(): Promise<AddOIDCSettingsResponse.AsObject> {
    const req = new AddOIDCSettingsRequest();

    const accessToken = new Duration().setSeconds((this.accessTokenLifetime?.value ?? 0) * 60 * 60);
    req.setAccessTokenLifetime(accessToken);

    const idToken = new Duration().setSeconds((this.idTokenLifetime?.value ?? 0) * 60 * 60);
    req.setIdTokenLifetime(idToken);

    const refreshToken = new Duration().setSeconds((this.refreshTokenExpiration?.value ?? 0) * 60 * 60 * 24);
    req.setRefreshTokenExpiration(refreshToken);

    const refreshIdleToken = new Duration().setSeconds((this.refreshTokenIdleExpiration?.value ?? 0) * 60 * 60 * 24);
    req.setRefreshTokenIdleExpiration(refreshIdleToken);

    return (this.service as AdminService).addOIDCSettings(req);
  }

  public savePolicy(): void {
    if (this.settingsSet) {
      this.updateData()
        .then(() => {
          this.toast.showInfo('SETTING.SMTP.SAVED', true);
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else {
      this.addData()
        .then(() => {
          this.toast.showInfo('SETTING.SMTP.SAVED', true);
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
