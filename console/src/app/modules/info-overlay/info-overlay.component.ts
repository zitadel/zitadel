import { Component, Input, OnDestroy } from '@angular/core';
import { Subject, takeUntil } from 'rxjs';
import { OverlayService } from 'src/app/services/overlay.service';

export enum InfoOverlayArrowType {
  TOP_LEFT,
  TOP_RIGHT,
}

@Component({
  selector: 'cnsl-info-overlay',
  templateUrl: './info-overlay.component.html',
  styleUrls: ['./info-overlay.component.scss'],
})
export class InfoOverlayComponent implements OnDestroy {
  @Input() workflowId: string = '';
  @Input() workflowStepId: string = '';

  @Input() public idsToHighlight: string[] = [];
  @Input() public inset: string = '';
  @Input() public show: boolean = false;
  @Input() arrowType: InfoOverlayArrowType = InfoOverlayArrowType.TOP_LEFT;
  InfoOverlayArrowType: any = InfoOverlayArrowType;

  private destroy$: Subject<void> = new Subject();
  private previousZIndex: string = 'auto';

  constructor(public overlayService: OverlayService) {
    this.overlayService.currentOverlayId$.pipe(takeUntil(this.destroy$)).subscribe((overlayStepId) => {
      console.log(overlayStepId);
      if (this.workflowStepId && this.workflowStepId === overlayStepId) {
        this.show = true;
        this.highlightIds();
      } else {
        this.show = false;
        this.resetIds();
      }
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private highlightIds(): void {
    this.idsToHighlight.forEach((id) => {
      const element = document.getElementById(id);
      if (element) {
        this.previousZIndex = element!.style.zIndex ?? 'auto'; // use id map for multiple
        element!.style.zIndex = '502';
      }
    });
  }

  private resetIds(): void {
    this.idsToHighlight.forEach((id) => {
      const element = document.getElementById(id);
      if (element) {
        element!.style.zIndex = this.show ? '502' : this.previousZIndex; // use id map for multiple
      }
    });
  }
}
