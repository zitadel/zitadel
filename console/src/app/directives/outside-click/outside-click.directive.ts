import { Directive, ElementRef, EventEmitter, HostListener, Output } from '@angular/core';

@Directive({
  selector: '[cnslOutsideClick]',
})
export class OutsideClickDirective {
  constructor(private elementRef: ElementRef) { }

  @Output() public clickOutside: EventEmitter<HTMLElement> = new EventEmitter();

  @HostListener('document:click', ['$event.target']) onMouseEnter(targetElement: HTMLElement): void {
    const clickedInside = this.elementRef.nativeElement.contains(targetElement);
    if (!clickedInside) {
      this.clickOutside.emit(targetElement);
    }
  }
}
