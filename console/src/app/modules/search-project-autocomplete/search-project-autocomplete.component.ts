import { COMMA, ENTER } from '@angular/cdk/keycodes';
import { Component, ElementRef, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import {
  MatLegacyAutocomplete as MatAutocomplete,
  MatLegacyAutocompleteSelectedEvent as MatAutocompleteSelectedEvent,
} from '@angular/material/legacy-autocomplete';
import { MatLegacyChipInputEvent as MatChipInputEvent } from '@angular/material/legacy-chips';
import { forkJoin, from, Subject } from 'rxjs';
import { debounceTime, switchMap, takeUntil, tap } from 'rxjs/operators';
import { ListProjectGrantsResponse, ListProjectsResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { GrantedProject, Project, ProjectNameQuery, ProjectQuery } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

import { ProjectType } from '../project-members/project-members-datasource';

export enum ProjectAutocompleteType {
  PROJECT_OWNED = 0,
  PROJECT_GRANTED = 1,
}

@Component({
  selector: 'cnsl-search-project-autocomplete',
  templateUrl: './search-project-autocomplete.component.html',
  styleUrls: ['./search-project-autocomplete.component.scss'],
})
export class SearchProjectAutocompleteComponent implements OnInit, OnDestroy {
  public selectable: boolean = true;
  public removable: boolean = true;
  public addOnBlur: boolean = true;
  public separatorKeysCodes: number[] = [ENTER, COMMA];
  public myControl: UntypedFormControl = new UntypedFormControl();
  public names: string[] = [];
  public projects: Array<GrantedProject.AsObject | Project.AsObject | any> = [];
  public filteredProjects: Array<GrantedProject.AsObject | Project.AsObject | any> = [];
  public isLoading: boolean = false;
  @ViewChild('nameInput') public nameInput!: ElementRef<HTMLInputElement>;
  @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
  @Input() public autocompleteType!: ProjectAutocompleteType;
  @Output() public selectionChanged: EventEmitter<{
    project: Project.AsObject | GrantedProject.AsObject;
    type: ProjectType;
  }> = new EventEmitter();

  private unsubscribed$: Subject<void> = new Subject();
  constructor(private mgmtService: ManagementService) {
    this.myControl.valueChanges
      .pipe(
        takeUntil(this.unsubscribed$),
        debounceTime(200),
        tap(() => (this.isLoading = true)),
        switchMap((value) => {
          const query = new ProjectQuery();
          const nameQuery = new ProjectNameQuery();
          nameQuery.setName(value);
          nameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          query.setNameQuery(nameQuery);

          switch (this.autocompleteType) {
            case ProjectAutocompleteType.PROJECT_GRANTED:
              return from(this.mgmtService.listGrantedProjects(10, 0, [query]));
            case ProjectAutocompleteType.PROJECT_OWNED:
              return from(this.mgmtService.listProjects(10, 0, [query]));
            default:
              return forkJoin([
                from(this.mgmtService.listGrantedProjects(10, 0, [query])),
                from(this.mgmtService.listProjects(10, 0, [query])),
              ]);
          }
        }),
      )
      .subscribe((returnValue) => {
        switch (this.autocompleteType) {
          case ProjectAutocompleteType.PROJECT_GRANTED:
            this.isLoading = false;
            this.filteredProjects = [...(returnValue as ListProjectGrantsResponse.AsObject).resultList];
            break;
          case ProjectAutocompleteType.PROJECT_OWNED:
            this.isLoading = false;
            this.filteredProjects = [...(returnValue as ListProjectsResponse.AsObject).resultList];
            break;
          default:
            this.isLoading = false;
            this.filteredProjects = [
              ...(returnValue as (ListProjectsResponse.AsObject | ListProjectGrantsResponse.AsObject)[])[0].resultList,
              ...(returnValue as (ListProjectsResponse.AsObject | ListProjectGrantsResponse.AsObject)[])[1].resultList,
            ];
            break;
        }
      });
  }

  public ngOnInit(): void {
    // feat-3916 show projects as soon as I am in the input field of the project
    const query = new ProjectQuery();
    const nameQuery = new ProjectNameQuery();
    nameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
    query.setNameQuery(nameQuery);

    switch (this.autocompleteType) {
      case ProjectAutocompleteType.PROJECT_GRANTED:
        this.mgmtService.listGrantedProjects(10, 0, [query]).then((projects) => {
          this.filteredProjects = projects.resultList;
        });
        break;
      case ProjectAutocompleteType.PROJECT_OWNED:
        this.mgmtService.listProjects(10, 0, [query]).then((projects) => {
          this.filteredProjects = projects.resultList;
        });
        break;
      default:
        Promise.all([
          this.mgmtService.listGrantedProjects(10, 0, [query]),
          this.mgmtService.listProjects(10, 0, [query]),
        ]).then((values) => {
          this.filteredProjects = values[0].resultList;
          this.filteredProjects = this.filteredProjects.concat(values[1].resultList);
        });
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribed$.next();
  }

  public displayFn(project?: any): string {
    return project && project.projectName ? `${project.projectName}` : project && project.name ? `${project.name}` : '';
  }

  public add(event: MatChipInputEvent): void {
    if (!this.matAutocomplete.isOpen) {
      const input = event.chipInput?.inputElement;
      const value = event.value;

      if ((value || '').trim()) {
        const index = this.filteredProjects.findIndex((project) => {
          if (project?.projectName) {
            return project.projectName === value;
          } else if (project?.name) {
            return project.name === value;
          } else {
            return false;
          }
        });
        if (index > -1) {
          if (this.projects && this.projects.length > 0) {
            this.projects.push(this.filteredProjects[index]);
          } else {
            this.projects = [this.filteredProjects[index]];
          }
        }
      }

      if (input) {
        input.value = '';
      }
    }
  }

  public remove(project: GrantedProject.AsObject): void {
    const index = this.projects.indexOf(project);

    if (index >= 0) {
      this.projects.splice(index, 1);
    }
  }

  public selected(event: MatAutocompleteSelectedEvent): void {
    const p: Project.AsObject | GrantedProject.AsObject = event.option.value;
    const type = (p as Project.AsObject).id
      ? ProjectType.PROJECTTYPE_OWNED
      : (p as GrantedProject.AsObject).projectId
      ? ProjectType.PROJECTTYPE_GRANTED
      : ProjectType.PROJECTTYPE_OWNED;

    this.selectionChanged.emit({
      project: p,
      type: type,
    });
  }
}
