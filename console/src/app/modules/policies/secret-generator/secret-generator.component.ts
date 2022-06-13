import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { UpdateSecretGeneratorRequest, UpdateSecretGeneratorResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { OIDCSettings, SecretGenerator, SecretGeneratorType } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { DialogAddSecretGeneratorComponent } from './dialog-add-secret-generator/dialog-add-secret-generator.component';

@Component({
  selector: 'cnsl-secret-generator',
  templateUrl: './secret-generator.component.html',
  styleUrls: ['./secret-generator.component.scss'],
})
export class SecretGeneratorComponent implements OnInit {
  public generators: SecretGenerator.AsObject[] = [];
  public oidcSettings!: OIDCSettings.AsObject;

  public loading: boolean = false;

  public readonly AVAILABLEGENERATORS: SecretGeneratorType[] = [
    SecretGeneratorType.SECRET_GENERATOR_TYPE_INIT_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
  ];
  constructor(private service: AdminService, private toast: ToastService, private dialog: MatDialog) {}

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.service
      .listSecretGenerators()
      .then((generators) => {
        if (generators.resultList) {
          this.generators = generators.resultList;
          console.log(this.generators);
        }
      })
      .catch((error) => {
        if (error.code === 5) {
        } else {
          this.toast.showError(error);
        }
      });
  }

  private updateData(): Promise<UpdateSecretGeneratorResponse.AsObject> | void {
    const dialogRef = this.dialog.open(DialogAddSecretGeneratorComponent, {
      data: {},
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: UpdateSecretGeneratorRequest) => {
      if (req) {
        return (this.service as AdminService).updateSecretGenerator(req);
      } else {
        return;
      }
    });
  }

  public openGeneratorDialog(generatorType: SecretGeneratorType): void {
    const dialogRef = this.dialog.open(DialogAddSecretGeneratorComponent, {
      data: {
        type: generatorType,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((req: UpdateSecretGeneratorRequest) => {
      if (req) {
        return (this.service as AdminService)
          .updateSecretGenerator(req)
          .then(() => {
            this.toast.showInfo('SETTING.SECRETS.UPDATED', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      } else {
        return;
      }
    });
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
}
