import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { UpdateSecretGeneratorResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { OIDCSettings, SecretGenerator } from 'src/app/proto/generated/zitadel/settings_pb';
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
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'IDP.DELETE_TITLE',
        descriptionKey: 'IDP.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    // dialogRef.afterClosed().subscribe((req: UpdateSecretGeneratorRequest) => {
    //   if (req) {
    //     return (this.service as AdminService).updateSecretGenerator(req);
    //   } else {
    //     return;
    //   }
    // });
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
