import { Component, Input, OnChanges, OnDestroy, OnInit, SimpleChanges } from '@angular/core';

export enum InfoOverlayArrowType {
  TOP_LEFT,
  TOP_RIGHT,
}

@Component({
  selector: 'cnsl-info-overlay',
  templateUrl: './info-overlay.component.html',
  styleUrls: ['./info-overlay.component.scss'],
})
export class InfoOverlayComponent implements OnChanges, OnInit, OnDestroy {
  @Input() workflowId: string = '';
  @Input() workflowStepId: string = '';

  @Input() public idsToHighlight: string[] = [];
  @Input() public inset: string = '';
  @Input() public show: boolean = false;
  @Input() arrowType: InfoOverlayArrowType = InfoOverlayArrowType.TOP_LEFT;
  InfoOverlayArrowType: any = InfoOverlayArrowType;

  private previousZIndex: string = 'auto';

  constructor() {}

  ngOnInit(): void {
    setTimeout(() => {
      this.highlightIds();
    }, 2000);
  }

  ngOnChanges(changes: SimpleChanges): void {
    console.log(changes);
    if (changes['idsToHighlight'].currentValue) {
      this.highlightIds();
    }
  }

  ngOnDestroy(): void {
    console.log(this.idsToHighlight);
  }

  private highlightIds(): void {
    console.log(this.idsToHighlight);

    this.idsToHighlight.forEach((id) => {
      const element = document.getElementById(id);
      if (element) {
        if (this.show) {
          this.previousZIndex = element!.style.zIndex ?? 'auto';
        }
        element!.style.zIndex = this.show ? '502' : this.previousZIndex; // black overlay is 500;
      }
    });
  }

  public dismiss() {
    this.show = false;
    this.highlightIds();
  }
}
