import { Directive, EventEmitter, HostListener, Output } from '@angular/core';

@Directive({
  selector: '[cnslDropzone]',
})
export class DropzoneDirective {
  @Output() dropped: EventEmitter<FileList> = new EventEmitter<FileList>();
  @Output() hovered: EventEmitter<boolean> = new EventEmitter<boolean>();

  @HostListener('drop', ['$event'])
  onDrop($event: DragEvent): void {
    $event.preventDefault();
    this.dropped.emit($event.dataTransfer?.files);
    this.hovered.emit(false);
  }

  @HostListener('dragover', ['$event'])
  onDragOver($event: any): void {
    $event.preventDefault();
    this.hovered.emit(true);
  }

  @HostListener('dragleave', ['$event'])
  onDragLeave($event: any): void {
    $event.preventDefault();
    this.hovered.emit(false);
  }
}
