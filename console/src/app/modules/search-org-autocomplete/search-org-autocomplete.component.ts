import { ChangeDetectionStrategy, Component, effect, inject, output, signal } from '@angular/core';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { injectInfiniteQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { UserService } from 'src/app/services/user.service';
import { AuthService, ListMyProjectOrgsRequestSchema } from '@zitadel/proto/zitadel/auth_pb';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { InputModule } from '../input/input.module';
import { TranslatePipe } from '@ngx-translate/core';
import { TextQueryMethod } from '@zitadel/proto/zitadel/object_pb';
import { ScrollingModule } from '@angular/cdk/scrolling';
import { Org, OrgFieldName, OrgState } from '@zitadel/proto/zitadel/org_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { MatTooltip } from '@angular/material/tooltip';
import { requiredValidator } from '../form-field/validators/validators';
import { ScrollableDirective } from 'src/app/directives/scrollable/scrollable.directive';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { debounceTime } from 'rxjs';
import { MatProgressSpinner } from '@angular/material/progress-spinner';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-search-org-autocomplete',
  templateUrl: './search-org-autocomplete.component.html',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    ScrollingModule,
    MatInputModule,
    MatAutocompleteModule,
    ReactiveFormsModule,
    InputModule,
    TranslatePipe,
    FormsModule,
    MatTooltip,
    ScrollableDirective,
    MatProgressSpinner,
  ],
})
export class SearchOrgAutocompleteComponent {
  public readonly selectionChanged = output<Org>();

  protected readonly searchControl = new FormControl<Org | null>(null, { validators: [requiredValidator] });
  protected readonly searchSignal = signal('');
  protected readonly debouncedSearchSignal = toSignal(toObservable(this.searchSignal).pipe(debounceTime(250)), {
    initialValue: this.searchSignal(),
  });

  private readonly authService = inject(NewAuthService);
  private readonly userService = inject(UserService);

  protected readonly query = injectInfiniteQuery(() => {
    const search = this.debouncedSearchSignal();

    const stateQuery = [
      {
        query: {
          case: 'stateQuery',
          value: {
            state: OrgState.ACTIVE,
          },
        },
      },
    ] satisfies MessageInitShape<typeof ListMyProjectOrgsRequestSchema>['queries'];

    return {
      queryKey: [this.userService.userId(), AuthService.name, AuthService.method.listMyProjectOrgs.name, 'infinite', search],
      placeholderData: keepPreviousData,
      queryFn: async ({ pageParam, signal }) =>
        this.authService.listMyProjectOrgs(
          {
            query: {
              limit: 20,
              offset: BigInt(pageParam.offset),
            },
            sortingColumn: OrgFieldName.NAME,
            queries: search
              ? [
                  ...stateQuery,
                  {
                    query: {
                      case: 'nameQuery',
                      value: {
                        name: search,
                        method: TextQueryMethod.CONTAINS_IGNORE_CASE,
                      },
                    },
                  },
                ]
              : stateQuery,
          },
          signal,
        ),
      initialPageParam: {
        offset: 0,
      },
      getNextPageParam: (lastPage, allPages, { offset }) => {
        const nextPageParam = { offset: offset + lastPage.result.length };
        const loadedCount = allPages.reduce((count, page) => count + page.result.length, 0);

        if (loadedCount < lastPage.details!.totalResult) {
          return nextPageParam;
        }

        return undefined;
      },
    };
  });

  // used to only open matTooltip if the text gets truncated
  protected readonly showTooltip = signal(false);

  constructor() {
    const toastService = inject(ToastService);
    effect(() => {
      const error = this.query.error();
      if (error) {
        toastService.showError(error);
      }
    });
  }

  protected async onScroll(position: 'bottom' | 'top') {
    if (this.query.isFetchingNextPage() || !this.query.hasNextPage() || position === 'top') {
      return;
    }

    await this.query.fetchNextPage();
  }

  protected displayOrg(org: Org | null) {
    return org?.name ?? '';
  }
}
