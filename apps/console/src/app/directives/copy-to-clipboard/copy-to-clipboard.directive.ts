import { Directive, EventEmitter, HostListener, Input, Output } from '@angular/core';

@Directive({
  selector: '[cnslCopyToClipboard]',
})
export class CopyToClipboardDirective {
  @Input() valueToCopy: string = '';
  @Output() copiedValue: EventEmitter<string> = new EventEmitter();

  @HostListener('click', ['$event']) onClick($event: any): void {
    $event.preventDefault();
    $event.stopPropagation();
    this.copytoclipboard(this.valueToCopy);
  }

  public copytoclipboard(value: string): void {
    const selBox = document.createElement('textarea');
    selBox.style.position = 'fixed';
    selBox.style.left = '0';
    selBox.style.top = '0';
    selBox.style.opacity = '0';
    selBox.value = value;
    document.body.appendChild(selBox);
    selBox.focus();
    selBox.select();
    document.execCommand('copy');
    document.body.removeChild(selBox);
    this.copiedValue.emit(value);
    setTimeout(() => {
      this.copiedValue.emit('');
    }, 3000);
  }
}
