import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SearchRolesAutocompleteComponent } from './search-roles-autocomplete.component';



describe('SearchProjectComponent', () => {
    let component: SearchRolesAutocompleteComponent;
    let fixture: ComponentFixture<SearchRolesAutocompleteComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [SearchRolesAutocompleteComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(SearchRolesAutocompleteComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
