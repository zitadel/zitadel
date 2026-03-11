import { CommonModule } from '@angular/common';
import { Component, DestroyRef, OnInit } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import type { Duration } from '@bufbuild/protobuf/wkt';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TranslateModule } from '@ngx-translate/core';
import { DetectionRule, DetectionRuleEngine } from '@zitadel/proto/zitadel/settings/v2/detection_rules_pb';

import { CardModule } from '../../card/card.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { TimestampToDatePipeModule } from '../../../pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { GrpcAuthService } from '../../../services/grpc-auth.service';
import { NewSettingsService } from '../../../services/new-settings.service';
import { ToastService } from '../../../services/toast.service';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { DetectionRuleDialogComponent, DetectionRuleEditorResult } from './detection-rule-dialog.component';

@Component({
  selector: 'cnsl-detection-rules',
  standalone: true,
  templateUrl: './detection-rules.component.html',
  styleUrls: ['./detection-rules.component.scss'],
  imports: [
    CommonModule,
    TranslateModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    CardModule,
    InfoSectionModule,
    TimestampToDatePipeModule,
  ],
})
export class DetectionRulesComponent implements OnInit {
  protected readonly canWrite$ = this.authService.isAllowed(['iam.policy.write']);
  protected readonly canDelete$ = this.authService.isAllowed(['iam.policy.delete']);

  protected rules: DetectionRule[] = [];
  protected loading = true;
  protected mutating = false;

  constructor(
    private readonly dialog: MatDialog,
    private readonly authService: GrpcAuthService,
    private readonly settingsService: NewSettingsService,
    private readonly toast: ToastService,
    private readonly destroyRef: DestroyRef,
  ) {}

  public ngOnInit(): void {
    void this.loadRules();
  }

  protected openCreateDialog(): void {
    const dialogRef = this.dialog.open(DetectionRuleDialogComponent, {
      width: '860px',
    });

    dialogRef
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(async (result) => {
        if (!result) {
          return;
        }
        await this.createRule(result);
      });
  }

  protected openEditDialog(rule: DetectionRule): void {
    const dialogRef = this.dialog.open(DetectionRuleDialogComponent, {
      width: '860px',
      data: { rule },
    });

    dialogRef
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(async (result) => {
        if (!result) {
          return;
        }
        await this.updateRule(rule.id, result);
      });
  }

  protected confirmDelete(rule: DetectionRule): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      width: '400px',
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'SETTING.DETECTION_RULES.DELETE_TITLE',
        descriptionKey: 'SETTING.DETECTION_RULES.DELETE_DESCRIPTION',
      },
    });

    dialogRef
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(async (confirmed) => {
        if (!confirmed) {
          return;
        }
        await this.deleteRule(rule.id);
      });
  }

  protected translationKeyForAction(engine: DetectionRuleEngine): string {
    switch (engine) {
      case DetectionRuleEngine.BLOCK:
        return 'SETTING.DETECTION_RULES.ENGINES.BLOCK';
      case DetectionRuleEngine.RATE_LIMIT:
        return 'SETTING.DETECTION_RULES.ENGINES.RATE_LIMIT';
      case DetectionRuleEngine.LLM:
        return 'SETTING.DETECTION_RULES.ENGINES.LLM';
      case DetectionRuleEngine.LOG:
        return 'SETTING.DETECTION_RULES.ENGINES.LOG';
      case DetectionRuleEngine.CAPTCHA:
        return 'SETTING.DETECTION_RULES.ENGINES.CAPTCHA';
      default:
        return 'SETTING.DETECTION_RULES.ENGINES.UNSPECIFIED';
    }
  }

  protected durationToSeconds(duration?: Duration): number {
    if (!duration) {
      return 0;
    }
    return Number(duration.seconds ?? 0n) + (duration.nanos ?? 0) / 1_000_000_000;
  }

  protected trackRule(_: number, rule: DetectionRule): string {
    return rule.id;
  }

  private async loadRules(): Promise<void> {
    this.loading = true;
    try {
      const response = await this.settingsService.listDetectionRules();
      this.rules = [...response.rules].sort((left, right) => left.id.localeCompare(right.id));
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.loading = false;
    }
  }

  private async createRule(result: DetectionRuleEditorResult): Promise<void> {
    this.mutating = true;
    try {
      await this.settingsService.createDetectionRule({ rule: result });
      this.toast.showInfo('SETTING.DETECTION_RULES.CREATED', true);
      await waitForProjection();
      await this.loadRules();
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.mutating = false;
    }
  }

  private async updateRule(ruleId: string, result: DetectionRuleEditorResult): Promise<void> {
    this.mutating = true;
    try {
      await this.settingsService.updateDetectionRule({ ruleId, rule: result });
      this.toast.showInfo('SETTING.DETECTION_RULES.UPDATED', true);
      await waitForProjection();
      await this.loadRules();
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.mutating = false;
    }
  }

  private async deleteRule(ruleId: string): Promise<void> {
    this.mutating = true;
    try {
      await this.settingsService.deleteDetectionRule(ruleId);
      this.toast.showInfo('SETTING.DETECTION_RULES.DELETED', true);
      await waitForProjection();
      await this.loadRules();
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.mutating = false;
    }
  }
}

async function waitForProjection(): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 1000));
}
