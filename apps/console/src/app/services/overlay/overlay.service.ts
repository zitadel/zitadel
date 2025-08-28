import { ConnectedPosition, ConnectionPositionPair, Overlay, OverlayConfig, OverlayRef } from '@angular/cdk/overlay';
import { ComponentPortal, PortalInjector } from '@angular/cdk/portal';
import { ComponentRef, Injectable, Injector } from '@angular/core';

import { InfoOverlayComponent, OVERLAY_DATA } from '../../modules/info-overlay/info-overlay.component';
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
  constructor(
    private overlay: Overlay,
    private injector: Injector,
  ) {}

  public open(overlay: CnslOverlay) {
    const dialogConfig: InfoOverlayConfig = { ...DEFAULT_CONFIG, ...overlay };
    const overlayRef = this.createOverlay(dialogConfig);

    const dialogRef = new CnslOverlayRef(overlayRef);

    const overlayComponent = this.attachOverlayContainer(overlayRef, dialogConfig, overlayRef);

    overlayRef.backdropClick().subscribe((_) => dialogRef.close());

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

    const htmlOrigin: HTMLElement | null = document.getElementById(config.origin);

    const positions: ConnectedPosition[] = [
      new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }, 0, 10),
      new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
    ];

    let positionStrategy;
    if (htmlOrigin) {
      positionStrategy = this.overlay
        .position()
        .flexibleConnectedTo(htmlOrigin)
        .withPositions(positions)
        .withFlexibleDimensions(true)
        .withPush(false);
    } else {
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
