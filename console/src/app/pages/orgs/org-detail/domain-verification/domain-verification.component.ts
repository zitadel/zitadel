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

    public validating: boolean = false;
    constructor(
        private toast: ToastService,
        public dialogRef: MatDialogRef<DomainVerificationComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
        private mgmtService: ManagementService,
    ) {
        this.domain = data.domain;
        if (this.domain.validationType === OrgDomainValidationType.ORGDOMAINVALIDATIONTYPE_UNSPECIFIED) {
            this.showNew = true;
        }
    }

    async loadHttpToken(): Promise<void> {
        this.mgmtService.GenerateMyOrgDomainValidation(
            this.domain.domain,
            OrgDomainValidationType.ORGDOMAINVALIDATIONTYPE_HTTP).then((http) => {
                this.http = http.toObject();
            });
    }

    async loadDnsToken(): Promise<void> {
        this.mgmtService.GenerateMyOrgDomainValidation(
            this.domain.domain,
            OrgDomainValidationType.ORGDOMAINVALIDATIONTYPE_DNS).then((dns) => {
                this.dns = dns.toObject();
            });
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public validate(): void {
        this.validating = true;
        this.mgmtService.ValidateMyOrgDomain(this.domain.domain).then(() => {
            this.dialogRef.close(true);
            this.toast.showInfo('ORG.PAGES.ORGDOMAIN.VERIFICATION_SUCCESSFUL', true);
            this.validating = false;
        }).catch((error) => {
            this.toast.showError(error);
            this.validating = false;
        });
    }

    public saveFile(): void {
        const blob = new Blob([this.http.token], { type: 'text/plain;charset=utf-8' });
        saveAs(blob, this.http.token + '.txt');
    }
}
