import { Location } from '@angular/common';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { v4 as uuidv4 } from 'uuid';

export abstract class StatehandlerProcessorService {
  public abstract createState(url: string): string;
  public abstract restoreState(state?: string): void;
}

@Injectable()
export class StatehandlerProcessorServiceImpl implements StatehandlerProcessorService {
  constructor(
    private location: Location,
    private router: Router,
  ) {}

  public createState(url: string): string {
    const externalUrl = url;
    const urlId = uuidv4();
    sessionStorage.setItem(urlId, externalUrl);
    return urlId;
  }

  public restoreState(state?: string): void {
    if (state === undefined) {
      return;
    } else {
      const url = sessionStorage.getItem(state);
      if (url === null) {
        return;
      } else {
        sessionStorage.removeItem(state);
        this.router.navigateByUrl(url);
      }
    }
  }
}
