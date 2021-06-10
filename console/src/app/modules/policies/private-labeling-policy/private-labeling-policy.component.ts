import { HttpErrorResponse } from '@angular/common/http';
import { Component, EventEmitter, Injector, OnDestroy, Type } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
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
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { AssetEndpoint, AssetService, AssetType } from 'src/app/services/asset.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { CnslLinks } from '../../links/links.component';
import { IAM_COMPLEXITY_LINK, IAM_LOGIN_POLICY_LINK, IAM_POLICY_LINK } from '../../policy-grid/policy-links';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

export enum Theme {
  DARK,
  LIGHT,
}

export enum Preview {
  CURRENT,
  PREVIEW,
}

export enum ColorType {
  BACKGROUND,
  PRIMARY,
  WARN,
  FONTDARK,
  FONTLIGHT,
  BACKGROUNDDARK,
  BACKGROUNDLIGHT,
}

const ORG_STORAGE_KEY = 'organization';

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

  public images: { [key: string]: any; } = {};

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

  public Theme: any = Theme;
  public Preview: any = Preview;
  public ColorType: any = ColorType;
  public AssetType: any = AssetType;

  public refreshPreview: EventEmitter<void> = new EventEmitter();
  public loadingImages: boolean = false;
  private org!: Org.AsObject;

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private injector: Injector,
    private assetService: AssetService,
    private sanitizer: DomSanitizer,
    private storageService: StorageService,
  ) {
    const org: Org.AsObject | null = (this.storageService.getItem(ORG_STORAGE_KEY));

    if (org) {
      this.org = org;
    }

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

      const formData = new FormData();
      formData.append('file', file);
      if (theme === Theme.DARK) {
        switch (this.serviceType) {
          case PolicyComponentServiceType.MGMT:
            return this.handleUploadPromise(this.assetService.upload(AssetEndpoint.MGMTDARKLOGO, formData, this.org.id));
          case PolicyComponentServiceType.ADMIN:
            return this.handleUploadPromise(this.assetService.upload(AssetEndpoint.IAMDARKLOGO, formData, this.org.id));
        }
      }
      if (theme === Theme.LIGHT) {
        switch (this.serviceType) {
          case PolicyComponentServiceType.MGMT:
            return this.handleUploadPromise(this.assetService.upload(AssetEndpoint.MGMTLOGO, formData, this.org.id));
          case PolicyComponentServiceType.ADMIN:
            return this.handleUploadPromise(this.assetService.upload(AssetEndpoint.IAMLOGO, formData, this.org.id));
        }
      }

    }
  }

  public onDropFont(filelist: FileList): Promise<any> | void {
    const file = filelist.item(0);
    if (file) {
      const formData = new FormData();
      formData.append('file', file);
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          return this.handleFontUploadPromise(this.assetService.upload(AssetEndpoint.MGMTFONT, formData, this.org.id));
        case PolicyComponentServiceType.ADMIN:
          return this.handleFontUploadPromise(this.assetService.upload(AssetEndpoint.IAMFONT, formData, this.org.id));
      }
    }
  }

  public deleteFont(): Promise<any> {
    const handler = (prom: Promise<any>) => prom.then(() => {
      this.toast.showInfo('POLICY.TOAST.DELETESUCCESS', true);
      setTimeout(() => {
        this.loadingImages = true;
        this.getPreviewData().then(data => {

          if (data.policy) {
            this.previewData = data.policy;
            this.loadPreviewImages();
          }
        });
      }, 1000);
    }).catch(error => this.toast.showError(error));

    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return handler((this.service as ManagementService).removeLabelPolicyFont());
      case PolicyComponentServiceType.ADMIN:
        return handler((this.service as AdminService).removeLabelPolicyFont());
    }
  }

  public deleteAsset(type: AssetType, theme: Theme): any {
    const previewHandler = (prom: Promise<any>) => {
      return prom.then(() => {
        this.toast.showInfo('POLICY.TOAST.DELETESUCCESS', true);
        setTimeout(() => {
          this.loadingImages = true;
          this.getPreviewData().then(data => {

            if (data.policy) {
              this.previewData = data.policy;
              this.loadPreviewImages();
            }
          });
        }, 1000);
      }).catch(error => this.toast.showError(error));
    };

    switch (this.serviceType) {
      case PolicyComponentServiceType.ADMIN:
        if (type === AssetType.LOGO) {
          if (theme === Theme.DARK) {
            return previewHandler(this.service.removeLabelPolicyLogoDark());
          } else if (theme === Theme.LIGHT) {
            return previewHandler(this.service.removeLabelPolicyLogo());
          }
        } else if (type === AssetType.ICON) {
          if (theme === Theme.DARK) {
            return previewHandler(this.service.removeLabelPolicyIconDark());
          } else if (theme === Theme.LIGHT) {
            return previewHandler(this.service.removeLabelPolicyIcon());
          }
        }
        break;
      case PolicyComponentServiceType.MGMT:
        if (type === AssetType.LOGO) {
          if (theme === Theme.DARK) {
            return previewHandler(this.service.removeLabelPolicyLogoDark());
          } else if (theme === Theme.LIGHT) {
            return previewHandler(this.service.removeLabelPolicyLogo());
          }
        } else if (type === AssetType.ICON) {
          if (theme === Theme.DARK) {
            return previewHandler(this.service.removeLabelPolicyIconDark());
          } else if (theme === Theme.LIGHT) {
            return previewHandler(this.service.removeLabelPolicyIcon());
          }
        }
        break;
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
    const file = filelist.item(0);
    if (file) {

      const formData = new FormData();
      formData.append('file', file);
      if (theme === Theme.DARK) {
        switch (this.serviceType) {
          case PolicyComponentServiceType.MGMT:
            this.handleUploadPromise(this.assetService.upload(AssetEndpoint.MGMTDARKICON, formData, this.org.id));
            break;
          case PolicyComponentServiceType.ADMIN:
            this.handleUploadPromise(this.assetService.upload(AssetEndpoint.IAMDARKICON, formData, this.org.id));
            break;
        }
      }
      if (theme === Theme.LIGHT) {
        switch (this.serviceType) {
          case PolicyComponentServiceType.MGMT:
            this.handleUploadPromise(this.assetService.upload(AssetEndpoint.MGMTICON, formData, this.org.id));
            break;
          case PolicyComponentServiceType.ADMIN:
            this.handleUploadPromise(this.assetService.upload(AssetEndpoint.IAMICON, formData, this.org.id));
            break;
        }
      }
    }
  }

  private handleFontUploadPromise(task: Promise<any>): Promise<any> {
    return task.then(() => {
      this.toast.showInfo('POLICY.TOAST.UPLOADSUCCESS', true);
      setTimeout(() => {
        this.getPreviewData().then(data => {
          if (data.policy) {
            this.previewData = data.policy;
          }
        });
      }, 1000);
    }).catch(error => this.toast.showError(error));
  }

  private handleUploadPromise(task: Promise<any>): Promise<any> {
    const enhTask = task.then(() => {
      this.toast.showInfo('POLICY.TOAST.UPLOADSUCCESS', true);
      setTimeout(() => {
        this.loadingImages = true;
        this.getPreviewData().then(data => {

          if (data.policy) {
            this.previewData = data.policy;
            this.loadPreviewImages();
          }
        });
      }, 1000);
    }).catch(error => this.toast.showError(error));

    if (this.serviceType == PolicyComponentServiceType.MGMT && ((this.previewData as LabelPolicy.AsObject).isDefault)) {
      return this.savePolicy().then(() => enhTask);
    } else {
      return enhTask;
    }
  }

  public fetchData(): void {
    this.loading = true;

    this.getPreviewData().then(data => {
      console.log('preview', data);
      this.loadingImages = true;

      if (data.policy) {
        this.previewData = data.policy;
        this.loading = false;

        this.loadPreviewImages();
      }
    }).catch(error => {
      this.toast.showError(error);
    });

    this.getData().then(data => {
      console.log('data', data);

      if (data.policy) {
        this.data = data.policy;
        this.loading = false;

        this.loadImages();
      }
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  private loadImages(): void {
    const promises: Promise<any>[] = [];
    if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      if (this.data.logoUrlDark) {
        promises.push(this.loadAsset('darkLogo', AssetEndpoint.IAMDARKLOGO));
      }
      if (this.data.iconUrlDark) {
        promises.push(this.loadAsset('darkIcon', AssetEndpoint.IAMDARKICON));
      }
      if (this.data.logoUrl) {
        promises.push(this.loadAsset('logo', AssetEndpoint.IAMLOGO));
      }
      if (this.data.iconUrl) {
        promises.push(this.loadAsset('icon', AssetEndpoint.IAMICON));
      }
    } else if (this.serviceType === PolicyComponentServiceType.MGMT) {
      if (this.data.logoUrlDark) {
        promises.push(this.loadAsset('darkLogo', AssetEndpoint.MGMTDARKLOGO));
      }
      if (this.data.iconUrlDark) {
        promises.push(this.loadAsset('darkIcon', AssetEndpoint.MGMTDARKICON));
      }
      if (this.data.logoUrl) {
        promises.push(this.loadAsset('logo', AssetEndpoint.MGMTLOGO));
      }
      if (this.data.iconUrl) {
        promises.push(this.loadAsset('icon', AssetEndpoint.MGMTICON));
      }
    }

    if (promises.length) {
      Promise.all(promises).then(() => {
        this.loadingImages = false;
      }).catch(error => {
        this.loadingImages = false;
      });
    } else {
      this.loadingImages = false;
    }
  }

  private loadPreviewImages(): void {
    const promises: Promise<any>[] = [];

    if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      if (this.previewData.logoUrlDark) {
        promises.push(this.loadAsset('previewDarkLogo', AssetEndpoint.IAMDARKLOGOPREVIEW));
      }
      if (this.previewData.iconUrlDark) {
        promises.push(this.loadAsset('previewDarkIcon', AssetEndpoint.IAMDARKICONPREVIEW));
      }
      if (this.previewData.logoUrl) {
        promises.push(this.loadAsset('previewLogo', AssetEndpoint.IAMLOGOPREVIEW));
      }
      if (this.previewData.iconUrl) {
        promises.push(this.loadAsset('previewIcon', AssetEndpoint.IAMICONPREVIEW));
      }
    } else if (this.serviceType === PolicyComponentServiceType.MGMT) {
      if (this.previewData.logoUrlDark) {
        promises.push(this.loadAsset('previewDarkLogo', AssetEndpoint.MGMTDARKLOGOPREVIEW));
      }
      if (this.previewData.iconUrlDark) {
        promises.push(this.loadAsset('previewDarkIcon', AssetEndpoint.MGMTDARKICONPREVIEW));
      }
      if (this.previewData.logoUrl) {
        promises.push(this.loadAsset('previewLogo', AssetEndpoint.MGMTLOGOPREVIEW));
      }
      if (this.previewData.iconUrl) {
        promises.push(this.loadAsset('previewIcon', AssetEndpoint.MGMTICONPREVIEW));
      }
    }

    if (promises.length) {
      Promise.all(promises).then(() => {
        this.loadingImages = false;
      }).catch(error => {
        this.loadingImages = false;
      });
    } else {
      this.loadingImages = false;
    }
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

  private loadAsset(imagekey: string, url: string): Promise<any> {
    return this.assetService.load(`${url}`, this.org.id).then(data => {
      const objectURL = URL.createObjectURL(data);
      this.images[imagekey] = this.sanitizer.bypassSecurityTrustUrl(objectURL);
      this.refreshPreview.emit();
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public removePolicy(): void {
    if (this.service instanceof ManagementService) {
      this.service.resetLabelPolicyToDefault().then(() => {
        this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
        setTimeout(() => {
          this.fetchData();
        }, 1000);
      }).catch(error => {
        this.toast.showError(error);
      });
    }
  }

  public savePolicy(): Promise<any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if ((this.previewData as LabelPolicy.AsObject).isDefault) {
          const req0 = new AddCustomLabelPolicyRequest();
          this.overwriteValues(req0);

          return (this.service as ManagementService).addCustomLabelPolicy(req0).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          }).catch((error: HttpErrorResponse) => {
            this.toast.showError(error);
          });
        } else {
          const req1 = new UpdateCustomLabelPolicyRequest();
          this.overwriteValues(req1);

          return (this.service as ManagementService).updateCustomLabelPolicy(req1).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          }).catch(error => {
            this.toast.showError(error);
          });
        }
      case PolicyComponentServiceType.ADMIN:
        const req = new UpdateLabelPolicyRequest();
        this.overwriteValues(req);
        return (this.service as AdminService).updateLabelPolicy(req).then(() => {
          this.toast.showInfo('POLICY.TOAST.SET', true);
        }).catch(error => {
          this.toast.showError(error);
        });
    }
  }

  public saveWatermark(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if ((this.previewData as LabelPolicy.AsObject).isDefault) {
          const req0 = new AddCustomLabelPolicyRequest();
          req0.setDisableWatermark(this.previewData.disableWatermark);

          (this.service as ManagementService).addCustomLabelPolicy(req0).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          }).catch((error: HttpErrorResponse) => {
            this.toast.showError(error);
          });
        } else {
          const req1 = new UpdateCustomLabelPolicyRequest();
          req1.setDisableWatermark(this.previewData.disableWatermark);

          (this.service as ManagementService).updateCustomLabelPolicy(req1).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          }).catch(error => {
            this.toast.showError(error);
          });
        }
        break;
      case PolicyComponentServiceType.ADMIN:
        const req = new UpdateLabelPolicyRequest();
        req.setDisableWatermark(this.data.disableWatermark);

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
    req.setBackgroundColorDark(this.previewData.backgroundColorDark);
    req.setBackgroundColor(this.previewData.backgroundColor);

    req.setFontColorDark(this.previewData.fontColorDark);
    req.setFontColor(this.previewData.fontColor);

    req.setPrimaryColorDark(this.previewData.primaryColorDark);
    req.setPrimaryColor(this.previewData.primaryColor);

    req.setWarnColorDark(this.previewData.warnColorDark);
    req.setWarnColor(this.previewData.warnColor);

    req.setDisableWatermark(this.previewData.disableWatermark);
    req.setHideLoginNameSuffix(this.previewData.hideLoginNameSuffix);
  }

  public activatePolicy(): Promise<any> {
    // dialog warning
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).activateCustomLabelPolicy().then(() => {
          this.toast.showInfo('POLICY.PRIVATELABELING.ACTIVATED', true);
          setTimeout(() => {
            this.loadingImages = true;
            this.getData().then(data => {

              if (data.policy) {
                this.data = data.policy;
                this.loadImages();
              }
            });
          }, 1000);
        }).catch(error => {
          this.toast.showError(error);
        });
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).activateLabelPolicy().then(() => {
          this.toast.showInfo('POLICY.PRIVATELABELING.ACTIVATED', true);
          setTimeout(() => {
            this.loadingImages = true;
            this.getData().then(data => {

              if (data.policy) {
                this.data = data.policy;
                this.loadImages();
              }
            });
          }, 1000);
        }).catch(error => {
          this.toast.showError(error);
        });
    }
  }

  public resetPolicy(): Promise<any> {
    return (this.service as ManagementService).resetLabelPolicyToDefault().then(() => {
      this.toast.showInfo('POLICY.PRIVATELABELING.RESET', true);
      setTimeout(() => {
        this.fetchData();
      });
    }).catch(error => {
      this.toast.showError(error);
    });
  }
}
