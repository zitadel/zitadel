import { Component, Input, OnInit } from '@angular/core';

@Component({
    selector: 'app-avatar',
    templateUrl: './avatar.component.html',
    styleUrls: ['./avatar.component.scss'],
})
export class AvatarComponent implements OnInit {
    @Input() name: string = '';
    @Input() credentials: string = '';
    @Input() size: number = 24;
    @Input() fontSize: number = 16;
    @Input() active: boolean = false;
    @Input() color: string = '';
    constructor() { }

    ngOnInit(): void {
        if (!this.credentials) {
            const split: string[] = this.name.split(' ');
            this.credentials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : '');
            if (!this.color) {
                this.color = this.getColor(this.name);
            }
        }

        if (this.size > 50) {
            this.fontSize = 32;
        }
    }

    getColor(userName: string): string {
        const colors = [
            '#e51c23',
            '#e91e63',
            '#9c27b0',
            '#673ab7',
            '#3f51b5',
            '#5677fc',
            '#03a9f4',
            '#00bcd4',
            '#009688',
            '#259b24',
            '#8bc34a',
            '#afb42b',
            '#ff9800',
            '#ff5722',
            '#795548',
            '#607d8b',
        ];

        let hash = 0;
        if (userName.length === 0) {
            return colors[hash];
        }
        for (let i = 0; i < userName.length; i++) {
            // tslint:disable-next-line: no-bitwise
            hash = userName.charCodeAt(i) + ((hash << 5) - hash);
            // tslint:disable-next-line: no-bitwise
            hash = hash & hash;
        }
        hash = ((hash % colors.length) + colors.length) % colors.length;
        return colors[hash];
    }
}
