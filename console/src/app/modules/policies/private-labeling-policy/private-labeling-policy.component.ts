import { HttpErrorResponse } from '@angular/common/http';
import { Component, EventEmitter, Injector, Input, OnDestroy, OnInit, Type } from '@angular/core';
import { Subject, Subscription } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
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
import { StorageKey, StorageLocation, StorageService } from 'src/app/services/storage.service';
import { ThemeService } from 'src/app/services/theme.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

export enum Theme {
  DARK,
  LIGHT,
}

export enum View {
  PREVIEW,
  CURRENT,
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

const MAX_ALLOWED_SIZE = 0.5 * 1024 * 1024;

@Component({
  selector: 'cnsl-private-labeling-policy',
  templateUrl: './private-labeling-policy.component.html',
  styleUrls: ['./private-labeling-policy.component.scss'],
})
export class PrivateLabelingPolicyComponent implements OnInit, OnDestroy {
  public theme: Theme = Theme.LIGHT;

  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
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

  public Theme: any = Theme;
  public View: any = View;
  public ColorType: any = ColorType;
  public AssetType: any = AssetType;

  public refreshPreview: EventEmitter<void> = new EventEmitter();
  public org!: Org.AsObject;
  public InfoSectionType: any = InfoSectionType;

  private destroy$: Subject<void> = new Subject();
  public view: View = View.PREVIEW;
  constructor(
    private toast: ToastService,
    private injector: Injector,
    private assetService: AssetService,
    private storageService: StorageService,
    private themeService: ThemeService,
  ) {}

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
      if (file.size > MAX_ALLOWED_SIZE) {
        this.toast.showInfo('POLICY.PRIVATELABELING.MAXSIZEEXCEEDED', true);
      } else {
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
  }

  public ngOnInit(): void {
    this.themeService.isDarkTheme.pipe(takeUntil(this.destroy$)).subscribe((isDark) => {
      if (isDark) {
        this.theme = Theme.DARK;
      } else {
        this.theme = Theme.LIGHT;
      }
    });

    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        this.service = this.injector.get(ManagementService as Type<ManagementService>);

        const org: Org.AsObject | null = this.storageService.getItem(StorageKey.organization, StorageLocation.session);

        if (org) {
          this.org = org;
        }
        break;
      case PolicyComponentServiceType.ADMIN:
        this.service = this.injector.get(AdminService as Type<AdminService>);

        break;
    }

    this.fetchData();
  }

  public onDropFont(filelist: FileList | null): Promise<any> | void {
    if (filelist) {
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
  }

  public deleteFont(): Promise<any> {
    const handler = (prom: Promise<any>) =>
      prom
        .then(() => {
          this.toast.showInfo('POLICY.TOAST.DELETESUCCESS', true);
          setTimeout(() => {
            this.getPreviewData().then((data) => {
              if (data.policy) {
                this.previewData = data.policy;
              }
            });
          }, 1000);
        })
        .catch((error) => this.toast.showError(error));

    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return handler((this.service as ManagementService).removeLabelPolicyFont());
      case PolicyComponentServiceType.ADMIN:
        return handler((this.service as AdminService).removeLabelPolicyFont());
    }
  }

