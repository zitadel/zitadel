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
    selector: 'app-owned-project-detail',
    templateUrl: './owned-project-detail.component.html',
    styleUrls: ['./owned-project-detail.component.scss'],

})
export class OwnedProjectDetailComponent implements OnInit, OnDestroy {
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
    private zitadelsub: Subscription = new Subscription();

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
        this.zitadelsub.unsubscribe();
    }

    private async getData({ id }: Params): Promise<void> {
        this.projectId = id;

        if (this.projectId) {
            this.projectService.GetProjectById(id).then(proj => {
                this.project = proj.toObject();
                this.isZitadel$ = from(this.projectService.SearchApplications(this.project.id, 100, 0).then(appsResp => {
                    const ret = appsResp.toObject().resultList
                        .filter(app => app.oidcConfig?.clientId === this.grpcService.clientid).length > 0;
                    return ret;
                }));
            }).catch(error => {
                this.toast.showError(error.message);
            });
        }

        this.zitadelsub = this.isZitadel$.subscribe(isZita => console.log(`zitade: ${isZita}`));
    }

    public changeState(newState: ProjectState): void {
        if (newState === ProjectState.PROJECTSTATE_ACTIVE) {
            this.projectService.ReactivateProject(this.projectId).then(() => {
                this.toast.showInfo('Reactivated Project');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        } else if (newState === ProjectState.PROJECTSTATE_INACTIVE) {
            this.toast.showInfo('You cant update this project.');
        }
    }

    public saveProject(): void {
        this.projectService.UpdateProject(this.project as Project.AsObject).then(() => {
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
