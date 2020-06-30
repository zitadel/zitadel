import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { IamContributorsComponent } from './iam-contributors.component';

describe('OrgContributorsComponent', () => {
    let component: IamContributorsComponent;
    let fixture: ComponentFixture<IamContributorsComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [IamContributorsComponent],
            imports: [
                NoopAnimationsModule,
                MatPaginatorModule,
                MatSortModule,
                MatTableModule,
            ],
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(IamContributorsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should compile', () => {
        expect(component).toBeTruthy();
    });
});
