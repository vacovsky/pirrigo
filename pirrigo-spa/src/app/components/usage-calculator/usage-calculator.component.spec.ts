import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UsageCalculatorComponent } from './usage-calculator.component';

describe('UsageCalculatorComponent', () => {
  let component: UsageCalculatorComponent;
  let fixture: ComponentFixture<UsageCalculatorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ UsageCalculatorComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UsageCalculatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
