import { Component, Inject, InjectionToken } from '@angular/core';
import { OverlayWorkflowService } from 'src/app/services/overlay-workflow.service';

export enum InfoOverlayArrowType {
  TOP_LEFT,
  TOP_RIGHT,
}

export const OVERLAY_DATA = new InjectionToken<any>('OVERLAY_DATA');

@Component({
  selector: 'cnsl-info-overlay',
  templateUrl: './info-overlay.component.html',
  styleUrls: ['./info-overlay.component.scss'],
})
export class InfoOverlayComponent {
  InfoOverlayArrowType: any = InfoOverlayArrowType;

  constructor(public workflowService: OverlayWorkflowService, @Inject(OVERLAY_DATA) public data: any) {}
}
