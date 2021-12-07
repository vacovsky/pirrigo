import { TestBed } from '@angular/core/testing';

import { ChartTransformService } from './chart-transform.service';

describe('ChartTransformService', () => {
  let service: ChartTransformService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ChartTransformService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
