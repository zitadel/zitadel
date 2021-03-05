import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
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
    public rolesKeyList: string[] = [];

    public STEPS: number = 2;
    public currentCreateStep: number = 1;

    private routeSubscription: Subscription = new Subscription();
    constructor(
        private route: ActivatedRoute,
        private toast: ToastService,
        private mgmtService: ManagementService,
        private authService: GrpcAuthService,
        private _location: Location,
    ) { }

    public ngOnInit(): void {
        this.routeSubscription = this.route.params.subscribe(params => {
            this.projectId = params.projectid;
        });
    }

    public ngOnDestroy(): void {
        this.routeSubscription.unsubscribe();
    }

    public searchOrg(domain: string): void {
        this.mgmtService.getOrgByDomainGlobal(domain).then((ret) => {
            if (ret.org) {
                const tmp = ret.org;
                this.authService.getActiveOrg().then((org) => {
                    if (tmp !== org) {
                        this.org = tmp;
                    }
                });
                this.org = ret.org;
            }
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public close(): void {
        this._location.back();
    }

    public addGrant(): void {
        this.mgmtService
            .addProjectGrant(this.org.id, this.projectId, this.rolesKeyList)
            .then(() => {
                this.close();
            })
            .catch(error => {
                this.toast.showError(error);
            });
    }

    public selectRoles(roles: Role.AsObject[]): void {
        this.rolesKeyList = roles.map(role => role.key);
    }

    public next(): void {
        this.currentCreateStep++;
    }

    public previous(): void {
        this.currentCreateStep--;
    }
}

