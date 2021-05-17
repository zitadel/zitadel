import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ZitadelTierComponent } from './zitadel-tier.component';

describe('ZitadelTierComponent', () => {
  let component: ZitadelTierComponent;
  let fixture: ComponentFixture<ZitadelTierComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ZitadelTierComponent],
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ZitadelTierComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
