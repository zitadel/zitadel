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

enum Preview {
  CURRENT,
  PREVIEW,
}

@Component({
  selector: 'app-private-labeling-policy',
  templateUrl: './private-labeling-policy.component.html',
  styleUrls: ['./private-labeling-policy.component.scss'],
})
export class PrivateLabelingPolicyComponent implements OnDestroy {
  public theme: Theme = Theme.LIGHT;
  public preview: Preview = Preview.PREVIEW;

  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  public service!: ManagementService | AdminService;

  public previewData!: LabelPolicy.AsObject;
  public data!: LabelPolicy.AsObject;

  public panelOpenState: boolean = false;
  public isHoveringOverDarkLogo: boolean = false;
  public isHoveringOverDarkIcon: boolean = false;
  public isHoveringOverLightLogo: boolean = false;
  public isHoveringOverLightIcon: boolean = false;
  public isHoveringOverFont: boolean = false;

  private sub: Subscription = new Subscription();
  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public loading: boolean = false;
  public nextLinks: CnslLinks[] = [
    IAM_COMPLEXITY_LINK,
    IAM_POLICY_LINK,
    IAM_LOGIN_POLICY_LINK,
  ];

  public logoFile!: File;
  public fontFile!: File;
  public logoURL: any = '';

  public Theme: any = Theme;
  public Preview: any = Preview;

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

  public toggleHoverFont(isHovering: boolean): void {
    this.isHoveringOverFont = isHovering;
  }

  public onDropLogo(theme: Theme, filelist: FileList): Promise<any> | void {
    const file = filelist.item(0);
    if (file) {
      this.logoFile = file;

      const reader = new FileReader();
      reader.readAsDataURL(this.logoFile);
      reader.onload = (event) => { // called once readAsDataURL is completed
        console.log(event.target?.result);
        this.logoURL = event.target?.result;

        const formData = new FormData();
        formData.append('file', file);
        if (theme === Theme.DARK) {
          switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
              return this.uploadService.upload(UploadEndpoint.MGMTDARKLOGO, formData);
            case PolicyComponentServiceType.ADMIN:
              return this.uploadService.upload(UploadEndpoint.IAMDARKLOGO, formData);
          }
        }
        if (theme === Theme.LIGHT) {
          switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
              return this.uploadService.upload(UploadEndpoint.MGMTDARKLOGO, formData);
            case PolicyComponentServiceType.ADMIN:
              return this.uploadService.upload(UploadEndpoint.IAMDARKLOGO, formData);
          }
        }
      };
    }
  }

  public onDropFont(filelist: FileList): Promise<any> | void {
    const file = filelist.item(0);
    if (file) {
      this.fontFile = file;

      // const reader = new FileReader();
      // reader.readAsDataURL(this.fontFile);
      // reader.onload = (event) => { // called once readAsDataURL is completed
      //   console.log(event.target?.result);

      const formData = new FormData();
      formData.append('file', file);
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          return this.uploadService.upload(UploadEndpoint.MGMTFONT, formData);
        case PolicyComponentServiceType.ADMIN:
          return this.uploadService.upload(UploadEndpoint.IAMFONT, formData);
      }
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

  public changeColor(attrToSet: string, valueToSet: string) {
    // attrToSet = valueToSet;
    this.savePolicy();
  }

  public onDropIcon(theme: Theme, filelist: FileList): Promise<any> | void {
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
        formData.append('file', file);
        if (theme === Theme.DARK) {
          switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
              return this.uploadService.upload(UploadEndpoint.MGMTDARKICON, formData);
            case PolicyComponentServiceType.ADMIN:
              return this.uploadService.upload(UploadEndpoint.IAMDARKICON, formData);
          }
        }
        if (theme === Theme.LIGHT) {
          switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
              return this.uploadService.upload(UploadEndpoint.MGMTLIGHTICON, formData);
            case PolicyComponentServiceType.ADMIN:
              return this.uploadService.upload(UploadEndpoint.IAMLIGHTICON, formData);
          }
        }
      };
    }
  }

  public fetchData(): void {
    this.loading = true;

    this.getPreviewData().then(data => {
      console.log('preview', data);

      if (data.policy) {
        this.previewData = data.policy;
        this.loading = false;
      }
    }).catch(error => {
      this.toast.showError(error);
    });

    this.getData().then(data => {
      console.log('data', data);

      if (data.policy) {
        this.data = data.policy;
        this.loading = false;
      }
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  private async getPreviewData():
    Promise<MgmtGetPreviewLabelPolicyResponse.AsObject |
      AdminGetPreviewLabelPolicyResponse.AsObject |
      MgmtGetLabelPolicyResponse.AsObject |
      AdminGetLabelPolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getPreviewLabelPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getPreviewLabelPolicy();
    }
  }

  private async getData():
    Promise<MgmtGetPreviewLabelPolicyResponse.AsObject |
      AdminGetPreviewLabelPolicyResponse.AsObject |
      MgmtGetLabelPolicyResponse.AsObject |
      AdminGetLabelPolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getLabelPolicy();
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
        if ((this.previewData as LabelPolicy.AsObject).isDefault) {
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
        (this.service as AdminService).updateLabelPolicy(req).then(() => {
          this.toast.showInfo('POLICY.TOAST.SET', true);
        }).catch(error => {
          this.toast.showError(error);
        });
        break;
    }
  }

  public get isDefault(): boolean {
    if (this.previewData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.previewData as LabelPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }

  public overwriteValues(req: AddCustomLabelPolicyRequest | UpdateCustomLabelPolicyRequest): void {
    req.setPrimaryColorDark(this.previewData.primaryColorDark);
    req.setPrimaryColor(this.previewData.primaryColor);
    req.setWarnColorDark(this.previewData.warnColorDark);
    req.setWarnColor(this.previewData.warnColor);

    req.setDisableWatermark(this.previewData.disableWatermark);
    req.setHideLoginNameSuffix(this.previewData.hideLoginNameSuffix);
  }

  public activatePolicy(): Promise<any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).activateCustomLabelPolicy().then(() => {
          this.toast.showInfo('POLICY.PRIVATELABELING.ACTIVATED', true);
        }).catch(error => {
          this.toast.showError(error);
        });
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).activateLabelPolicy().then(() => {
          this.toast.showInfo('POLICY.PRIVATELABELING.ACTIVATED', true);
        }).catch(error => {
          this.toast.showError(error);
        });
    }
  }

  public resetPolicy(): Promise<any> {
    return (this.service as ManagementService).resetLabelPolicyToDefault().then(() => {
      this.toast.showInfo('POLICY.PRIVATELABELING.RESET', true);
    }).catch(error => {
      this.toast.showError(error);
    });
  }
}
