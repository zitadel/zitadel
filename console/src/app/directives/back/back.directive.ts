import { Directive, ElementRef, HostListener, Renderer2 } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NavigationService } from 'src/app/services/navigation.service';

@Directive({
  selector: '[cnslBack]',
})
export class BackDirective {
  new: Boolean = false;
  @HostListener('click')
  onClick(): void {
    this.navigation.back();
    if (this.new) {
      this.navigation.back();
    }
  }

  constructor(
    private navigation: NavigationService,
    private elRef: ElementRef,
    private renderer2: Renderer2,
    private route: ActivatedRoute,
  ) {
    this.route.queryParams.subscribe((params) => {
      this.new = params['new'];
    });

    if (navigation.isBackPossible) {
      // this.renderer2.removeStyle(this.elRef.nativeElement, 'visibility');
    } else {
      this.renderer2.setStyle(this.elRef.nativeElement, 'display', 'none');
    }
  }
}
