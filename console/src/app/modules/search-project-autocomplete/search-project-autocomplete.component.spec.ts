import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { SearchProjectAutocompleteComponent } from './search-project-autocomplete.component';


describe('SearchProjectComponent', () => {
    let component: SearchProjectAutocompleteComponent;
    let fixture: ComponentFixture<SearchProjectAutocompleteComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [SearchProjectAutocompleteComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(SearchProjectAutocompleteComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
