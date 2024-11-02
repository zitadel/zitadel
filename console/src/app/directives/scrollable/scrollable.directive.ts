import { Directive, ElementRef, EventEmitter, HostListener, Output } from '@angular/core';

@Directive({
  standalone: true,
  selector: '[cnslScrollable]',
})
export class ScrollableDirective {
  // when using this directive, add overflow-y scroll to css
  @Output() scrollPosition: EventEmitter<any> = new EventEmitter();

  constructor(public el: ElementRef) {}

  @HostListener('scroll', ['$event'])
  public onScroll(event: any): void {
    try {
      const top = event.target.scrollTop;
      const height = this.el.nativeElement.scrollHeight;
      const offset = this.el.nativeElement.offsetHeight;
      // emit bottom event
      if (top > height - offset - 1) {
        this.scrollPosition.emit('bottom');
      }

      // emit top event
      if (top === 0) {
        this.scrollPosition.emit('top');
      }
    } catch (err) {}
  }
}
