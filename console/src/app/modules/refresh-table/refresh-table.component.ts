import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
    selector: 'app-refresh-table',
    templateUrl: './refresh-table.component.html',
    styleUrls: ['./refresh-table.component.scss'],
})
export class RefreshTableComponent implements OnInit {
    @Input() public selection: SelectionModel<any> = new SelectionModel<any>(true, []);
    @Input() public dataSize: number = 0;
    @Output() public refreshed: EventEmitter<void> = new EventEmitter();
    constructor() { }

    ngOnInit(): void {
    }

    emitRefresh(): void {
        return this.refreshed.emit();
    }
}
