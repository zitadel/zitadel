import 'codemirror/mode/yaml/yaml';

import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, Inject, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CodemirrorComponent, CodemirrorModule } from '@ctrl/ngx-codemirror';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import { MessageInitShape } from '@bufbuild/protobuf';
import {
  DetectionRule,
  DetectionRuleEngine,
  DetectionRuleSchema,
} from '@zitadel/proto/zitadel/settings/v2/detection_rules_pb';
import { stringify, parse } from 'yaml';

export type DetectionRuleEditorResult = MessageInitShape<typeof DetectionRuleSchema>;

interface DetectionRuleDialogData {
  rule?: DetectionRule;
}

/** Template shown when creating a new rule. */
const NEW_RULE_TEMPLATE = `id: my-rule
description: ""
expr: "true"
engine: log  # block | rate_limit | llm | log | captcha
priority: 0
stop_on_match: false
finding:
  name: my_finding
  message: ""
  block: false
# rate_limit:  # uncomment when action is rate_limit
#   key: "user:{{.Current.UserID}}"
#   window: 5m
#   max: 20
# context_template: ""  # uncomment when action is llm
`;

@Component({
  selector: 'cnsl-detection-rule-dialog',
  standalone: true,
  templateUrl: './detection-rule-dialog.component.html',
  styleUrls: ['./detection-rule-dialog.component.scss'],
  imports: [CommonModule, FormsModule, TranslateModule, MatButtonModule, MatDialogModule, CodemirrorModule],
})
export class DetectionRuleDialogComponent implements AfterViewInit {
  @ViewChild(CodemirrorComponent) private readonly codeMirror?: CodemirrorComponent;

  protected yamlContent: string;
  protected parseError: string | null = null;

  protected readonly codemirrorOptions = {
    lineNumbers: true,
    theme: 'material',
    mode: 'yaml',
    tabSize: 2,
    indentWithTabs: false,
    lineWrapping: true,
  };

  constructor(
    private readonly dialogRef: MatDialogRef<DetectionRuleDialogComponent, DetectionRuleEditorResult>,
    @Inject(MAT_DIALOG_DATA) protected readonly data: DetectionRuleDialogData | null,
  ) {
    this.yamlContent = data?.rule ? ruleToYaml(data.rule) : NEW_RULE_TEMPLATE;
  }

  ngAfterViewInit(): void {
    // Dialog open animation displaces CodeMirror's internal cursor calculations.
    // Refresh once the animation has settled so selection/cursor render correctly.
    setTimeout(() => this.codeMirror?.codeMirror?.refresh(), 200);
  }

  protected get isEdit(): boolean {
    return !!this.data?.rule;
  }

  protected onYamlChange(content: string): void {
    this.yamlContent = content;
    this.parseError = null;
    try {
      yamlToRule(content);
    } catch (e) {
      this.parseError = e instanceof Error ? e.message : String(e);
    }
  }

  protected get hasError(): boolean {
    return this.parseError !== null;
  }

  protected closeWithResult(): void {
    if (this.hasError) {
      return;
    }
    try {
      const rule = yamlToRule(this.yamlContent);
      // In edit mode, preserve the original rule ID — the update API uses it as a path parameter.
      if (this.isEdit && this.data?.rule) {
        rule.id = this.data.rule.id;
      }
      this.dialogRef.close(rule);
    } catch (e) {
      this.parseError = e instanceof Error ? e.message : String(e);
    }
  }
}

// ---------------------------------------------------------------------------
// YAML ↔ Rule conversion
// ---------------------------------------------------------------------------

interface RuleYaml {
  id?: string;
  description?: string;
  expr?: string;
  engine?: string;
  priority?: number;
  stop_on_match?: boolean;
  finding?: { name?: string; message?: string; block?: boolean };
  rate_limit?: { key?: string; window?: string; max?: number };
  context_template?: string;
}

