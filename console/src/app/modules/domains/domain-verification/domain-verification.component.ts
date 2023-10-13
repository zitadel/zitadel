import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { saveAs } from 'file-saver';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { GenerateOrgDomainValidationResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Domain, DomainValidationType } from 'src/app/proto/generated/zitadel/org_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-domain-verification',
  templateUrl: './domain-verification.component.html',
  styleUrls: ['./domain-verification.component.scss'],
})
export class DomainVerificationComponent {
  public domain?: Domain.AsObject;

  public DomainValidationType: any = DomainValidationType;

  public http?: GenerateOrgDomainValidationResponse.AsObject;
  public dns?: GenerateOrgDomainValidationResponse.AsObject;

  public copied: string = '';

  public showNew: boolean = false;

  public validating: boolean = false;
  public InfoSectionType: any = InfoSectionType;

  constructor(
    private toast: ToastService,
    public dialogRef: MatDialogRef<DomainVerificationComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private mgmtService: ManagementService,
  ) {
    this.domain = data.domain;

    if (this.domain?.validationType === DomainValidationType.DOMAIN_VALIDATION_TYPE_UNSPECIFIED) {
      this.showNew = true;
    }
  }

  async loadHttpToken(): Promise<void> {
    if (this.domain) {
      this.mgmtService
        .generateOrgDomainValidation(this.domain.domainName, DomainValidationType.DOMAIN_VALIDATION_TYPE_HTTP)
        .then((http) => {
          this.http = http;
        });
    }
  }

  async loadDnsToken(): Promise<void> {
    if (this.domain) {
      this.mgmtService
        .generateOrgDomainValidation(this.domain.domainName, DomainValidationType.DOMAIN_VALIDATION_TYPE_DNS)
        .then((dns) => {
          this.dns = dns;
        });
    }
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public validate(): void {
    this.validating = true;
    if (this.domain) {
      this.mgmtService
        .validateOrgDomain(this.domain.domainName)
        .then(() => {
          this.dialogRef.close(true);
          this.toast.showInfo('ORG.PAGES.ORGDOMAIN.VERIFICATION_SUCCESSFUL', true);
          this.validating = false;
        })
        .catch((error) => {
          this.toast.showError(error);
          this.validating = false;
        });
    }
  }

  public saveFile(): void {
    if (this.http) {
      const blob = new Blob([this.http.token], { type: 'text/plain;charset=utf-8' });
      saveAs(blob, this.http.token + '.txt');
    }
  }
}
