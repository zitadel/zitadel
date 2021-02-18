import { Component, EventEmitter, Input, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';

@Component({
    selector: 'cnsl-redirect-uris',
    templateUrl: './redirect-uris.component.html',
    styleUrls: ['./redirect-uris.component.scss']
})
export class RedirectUrisComponent implements OnInit {
    @Input() title: string = '';
    @Input() devMode: boolean = false;
    @Input() canWrite: boolean = false;
    @Input() public urisList: string[] = [];
    @Input() public redirectControl: FormControl = new FormControl({ value: '', disabled: true });
    @Input() public changedUris: EventEmitter<string[]> = new EventEmitter();
    constructor() { }

    ngOnInit(): void {
        if (this.canWrite) {
            this.redirectControl.enable();
        }
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
        console.log(redirect);
        const index = this.urisList.indexOf(redirect);

        if (index >= 0) {
            this.urisList.splice(index, 1);
        }
    }
}
