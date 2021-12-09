import { Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { BehaviorSubject, catchError, debounceTime, finalize, from, map, Observable, of } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Org, OrgNameQuery, OrgQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-org-context',
  templateUrl: './org-context.component.html',
  styleUrls: ['./org-context.component.scss'],
})
export class OrgContextComponent implements OnInit {
  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);
  public orgs$: Observable<Org.AsObject[]> = of([]);
  public filterControl: FormControl = new FormControl('');
  @Input() public org!: Org.AsObject;
  @ViewChild('input', { static: false }) input!: ElementRef;
  @Output() public closedCard: EventEmitter<void> = new EventEmitter();
  @Output() public setOrg: EventEmitter<Org.AsObject> = new EventEmitter();

  constructor(public authService: AuthenticationService, private auth: GrpcAuthService) {
    this.filterControl.valueChanges.pipe(debounceTime(500)).subscribe((value) => {
      this.loadOrgs(value.trim().toLowerCase());
    });
  }

  public ngOnInit(): void {
    this.focusFilter();
    this.loadOrgs();
  }

  public setActiveOrg(org: Org.AsObject) {
    this.setOrg.emit(org);
    this.closedCard.emit();
  }

  public loadOrgs(filter?: string): void {
    if (!filter) {
      const value = this.input?.nativeElement?.value;
      if (value) {
        filter = value;
      }
    }

    let query;
    if (filter) {
      query = new OrgQuery();
      const orgNameQuery = new OrgNameQuery();
      orgNameQuery.setName(filter);
      orgNameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
      query.setNameQuery(orgNameQuery);
    }

    this.orgLoading$.next(true);
    this.orgs$ = from(this.auth.listMyProjectOrgs(10, 0, query ? [query] : undefined)).pipe(
      map((resp) => {
        return resp.resultList.sort((left, right) => left.name.localeCompare(right.name));
      }),
      catchError(() => of([])),
      finalize(() => {
        this.orgLoading$.next(false);
      }),
    );
  }

  public closeCard(element: HTMLElement): void {
    if (!element.classList.contains('dontcloseonclick') && !element.classList.contains('mat-button-wrapper')) {
      this.closedCard.emit();
    }
  }

  private focusFilter(): void {
    setTimeout(() => {
      this.input.nativeElement.focus();
    }, 0);
  }
}
