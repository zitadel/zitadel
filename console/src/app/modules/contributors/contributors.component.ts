import { Component, EventEmitter, Input, Output } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Component({
    selector: 'app-contributors',
    templateUrl: './contributors.component.html',
    styleUrls: ['./contributors.component.scss'],
})
export class ContributorsComponent {
    @Input() title: string = '';
    @Input() description: string = '';
    @Input() disabled: boolean = false;
    @Input() totalResult: number = 0;
    @Input() loading: boolean = false;
    @Input() membersSubject!: BehaviorSubject<any[]>;
    @Output() addClicked: EventEmitter<void> = new EventEmitter();
    @Output() showDetailClicked: EventEmitter<void> = new EventEmitter();

    public emitAddMember(): void {
        this.addClicked.emit();
    }

    public emitShowDetail(): void {
        this.showDetailClicked.emit();
    }
}
