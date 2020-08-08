import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { OrgDomainValidationType, OrgDomainView } from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';

@Component({
    selector: 'app-domain-verification',
    templateUrl: './domain-verification.component.html',
    styleUrls: ['./domain-verification.component.scss'],
})
export class DomainVerificationComponent implements OnInit {
    public domain!: OrgDomainView.AsObject;

    public OrgDomainValidationType: any = OrgDomainValidationType;
    constructor(
        public dialogRef: MatDialogRef<DomainVerificationComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        private orgService: OrgService,
    ) {
        this.domain = data.domain;
    }

    ngOnInit(): void {
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public verify(type: OrgDomainValidationType): void {
        this.orgService.GenerateMyOrgDomainValidation(this.domain.domain, type);
    }
}
