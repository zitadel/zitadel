import { ChangeDetectionStrategy, Component, EventEmitter, Output, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';
import { Router, RouterLink } from '@angular/router';
import { InstanceDetail } from '@zitadel/proto/zitadel/instance_pb';

@Component({
  selector: 'cnsl-instance-selector',
  templateUrl: './instance-selector.component.html',
  styleUrls: ['./instance-selector.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [TranslateModule, MatButtonModule, RouterLink],
})
export class InstanceSelectorComponent {
  @Output() public instanceChanged = new EventEmitter<string>();
  @Output() public settingsClicked = new EventEmitter<void>();

  @Input({ required: true })
  public instance!: InstanceDetail;

  constructor(private readonly router: Router) {}

  protected async setInstance({ id }: InstanceDetail) {
    this.instanceChanged.emit(id);
    await this.router.navigate(['/']);
  }
}
