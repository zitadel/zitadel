import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class EnvService {
  env: any;
  constructor(private http: HttpClient) { }

  public loadEnvironment(): Observable<any> {
    if (this.env) {
      console.log('loaded env from cache');
      return of(this.env);
    } else {
      return this.http.get('./assets/environment.json');
    }
  }
}
