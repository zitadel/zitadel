import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, Input, signal } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';

@Component({
  standalone: true,
  selector: 'cnsl-copy-row',
  templateUrl: './copy-row.component.html',
  styleUrls: ['./copy-row.component.scss'],
  imports: [CommonModule, TranslateModule, MatButtonModule, MatTooltipModule, CopyToClipboardModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CopyRowComponent {
  @Input({ required: true }) public label!: string;
  @Input({ required: true }) public value!: string;
  @Input() public labelMinWidth = '';

  protected readonly copied = signal('');
}
