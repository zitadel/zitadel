import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {
    SearchProjectAutocompleteComponent,
} from '../../../modules/search-project-autocomplete/search-project-autocomplete.component';

describe('SearchProjectComponent', () => {
    let component: SearchProjectAutocompleteComponent;
    let fixture: ComponentFixture<SearchProjectAutocompleteComponent>;

    beforeEach(async(() => {
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
