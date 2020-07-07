import { Component, ElementRef, EventEmitter, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatAutocomplete } from '@angular/material/autocomplete';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-org-member-roles-autocomplete',
    templateUrl: './org-member-roles-autocomplete.component.html',
    styleUrls: ['./org-member-roles-autocomplete.component.scss'],
})
export class OrgMemberRolesAutocompleteComponent {
    public myControl: FormControl = new FormControl();
    public names: string[] = [];
    public roles: string[] = [];
    public allRoles: string[] = [];
    public isLoading: boolean = false;
    @ViewChild('nameInput') public nameInput!: ElementRef<HTMLInputElement>;
    @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
    @Output() public selectionChanged: EventEmitter<string[]> = new EventEmitter();
    constructor(private orgService: OrgService, private toast: ToastService) {
        this.orgService.GetOrgMemberRoles().then(resp => {
            this.allRoles = resp.toObject().rolesList;
        }).catch(error => {
            this.toast.showError(error.message);
        });

        this.myControl.valueChanges.subscribe(change => {
            this.selectionChanged.emit(change);
        });
    }
}
