import { HttpErrorResponse } from '@angular/common/http';
import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
  GetLabelPolicyResponse as AdminGetLabelPolicyResponse,
  GetPreviewLabelPolicyResponse as AdminGetPreviewLabelPolicyResponse,
  UpdateLabelPolicyRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  AddCustomLabelPolicyRequest,
  GetLabelPolicyResponse as MgmtGetLabelPolicyResponse,
  GetPreviewLabelPolicyResponse as MgmtGetPreviewLabelPolicyResponse,
  UpdateCustomLabelPolicyRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { UploadEndpoint, UploadService } from 'src/app/services/upload.service';

import { CnslLinks } from '../../links/links.component';
import { IAM_COMPLEXITY_LINK, IAM_LOGIN_POLICY_LINK, IAM_POLICY_LINK } from '../../policy-grid/policy-links';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

enum Theme {
  DARK,
  LIGHT,
}

@Component({
  selector: 'app-private-labeling-policy',
  templateUrl: './private-labeling-policy.component.html',
  styleUrls: ['./private-labeling-policy.component.scss'],
})
export class PrivateLabelingPolicyComponent implements OnDestroy {
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  public service!: ManagementService | AdminService;

  public data!: LabelPolicy.AsObject;
  public panelOpenState: boolean = false;
  public isHoveringOverDarkLogo: boolean = false;
  public isHoveringOverDarkIcon: boolean = false;
  public isHoveringOverLightLogo: boolean = false;
  public isHoveringOverLightIcon: boolean = false;

  private sub: Subscription = new Subscription();
  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public loading: boolean = false;
  public nextLinks: CnslLinks[] = [
    IAM_COMPLEXITY_LINK,
    IAM_POLICY_LINK,
    IAM_LOGIN_POLICY_LINK,
  ];

  public logoFile!: File;
  public logoURL: any = '';

  public font: string = 'Lato';
  public fontCssRule: string = '\'Lato\', sans-serif';
  public fonts: Array<{ name: string; css: string; }> = [
    { name: 'Lato', css: '\'Lato\', sans-serif' },
    { name: 'Merriweather', css: '\'Merriweather\', sans-serif' },
    { name: 'System', css: 'ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji;' },
  ];

  public colors: Array<{ name: string; color: string; }> = [
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
    { name: 'indigo', color: '#3f51b5' },
    { name: 'blue', color: '#2196f3' },
    { name: 'lightblue', color: '#03a9f4' },
    { name: 'cyan', color: '#00bcd4' },
    { name: 'teal', color: '#009688' },
    { name: 'green', color: '#4caf50' },
    { name: 'lightgreen', color: '#8bc34a' },
    { name: 'lime', color: '#cddc39' },
    { name: 'yellow', color: '#ffeb3b' },
    { name: 'amber', color: '#ffc107' },
    { name: 'orange', color: '#fb8c00' },
    { name: 'deeporange', color: '#ff5722' },
    { name: 'brown', color: '#795548' },
    { name: 'grey', color: '#9e9e9e' },
    { name: 'bluegrey', color: '#607d8b' },
    { name: 'white', color: '#ffffff' },
  ];

  public warncolors: Array<{ name: string; color: string; }> = [
    { name: 'red', color: '#f44336' },
    { name: 'pink', color: '#e91e63' },
    { name: 'purple', color: '#9c27b0' },
    { name: 'deeppurple', color: '#673ab7' },
  ];

