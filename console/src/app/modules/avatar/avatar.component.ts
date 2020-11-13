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
    @Input() fontSize: number = 14;
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
            '#B44D51',
            '#B75073',
            '#84498E',
            '#705998',
            '#5C6598',
            '#7F90D3',
            '#3E93B9',
            '#3494A0',
            '#25716A',
            '#427E41',
            '#89A568',
            '#90924D',
            '#E2B032',
            '#C97358',
            '#6D5B54',
            '#6B7980',
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
