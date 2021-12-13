import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OrgCreateComponent } from './org-create.component';

describe('OrgCreateComponent', () => {
    let component: OrgCreateComponent;
    let fixture: ComponentFixture<OrgCreateComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [OrgCreateComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OrgCreateComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
