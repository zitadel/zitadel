import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Component({
    selector: 'app-contributors',
    templateUrl: './contributors.component.html',
    styleUrls: ['./contributors.component.scss'],
    animations: [
        trigger('list', [
            transition(':enter', [
                query('@animate',
                    stagger(80, animateChild()),
                ),
            ]),
        ]),
        trigger('animate', [
            transition(':enter', [
                style({ opacity: 0, transform: 'translateX(100%)' }),
                animate('100ms', style({ opacity: 1, transform: 'translateX(0)' })),
            ]),
        ]),
    ],
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
