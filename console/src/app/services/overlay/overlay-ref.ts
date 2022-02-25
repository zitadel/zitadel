import { OverlayRef } from '@angular/cdk/overlay';

export class CnslOverlayRef {
  constructor(private overlayRef: OverlayRef) {}

  close(): void {
    this.overlayRef.dispose();
  }
}
