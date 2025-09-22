import {
  ChangeDetectionStrategy,
  Component,
  computed,
  DestroyRef,
  effect,
  ElementRef,
  EventEmitter,
  Input,
  Output,
  signal,
  Signal,
  ViewChild,
} from '@angular/core';
import { injectInfiniteQuery, injectMutation, keepPreviousData, QueryClient } from '@tanstack/angular-query-experimental';
import { NewOrganizationService } from 'src/app/services/new-organization.service';
import { AsyncPipe, NgForOf, NgIf } from '@angular/common';
import { ToastService } from 'src/app/services/toast.service';
import { FormBuilder, FormControl, ReactiveFormsModule } from '@angular/forms';
import { ListOrganizationsRequestSchema, ListOrganizationsResponse } from '@zitadel/proto/zitadel/org/v2/org_service_pb';
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
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import { heroCheck, heroMagnifyingGlass } from '@ng-icons/heroicons/outline';
import { heroArrowLeftCircleSolid } from '@ng-icons/heroicons/solid';
import { UserService } from 'src/app/services/user.service';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { NewAuthService } from 'src/app/services/new-auth.service';

type NameQuery = Extract<
  NonNullable<MessageInitShape<typeof ListOrganizationsRequestSchema>['queries']>[number]['query'],
  { case: 'nameQuery' }
>;

const QUERY_LIMIT = 20;

@Component({
  selector: 'cnsl-organization-selector',
  templateUrl: './organization-selector.component.html',
  styleUrls: ['./organization-selector.component.scss'],
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
    NgIconComponent,
    HasRolePipeModule,
    AsyncPipe,
  ],
  providers: [provideIcons({ heroCheck, heroMagnifyingGlass, heroArrowLeftCircleSolid })],
})
export class OrganizationSelectorComponent {
  @Input()
  public backButton = '';

  @Output()
  public backButtonPressed = new EventEmitter<void>();

  @Output()
  public orgChanged = new EventEmitter<Organization>();

  @ViewChild('moreButton', { static: false, read: ElementRef })
  public set moreButton(button: ElementRef<HTMLButtonElement>) {
    this.moreButtonSignal.set(button);
  }

  private moreButtonSignal = signal<ElementRef<HTMLButtonElement> | undefined>(undefined);

  protected setOrgId = injectMutation(() => ({
    mutationFn: (orgId: string) => this.newOrganizationService.setOrgId(orgId),
  }));

  protected readonly form: ReturnType<typeof this.buildForm>;
  private readonly nameQuery: Signal<NameQuery | undefined>;
  protected readonly organizationsQuery: ReturnType<typeof this.getOrganizationsQuery>;
  protected readonly activeOrg = this.newOrganizationService.activeOrganizationQuery();
  protected readonly activeOrgIfSearchMatches: Signal<Organization | undefined>;
  private readonly listMyZitadelPermissionsQuery = this.newAuthService.listMyZitadelPermissionsQuery();

  constructor(
    private readonly newOrganizationService: NewOrganizationService,
    private readonly formBuilder: FormBuilder,
    private readonly router: Router,
    private readonly destroyRef: DestroyRef,
    private readonly userService: UserService,
    private readonly newAuthService: NewAuthService,
    private readonly queryClient: QueryClient,
    toast: ToastService,
  ) {
    this.form = this.buildForm();
    this.nameQuery = this.getNameQuery(this.form);
    this.organizationsQuery = this.getOrganizationsQuery(this.nameQuery);
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
      const orgs = this.organizationsQuery.data()?.orgs;

      // orgs not yet loaded or user has no orgs
      if (!orgs || orgs.length === 0) {
        return;
      }

      // no orgId set so we set it to the first org
      if (!orgId) {
        newOrganizationService.setOrgId(orgs[0].id).then();
        return;
      }

      // user has a selected org and it was found
      if (orgs.some((org) => org.id === orgId)) {
        return;
      }

      // maybe the org is not yet loaded in the org selector so we try to fetch it
      // if the user has permission to the org this will succeed and we do nothing
      this.queryClient
        .fetchQuery(this.newOrganizationService.organizationByIdQueryOptions(orgId))
        .then((org) => {
          if (org) {
            return;
          }
          throw new Error('org not found');
        })
        .catch((_) => {
          // user has no org selected or no permission for said org so we default to first org
          return newOrganizationService.setOrgId(orgs[0].id);
        });
    });

    this.infiniteScrollLoading();
  }

  private infiniteScrollLoading() {
    const intersection = new IntersectionObserver(async (entries) => {
      if (!entries[0]?.isIntersecting) {
        return;
      }
      await this.organizationsQuery.fetchNextPage();
    });
    this.destroyRef.onDestroy(() => {
      intersection.disconnect();
    });

    effect((onCleanup) => {
      const moreButton = this.moreButtonSignal()?.nativeElement;
      const permissions = this.listMyZitadelPermissionsQuery.data();

      if (!moreButton || !permissions) {
        return;
      }

      // only do infinite scrolling when user has access to all orgs
      if (!permissions.includes('iam.read')) {
        return;
      }

      intersection.observe(moreButton);
      onCleanup(() => {
        intersection.unobserve(moreButton);
      });
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
      const isExpired = this.userService.isExpired();
      return {
        queryKey: [this.userService.userId(), 'organization', 'listOrganizationsInfinite', query],
        queryFn: ({ pageParam, signal }) => this.newOrganizationService.listOrganizations(pageParam, signal),
        enabled: !isExpired,
        initialPageParam: {
          query: {
            limit: QUERY_LIMIT,
            offset: BigInt(0),
          },
          queries: query ? [{ query }] : undefined,
        },
        placeholderData: keepPreviousData,
        getNextPageParam: (lastPage, pages, pageParam) =>
          this.countLoadedOrgs(pages) < (lastPage.details?.totalResult ?? BigInt(Number.MAX_SAFE_INTEGER))
            ? {
                ...pageParam,
                query: {
                  ...pageParam.query,
                  offset: pageParam.query.offset + BigInt(lastPage.result.length),
                },
              }
            : undefined,
        select: (data) => ({
          orgs: data.pages.flatMap((page) => page.result),
          totalResult: Number(data.pages[data.pages.length - 1]?.details?.totalResult ?? 0),
        }),
      };
    });
  }

  private countLoadedOrgs(pages?: ListOrganizationsResponse[]) {
    if (!pages) {
      return BigInt(0);
    }
    return pages.reduce((acc, page) => acc + BigInt(page.result.length), BigInt(0));
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

  protected trackOrgResponse(_: number, { id }: Organization): string {
    return id;
  }

  protected readonly QUERY_LIMIT = QUERY_LIMIT;
}
