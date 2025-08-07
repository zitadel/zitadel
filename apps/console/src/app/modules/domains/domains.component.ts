import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Domain, DomainValidationType } from 'src/app/proto/generated/zitadel/org_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { AddDomainDialogComponent } from './add-domain-dialog/add-domain-dialog.component';
import { DomainVerificationComponent } from './domain-verification/domain-verification.component';

@Component({
  selector: 'cnsl-domains',
  templateUrl: './domains.component.html',
  styleUrls: ['./domains.component.scss'],
})
export class DomainsComponent implements OnInit {
  public domains: Domain.AsObject[] = [];
  public primaryDomain: string = '';
  public InfoSectionType: any = InfoSectionType;
  public verifyOrgDomains: boolean | undefined;

  constructor(
    private mgmtService: ManagementService,
    private toast: ToastService,
    private dialog: MatDialog,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }

  ngOnInit(): void {
    this.loadDomains();
  }

  public loadDomains(): void {
    this.mgmtService.getDomainPolicy().then((result) => {
      this.verifyOrgDomains = result.policy?.validateOrgDomains;
    });

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
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((domainName) => {
      if (domainName) {
        this.mgmtService
          .addOrgDomain(domainName)
          .then(() => {
            this.toast.showInfo('ORG.TOAST.DOMAINADDED', true);
            if (this.verifyOrgDomains) {
              this.verifyDomain({
                domainName: domainName,
                validationType: DomainValidationType.DOMAIN_VALIDATION_TYPE_UNSPECIFIED,
              });
            } else {
              this.loadDomains();
            }
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

    dialogRef.afterClosed().subscribe((del) => {
      if (del) {
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

  public verifyDomain(domain: Partial<Domain.AsObject>): void {
    const dialogRef = this.dialog.open(DomainVerificationComponent, {
      data: {
        domain: domain,
      },
      width: '500px',
    });

    dialogRef.afterClosed().subscribe(() => {
      this.loadDomains();
    });
  }
}
