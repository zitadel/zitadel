import { Component, Input, OnInit, OnChanges, SimpleChanges } from '@angular/core';
import { Observable, BehaviorSubject, combineLatest, map } from 'rxjs';

export type EnvFormat = 'env' | 'json' | 'yaml';

export interface EnvVar {
  key: string;
  value: string;
  description?: string;
}

@Component({
  selector: 'cnsl-env-vars-block',
  templateUrl: './env-vars-block.component.html',
  styleUrls: ['./env-vars-block.component.scss'],
  standalone: false,
})
export class EnvVarsBlockComponent implements OnInit, OnChanges {
  @Input() envVars: EnvVar[] = [];
  @Input() title: string = 'Environment Variables';
  @Input() description?: string;
  @Input() defaultFormat: EnvFormat = 'env';
  @Input() availableFormats: EnvFormat[] = ['env', 'json', 'yaml'];
  @Input() showDownload: boolean = true;
  @Input() showCopy: boolean = true;

  public selectedFormat$ = new BehaviorSubject<EnvFormat>('env');
  public formattedContent$!: Observable<string>;
  public copied: string | null = null;

  public formatOptions: { value: EnvFormat; label: string; icon: string }[] = [
    { value: 'env', label: '.env', icon: 'file-text' },
    { value: 'json', label: 'JSON', icon: 'code' },
    { value: 'yaml', label: 'YAML', icon: 'file-code' },
  ];

  ngOnInit(): void {
    this.selectedFormat$.next(this.defaultFormat);
    this.setupFormattedContent();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['envVars'] && this.formattedContent$) {
      this.setupFormattedContent();
    }
  }

  private setupFormattedContent(): void {
    this.formattedContent$ = combineLatest([this.selectedFormat$]).pipe(map(([format]) => this.formatEnvVars(format)));
  }

  public switchFormat(format: EnvFormat): void {
    this.selectedFormat$.next(format);
  }

  public copyToClipboard(): void {
    this.formattedContent$.pipe().subscribe((content) => {
      navigator.clipboard.writeText(content).then(() => {
        this.copied = 'copy';
        setTimeout(() => (this.copied = null), 2000);
      });
    });
  }

  public downloadFile(): void {
    this.formattedContent$.pipe().subscribe((content) => {
      const format = this.selectedFormat$.value;
      const filename = this.getFilename(format);
      const mimeType = this.getMimeType(format);

      const blob = new Blob([content], { type: mimeType });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      link.click();
      window.URL.revokeObjectURL(url);

      this.copied = 'download';
      setTimeout(() => (this.copied = null), 2000);
    });
  }

  private formatEnvVars(format: EnvFormat): string {
    switch (format) {
      case 'env':
        return this.formatAsEnv();
      case 'json':
        return this.formatAsJson();
      case 'yaml':
        return this.formatAsYaml();
      default:
        return this.formatAsEnv();
    }
  }

  private formatAsEnv(): string {
    return this.envVars
      .map((envVar) => {
        let line = '';
        if (envVar.description) {
          line += `# ${envVar.description}\n`;
        }
        line += `${envVar.key}="${envVar.value}"`;
        return line;
      })
      .join('\n');
  }

  private formatAsJson(): string {
    const envObject = this.envVars.reduce(
      (acc, envVar) => {
        acc[envVar.key] = envVar.value;
        return acc;
      },
      {} as Record<string, string>,
    );

    return JSON.stringify(envObject, null, 2);
  }

  private formatAsYaml(): string {
    return this.envVars
      .map((envVar) => {
        let yaml = '';
        if (envVar.description) {
          yaml += `# ${envVar.description}\n`;
        }
        yaml += `${envVar.key}: "${envVar.value}"`;
        return yaml;
      })
      .join('\n');
  }

  private getFilename(format: EnvFormat): string {
    switch (format) {
      case 'env':
        return '.env';
      case 'json':
        return 'environment.json';
      case 'yaml':
        return 'environment.yaml';
      default:
        return '.env';
    }
  }

  private getMimeType(format: EnvFormat): string {
    switch (format) {
      case 'env':
        return 'text/plain';
      case 'json':
        return 'application/json';
      case 'yaml':
        return 'text/yaml';
      default:
        return 'text/plain';
    }
  }

  public get availableFormatOptions() {
    return this.formatOptions.filter((option) => this.availableFormats.includes(option.value));
  }
}
