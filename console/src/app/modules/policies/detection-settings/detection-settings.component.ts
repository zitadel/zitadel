import { CommonModule } from '@angular/common';
import { Component, DestroyRef, OnInit } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MessageInitShape } from '@bufbuild/protobuf';
import { FormBuilder, FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { DurationSchema, type Duration } from '@bufbuild/protobuf/wkt';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TranslateModule } from '@ngx-translate/core';

import { CardModule } from '../../card/card.module';
import { InputModule } from '../../input/input.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { requiredValidator } from '../../form-field/validators/validators';
import { NewSettingsService } from '../../../services/new-settings.service';
import { GrpcAuthService } from '../../../services/grpc-auth.service';
import { ToastService } from '../../../services/toast.service';

type DetectionSettingsFormValue = {
  enabled: boolean;
  failOpen: boolean;
  failureBurstThreshold: number;
  historyWindowSeconds: number;
  contextChangeWindowSeconds: number;
  maxSignalsPerUser: number;
  maxSignalsPerSession: number;
};

@Component({
  selector: 'cnsl-detection-settings',
  standalone: true,
  templateUrl: './detection-settings.component.html',
  styleUrls: ['./detection-settings.component.scss'],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    TranslateModule,
    MatButtonModule,
    MatCheckboxModule,
    MatProgressSpinnerModule,
    CardModule,
    InputModule,
    InfoSectionModule,
  ],
})
export class DetectionSettingsComponent implements OnInit {
  protected readonly canWrite$ = this.authService.isAllowed(['iam.policy.write']);
  protected readonly form = this.fb.group({
    enabled: new FormControl(false, { nonNullable: true }),
    failOpen: new FormControl(false, { nonNullable: true }),
    failureBurstThreshold: new FormControl(1, {
      nonNullable: true,
      validators: [requiredValidator, Validators.min(1)],
    }),
    historyWindowSeconds: new FormControl(1, {
      nonNullable: true,
      validators: [requiredValidator, Validators.min(1)],
    }),
    contextChangeWindowSeconds: new FormControl(1, {
      nonNullable: true,
      validators: [requiredValidator, Validators.min(1)],
    }),
    maxSignalsPerUser: new FormControl(1, {
      nonNullable: true,
      validators: [requiredValidator, Validators.min(1)],
    }),
    maxSignalsPerSession: new FormControl(1, {
      nonNullable: true,
      validators: [requiredValidator, Validators.min(1)],
    }),
  });

  protected loading = true;
  protected saving = false;

  private lastLoadedValue?: DetectionSettingsFormValue;

  constructor(
    private readonly fb: FormBuilder,
    private readonly settingsService: NewSettingsService,
    private readonly authService: GrpcAuthService,
    private readonly toast: ToastService,
    private readonly destroyRef: DestroyRef,
  ) {
    this.canWrite$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((canWrite) => {
      if (canWrite) {
        this.form.enable({ emitEvent: false });
      } else {
        this.form.disable({ emitEvent: false });
      }
    });
  }

  public ngOnInit(): void {
    void this.loadSettings();
  }

  protected async save(): Promise<void> {
    if (this.form.invalid || this.saving) {
      return;
    }

    this.saving = true;
    try {
      const value = this.form.getRawValue();
      await this.settingsService.setDetectionSettings({
        settings: {
          enabled: value.enabled,
          failOpen: value.failOpen,
          failureBurstThreshold: value.failureBurstThreshold,
          historyWindow: secondsToDuration(value.historyWindowSeconds),
          contextChangeWindow: secondsToDuration(value.contextChangeWindowSeconds),
          maxSignalsPerUser: value.maxSignalsPerUser,
          maxSignalsPerSession: value.maxSignalsPerSession,
        },
      });
      this.toast.showInfo('SETTING.DETECTION.SAVED', true);
      await waitForProjection();
      await this.loadSettings();
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.saving = false;
    }
  }

  protected discard(): void {
    if (!this.lastLoadedValue) {
      return;
    }
    this.form.reset(this.lastLoadedValue);
    this.form.markAsPristine();
    this.form.markAsUntouched();
  }

  private async loadSettings(): Promise<void> {
    this.loading = true;
    try {
      const response = await this.settingsService.getDetectionSettings();
      const nextValue = {
        enabled: response.settings?.enabled ?? false,
        failOpen: response.settings?.failOpen ?? false,
        failureBurstThreshold: response.settings?.failureBurstThreshold ?? 1,
        historyWindowSeconds: durationToSeconds(response.settings?.historyWindow),
        contextChangeWindowSeconds: durationToSeconds(response.settings?.contextChangeWindow),
        maxSignalsPerUser: response.settings?.maxSignalsPerUser ?? 1,
        maxSignalsPerSession: response.settings?.maxSignalsPerSession ?? 1,
      } satisfies DetectionSettingsFormValue;
      this.lastLoadedValue = nextValue;
      this.form.reset(nextValue);
      this.form.markAsPristine();
      this.form.markAsUntouched();
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.loading = false;
    }
  }
}

function durationToSeconds(duration?: Duration): number {
  if (!duration) {
    return 1;
  }
  return Math.max(1, Number(duration.seconds ?? 0n) + (duration.nanos ?? 0) / 1_000_000_000);
}

function secondsToDuration(seconds: number) {
  return {
    seconds: BigInt(Math.max(1, Math.trunc(seconds))),
    nanos: 0,
  } satisfies MessageInitShape<typeof DurationSchema>;
}

async function waitForProjection(): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 1000));
}
