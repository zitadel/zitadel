import { LiveAnnouncer } from '@angular/cdk/a11y';
import { Component, computed, effect, signal } from '@angular/core';
import { Sort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { OrgQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { NewOrganizationService } from '../../services/new-organization.service';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { MessageInitShape } from '@bufbuild/protobuf';
import { ListOrganizationsRequestSchema } from '@zitadel/proto/zitadel/org/v2/org_service_pb';
import { PageEvent } from '@angular/material/paginator';
import { OrganizationFieldName } from '@zitadel/proto/zitadel/org/v2/query_pb';
import { Organization, OrganizationState } from '@zitadel/proto/zitadel/org/v2/org_pb';
import { PaginatorComponent } from '../paginator/paginator.component';

type ListQuery = NonNullable<MessageInitShape<typeof ListOrganizationsRequestSchema>['query']>;
type SearchQuery = NonNullable<MessageInitShape<typeof ListOrganizationsRequestSchema>['queries']>[number];

@Component({
  selector: 'cnsl-org-table',
  templateUrl: './org-table.component.html',
  styleUrls: ['./org-table.component.scss'],
})
export class OrgTableComponent {
  public displayedColumns: string[] = ['name', 'state', 'primaryDomain', 'creationDate', 'changeDate', 'actions'];
  public copied: string = '';

  public defaultOrgId: string = '';

  protected readonly listQuery = signal<ListQuery & { limit: number }>({ limit: 20, offset: BigInt(0) });
  private readonly searchQueries = signal<SearchQuery[]>([]);
  private readonly sortingColumn = signal<OrganizationFieldName | undefined>(undefined);

  private readonly req = computed<MessageInitShape<typeof ListOrganizationsRequestSchema>>(() => ({
    query: this.listQuery(),
    queries: this.searchQueries().length ? this.searchQueries() : undefined,
    sortingColumn: this.sortingColumn(),
  }));

  protected listOrganizationsQuery = injectQuery(() => ({
    ...this.newOrganizationService.listOrganizationsQueryOptions(this.req()),
    placeholderData: keepPreviousData,
  }));

  protected readonly dataSource = this.getDataSource();

  constructor(
    private readonly authService: GrpcAuthService,
    private readonly mgmtService: ManagementService,
    private readonly adminService: AdminService,
    protected readonly router: Router,
    private readonly toast: ToastService,
    private readonly liveAnnouncer: LiveAnnouncer,
    private readonly translate: TranslateService,
    private readonly newOrganizationService: NewOrganizationService,
  ) {
    this.mgmtService.getIAM().then((iam) => {
      this.defaultOrgId = iam.defaultOrgId;
    });

    effect(() => {
      if (this.listOrganizationsQuery.isError()) {
        this.toast.showError(this.listOrganizationsQuery.error());
      }
    });
  }

  private getDataSource() {
    const dataSource = new MatTableDataSource<Organization>();
    effect(() => {
      const organizations = this.listOrganizationsQuery.data()?.result ?? [];
      if (dataSource.data != organizations) {
        dataSource.data = organizations;
      }
    });

    return dataSource;
  }

  public async sortChange(sortState: Sort) {
    this.sortingColumn.set(sortState.active === 'name' ? OrganizationFieldName.NAME : undefined);

    const listQuery = { ...this.listQuery() };
    if (sortState.direction === 'asc') {
      this.listQuery.set({ ...listQuery, asc: true });
    } else {
      delete listQuery.asc;
      this.listQuery.set(listQuery);
    }

    if (sortState.direction && sortState.active) {
      await this.liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
    } else {
      await this.liveAnnouncer.announce('Sorting cleared');
    }
  }

  public setDefaultOrg(org: Organization) {
    this.adminService
      .setDefaultOrg(org.id)
      .then(() => {
        this.toast.showInfo('ORG.PAGES.DEFAULTORGSET', true);
        this.defaultOrgId = org.id;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public applySearchQuery(searchQueries: OrgQuery[], paginator: PaginatorComponent): void {
    if (this.searchQueries().length === 0 && searchQueries.length === 0) {
      return;
    }
    paginator.pageIndex = 0;
    this.searchQueries.set(searchQueries.map((q) => ({ query: this.oldQueryToNewQuery(q.toObject()) })));
  }

  private oldQueryToNewQuery(query: OrgQuery.AsObject): SearchQuery['query'] {
    if (query.idQuery) {
      return {
        case: 'idQuery' as const,
        value: {
          id: query.idQuery.id,
        },
      };
    }
    if (query.stateQuery) {
      return {
        case: 'stateQuery' as const,
        value: {
          state: query.stateQuery.state as unknown as any,
        },
      };
    }
    if (query.domainQuery) {
      return {
        case: 'domainQuery' as const,
        value: {
          domain: query.domainQuery.domain,
          method: query.domainQuery.method as unknown as any,
        },
      };
    }
    if (query.nameQuery) {
      return {
        case: 'nameQuery' as const,
        value: {
          name: query.nameQuery.name,
          method: query.nameQuery.method as unknown as any,
        },
      };
    }
    throw new Error('Invalid query');
  }

  public async setAndNavigateToOrg(org: Organization): Promise<void> {
    if (org.state !== OrganizationState.REMOVED) {
      await this.newOrganizationService.setOrgId(org.id);
      await this.router.navigate(['/org']);
    } else {
      this.translate.get('ORG.TOAST.ORG_WAS_DELETED').subscribe((data) => {
        this.toast.showInfo(data);
      });
    }
  }

  protected pageChanged(event: PageEvent) {
    this.listQuery.set({
      limit: event.pageSize,
      offset: BigInt(event.pageSize) * BigInt(event.pageIndex),
    });
  }

  protected readonly Number = Number;
  protected readonly OrganizationState = OrganizationState;
}
