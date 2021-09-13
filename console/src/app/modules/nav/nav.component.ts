import { Component, EventEmitter, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { debounceTime } from 'rxjs/operators';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';

@Component({
  selector: 'cnsl-nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss']
})
export class NavComponent implements OnInit {
  public org!: Org.AsObject;
  public setActiveOrg: EventEmitter<Org.AsObject> = new EventEmitter();
  public filterControl: FormControl = new FormControl('');

  constructor(public authenticationService: AuthenticationService,
  ) {
    this.filterControl.valueChanges.pipe(debounceTime(300)).subscribe(value => {
      this.loadOrgs(
        value.trim().toLowerCase(),
      );
    });
  }

  ngOnInit(): void {
  }


  public loadOrgs(filter?: string): void {
    let query;
    if (filter) {
      query = new OrgQuery();
      const orgNameQuery = new OrgNameQuery();
      orgNameQuery.setName(filter);
      orgNameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
      query.setNameQuery(orgNameQuery);
    }

    this.orgLoading$.next(true);
    this.orgs$ = from(this.authService.listMyProjectOrgs(10, 0, query ? [query] : undefined)).pipe(
      map(resp => {
        return resp.resultList;
      }),
      catchError(() => of([])),
      finalize(() => {
        this.orgLoading$.next(false);
        this.focusFilter();
      }),
    );
  }

  public closeAccountCard(): void {
    if (this.showAccount) {
      this.showAccount = false;
    }
  }

}
