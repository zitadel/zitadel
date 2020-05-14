import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchUserAutocompleteComponent } from './search-user-autocomplete.component';

describe('SearchUserAutocompleteComponent', () => {
    let component: SearchUserAutocompleteComponent;
    let fixture: ComponentFixture<SearchUserAutocompleteComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [SearchUserAutocompleteComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(SearchUserAutocompleteComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