  public deleteAsset(type: AssetType, theme: Theme): any {
    const previewHandler = (prom: Promise<any>) => {
      return prom
        .then(() => {
          this.toast.showInfo('POLICY.TOAST.DELETESUCCESS', true);
          setTimeout(() => {
            this.getPreviewData().then((data) => {
              if (data.policy) {
                this.previewData = data.policy;
              }
            });
          }, 1000);
        })
        .catch((error) => this.toast.showError(error));
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
      if (file.size > MAX_ALLOWED_SIZE) {
        this.toast.showInfo('POLICY.PRIVATELABELING.MAXSIZEEXCEEDED', true);
      } else {
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
  }

  private handleFontUploadPromise(task: Promise<any>): Promise<any> {
    const enhTask = task
      .then(() => {
        this.toast.showInfo('POLICY.TOAST.UPLOADSUCCESS', true);
        setTimeout(() => {
          this.getPreviewData().then((data) => {
            if (data.policy) {
              this.previewData = data.policy;
            }
          });
        }, 1000);
      })
      .catch((error) => this.toast.showError(error));

    if (this.serviceType === PolicyComponentServiceType.MGMT && (this.previewData as LabelPolicy.AsObject).isDefault) {
      return this.savePolicy().then(() => enhTask);
    } else {
      return enhTask;
    }
  }

  private handleUploadPromise(task: Promise<any>): Promise<any> {
    const enhTask = task
      .then(() => {
        this.toast.showInfo('POLICY.TOAST.UPLOADSUCCESS', true);
        setTimeout(() => {
          this.getPreviewData().then((data) => {
            if (data.policy) {
              this.previewData = data.policy;
            }
          });
        }, 1000);
      })
      .catch((error) => this.toast.showError(error));

    if (this.serviceType === PolicyComponentServiceType.MGMT && (this.previewData as LabelPolicy.AsObject).isDefault) {
      return this.savePolicy().then(() => enhTask);
    } else {
      return enhTask;
    }
  }

  public fetchData(): void {
    this.loading = true;

    this.getPreviewData()
      .then((data) => {
        if (data.policy) {
          this.previewData = data.policy;
          this.loading = false;
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });

    this.getData()
      .then((data) => {
        if (data.policy) {
          this.data = data.policy;
          this.loading = false;
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
    this.destroy$.next();
    this.destroy$.complete();
  }

  private async getPreviewData(): Promise<
    | MgmtGetPreviewLabelPolicyResponse.AsObject
    | AdminGetPreviewLabelPolicyResponse.AsObject
    | MgmtGetLabelPolicyResponse.AsObject
    | AdminGetLabelPolicyResponse.AsObject
  > {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getPreviewLabelPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getPreviewLabelPolicy();
    }
  }

  private async getData(): Promise<
    | MgmtGetPreviewLabelPolicyResponse.AsObject
    | AdminGetPreviewLabelPolicyResponse.AsObject
    | MgmtGetLabelPolicyResponse.AsObject
    | AdminGetLabelPolicyResponse.AsObject
  > {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getLabelPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getLabelPolicy();
    }
  }

  public removePolicy(): void {
    if (this.service instanceof ManagementService) {
      this.service
        .resetLabelPolicyToDefault()
        .then(() => {
          this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
          setTimeout(() => {
            this.fetchData();
          }, 1000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public savePolicy(): Promise<any> {
    const reloadPolicy = () => {
      setTimeout(() => {
        this.getData().then((data) => {
          if (data.policy) {
            this.data = data.policy;
          }
        });
      }, 500);
    };
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if ((this.previewData as LabelPolicy.AsObject).isDefault) {
          const req0 = new AddCustomLabelPolicyRequest();
          this.overwriteValues(req0);

          return (this.service as ManagementService)
            .addCustomLabelPolicy(req0)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);

              reloadPolicy();
            })
            .catch((error: HttpErrorResponse) => {
              this.toast.showError(error);
            });
        } else {
          const req1 = new UpdateCustomLabelPolicyRequest();
          this.overwriteValues(req1);

          return (this.service as ManagementService)
            .updateCustomLabelPolicy(req1)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);

              reloadPolicy();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      case PolicyComponentServiceType.ADMIN:
        const req = new UpdateLabelPolicyRequest();
        this.overwriteValues(req);
        return (this.service as AdminService)
          .updateLabelPolicy(req)
          .then(() => {
            reloadPolicy();
            this.toast.showInfo('POLICY.TOAST.SET', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
    }
  }

  public get isDefault(): boolean {
    if (this.previewData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.previewData as LabelPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }

  public setDarkBackgroundColorAndSave($event: string): void {
    this.previewData.backgroundColorDark = $event;
    this.savePolicy();
  }

  public setDarkPrimaryColorAndSave($event: string): void {
    this.previewData.primaryColorDark = $event;
    this.savePolicy();
  }

  public setDarkWarnColorAndSave($event: string): void {
    this.previewData.warnColorDark = $event;
    this.savePolicy();
  }

  public setDarkFontColorAndSave($event: string): void {
    this.previewData.fontColorDark = $event;
    this.savePolicy();
  }

  public setBackgroundColorAndSave($event: string): void {
    this.previewData.backgroundColor = $event;
    this.savePolicy();
  }

  public setPrimaryColorAndSave($event: string): void {
    this.previewData.primaryColor = $event;
    this.savePolicy();
  }

  public setWarnColorAndSave($event: string): void {
    this.previewData.warnColor = $event;
    this.savePolicy();
  }

  public setFontColorAndSave($event: string): void {
    this.previewData.fontColor = $event;
    this.savePolicy();
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
        return (this.service as ManagementService)
          .activateCustomLabelPolicy()
          .then(() => {
            this.toast.showInfo('POLICY.PRIVATELABELING.ACTIVATED', true);
            setTimeout(() => {
              this.getData().then((data) => {
                if (data.policy) {
                  this.data = data.policy;
                  this.applyToConsole(data.policy);
                }
              });
            }, 1000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService)
          .activateLabelPolicy()
          .then(() => {
            this.toast.showInfo('POLICY.PRIVATELABELING.ACTIVATED', true);
            setTimeout(() => {
              this.getData().then((data) => {
                if (data.policy) {
                  this.data = data.policy;
                  this.applyToConsole(data.policy);
                }
              });
            }, 1000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
    }
  }

  private applyToConsole(labelpolicy: LabelPolicy.AsObject): void {
    const darkPrimary = labelpolicy?.primaryColorDark || '#bbbafa';
    const lightPrimary = labelpolicy?.primaryColor || '#5469d4';

    const darkWarn = labelpolicy?.warnColorDark || '#ff3b5b';
    const lightWarn = labelpolicy?.warnColor || '#cd3d56';

    const darkBackground = labelpolicy?.backgroundColorDark || '#111827';
    const lightBackground = labelpolicy?.backgroundColor || '#fafafa';

    this.themeService.savePrimaryColor(darkPrimary, true);
    this.themeService.savePrimaryColor(lightPrimary, false);

    this.themeService.saveWarnColor(darkWarn, true);
    this.themeService.saveWarnColor(lightWarn, false);

    this.themeService.saveBackgroundColor(darkBackground, true);
    this.themeService.saveBackgroundColor(lightBackground, false);
  }

  public resetPolicy(): Promise<any> {
    return (this.service as ManagementService)
      .resetLabelPolicyToDefault()
      .then(() => {
        this.toast.showInfo('POLICY.PRIVATELABELING.RESET', true);
        setTimeout(() => {
          this.fetchData();
        });
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  // /**
  //  *  defaults to false because urls are distinct anyway
  //  */
  // public get previewEqualsCurrentPolicy(): boolean {
  //   const getComparable = (policy: LabelPolicy.AsObject): Partial<LabelPolicy.AsObject> => {
  //     return Object.assign({
  //       primaryColor: policy.primaryColor,
  //       hideLoginNameSuffix: policy.primaryColor,
  //       warnColor: policy.warnColor,
  //       backgroundColor: policy.backgroundColor,
  //       fontColor: policy.fontColor,
  //       primaryColorDark: policy.primaryColorDark,
  //       backgroundColorDark: policy.backgroundColorDark,
  //       warnColorDark: policy.warnColorDark,
  //       fontColorDark: policy.fontColorDark,
  //       disableWatermark: policy.disableWatermark,
  //       logoUrl: policy.logoUrl,
  //       iconUrl: policy.iconUrl,
  //       logoUrlDark: policy.logoUrlDark,
  //       iconUrlDark: policy.iconUrlDark,
  //       fontUrl: policy.fontUrl,
  //     });
  //   };

  //   const c = getComparable(this.data);
  //   const p = getComparable(this.previewData);

  //   return JSON.stringify(p) === JSON.stringify(c);
  // }
}
