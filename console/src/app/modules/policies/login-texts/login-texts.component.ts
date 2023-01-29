import { Component, Injector, Input, OnDestroy, OnInit, Type } from '@angular/core';
import { UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, interval, Observable, of, Subject, Subscription } from 'rxjs';
import { map, pairwise, startWith, takeUntil } from 'rxjs/operators';
import {
  GetCustomLoginTextsRequest as AdminGetCustomLoginTextsRequest,
  GetDefaultLoginTextsRequest as AdminGetDefaultLoginTextsRequest,
  SetCustomLoginTextsRequest as AdminSetCustomLoginTextsRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  GetCustomLoginTextsRequest,
  GetDefaultLoginTextsRequest,
  SetCustomLoginTextsRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { mapRequestValues } from './helper';

const MIN_INTERVAL_SECONDS = 10; // if the difference of a newer version to the current exceeds this time, a refresh button is shown.

/* eslint-disable */
const KeyNamesArray = [
  'emailVerificationDoneText',
  'emailVerificationText',
  'externalUserNotFoundText',
  'footerText',
  'initMfaDoneText',
  'initMfaOtpText',
  'initMfaPromptText',
  'initMfaU2fText',
  'initPasswordDoneText',
  'initPasswordText',
  'initializeDoneText',
  'initializeUserText',
  'linkingUserDoneText',
  'loginText',
  'logoutText',
  'mfaProvidersText',
  'passwordChangeDoneText',
  'passwordChangeText',
  'passwordResetDoneText',
  'passwordText',
  'registrationOptionText',
  'registrationOrgText',
  'registrationUserText',
  'selectAccountText',
  'successLoginText',
  'usernameChangeDoneText',
  'usernameChangeText',
  'verifyMfaOtpText',
  'verifyMfaU2fText',
  'passwordlessPromptText',
  'passwordlessRegistrationDoneText',
  'passwordlessRegistrationText',
  'passwordlessText',
  'externalRegistrationUserOverviewText',
];
/* eslint-enable */

const REQUESTMAP = {
  [PolicyComponentServiceType.MGMT]: {
    get: new GetCustomLoginTextsRequest(),
    set: new SetCustomLoginTextsRequest(),
    getDefault: new GetDefaultLoginTextsRequest(),
    setFcn: (mgmtmap: Partial<SetCustomLoginTextsRequest.AsObject>): SetCustomLoginTextsRequest => {
      let req = new SetCustomLoginTextsRequest();
      req.setLanguage(mgmtmap.language ?? '');
      req = mapRequestValues(mgmtmap, req);
      return req;
    },
  },
  [PolicyComponentServiceType.ADMIN]: {
    get: new AdminGetCustomLoginTextsRequest(),
    set: new AdminSetCustomLoginTextsRequest(),
    getDefault: new AdminGetDefaultLoginTextsRequest(),
    setFcn: (adminmap: Partial<AdminSetCustomLoginTextsRequest.AsObject>): AdminSetCustomLoginTextsRequest => {
      let req = new AdminSetCustomLoginTextsRequest();
      req.setLanguage(adminmap.language ?? '');
      req = mapRequestValues(adminmap, req);
      return req;
    },
  },
};
@Component({
  selector: 'cnsl-login-texts',
  templateUrl: './login-texts.component.html',
  styleUrls: ['./login-texts.component.scss'],
})
export class LoginTextsComponent implements OnInit, OnDestroy {
  public loading: boolean = false;
  public currentPolicyChangeDate!: Timestamp.AsObject | undefined;
  public newerPolicyChangeDate!: Timestamp.AsObject | undefined;

  public totalCustomPolicy?: { [key: string]: { [key: string]: string } | boolean } = {}; // LoginCustomText.AsObject

  public getDefaultInitMessageTextMap$: Observable<{ [key: string]: string }> = of({});
  public getCustomInitMessageTextMap$: BehaviorSubject<{ [key: string]: string | boolean }> = new BehaviorSubject({});

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public KeyNamesArray: string[] = KeyNamesArray;
  public LOCALES: string[] = ['en', 'de', 'it', 'fr', 'pl', 'zh'];

  private sub: Subscription = new Subscription();

  public updateRequest!: SetCustomLoginTextsRequest;

  public destroy$: Subject<void> = new Subject();
  public InfoSectionType: any = InfoSectionType;
  public form: UntypedFormGroup = new UntypedFormGroup({
    currentSubMap: new UntypedFormControl('emailVerificationDoneText'),
    locale: new UntypedFormControl('en'),
  });

  public isDefault: boolean = false;

  public canWrite$: Observable<boolean> = this.authService.isAllowed([
    this.serviceType === PolicyComponentServiceType.ADMIN
      ? 'iam.policy.write'
      : this.serviceType === PolicyComponentServiceType.MGMT
      ? 'policy.write'
      : '',
  ]);
  constructor(
    private authService: GrpcAuthService,
    private injector: Injector,
    private dialog: MatDialog,
    private toast: ToastService,
  ) {
    this.form.valueChanges
      .pipe(startWith({ currentSubMap: 'emailVerificationDoneText', locale: 'en' }), pairwise(), takeUntil(this.destroy$))
      .subscribe((pair) => {
        this.checkForUnsaved(pair[0].currentSubMap).then((wantsToSave) => {
          if (wantsToSave) {
            this.saveCurrentTexts()
              .then(() => {
                this.loadData();
              })
              .catch(() => {
                // load even if save failed
                this.loadData();
              });
          } else {
            this.loadData();
          }
        });
      });
  }

  ngOnInit(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        this.service = this.injector.get(ManagementService as Type<ManagementService>);

        this.service.getSupportedLanguages().then((lang) => {
          this.LOCALES = lang.languagesList;
        });

        this.loadData();
        break;
      case PolicyComponentServiceType.ADMIN:
        this.service = this.injector.get(AdminService as Type<AdminService>);

        this.service.getSupportedLanguages().then((lang) => {
          this.LOCALES = lang.languagesList;
        });

        this.loadData();
        break;
    }

    interval(10000)
      .pipe(
        // debounceTime(5000),
        takeUntil(this.destroy$),
      )
      .subscribe((x) => {
        this.checkForChanges();
      });
  }

  public getDefaultValues(req: any): Promise<any> {
    return this.service.getDefaultLoginTexts(req).then((res) => {
      if (res.customText) {
        // delete res.customText.details;
        return Object.assign({}, res.customText);
      } else {
        return {};
      }
    });
  }

  public getCurrentValues(req: any): Promise<any> {
    return (this.service as ManagementService).getCustomLoginTexts(req).then((res) => {
      if (res.customText) {
        this.currentPolicyChangeDate = res.customText.details?.changeDate;
        return Object.assign({}, res.customText);
      } else {
        return {};
      }
    });
  }

  public async loadData(): Promise<any> {
    this.loading = true;
    const reqDefaultInit = REQUESTMAP[this.serviceType].getDefault;
    reqDefaultInit.setLanguage(this.locale);
    this.getDefaultInitMessageTextMap$ = from(this.getDefaultValues(reqDefaultInit)).pipe(map((m) => m[this.currentSubMap]));

    const reqCustomInit = REQUESTMAP[this.serviceType].get.setLanguage(this.locale);
    return this.getCurrentValues(reqCustomInit)
      .then((policy) => {
        this.loading = false;
        if (policy) {
          this.isDefault = policy.isDefault ?? false;

          this.totalCustomPolicy = policy;
          this.getCustomInitMessageTextMap$.next(policy[this.currentSubMap]);
        }
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
  }

  private async patchSingleCurrentMap(): Promise<any> {
    const reqCustomInit = REQUESTMAP[this.serviceType].get.setLanguage(this.locale);
    this.getCurrentValues(reqCustomInit).then((policy) => {
      this.getCustomInitMessageTextMap$.next(policy[this.currentSubMap]);
    });
  }

  public checkForChanges(): void {
    const reqCustomInit = REQUESTMAP[this.serviceType].get.setLanguage(this.locale);

    (this.service as ManagementService).getCustomLoginTexts(reqCustomInit).then((policy) => {
      this.newerPolicyChangeDate = policy.customText?.details?.changeDate;
    });
  }

  /**
   *
   * @param oldkey which was potentially unsaved
   * @returns a boolean if saving is desired
   */
  public checkForUnsaved(oldkey: string): Promise<boolean> {
    const old = this.getCustomInitMessageTextMap$.getValue();
    const unsaved = this.totalCustomPolicy ? this.totalCustomPolicy[oldkey] : undefined;

    if (old && unsaved && JSON.stringify(old) !== JSON.stringify(unsaved)) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.SAVE',
          cancelKey: 'ACTIONS.CONTINUEWITHOUTSAVE',
          titleKey: 'POLICY.LOGIN_TEXTS.UNSAVED_TITLE',
          descriptionKey: 'POLICY.LOGIN_TEXTS.UNSAVED_DESCRIPTION',
        },
        width: '400px',
      });

      return dialogRef.afterClosed().toPromise();
    } else {
      return Promise.resolve(false);
    }
  }

  public updateCurrentValues(values: { [key: string]: string }): void {
    if (this.totalCustomPolicy) {
      const setFcn = REQUESTMAP[this.serviceType].setFcn;
      this.totalCustomPolicy[this.currentSubMap] = values;

      this.updateRequest = setFcn(this.totalCustomPolicy);
      this.updateRequest.setLanguage(this.locale);
    }
  }

  public saveCurrentTexts(): Promise<any> {
    const entirePayload = this.updateRequest.toObject();
    this.getCustomInitMessageTextMap$.next((entirePayload as any)[this.currentSubMap]);

    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService)
          .setCustomLoginText(this.updateRequest)
          .then(() => {
            this.updateCurrentPolicyDate();
            this.isDefault = false;
            this.toast.showInfo('POLICY.MESSAGE_TEXTS.TOAST.UPDATED', true);
            setTimeout(() => {
              this.patchSingleCurrentMap();
            }, 1000);
          })
          .catch((error) => this.toast.showError(error));
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService)
          .setCustomLoginText(this.updateRequest)
          .then(() => {
            this.updateCurrentPolicyDate();
            this.isDefault = false;
            this.toast.showInfo('POLICY.MESSAGE_TEXTS.TOAST.UPDATED', true);
          })
          .catch((error) => this.toast.showError(error));
    }
  }

  private updateCurrentPolicyDate(): void {
    const ts = new Timestamp();
    const milliseconds = new Date().getTime();
    const seconds = Math.abs(milliseconds / 1000);
    const nanos = (milliseconds - seconds * 1000) * 1000 * 1000;
    ts.setSeconds(seconds);
    ts.setNanos(nanos);

    if (this.currentPolicyChangeDate) {
      const oldDate = new Date(
        this.currentPolicyChangeDate.seconds * 1000 + this.currentPolicyChangeDate.nanos / 1000 / 1000,
      );
      const newDate = ts.toDate();
      if (newDate.getTime() > oldDate.getTime()) {
        this.currentPolicyChangeDate = ts.toObject();
      }
    }
  }

  public resetDefault(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        icon: 'las la-history',
        confirmKey: 'ACTIONS.RESTORE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'POLICY.LOGIN_TEXTS.RESET_TITLE',
        descriptionKey: 'POLICY.LOGIN_TEXTS.RESET_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
          (this.service as ManagementService)
            .resetCustomLoginTextToDefault(this.locale)
            .then(() => {
              this.updateCurrentPolicyDate();
              this.isDefault = true;
              setTimeout(() => {
                this.loadData();
              }, 1000);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
          (this.service as AdminService)
            .resetCustomLoginTextToDefault(this.locale)
            .then(() => {
              this.updateCurrentPolicyDate();
              setTimeout(() => {
                this.loadData();
              }, 1000);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    });
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
    this.destroy$.next();
    this.destroy$.complete();
  }

  public get newerVersionExists(): boolean {
    const toDate = (ts: Timestamp.AsObject) => {
      return new Date(ts.seconds * 1000 + ts.nanos / 1000 / 1000);
    };
    if (this.newerPolicyChangeDate && this.currentPolicyChangeDate) {
      const ms = toDate(this.newerPolicyChangeDate).getTime() - toDate(this.currentPolicyChangeDate).getTime();
      // show button if changes are newer than 10s
      return ms / 1000 > MIN_INTERVAL_SECONDS;
    } else {
      return false;
    }
  }

  public get locale(): string {
    return this.form.get('locale')?.value;
  }

  public get currentSubMap(): string {
    return this.form.get('currentSubMap')?.value;
  }
}
