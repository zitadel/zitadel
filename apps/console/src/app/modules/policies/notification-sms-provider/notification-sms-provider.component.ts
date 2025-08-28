import { Component, Input, OnInit } from '@angular/core';
import { AddSMSProviderTwilioRequest, UpdateSMSProviderTwilioRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { SMSProvider, SMSProviderConfigState } from 'src/app/proto/generated/zitadel/settings_pb';

import { MatDialog } from '@angular/material/dialog';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { DialogAddSMSProviderComponent } from './dialog-add-sms-provider/dialog-add-sms-provider.component';

@Component({
  selector: 'cnsl-notification-sms-provider',
  templateUrl: './notification-sms-provider.component.html',
  styleUrls: ['./notification-sms-provider.component.scss'],
})
export class NotificationSMSProviderComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  public smsProviders: SMSProvider.AsObject[] = [];

  public smsProvidersLoading: boolean = false;

  public SMSProviderConfigState: any = SMSProviderConfigState;
  public InfoSectionType: any = InfoSectionType;

  constructor(
    private service: AdminService,
    private dialog: MatDialog,
    private toast: ToastService,
  ) {}

  ngOnInit(): void {
    this.fetchData();
  }

  private fetchData(): void {
    this.smsProvidersLoading = true;
    this.service
      .listSMSProviders()
      .then((smsProviders) => {
        this.smsProvidersLoading = false;
        if (smsProviders.resultList) {
          this.smsProviders = smsProviders.resultList;
        }
      })
      .catch((error) => {
        this.smsProvidersLoading = false;
        this.toast.showError(error);
      });
  }

  public editSMSProvider(): void {
    const dialogRef = this.dialog.open(DialogAddSMSProviderComponent, {
      width: '400px',
      data: {
        smsProviders: this.smsProviders,
      },
    });

    dialogRef.afterClosed().subscribe((req: AddSMSProviderTwilioRequest | UpdateSMSProviderTwilioRequest) => {
      if (req) {
        if (!!this.twilio) {
          this.service
            .updateSMSProviderTwilio(req as UpdateSMSProviderTwilioRequest)
            .then(() => {
              this.toast.showInfo('SETTING.SMS.TWILIO.UPDATED', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else {
          this.service
            .addSMSProviderTwilio(req as AddSMSProviderTwilioRequest)
            .then(() => {
              this.toast.showInfo('SETTING.SMS.TWILIO.ADDED', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public toggleSMSProviderState(id: string): void {
    const provider = this.smsProviders.find((p) => p.id === id);
    if (provider) {
      if (provider.state === SMSProviderConfigState.SMS_PROVIDER_CONFIG_ACTIVE) {
        this.service
          .deactivateSMSProvider(id)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.DEACTIVATED', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      } else if (provider.state === SMSProviderConfigState.SMS_PROVIDER_CONFIG_INACTIVE) {
        this.service
          .activateSMSProvider(id)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.ACTIVATED', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    }
  }

  public removeSMSProvider(id: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'SETTING.SMS.REMOVEPROVIDER',
        descriptionKey: 'SETTING.SMS.REMOVEPROVIDER_DESC',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.service
          .removeSMSProvider(id)
          .then(() => {
            this.toast.showInfo('SETTING.SMS.TWILIO.REMOVED', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public get twilio(): SMSProvider.AsObject | undefined {
    return this.smsProviders.find((p) => p.twilio);
  }
}
