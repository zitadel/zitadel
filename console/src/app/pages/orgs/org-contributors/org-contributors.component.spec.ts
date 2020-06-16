import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { OrgContributorsComponent } from './org-contributors.component';

describe('OrgContributorsComponent', () => {
    let component: OrgContributorsComponent;
    let fixture: ComponentFixture<OrgContributorsComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [OrgContributorsComponent],
            imports: [
                NoopAnimationsModule,
                MatPaginatorModule,
                MatSortModule,
                MatTableModule,
            ],
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OrgContributorsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should compile', () => {
        expect(component).toBeTruthy();
    });
});
