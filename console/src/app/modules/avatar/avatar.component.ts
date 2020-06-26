import { Component, Input, OnInit } from '@angular/core';

@Component({
    selector: 'app-avatar',
    templateUrl: './avatar.component.html',
    styleUrls: ['./avatar.component.scss']
})
export class AvatarComponent implements OnInit {
    @Input() name: string = '';
    @Input() credentials: string = '';
    @Input() size: number = 24;
    @Input() fontSize: number = 16;
    @Input() active: boolean = false;
    constructor() { }

    ngOnInit(): void {
        if (!this.credentials) {
            console.log(this.name);
            const split: string[] = this.name.split(' ');
            this.credentials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : '');
        }

        if (this.size > 50) {
            this.fontSize = 32;
        }
    }
}
