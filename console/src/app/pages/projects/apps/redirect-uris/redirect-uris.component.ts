import { Component, EventEmitter, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Observable, Subscription } from 'rxjs';

@Component({
    selector: 'cnsl-redirect-uris',
    templateUrl: './redirect-uris.component.html',
    styleUrls: ['./redirect-uris.component.scss'],
})
export class RedirectUrisComponent implements OnInit, OnDestroy {
    @Input() title: string = '';
    @Input() devMode: boolean = false;
    @Input() canWrite: boolean = false;
    @Input() isNative!: boolean;
    @Input() public urisList: string[] = [];
    @Input() public redirectControl: FormControl = new FormControl({ value: '', disabled: true });
    @Input() public changedUris: EventEmitter<string[]> = new EventEmitter();
    @Input() public getValues: Observable<void> = new Observable();

    @ViewChild('redInput') input!: any;
    private sub: Subscription = new Subscription();
    constructor() { }

    ngOnInit(): void {
        if (this.canWrite) {
            this.redirectControl.enable();
        }

        this.sub = this.getValues.subscribe(() => {
            this.add(this.input.nativeElement);
        });
    }

    ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    public add(input: any): void {
        if (this.redirectControl.valid) {
            if (input.value !== '' && input.value !== ' ' && input.value !== '/') {
                this.urisList.push(input.value);
            }
            if (input) {
                input.value = '';
            }
        }
    }

    public remove(redirect: any): void {
        const index = this.urisList.indexOf(redirect);

        if (index >= 0) {
            this.urisList.splice(index, 1);
        }
    }
}
