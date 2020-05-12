import { SelectionModel } from '@angular/cdk/collections';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Params } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { from, Observable, of, Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import {
    Application,
    ApplicationSearchResponse,
    Project,
    ProjectMember,
    ProjectMemberSearchResponse,
    ProjectRole,
    ProjectRoleSearchResponse,
    ProjectState,
    ProjectType,
} from 'src/app/proto/generated/management_pb';
import { GrpcService } from 'src/app/services/grpc.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-project-detail',
    templateUrl: './project-detail.component.html',
    styleUrls: ['./project-detail.component.scss'],

})
export class ProjectDetailComponent implements OnInit, OnDestroy {
    public projectId: string = '';
    public project!: Project.AsObject;

    public pageSizeRoles: number = 10;
    public roleDataSource: MatTableDataSource<ProjectRole.AsObject> = new MatTableDataSource<ProjectRole.AsObject>();
    public roleResult!: ProjectRoleSearchResponse.AsObject;
    public roleColumns: string[] = ['name', 'displayname', 'group', 'actions'];

    public pageSizeMembers: number = 10;
    public projectDataSource: MatTableDataSource<ProjectMember.AsObject> = new MatTableDataSource<ProjectMember.AsObject>();
    public memberResult!: ProjectMemberSearchResponse.AsObject;
    public memberColumns: string[] = ['firstname', 'lastname', 'username', 'email', 'roles'];
    public selection: SelectionModel<ProjectMember.AsObject> = new SelectionModel<ProjectMember.AsObject>(true, []);

    public pageSizeApps: number = 10;
    public appsDataSource: MatTableDataSource<Application.AsObject> = new MatTableDataSource<Application.AsObject>();
    public appsResult!: ApplicationSearchResponse.AsObject;
    public appsColumns: string[] = ['name'];

    public ProjectState: any = ProjectState;
    public ProjectType: any = ProjectType;
    public ChangeType: any = ChangeType;

    public grid: boolean = true;
    private subscription?: Subscription;
    public editstate: boolean = false;

    public isZitadel$: Observable<boolean> = of(false);

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private projectService: ProjectService,
        private _location: Location,
        private grpcService: GrpcService,
    ) {
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    private async getData({ id }: Params): Promise<void> {
        this.projectId = id;
        this.projectService.GetProjectById(id).then(proj => {
            this.project = proj.toObject();
            if (this.project.type !== ProjectType.PROJECTTYPE_SELF ||
                this.project.state === ProjectState.PROJECTSTATE_INACTIVE ||
                this.project.state === ProjectState.PROJECTSTATE_UNSPECIFIED) {
            }

            this.isZitadel$ = from(this.projectService.SearchApplications(this.project.id, 100, 0).then(appsResp => {
                const ret = appsResp.toObject().resultList
                    .filter(app => app.oidcConfig?.clientId === this.grpcService.clientid).length > 0;
                return ret;
            })); // TODO: replace with prettier thing
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }


    public changeState(newState: ProjectState): void {
        if (newState === ProjectState.PROJECTSTATE_ACTIVE) {
            this.projectService.ReactivateProject(this.projectId).then(() => {
                this.toast.showInfo('Reactivated Project');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        } else if (newState === ProjectState.PROJECTSTATE_INACTIVE) {
            this.projectService.DeactivateProject(this.projectId).then(() => {
                this.toast.showInfo('Deactivated Project');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        }
    }

    public saveProject(): void {
        this.projectService.UpdateProject(this.project).then(() => {
            this.toast.showInfo('Project updated');
        }).catch(error => {
            this.toast.showInfo(error.message);
        });
    }

    public navigateBack(): void {
        this._location.back();
    }

    public updateName(): void {
        this.saveProject();
        this.editstate = false;
    }
}
