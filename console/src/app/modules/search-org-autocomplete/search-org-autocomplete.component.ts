import { Component, EventEmitter, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatLegacyAutocomplete as MatAutocomplete } from '@angular/material/legacy-autocomplete';
import { debounceTime, from, map, Subject, switchMap, takeUntil, tap } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Org, OrgNameQuery, OrgQuery, OrgState, OrgStateQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-search-org-autocomplete',
  templateUrl: './search-org-autocomplete.component.html',
  styleUrls: ['./search-org-autocomplete.component.scss'],
})
export class SearchOrgAutocompleteComponent implements OnInit, OnDestroy {
  public selectable: boolean = true;
  public myControl: UntypedFormControl = new UntypedFormControl();
  public filteredOrgs: Array<Org.AsObject> = [];
  public isLoading: boolean = false;
  @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
  @Output() public selectionChanged: EventEmitter<Org.AsObject> = new EventEmitter();

  private unsubscribed$: Subject<void> = new Subject();
  constructor(public authService: AuthenticationService, private auth: GrpcAuthService) {
    this.myControl.valueChanges
      .pipe(
        takeUntil(this.unsubscribed$),
        debounceTime(200),
        tap(() => (this.isLoading = true)),
        switchMap((value) => {
          const stateQuery = new OrgQuery();
          const orgStateQuery = new OrgStateQuery();
          orgStateQuery.setState(OrgState.ORG_STATE_ACTIVE);
          stateQuery.setStateQuery(orgStateQuery);

          let queries: OrgQuery[] = [stateQuery];

          if (value) {
            const nameQuery = new OrgQuery();
            const orgNameQuery = new OrgNameQuery();
            orgNameQuery.setName(value);
            orgNameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
            nameQuery.setNameQuery(orgNameQuery);
            queries = [stateQuery, nameQuery];
          }

          return from(this.auth.listMyProjectOrgs(undefined, 0, queries)).pipe(
            map((resp) => {
              return resp.resultList.sort((left, right) => left.name.localeCompare(right.name));
            }),
          );
        }),
      )
      .subscribe((returnValue) => {
        this.isLoading = false;
        this.filteredOrgs = returnValue;
      });
  }

  public ngOnInit(): void {
    const query = new OrgQuery();
    const orgStateQuery = new OrgStateQuery();
    orgStateQuery.setState(OrgState.ORG_STATE_ACTIVE);
    query.setStateQuery(orgStateQuery);

    this.auth.listMyProjectOrgs(undefined, 0, [query]).then((orgs) => {
      this.filteredOrgs = orgs.resultList;
    });
  }

  public ngOnDestroy(): void {
    this.unsubscribed$.next();
  }

  public displayFn(org?: Org.AsObject): string {
    return org ? `${org.name}` : '';
  }

  public selected(event: MatAutocompleteSelectedEvent): void {
    this.selectionChanged.emit(event.option.value);
  }
}