function ruleToYaml(rule: DetectionRule): string {
  const obj: RuleYaml = {
    id: rule.id,
    description: rule.description || undefined,
    expr: rule.expr,
    engine: engineToString(rule.engine),
    priority: rule.priority ?? 0,
    stop_on_match: rule.stopOnMatch ?? false,
    finding: {
      name: rule.finding?.name || undefined,
      message: rule.finding?.message || undefined,
      block: rule.finding?.block || undefined,
    },
  };

  if (rule.rateLimit) {
    obj.rate_limit = {
      key: rule.rateLimit.key,
      window: durationToString(rule.rateLimit.window),
      max: rule.rateLimit.max,
    };
  }

  if (rule.contextTemplate) {
    obj.context_template = rule.contextTemplate;
  }

  return stringify(obj, { lineWidth: 0 });
}

function yamlToRule(yamlStr: string): DetectionRuleEditorResult {
  const parsed = parse(yamlStr) as RuleYaml;

  if (!parsed || typeof parsed !== 'object') {
    throw new Error('Invalid YAML: must be a mapping');
  }
  if (!parsed.id || typeof parsed.id !== 'string') {
    throw new Error('Rule must have an "id" field');
  }
  if (!parsed.expr || typeof parsed.expr !== 'string') {
    throw new Error('Rule must have an "expr" field');
  }
  if (!parsed.engine || typeof parsed.engine !== 'string') {
    throw new Error('Rule must have an "engine" field (block | rate_limit | llm | log | captcha)');
  }

  const engine = stringToEngine(parsed.engine);
  const rule: DetectionRuleEditorResult = {
    id: parsed.id,
    description: parsed.description ?? '',
    expr: parsed.expr,
    engine,
    priority: parsed.priority ?? 0,
    stopOnMatch: parsed.stop_on_match ?? false,
    finding: {
      name: parsed.finding?.name ?? '',
      message: parsed.finding?.message ?? '',
      block: parsed.finding?.block ?? false,
    },
    contextTemplate: parsed.context_template ?? '',
  };

  if (engine === DetectionRuleEngine.RATE_LIMIT) {
    if (!parsed.rate_limit) {
      throw new Error('rate_limit engine requires a "rate_limit" section');
    }
    rule.rateLimit = {
      key: parsed.rate_limit.key ?? '',
      window: parseDuration(parsed.rate_limit.window ?? '5m'),
      max: parsed.rate_limit.max ?? 10,
    };
  }

  return rule;
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function engineToString(engine: DetectionRuleEngine): string {
  switch (engine) {
    case DetectionRuleEngine.BLOCK:
      return 'block';
    case DetectionRuleEngine.RATE_LIMIT:
      return 'rate_limit';
    case DetectionRuleEngine.LLM:
      return 'llm';
    case DetectionRuleEngine.LOG:
      return 'log';
    case DetectionRuleEngine.CAPTCHA:
      return 'captcha';
    default:
      return 'log';
  }
}

function stringToEngine(s: string): DetectionRuleEngine {
  switch (s.toLowerCase()) {
    case 'block':
      return DetectionRuleEngine.BLOCK;
    case 'rate_limit':
      return DetectionRuleEngine.RATE_LIMIT;
    case 'llm':
      return DetectionRuleEngine.LLM;
    case 'log':
      return DetectionRuleEngine.LOG;
    case 'captcha':
      return DetectionRuleEngine.CAPTCHA;
    default:
      throw new Error(`Unknown engine "${s}" — must be one of: block, rate_limit, llm, log, captcha`);
  }
}

function durationToString(d?: { seconds?: bigint; nanos?: number }): string {
  if (!d) {
    return '5m';
  }
  const totalSeconds = Number(d.seconds ?? 0n);
  if (totalSeconds > 0 && totalSeconds % 3600 === 0) {
    return `${totalSeconds / 3600}h`;
  }
  if (totalSeconds > 0 && totalSeconds % 60 === 0) {
    return `${totalSeconds / 60}m`;
  }
  return `${totalSeconds}s`;
}

function parseDuration(s: string): { seconds: bigint; nanos: number } {
  const match = s.match(/^(\d+)(h|m|s)$/);
  if (!match) {
    throw new Error(`Invalid duration "${s}" — use formats like "5m", "30s", "1h"`);
  }
  const n = parseInt(match[1], 10);
  let seconds = 0;
  switch (match[2]) {
    case 'h':
      seconds = n * 3600;
      break;
    case 'm':
      seconds = n * 60;
      break;
    case 's':
      seconds = n;
      break;
  }
  return { seconds: BigInt(seconds), nanos: 0 };
}
