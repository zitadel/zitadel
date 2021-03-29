import { Pipe, PipeTransform } from '@angular/core';
import { Observable } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Pipe({
    name: 'hasFeature',
})
export class HasFeaturePipe implements PipeTransform {
    constructor(private authService: GrpcAuthService) { }

    public transform(values: string[]): Observable<boolean> {
        return this.authService.canUseFeature(values);
    }
}
