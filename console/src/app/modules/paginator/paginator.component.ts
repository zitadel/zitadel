import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';

export interface PageEvent {
  length: number;
  pageSize: number;
  pageIndex: number;
  pageSizeOptions: Array<number>;
}

@Component({
  selector: 'cnsl-paginator',
  templateUrl: './paginator.component.html',
  styleUrls: ['./paginator.component.scss'],
})
export class PaginatorComponent {
  @Input() public timestamp: Timestamp.AsObject | undefined = undefined;
  @Input() public length: number = 0;
  @Input() public pageSize: number = 10;
  @Input() public pageIndex: number = 0;
  @Input() public pageSizeOptions: Array<number> = [10, 25, 50];
  @Input() public hidePagination: boolean = false;
  @Input() public showMoreButton: boolean = false;
  @Input() public disableShowMore: boolean | null = false;
  @Output() public moreRequested: EventEmitter<void> = new EventEmitter();
  @Output() public page: EventEmitter<PageEvent> = new EventEmitter();
  constructor() {}

  public previous(): void {
    if (this.previousPossible) {
      this.pageIndex = this.pageIndex - 1;
      this.emitChange();
    }
  }

  public next(): void {
    if (this.nextPossible) {
      this.pageIndex = this.pageIndex + 1;
      this.emitChange();
    }
  }

  get previousPossible(): boolean {
    const temp = this.pageIndex - 1;
    return temp >= 0;
  }

  get nextPossible(): boolean {
    const temp = this.pageIndex + 1;
    return temp <= this.length / this.pageSize;
  }

  get startIndex(): number {
    return this.pageIndex * this.pageSize;
  }

  get endIndex(): number {
    const max = this.startIndex + this.pageSize;
    return this.length < max ? this.length : max;
  }

  public emitChange(): void {
    this.page.emit({
      length: this.length,
      pageSize: this.pageSize,
      pageIndex: this.pageIndex,
      pageSizeOptions: this.pageSizeOptions,
    });
  }

  public updatePageSize(newSize: number): void {
    this.pageSize = newSize;
    this.pageIndex = 0;
    this.emitChange();
  }
}
