import { Directive, ElementRef, HostListener, Renderer2 } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs';
import { NavigationService } from 'src/app/services/navigation.service';

@Directive({
  selector: '[cnslBack]',
})
export class BackDirective {
  new: Boolean = false;
  @HostListener('click')
  onClick(): void {
    this.navigation.back();
    // Go back again to avoid create dialog starts again
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
    // Check if a new element was created using a create dialog
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      this.new = params['new'];
    });

    if (navigation.isBackPossible) {
      // this.renderer2.removeStyle(this.elRef.nativeElement, 'visibility');
    } else {
      this.renderer2.setStyle(this.elRef.nativeElement, 'display', 'none');
    }
  }
}
