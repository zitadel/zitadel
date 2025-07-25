import {
  ChangeDetectionStrategy,
  Component,
  computed,
  effect,
  EventEmitter,
  Injector,
  Input,
  OnInit,
  Output,
  runInInjectionContext,
  Signal,
  untracked,
} from '@angular/core';
import { CdkConnectedOverlay, CdkOverlayOrigin, FlexibleConnectedPositionStrategy, Overlay } from '@angular/cdk/overlay';
import { BreakpointObserver } from '@angular/cdk/layout';
import { map } from 'rxjs/operators';
import { toSignal } from '@angular/core/rxjs-interop';
import { AsyncPipe, NgIf } from '@angular/common';
import { ReplaySubject } from 'rxjs';

@Component({
  selector: 'cnsl-header-dropdown',
  templateUrl: './header-dropdown.component.html',
  styleUrls: ['./header-dropdown.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CdkConnectedOverlay, NgIf, AsyncPipe],
})
export class HeaderDropdownComponent implements OnInit {
  @Input({ required: true })
  public trigger!: CdkOverlayOrigin;

  @Input({ required: true })
  public set isOpen(isOpen: boolean) {
    this.isOpen$.next(isOpen);
  }

  @Output()
  public closed = new EventEmitter<void>();

  protected readonly isOpen$ = new ReplaySubject<boolean>(1);
  protected readonly isHandset: Signal<boolean>;
  protected readonly positionStrategy: Signal<FlexibleConnectedPositionStrategy>;
  protected readonly scrollStrategy = this.overlay.scrollStrategies.block();

  constructor(
    private readonly overlay: Overlay,
    private readonly breakpointObserver: BreakpointObserver,
    private readonly injector: Injector,
  ) {
    this.isHandset = this.getIsHandset();
    this.positionStrategy = this.getPositionStrategy(this.isHandset);
  }

  ngOnInit(): void {
    // because closeWhenResized accesses the input properties, we need to run it in ngOnInit
    // this method is used to close the dropdown when the screen is resized
    // to make sure the dropdown will be rendered in the correct position
    runInInjectionContext(this.injector, () => {
      const isOpen = toSignal(this.isOpen$, { requireSync: true });
      effect(
        () => {
          this.isHandset();
          if (untracked(() => isOpen())) {
            this.closed.emit();
          }
        },
        { allowSignalWrites: true },
      );
    });
  }

  private getIsHandset() {
    const mediaQuery = '(max-width: 599px)';
    const isHandset$ = this.breakpointObserver.observe(mediaQuery).pipe(map(({ matches }) => matches));
    return toSignal(isHandset$, { initialValue: this.breakpointObserver.isMatched(mediaQuery) });
  }

  private getPositionStrategy(isHandset: Signal<boolean>): Signal<FlexibleConnectedPositionStrategy> {
    return computed(() =>
      isHandset()
        ? this.overlay
            .position()
            .flexibleConnectedTo(document.body)
            .withPositions([
              {
                originX: 'start',
                originY: 'bottom',
                overlayX: 'start',
                overlayY: 'bottom',
              },
            ])
        : this.overlay
            .position()
            .flexibleConnectedTo(this.trigger.elementRef)
            .withPositions([
              {
                originX: 'start',
                originY: 'bottom',
                overlayX: 'start',
                overlayY: 'top',
                offsetY: 8, // 8px gap between trigger and overlay
              },
            ]),
    );
  }
}
