import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Subscription } from 'rxjs';
import { Org, ProjectRole } from 'src/app/proto/generated/management_pb';
import { AuthService } from 'src/app/services/auth.service';
import { OrgService } from 'src/app/services/org.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-project-grant-create',
    templateUrl: './project-grant-create.component.html',
    styleUrls: ['./project-grant-create.component.scss'],
})
export class ProjectGrantCreateComponent implements OnInit, OnDestroy {
    public org!: Org.AsObject;
    public projectId: string = '';
    public grantId: string = '';
    public rolesList: string[] = [];

    public STEPS: number = 2;
    public currentCreateStep: number = 1;

    private routeSubscription: Subscription = new Subscription();
    constructor(
        private orgService: OrgService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private projectService: ProjectService,
        private authService: AuthService,
        private _location: Location,
    ) { }

    public ngOnInit(): void {
        this.routeSubscription = this.route.params.subscribe(params => {
            this.projectId = params.projectid;
            console.log(params);
        });
    }

    public ngOnDestroy(): void {
        this.routeSubscription.unsubscribe();
    }

    public searchOrg(domain: any): void {
        this.orgService.getOrgByDomainGlobal(domain.value).then((ret) => {
            console.log(ret.toObject());
            const tmp = ret.toObject();
            this.authService.GetActiveOrg().then((org) => {
                if (tmp !== org) {
                    this.org = tmp;
                }
            });
            this.org = ret.toObject();
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }

    public close(): void {
        this._location.back();
    }

    public addGrant(): void {
        this.projectService
            .CreateProjectGrant(this.org.id, this.projectId, this.rolesList)
            .then((data) => {
                console.log(data);
                this.close();
            })
            .catch(error => {
                this.toast.showError(error.message);
                console.log(error);
            });
    }


    public dateFromTimestamp(date: Timestamp.AsObject): any {
        const ts: Date = new Date(date.seconds * 1000 + date.nanos / 1000);
        return ts;
    }

    public selectRoles(roles: ProjectRole.AsObject[]): void {
        this.rolesList = roles.map(role => role.name);
    }

    public next(): void {
        this.currentCreateStep++;
    }

    public previous(): void {
        this.currentCreateStep--;
    }
}

