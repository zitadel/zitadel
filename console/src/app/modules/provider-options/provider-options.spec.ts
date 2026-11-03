import { AutoLinkingOption, Options } from 'src/app/proto/generated/zitadel/idp_pb';

import { providerOptionsFromStored } from './provider-options';

describe('providerOptionsFromStored', () => {
  it('copies every stored option so a save that never expanded the optional section keeps it', () => {
    const stored: Options.AsObject = {
      isAutoCreation: true,
      isAutoUpdate: true,
      isCreationAllowed: false,
      isLinkingAllowed: false,
      autoLinking: AutoLinkingOption.AUTO_LINKING_OPTION_USERNAME,
    };

    const options = providerOptionsFromStored(stored);

    expect(options.getIsAutoCreation()).toBe(true);
    expect(options.getIsAutoUpdate()).toBe(true);
    expect(options.getIsCreationAllowed()).toBe(false);
    expect(options.getIsLinkingAllowed()).toBe(false);
    expect(options.getAutoLinking()).toBe(AutoLinkingOption.AUTO_LINKING_OPTION_USERNAME);
  });

  it('falls back to the defaults a new provider starts with when none are stored', () => {
    const options = providerOptionsFromStored(undefined);

    expect(options.getIsCreationAllowed()).toBe(true);
    expect(options.getIsLinkingAllowed()).toBe(true);
    expect(options.getAutoLinking()).toBe(AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED);
    expect(options.getIsAutoCreation()).toBe(false);
    expect(options.getIsAutoUpdate()).toBe(false);
  });
});
