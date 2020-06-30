import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { UserGrantView } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';

@Component({
    selector: 'app-user-grant',
    templateUrl: './user-grant.component.html',
    styleUrls: ['./user-grant.component.scss'],
})
export class UserGrantComponent implements OnInit {
    public userid: string = '';
    public grantid: string = '';

    public grantView!: UserGrantView.AsObject;
    constructor(
        private mgmtUserService: MgmtUserService,
        private route: ActivatedRoute) {
        this.route.params.subscribe(params => {
            this.userid = params.projectid;
            this.grantid = params.grantid;

            this.mgmtUserService.UserGrantByID(this.grantid, this.userid).then(resp => {
                this.grantView = resp.toObject();
                console.log(this.grantView);
            });
        });
    }

    ngOnInit(): void {
    }

}
