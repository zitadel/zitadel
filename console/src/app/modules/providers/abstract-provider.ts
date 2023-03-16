import { FormGroup } from '@angular/forms';
import { Provider, ProviderType } from 'src/app/proto/generated/zitadel/idp_pb';

export abstract class AbstractProvider {
  abstract form: FormGroup;
  abstract getData(id: string, providerType: ProviderType): Promise<Provider.AsObject | undefined>;
  abstract addProvider(form: FormGroup): boolean;
  abstract updateProvder(id: string, form: FormGroup): boolean;
  abstract navigateBack(): void;
}
