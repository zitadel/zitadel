import { Component, EventEmitter, inject, Input, Output } from '@angular/core';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';
import { Timestamp as ConnectTimestamp } from '@bufbuild/protobuf/wkt';
import { PaginationPreferenceService } from 'src/app/services/pagination-preference.service';

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
  standalone: false,
})
export class PaginatorComponent {
  @Input() public timestamp: Timestamp.AsObject | ConnectTimestamp | undefined = undefined;
  @Input() public length: number = 0;
  @Input() public pageSize: number = 10;
  @Input() public pageIndex: number = 0;
  @Input() public pageSizeOptions: Array<number> = [10, 25, 50];
  @Input() public hidePagination: boolean = false;
  @Input() public showMoreButton: boolean = false;
  @Input() public disableShowMore: boolean | null = false;
  /**
   * When set, the selected page size is remembered under this key across reloads.
   * The hosting table has to read the same key for its initial page size, otherwise
   * the first request would still use the hardcoded default.
   */
  @Input() public persistKey?: string;
  @Output() public moreRequested: EventEmitter<void> = new EventEmitter();
  @Output() public page: EventEmitter<PageEvent> = new EventEmitter();

  private readonly paginationPreference = inject(PaginationPreferenceService);

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
    // Comparing against length / pageSize used to leave "next" enabled on the last page
    // whenever length was an exact multiple of pageSize, leading to an empty page.
    return (this.pageIndex + 1) * this.pageSize < this.length;
  }

  /** Zero based offset of the first row on the current page, used for requests. */
  get startIndex(): number {
    return this.pageIndex * this.pageSize;
  }

  /** One based position of the first row on the current page, used for display. */
  get displayStartIndex(): number {
    return this.length === 0 ? 0 : this.startIndex + 1;
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
    if (this.persistKey) {
      this.paginationPreference.set(this.persistKey, newSize);
    }
    this.emitChange();
  }
}
