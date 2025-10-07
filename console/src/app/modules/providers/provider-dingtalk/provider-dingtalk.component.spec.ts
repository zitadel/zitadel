import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProviderDingTalkComponent } from './provider-dingtalk.component';

describe('ProviderDingTalkComponent', () => {
  let component: ProviderDingTalkComponent;
  let fixture: ComponentFixture<ProviderDingTalkComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ProviderDingTalkComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ProviderDingTalkComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});