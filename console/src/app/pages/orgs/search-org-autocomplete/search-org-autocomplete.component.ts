import { COMMA, ENTER } from '@angular/cdk/keycodes';
import { Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatAutocomplete, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatChipInputEvent } from '@angular/material/chips';
import { debounceTime, tap } from 'rxjs/operators';
import { Org } from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-search-org-autocomplete',
    templateUrl: './search-org-autocomplete.component.html',
    styleUrls: ['./search-org-autocomplete.component.scss'],
})
export class SearchOrgAutocompleteComponent {
    public selectable: boolean = true;
    public removable: boolean = true;
    public addOnBlur: boolean = true;
    public separatorKeysCodes: number[] = [ENTER, COMMA];
    public myControl: FormControl = new FormControl();
    public names: string[] = [];
    public orgs: Array<Org.AsObject> = [];
    public filteredOrgs: Array<Org.AsObject> = [];
    public isLoading: boolean = false;
    @ViewChild('domainInput') public domainInput!: ElementRef<HTMLInputElement>;
    @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
    @Input() public singleOutput: boolean = false;
    @Output() public selectionChanged: EventEmitter<Org.AsObject | Org.AsObject[]> = new EventEmitter();
    constructor(private orgService: OrgService, private toast: ToastService) {
        this.myControl.valueChanges.pipe(debounceTime(200), tap(() => this.isLoading = true)).subscribe(value => {
            return this.orgService.getOrgByDomainGlobal(value).then((org) => {
                this.isLoading = false;
                if (org) {
                    this.filteredOrgs = [org.toObject()];
                }
            }).catch(error => {
                this.isLoading = false;
                // this.toast.showInfo(error.message);
            });
        });
    }

    public displayFn(org?: Org.AsObject): string | undefined {
        return org ? `${org.name}` : undefined;
    }

    public add(event: MatChipInputEvent): void {
        if (!this.matAutocomplete.isOpen) {
            const input = event.input;
            const value = event.value;

            if ((value || '').trim()) {
                const index = this.filteredOrgs.findIndex((org) => {
                    if (org.name) {
                        return org.name === value;
                    }
                });
                if (index > -1) {
                    if (this.orgs && this.orgs.length > 0) {
                        this.orgs.push(this.filteredOrgs[index]);
                    } else {
                        this.orgs = [this.filteredOrgs[index]];
                    }
                }
            }

            if (input) {
                input.value = '';
            }
        }
    }

    public remove(org: Org.AsObject): void {
        const index = this.orgs.indexOf(org);

        if (index >= 0) {
            this.orgs.splice(index, 1);
        }
    }

    public selected(event: MatAutocompleteSelectedEvent): void {
        const index = this.filteredOrgs.findIndex((org) => org === event.option.value);
        if (index !== -1) {
            if (this.singleOutput) {
                this.selectionChanged.emit(this.filteredOrgs[index]);
            } else {
                if (this.orgs && this.orgs.length > 0) {
                    this.orgs.push(this.orgs[index]);
                    this.selectionChanged.emit(this.orgs);
                } else {
                    this.orgs = [this.filteredOrgs[index]];
                }

                this.domainInput.nativeElement.value = '';
                this.myControl.setValue(null);
            }
        }
    }
}
