import { animate, animateChild, keyframes, query, stagger, style, transition, trigger } from '@angular/animations';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';

@Component({
    selector: 'app-contributors',
    templateUrl: './contributors.component.html',
    styleUrls: ['./contributors.component.scss'],
    animations: [
        trigger('cardAnimation', [
            transition('* => *', [
                query('@animate', stagger('40ms', animateChild()), { optional: true }),
            ]),
        ]),
        trigger('animate', [
            transition(':enter', [
                animate('.2s ease-in', keyframes([
                    style({ opacity: 0, offset: 0 }),
                    style({ opacity: .5, transform: 'scale(1.05)', offset: 0.3 }),
                    style({ opacity: 1, transform: 'scale(1)', offset: 1 }),
                ])),
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
    @Input() membersSubject!: BehaviorSubject<Member.AsObject[]>;
    @Output() addClicked: EventEmitter<void> = new EventEmitter();
    @Output() showDetailClicked: EventEmitter<void> = new EventEmitter();
    @Output() refreshClicked: EventEmitter<void> = new EventEmitter();

    public emitAddMember(): void {
        this.addClicked.emit();
    }

    public emitShowDetail(): void {
        this.showDetailClicked.emit();
    }

    public emitRefresh(): void {
        this.refreshClicked.emit();
    }
}
