import { TestBed } from '@angular/core/testing';

import { ApiClientService } from './apiclient.service';

describe('ApiclientService', () => {
  let service: ApiClientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ApiClientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
