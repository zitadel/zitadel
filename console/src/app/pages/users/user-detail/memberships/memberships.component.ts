import { Component, Input, OnInit } from '@angular/core';
import {
    ProjectGrantMemberSearchKey,
    ProjectGrantMemberSearchQuery,
    ProjectGrantMemberView,
    ProjectMemberSearchKey,
    ProjectMemberSearchQuery,
    ProjectMemberView,
} from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';

@Component({
    selector: 'app-memberships',
    templateUrl: './memberships.component.html',
    styleUrls: ['./memberships.component.scss'],
})
export class MembershipsComponent implements OnInit {
    usergrants: ProjectGrantMemberView.AsObject[] = [];
    projectmembers: ProjectMemberView.AsObject[] = [];

    @Input() public userId: string = '';

    constructor(private mgmtUserService: MgmtUserService) { }

    ngOnInit(): void {
        this.loadManager(this.userId);
    }

    public async loadManager(userId: string): Promise<void> {
        console.log('load managers');
        // manager of granted project
        const projectGrantQuery = new ProjectGrantMemberSearchQuery();
        projectGrantQuery.setKey(ProjectGrantMemberSearchKey.PROJECTGRANTMEMBERSEARCHKEY_USER_ID);
        projectGrantQuery.setValue(userId);

        this.usergrants = (await this.mgmtUserService.SearchProjectGrantMembers(100, 0, [projectGrantQuery]))
            .toObject().resultList;
        console.log(this.usergrants);

        // manager of owned project
        const projectMemberQuery = new ProjectMemberSearchQuery();
        projectMemberQuery.setKey(ProjectMemberSearchKey.PROJECTMEMBERSEARCHKEY_USER_ID);
        projectMemberQuery.setValue(userId);

        this.projectmembers = (await this.mgmtUserService.SearchProjectMembers(100, 0, [projectMemberQuery]))
            .toObject().resultList;
        console.log(this.projectmembers);

        // manager of organization
        // const projectMemberQuery = new ProjectMemberSearchQuery();
        // projectMemberQuery.setKey(ProjectMemberSearchKey.PROJECTMEMBERSEARCHKEY_USER_ID);
        // projectMemberQuery.setValue(userId);

        // this.projectmembers = (await this.mgmtUserService.searchor(100, 0, [projectMemberQuery]))
        //     .toObject().resultList;
        // console.log(this.projectmembers);
    }
}
