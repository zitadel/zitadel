import { BreakpointObserver } from '@angular/cdk/layout';
import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Component({
    selector: 'app-meta-layout',
    templateUrl: './meta-layout.component.html',
    styleUrls: ['./meta-layout.component.scss'],
})
export class MetaLayoutComponent {

    constructor(private breakpointObserver: BreakpointObserver) {
        this.isSmallScreen$.subscribe(small => this.hidden = small);
    }
    public hidden: boolean = false;
    public isSmallScreen$: Observable<boolean> = this.breakpointObserver
        .observe('(max-width: 1000px)')
        .pipe(map(result => {
            return result.matches;
        }));
}
