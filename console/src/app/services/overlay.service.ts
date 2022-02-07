import { ConnectionPositionPair, Overlay, OverlayConfig, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal, PortalInjector } from '@angular/cdk/portal';
import { ComponentRef, Injectable, Injector } from '@angular/core';

import { InfoOverlayComponent, OVERLAY_DATA } from '../modules/info-overlay/info-overlay.component';
import { CnslOverlayRef } from './overlay-ref';
import { CnslOverlay } from './overlay-workflow.service';

interface InfoOverlayConfig extends CnslOverlay {
  backdropClass?: string;
  panelClass?: string;
}

const DEFAULT_CONFIG: Partial<InfoOverlayConfig> = {
  backdropClass: 'cnsl-overlay-backdrop',
  panelClass: 'cnsl-overlay-panel',
};

@Injectable({
  providedIn: 'root',
})
export class OverlayService {
  constructor(private overlay: Overlay, private injector: Injector) {}

  public open(overlay: CnslOverlay) {
    // Override default configuration
    const dialogConfig: InfoOverlayConfig = { ...DEFAULT_CONFIG, ...overlay };
    console.log(dialogConfig);

    // Returns an OverlayRef which is a PortalHost
    const overlayRef = this.createOverlay(dialogConfig);

    // Instantiate remote control
    const dialogRef = new CnslOverlayRef(overlayRef);

    // Create ComponentPortal that can be attached to a PortalHost
    // const filePreviewPortal = new ComponentPortal(InfoOverlayComponent);
    const overlayComponent = this.attachOverlayContainer(overlayRef, dialogConfig, overlayRef);

    overlayRef.backdropClick().subscribe((_) => dialogRef.close());

    // Attach ComponentPortal to PortalHost
    // overlayRef.attach(filePreviewPortal);

    return dialogRef;
  }

  private attachOverlayContainer(overlayRef: OverlayRef, config: InfoOverlayConfig, dialogRef: OverlayRef) {
    const injector = this.createInjector(config, dialogRef);

    const containerPortal = new ComponentPortal(InfoOverlayComponent, null, injector);
    const containerRef: ComponentRef<InfoOverlayComponent> = overlayRef.attach(containerPortal);

    return containerRef.instance;
  }

  private createInjector(config: InfoOverlayConfig, dialogRef: OverlayRef): PortalInjector {
    const injectionTokens = new WeakMap();

    injectionTokens.set(OverlayRef, dialogRef);
    injectionTokens.set(OVERLAY_DATA, config.content);

    return new PortalInjector(this.injector, injectionTokens);
  }

  private createOverlay(config: InfoOverlayConfig): OverlayRef {
    const overlayConfig = this.getOverlayConfig(config);
    return this.overlay.create(overlayConfig);
  }

  private getOverlayConfig(config: InfoOverlayConfig): OverlayConfig {
    // const positionStrategy = this.overlay.position().global().centerHorizontally().centerVertically();
    const positions = [
      new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }),
      new ConnectionPositionPair({ originX: 'start', originY: 'top' }, { overlayX: 'start', overlayY: 'bottom' }),
    ];

    const htmlOrigin: HTMLElement | null = document.getElementById(config.origin);

    let positionStrategy;
    if (htmlOrigin) {
      console.log(`use html origin: ${config.origin}`);
      positionStrategy = this.overlay
        .position()
        .flexibleConnectedTo(htmlOrigin)
        .withPositions(positions)
        .withFlexibleDimensions(false)
        .withPush(false);
    } else {
      console.log(`use central position strategy`);
      positionStrategy = this.overlay.position().global().centerHorizontally().centerVertically();
    }

    const overlayConfig = new OverlayConfig({
      hasBackdrop: true,
      backdropClass: config.backdropClass,
      panelClass: config.panelClass,
      scrollStrategy: this.overlay.scrollStrategies.block(),
      positionStrategy,
    });

    return overlayConfig;
  }
}
