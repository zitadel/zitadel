import { animate, animation, keyframes, style, transition, trigger, useAnimation } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { RefreshService } from 'src/app/services/refresh.service';

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
    selector: 'app-refresh-table',
    templateUrl: './refresh-table.component.html',
    styleUrls: ['./refresh-table.component.scss'],
    animations: [
        trigger('rotate', [
            transition('* => *', [useAnimation(rotate, { params: { time: '1s' } })]),
        ]),
    ],
})
export class RefreshTableComponent implements OnInit {
    @Input() public selection: SelectionModel<any> = new SelectionModel<any>(true, []);
    @Input() public timestamp!: Timestamp.AsObject;
    @Input() public dataSize: number = 0;
    @Input() public emitRefreshAfterTimeoutInMs: number = 0;
    @Input() public loading: boolean = false;
    @Input() public emitRefreshOnPreviousRoutes: string[] = [];
    @Output() public refreshed: EventEmitter<void> = new EventEmitter();

    constructor(private refreshService: RefreshService) { }

    ngOnInit(): void {
        if (this.emitRefreshAfterTimeoutInMs) {
            setTimeout(() => {
                this.emitRefresh();
            }, this.emitRefreshAfterTimeoutInMs);
        }

        if (this.emitRefreshOnPreviousRoutes.length && this.refreshService.previousUrls
            .some(url => this.emitRefreshOnPreviousRoutes.includes(url))) {
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
