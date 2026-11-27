import { AutoLinkingOption, Options } from 'src/app/proto/generated/zitadel/idp_pb';

/**
 * Builds the `Options` object a provider save request carries: the options the
 * loaded provider stored, so that an edit saved without ever expanding the
 * optional section of the form leaves them unchanged, instead of silently
 * overwriting them with the defaults below. A provider without stored options
 * keeps those defaults.
 */
export function providerOptionsFromStored(stored: Options.AsObject | undefined): Options {
  if (!stored) {
    return new Options()
      .setIsCreationAllowed(true)
      .setIsLinkingAllowed(true)
      .setAutoLinking(AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED);
  }
  return new Options()
    .setIsAutoCreation(stored.isAutoCreation)
    .setIsAutoUpdate(stored.isAutoUpdate)
    .setIsCreationAllowed(stored.isCreationAllowed)
    .setIsLinkingAllowed(stored.isLinkingAllowed)
    .setAutoLinking(stored.autoLinking);
}
