import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ProjectType } from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';

@Component({
    selector: 'app-grant',
    templateUrl: './grant.component.html',
    styleUrls: ['./grant.component.scss'],
})
export class GrantComponent {
    public projectid: string = '';
    public grantid: string = '';

    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;

    public isZitadel: boolean = false;

    constructor(
        private orgService: OrgService,
        private route: ActivatedRoute) {
        this.route.params.subscribe(params => {
            this.projectid = params.projectid;
            this.grantid = params.grantid;

            this.orgService.GetIam().then(iam => {
                this.isZitadel = iam.toObject().iamProjectId === this.projectid;
            });
        });
    }
}
