import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { OrgMembersComponent } from './org-members.component';

describe('OrgMembersComponent', () => {
    let component: OrgMembersComponent;
    let fixture: ComponentFixture<OrgMembersComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [OrgMembersComponent],
            imports: [
                NoopAnimationsModule,
                MatPaginatorModule,
                MatSortModule,
                MatTableModule,
            ],
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OrgMembersComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should compile', () => {
        expect(component).toBeTruthy();
    });
});
