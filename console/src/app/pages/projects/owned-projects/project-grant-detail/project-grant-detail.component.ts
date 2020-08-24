import { Component } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { ActivatedRoute } from '@angular/router';
import {
    ProjectGrant,
    ProjectGrantState,
    ProjectGrantView,
    ProjectRoleView,
    ProjectType,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-project-grant-detail',
    templateUrl: './project-grant-detail.component.html',
    styleUrls: ['./project-grant-detail.component.scss'],
})
export class ProjectGrantDetailComponent {
    public grant!: ProjectGrantView.AsObject;
    public projectid: string = '';
    public grantid: string = '';

    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;

    public isZitadel: boolean = false;
    ProjectGrantState: any = ProjectGrantState;

    public memberRoleOptions: ProjectRoleView.AsObject[] = [];

    constructor(
        private orgService: ManagementService,
        private projectService: ManagementService,
        private route: ActivatedRoute,
        private toast: ToastService,
    ) {
        this.route.params.subscribe(params => {
            this.projectid = params.projectid;
            this.grantid = params.grantid;

            this.orgService.GetIam().then(iam => {
                this.isZitadel = iam.toObject().iamProjectId === this.projectid;
            });

            this.getRoleOptions(params.projectid);

            this.projectService.ProjectGrantByID(this.grantid, this.projectid).then((grant) => {
                this.grant = grant.toObject();
            });
        });
    }

    public changeState(newState: ProjectGrantState): void {
        if (newState === ProjectGrantState.PROJECTGRANTSTATE_ACTIVE) {
            this.projectService.ReactivateProjectGrant(this.grantid, this.projectid).then(() => {
                this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
                this.grant.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (newState === ProjectGrantState.PROJECTGRANTSTATE_INACTIVE) {
            this.projectService.DeactivateProjectGrant(this.grantid, this.projectid).then(() => {
                this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
                this.grant.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public getRoleOptions(projectId: string): void {
        this.projectService.SearchProjectRoles(projectId, 100, 0).then(resp => {
            this.memberRoleOptions = resp.toObject().resultList;
        });
    }

    updateRoles(selectionChange: MatSelectChange): void {
        this.projectService.UpdateProjectGrant(this.grant.id, this.grant.projectId, selectionChange.value)
            .then((newgrant: ProjectGrant) => {
                this.toast.showInfo('PROJECT.TOAST.GRANTUPDATED');
            }).catch(error => {
                this.toast.showError(error);
            });
    }
}
