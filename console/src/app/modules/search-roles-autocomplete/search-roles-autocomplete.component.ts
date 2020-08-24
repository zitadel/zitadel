import { COMMA, ENTER } from '@angular/cdk/keycodes';
import { Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatAutocomplete, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatChipInputEvent } from '@angular/material/chips';
import { from } from 'rxjs';
import { debounceTime, switchMap, tap } from 'rxjs/operators';
import {
    ProjectRole,
    ProjectRoleSearchKey,
    ProjectRoleSearchQuery,
    SearchMethod,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';


@Component({
    selector: 'app-search-roles-autocomplete',
    templateUrl: './search-roles-autocomplete.component.html',
    styleUrls: ['./search-roles-autocomplete.component.scss'],
})
export class SearchRolesAutocompleteComponent {
    public selectable: boolean = true;
    public removable: boolean = true;
    public addOnBlur: boolean = true;
    public separatorKeysCodes: number[] = [ENTER, COMMA];
    public myControl: FormControl = new FormControl();
    public names: string[] = [];
    public roles: Array<ProjectRole.AsObject> = [];
    public filteredRoles: Array<ProjectRole.AsObject> = [];
    public isLoading: boolean = false;
    @ViewChild('nameInput') public nameInput!: ElementRef<HTMLInputElement>;
    @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
    @Input() public projectId: string = '';
    @Input() public singleOutput: boolean = false;
    @Output() public selectionChanged: EventEmitter<ProjectRole.AsObject[] | ProjectRole.AsObject> = new EventEmitter();
    constructor(private mgmtService: ManagementService, private toast: ToastService) {
        this.myControl.valueChanges
            .pipe(
                debounceTime(200),
                tap(() => this.isLoading = true),
                switchMap(value => {
                    const query = new ProjectRoleSearchQuery();
                    query.setKey(ProjectRoleSearchKey.PROJECTROLESEARCHKEY_DISPLAY_NAME);
                    query.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS_IGNORE_CASE);
                    query.setValue(value);
                    return from(this.mgmtService.SearchProjectRoles(this.projectId, 10, 0, [query]));
                }),
            ).subscribe((roles) => {
                this.isLoading = false;
                this.filteredRoles = roles.toObject().resultList;
            }, error => {
                this.isLoading = false;
            });
    }

    public displayFn(project?: ProjectRole.AsObject): string | undefined {
        return project ? `${project.displayName}` : undefined;
    }

    public add(event: MatChipInputEvent): void {
        if (!this.matAutocomplete.isOpen) {
            const input = event.input;
            const value = event.value;

            if ((value || '').trim()) {
                const index = this.filteredRoles.findIndex((role) => {
                    if (role.key) {
                        return role.key === value;
                    }
                });
                if (index > -1) {
                    if (this.roles && this.roles.length > 0) {
                        this.roles.push(this.filteredRoles[index]);
                    } else {
                        this.roles = [this.filteredRoles[index]];
                    }
                }
            }

            if (input) {
                input.value = '';
            }
        }
    }

    public remove(role: ProjectRole.AsObject): void {
        const index = this.roles.indexOf(role);

        if (index >= 0) {
            this.roles.splice(index, 1);
        }
    }

    public selected(event: MatAutocompleteSelectedEvent): void {
        const index = this.filteredRoles.findIndex((role) => role.key === event.option.value);
        if (index !== -1) {
            if (this.singleOutput) {
                this.selectionChanged.emit(this.filteredRoles[index]);
            } else {
                if (this.roles && this.roles.length > 0) {
                    this.roles.push(this.filteredRoles[index]);
                } else {
                    this.roles = [this.filteredRoles[index]];
                }
                this.selectionChanged.emit(this.roles);

                this.nameInput.nativeElement.value = '';
                this.myControl.setValue(null);
            }

        }
    }
}
