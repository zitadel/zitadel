import { ChangeDetectionStrategy, Component, computed, effect, EventEmitter, Input, Output, Signal } from '@angular/core';
import { injectInfiniteQuery, injectMutation, keepPreviousData } from '@tanstack/angular-query-experimental';
import { NewOrganizationService } from 'src/app/services/new-organization.service';
import { NgForOf, NgIf } from '@angular/common';
import { ToastService } from 'src/app/services/toast.service';
import { FormBuilder, FormControl, ReactiveFormsModule } from '@angular/forms';
import { ListOrganizationsRequestSchema } from '@zitadel/proto/zitadel/org/v2/org_service_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { debounceTime } from 'rxjs/operators';
import { toSignal } from '@angular/core/rxjs-interop';
import { TextQueryMethod } from '@zitadel/proto/zitadel/object/v2/object_pb';
import { A11yModule } from '@angular/cdk/a11y';
import { MatButtonModule } from '@angular/material/button';
import { Organization } from '@zitadel/proto/zitadel/org/v2/org_pb';
import { MatMenuModule } from '@angular/material/menu';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from '../../input/input.module';
import { MatOptionModule } from '@angular/material/core';
import { Router } from '@angular/router';

type NameQuery = Extract<
  NonNullable<MessageInitShape<typeof ListOrganizationsRequestSchema>['queries']>[number]['query'],
  { case: 'nameQuery' }
>;

const QUERY_LIMIT = 5;

@Component({
  selector: 'cnsl-organization-selector',
  templateUrl: './organization-selector.component.html',
  styleUrls: ['./organization-selector.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    NgForOf,
    NgIf,
    ReactiveFormsModule,
    A11yModule,
    MatButtonModule,
    TranslateModule,
    MatMenuModule,
    InputModule,
    MatOptionModule,
  ],
})
export class OrganizationSelectorComponent {
  @Input()
  public backButton = '';

  @Output()
  public backButtonPressed = new EventEmitter<void>();

  @Output()
  public orgChanged = new EventEmitter<Organization>();

  protected setOrgId = injectMutation(() => ({
    mutationFn: (orgId: string) => this.newOrganizationService.setOrgId(orgId),
  }));

  protected readonly form: ReturnType<typeof this.buildForm>;
  private readonly nameQuery: Signal<NameQuery | undefined>;
  protected readonly organizationsQuery: ReturnType<typeof this.getOrganizationsQuery>;
  protected loadedOrgsCount: Signal<bigint>;
  protected activeOrg = this.newOrganizationService.activeOrganizationQuery();
  protected activeOrgIfSearchMatches: Signal<Organization | undefined>;

  constructor(
    private readonly newOrganizationService: NewOrganizationService,
    private readonly formBuilder: FormBuilder,
    private readonly router: Router,
    toast: ToastService,
  ) {
    this.form = this.buildForm();
    this.nameQuery = this.getNameQuery(this.form);
    this.organizationsQuery = this.getOrganizationsQuery(this.nameQuery);
    this.loadedOrgsCount = this.getLoadedOrgsCount(this.organizationsQuery);
    this.activeOrgIfSearchMatches = this.getActiveOrgIfSearchMatches(this.nameQuery);

    effect(() => {
      if (this.organizationsQuery.isError()) {
        toast.showError(this.organizationsQuery.error());
      }
    });
    effect(() => {
      if (this.setOrgId.isError()) {
        toast.showError(this.setOrgId.error());
      }
    });
    effect(() => {
      if (this.activeOrg.isError()) {
        toast.showError(this.activeOrg.error());
      }
    });

    effect(() => {
      const orgId = newOrganizationService.orgId();
      const orgs = this.organizationsQuery.data()?.pages[0]?.result;
      if (orgId || !orgs || orgs.length === 0) {
        return;
      }
      const _ = newOrganizationService.setOrgId(orgs[0].id);
    });
  }

  private buildForm() {
    return this.formBuilder.group({
      name: new FormControl('', { nonNullable: true }),
    });
  }

  private getNameQuery(form: ReturnType<typeof this.buildForm>): Signal<NameQuery | undefined> {
    const name$ = form.controls.name.valueChanges.pipe(debounceTime(125));
    const nameSignal = toSignal(name$, { initialValue: form.controls.name.value });

    return computed(() => {
      const name = nameSignal();
      if (!name) {
        return undefined;
      }
      const nameQuery: NameQuery = {
        case: 'nameQuery' as const,
        value: {
          name,
          method: TextQueryMethod.CONTAINS_IGNORE_CASE,
        },
      };
      return nameQuery;
    });
  }

  private getOrganizationsQuery(nameQuery: Signal<NameQuery | undefined>) {
    return injectInfiniteQuery(() => {
      const query = nameQuery();
      return {
        queryKey: ['organization', 'listOrganizationsInfinite', query],
        queryFn: ({ pageParam, signal }) => this.newOrganizationService.listOrganizations(pageParam, signal),
        initialPageParam: {
          query: {
            limit: QUERY_LIMIT,
            offset: BigInt(0),
          },
          queries: query ? [{ query }] : undefined,
        },
        placeholderData: keepPreviousData,
        getNextPageParam: (lastPage, _, pageParam) =>
          // if we received less than the limit last time we are at the end
          lastPage.result.length < pageParam.query.limit
            ? undefined
            : {
                ...pageParam,
                query: {
                  ...pageParam.query,
                  offset: pageParam.query.offset + BigInt(lastPage.result.length),
                },
              },
      };
    });
  }

  private getLoadedOrgsCount(organizationsQuery: ReturnType<typeof this.getOrganizationsQuery>) {
    return computed(() => {
      const pages = organizationsQuery.data()?.pages;
      if (!pages) {
        return BigInt(0);
      }
      return pages.reduce((acc, page) => acc + BigInt(page.result.length), BigInt(0));
    });
  }

  private getActiveOrgIfSearchMatches(nameQuery: Signal<NameQuery | undefined>) {
    return computed(() => {
      const activeOrg = this.activeOrg.data() ?? undefined;
      const query = nameQuery();
      if (!activeOrg || !query?.value?.name) {
        return activeOrg;
      }
      return activeOrg.name.toLowerCase().includes(query.value.name.toLowerCase()) ? activeOrg : undefined;
    });
  }

  protected async changeOrg(orgId: string) {
    const org = await this.setOrgId.mutateAsync(orgId);
    this.orgChanged.emit(org);
    await this.router.navigate(['/org']);
  }

  protected trackOrg(_: number, { id }: Organization): string {
    return id;
  }
}
