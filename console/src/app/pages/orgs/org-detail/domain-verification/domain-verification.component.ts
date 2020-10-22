import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { saveAs } from 'file-saver';
import { OrgDomainValidationResponse, OrgDomainValidationType, OrgDomainView } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-domain-verification',
    templateUrl: './domain-verification.component.html',
    styleUrls: ['./domain-verification.component.scss'],
})
export class DomainVerificationComponent {
    public domain!: OrgDomainView.AsObject;

    public OrgDomainValidationType: any = OrgDomainValidationType;

    public http!: OrgDomainValidationResponse.AsObject;
    public dns!: OrgDomainValidationResponse.AsObject;
    public copied: string = '';

    public showNew: boolean = false;

    constructor(
        private toast: ToastService,
        public dialogRef: MatDialogRef<DomainVerificationComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        private mgmtService: ManagementService,
    ) {
        this.domain = data.domain;
        console.log(data);
        if (this.domain.validationType === OrgDomainValidationType.ORGDOMAINVALIDATIONTYPE_UNSPECIFIED) {
            this.showNew = true;
        }
    }

    async loadHttpToken(): Promise<void> {
        this.http = (await this.mgmtService.GenerateMyOrgDomainValidation(
            this.domain.domain,
            OrgDomainValidationType.ORGDOMAINVALIDATIONTYPE_HTTP)).toObject();
    }

    async loadDnsToken(): Promise<void> {
        this.dns = (await this.mgmtService.GenerateMyOrgDomainValidation(
            this.domain.domain,
            OrgDomainValidationType.ORGDOMAINVALIDATIONTYPE_DNS)).toObject();
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public validate(): void {
        this.mgmtService.ValidateMyOrgDomain(this.domain.domain).then(() => {
            this.dialogRef.close(false);
        }).catch((error) => {
            this.toast.showError(error);
        });
    }

    public saveFile(): void {
        const blob = new Blob([this.http.token], { type: 'text/plain;charset=utf-8' });
        saveAs(blob, this.http.token + '.txt');
    }
}
