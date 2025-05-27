import { animate, animation, keyframes, style, transition, trigger, useAnimation } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { RefreshService } from 'src/app/services/refresh.service';
import { Timestamp as ConnectTimestamp } from '@bufbuild/protobuf/wkt';

import { ActionKeysType } from '../action-keys/action-keys.component';

const rotate = animation([
  animate(
    '{{time}} cubic-bezier(0.785, 0.135, 0.15, 0.86)',
    keyframes([
      style({
        transform: 'rotate(0deg)',
      }),
      style({
        transform: 'rotate(360deg)',
      }),
    ]),
  ),
]);
@Component({
  selector: 'cnsl-refresh-table',
  templateUrl: './refresh-table.component.html',
  styleUrls: ['./refresh-table.component.scss'],
  animations: [trigger('rotate', [transition('* => *', [useAnimation(rotate, { params: { time: '1s' } })])])],
})
export class RefreshTableComponent implements OnInit {
  @Input() public selection: SelectionModel<any> = new SelectionModel<any>(true, []);
  @Input() public timestamp: Timestamp.AsObject | ConnectTimestamp | undefined = undefined;
  @Input() public emitRefreshAfterTimeoutInMs: number = 0;
  @Input() public loading: boolean | null = false;
  @Input() public emitRefreshOnPreviousRoutes: string[] = [];
  @Output() public refreshed: EventEmitter<void> = new EventEmitter();
  @Input() public hideRefresh: boolean = false;
  @Input() public showBorder: boolean = false;
  @Input() public showSelectionActionButton: boolean = true;

  public ActionKeysType: any = ActionKeysType;
  constructor(private refreshService: RefreshService) {}

  ngOnInit(): void {
    if (this.emitRefreshAfterTimeoutInMs) {
      setTimeout(() => {
        this.emitRefresh();
      }, this.emitRefreshAfterTimeoutInMs);
    }

    if (
      this.emitRefreshOnPreviousRoutes.length &&
      this.refreshService.previousUrls.some((url) => this.emitRefreshOnPreviousRoutes.includes(url))
    ) {
      setTimeout(() => {
        this.emitRefresh();
      }, 1000);
    }
  }

  emitRefresh(): void {
    this.selection.clear();
    return this.refreshed.emit();
  }
}
