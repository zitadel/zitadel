import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { CopyToClipboardModule } from '../../directives/copy-to-clipboard/copy-to-clipboard.module';

@Component({
  standalone: true,
  selector: 'cnsl-copy-row',
  templateUrl: './copy-row.component.html',
  styleUrls: ['./copy-row.component.scss'],
  imports: [CommonModule, TranslateModule, MatButtonModule, MatTooltipModule, CopyToClipboardModule],
})
export class CopyRowComponent {
  @Input({ required: true }) public label = '';
  @Input({ required: true }) public value = '';
  @Input() public labelMinWidth = '';

  public copied = '';
}
