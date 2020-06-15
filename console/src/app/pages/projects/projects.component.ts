import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';

@Component({
    selector: 'app-projects',
    templateUrl: './projects.component.html',
    styleUrls: ['./projects.component.scss'],
})
export class ProjectsComponent implements OnInit, OnDestroy {
    // public projectId: string = '';
    // public grantId: string = '';
    private sub: Subscription = new Subscription();
    constructor(private route: ActivatedRoute,
    ) {
        // this.route.params.subscribe((params) => {
        //     this.projectId = params.projectId;
        //     this.grantId = params.grantId;
        // });
    }

    ngOnInit(): void {
    }


    public ngOnDestroy(): void {
        // this.sub.unsubscribe();
    }
}