  public Theme: any = Theme;

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private uploadService: UploadService,
  ) {
    this.sub = this.route.data.pipe(switchMap(data => {
      this.serviceType = data.serviceType;

      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.service = this.injector.get(ManagementService as Type<ManagementService>);
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);
          break;
      }

      return this.route.params;
    })).subscribe(() => {
      this.fetchData();
    });
  }

  public toggleHoverLogo(theme: Theme, isHovering: boolean): void {
    if (theme === Theme.DARK) {
      this.isHoveringOverDarkLogo = isHovering;
    }
    if (theme === Theme.LIGHT) {
      this.isHoveringOverLightLogo = isHovering;
    }
  }

  public onDropLogo(theme: Theme, filelist: FileList): void {
    const file = filelist.item(0);
    if (file) {
      console.log(filelist.item(0));
      this.logoFile = file;

      const reader = new FileReader();
      reader.readAsDataURL(this.logoFile);
      reader.onload = (event) => { // called once readAsDataURL is completed
        console.log(event.target?.result);
        this.logoURL = event.target?.result;

        const formData = new FormData();
        formData.append('file', this.logoURL);
        if (theme === Theme.DARK) {
          this.uploadService.upload(UploadEndpoint.DARKLOGO, formData);
        }
        if (theme === Theme.LIGHT) {
          this.uploadService.upload(UploadEndpoint.LIGHTLOGO, formData);
        }
      };
    }
  }

  public toggleHoverIcon(theme: Theme, isHovering: boolean): void {
    if (theme === Theme.DARK) {
      this.isHoveringOverDarkIcon = isHovering;
    }
    if (theme === Theme.LIGHT) {
      this.isHoveringOverLightIcon = isHovering;
    }
  }

  public onDropIcon(theme: Theme, filelist: FileList): void {
    console.log(filelist);
    const file = filelist.item(0);
    if (file) {
      console.log(filelist.item(0));
      this.logoFile = file;

      const reader = new FileReader();
      reader.readAsDataURL(this.logoFile);
      reader.onload = (event) => { // called once readAsDataURL is completed
        console.log(event.target?.result);
        this.logoURL = event.target?.result;

        const formData = new FormData();
        formData.append('file', this.logoURL);
        if (theme === Theme.DARK) {
          this.uploadService.upload(UploadEndpoint.DARKLOGO, formData);
        }
        if (theme === Theme.LIGHT) {
          this.uploadService.upload(UploadEndpoint.LIGHTLOGO, formData);
        }
      };
    }
  }

  public fetchData(): void {
    this.loading = true;

    this.getData().then(data => {
      console.log(data);

      if (data.policy) {
        this.data = data.policy;
        this.loading = false;
      }
    });
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  private async getData():
    Promise<MgmtGetPreviewLabelPolicyResponse.AsObject |
      AdminGetPreviewLabelPolicyResponse.AsObject |
      MgmtGetLabelPolicyResponse.AsObject |
      AdminGetLabelPolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getPreviewLabelPolicy(); // .getLabelPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getLabelPolicy();
    }
  }

  public removePolicy(): void {
    if (this.service instanceof ManagementService) {
      this.service.resetPasswordComplexityPolicyToDefault().then(() => {
        this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
        setTimeout(() => {
          this.fetchData();
        }, 1000);
      }).catch(error => {
        this.toast.showError(error);
      });
    }
  }

  public savePolicy(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if ((this.data as LabelPolicy.AsObject).isDefault) {
          const req0 = new AddCustomLabelPolicyRequest();
          this.overwriteValues(req0);

          (this.service as ManagementService).addCustomLabelPolicy(req0).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          }).catch((error: HttpErrorResponse) => {
            this.toast.showError(error);
          });
        } else {
          const req1 = new UpdateCustomLabelPolicyRequest();
          this.overwriteValues(req1);

          (this.service as ManagementService).updateCustomLabelPolicy(req1).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          }).catch(error => {
            this.toast.showError(error);
          });
        }
        break;
      case PolicyComponentServiceType.ADMIN:
        const req = new UpdateLabelPolicyRequest();
        this.overwriteValues(req);
        console.log(req.toObject());
        (this.service as AdminService).updateLabelPolicy(req).then(() => {
          this.toast.showInfo('POLICY.TOAST.SET', true);
        }).catch(error => {
          this.toast.showError(error);
        });
        break;
    }
  }

  public get isDefault(): boolean {
    if (this.data && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.data as LabelPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }

  public overwriteValues(req: AddCustomLabelPolicyRequest | UpdateCustomLabelPolicyRequest): void {
    req.setPrimaryColor(this.data.primaryColorDark);
    req.setPrimaryColor(this.data.primaryColor);
    req.setWarnColorDark(this.data.warnColorDark);
    req.setWarnColor(this.data.warnColor);

    req.setDisableWatermark(this.data.disableWatermark);
    req.setHideLoginNameSuffix(this.data.hideLoginNameSuffix);
    req.setErrorMsgPopup(this.data.errorMsgPopup);
  }
}
