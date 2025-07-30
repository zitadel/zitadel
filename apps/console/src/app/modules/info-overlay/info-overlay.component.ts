import { Component, Inject, InjectionToken } from '@angular/core';
import { OverlayWorkflowService } from 'src/app/services/overlay/overlay-workflow.service';

export const OVERLAY_DATA = new InjectionToken<any>('OVERLAY_DATA');

@Component({
  selector: 'cnsl-info-overlay',
  templateUrl: './info-overlay.component.html',
  styleUrls: ['./info-overlay.component.scss'],
})
export class InfoOverlayComponent {
  constructor(
    public workflowService: OverlayWorkflowService,
    @Inject(OVERLAY_DATA) public data: any,
  ) {}
}
