import { Component, OnInit } from '@angular/core';
import { AdminService } from 'src/app/services/admin.service';

@Component({
    selector: 'app-iam',
    templateUrl: './iam.component.html',
    styleUrls: ['./iam.component.scss']
})
export class IamComponent implements OnInit {

    constructor(private adminService: AdminService) {

    }

    ngOnInit(): void {
    }

}
