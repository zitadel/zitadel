import { Component, OnInit } from '@angular/core';
import { View } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';

@Component({
    selector: 'app-iam-views',
    templateUrl: './iam-views.component.html',
    styleUrls: ['./iam-views.component.scss']
})
export class IamViewsComponent implements OnInit {
    public views: View.AsObject[] = [];
    constructor(private adminService: AdminService) {
        this.getViews();
    }

    ngOnInit(): void {
    }

    public getViews(): void {
        this.adminService.GetViews().then(views => {
            this.views = views.toObject().viewsList;
        });
    }

}
