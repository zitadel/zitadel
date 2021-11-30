import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Domain } from 'src/app/proto/generated/zitadel/org_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddDomainDialogComponent } from '../org-detail/add-domain-dialog/add-domain-dialog.component';
import { DomainVerificationComponent } from '../org-detail/domain-verification/domain-verification.component';

@Component({
  selector: 'cnsl-org-domains',
  templateUrl: './org-domains.component.html',
  styleUrls: ['./org-domains.component.scss'],
})
export class OrgDomainsComponent implements OnInit {
  public domains: Domain.AsObject[] = [];
  public primaryDomain: string = '';
  public InfoSectionType: any = InfoSectionType;

  constructor(private mgmtService: ManagementService, private toast: ToastService, private dialog: MatDialog) {}

  ngOnInit(): void {
    this.loadDomains();
  }

  public loadDomains(): void {
    this.mgmtService.listOrgDomains().then((result) => {
      this.domains = result.resultList;
      this.primaryDomain = this.domains.find((domain) => domain.isPrimary)?.domainName ?? '';
    });
  }

  public setPrimary(domain: Domain.AsObject): void {
    this.mgmtService
      .setPrimaryOrgDomain(domain.domainName)
      .then(() => {
        this.toast.showInfo('ORG.TOAST.SETPRIMARY', true);
        this.loadDomains();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public addNewDomain(): void {
    const dialogRef = this.dialog.open(AddDomainDialogComponent, {
      data: {},
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .addOrgDomain(resp)
          .then(() => {
            this.toast.showInfo('ORG.TOAST.DOMAINADDED', true);

            setTimeout(() => {
              this.loadDomains();
            }, 1000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public removeDomain(domain: string): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'ORG.DOMAINS.DELETE.TITLE',
        descriptionKey: 'ORG.DOMAINS.DELETE.DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .removeOrgDomain(domain)
          .then(() => {
            this.toast.showInfo('ORG.TOAST.DOMAINREMOVED', true);
            const index = this.domains.findIndex((d) => d.domainName === domain);
            if (index > -1) {
              this.domains.splice(index, 1);
            }
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public verifyDomain(domain: Domain.AsObject): void {
    const dialogRef = this.dialog.open(DomainVerificationComponent, {
      data: {
        domain: domain,
      },
      width: '500px',
    });

    dialogRef.afterClosed().subscribe((reload: boolean) => {
      if (reload) {
        this.loadDomains();
      }
    });
  }
}
