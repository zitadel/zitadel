import { Component, Input, OnInit } from '@angular/core';


export interface CnslLinks {
    i18nTitle: string;
    i18nDesc: string;
    routerLink?: any;
    href?: string;
    withRole?: Array<string | RegExp>;
}

@Component({
    selector: 'cnsl-links',
    templateUrl: './links.component.html',
    styleUrls: ['./links.component.scss'],
})
export class LinksComponent implements OnInit {
    @Input() links: Array<CnslLinks> = [];
    constructor() { }

    ngOnInit(): void {
    }

}
