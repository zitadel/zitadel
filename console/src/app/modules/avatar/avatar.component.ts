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
            'linear-gradient(40deg, #B44D51 30%, rgb(241,138,138))',
            'linear-gradient(40deg, #B75073 30%, rgb(234,96,143))',
            'linear-gradient(40deg, #84498E 30%, rgb(214,116,230))',
            'linear-gradient(40deg, #705998 30%, rgb(163,131,220))',
            'linear-gradient(40deg, #5C6598 30%, rgb(135,148,222))',
            'linear-gradient(40deg, #7F90D3 30%, rgb(181,196,247))',
            'linear-gradient(40deg, #3E93B9 30%, rgb(150,215,245))',
            'linear-gradient(40deg, #3494A0 30%, rgb(71,205,222))',
            'linear-gradient(40deg, #25716A 30%, rgb(58,185,173))',
            'linear-gradient(40deg, #427E41 30%, rgb(97,185,96))',
            'linear-gradient(40deg, #89A568 30%, rgb(176,212,133))',
            'linear-gradient(40deg, #90924D 30%, rgb(187,189,98))',
            'linear-gradient(40deg, #E2B032 30%, rgb(245,203,99))',
            'linear-gradient(40deg, #C97358 30%, rgb(245,148,118))',
            'linear-gradient(40deg, #6D5B54 30%, rgb(152,121,108))',
            'linear-gradient(40deg, #6B7980 30%, rgb(134,163,177))',
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
