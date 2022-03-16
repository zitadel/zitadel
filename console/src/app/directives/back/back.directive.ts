import { Directive, ElementRef, HostListener, Renderer2 } from '@angular/core';
import { NavigationService } from 'src/app/services/navigation.service';

@Directive({
  selector: '[cnslBack]',
})
export class BackDirective {
  @HostListener('click')
  onClick(): void {
    this.navigation.back();
  }

  constructor(private navigation: NavigationService, private elRef: ElementRef, private renderer2: Renderer2) {
    if (navigation.isBackPossible) {
      // this.renderer2.removeStyle(this.elRef.nativeElement, 'visibility');
    } else {
      this.renderer2.setStyle(this.elRef.nativeElement, 'display', 'none');
    }
  }
}
