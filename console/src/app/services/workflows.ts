import { CnslOverlay } from './overlay-workflow.service';

export const IntroWorkflowOverlays: CnslOverlay[] = [
  {
    id: 'orgswitcher',
    origin: 'orgswitchbutton',
    toHighlight: ['orgswitchbutton'],
    content: {
      i18nText: 'OVERLAYS.ORGSWITCHER.TEXT',
    },
    requirements: {
      permission: ['org.read'],
    },
  },
  {
    id: 'systembutton',
    origin: 'systembutton',
    toHighlight: ['systembutton'],
    content: {
      i18nText: 'OVERLAYS.SYSTEM.TEXT',
    },
    requirements: {
      permission: ['iam.read'],
    },
  },
  {
    id: 'profilebutton',
    origin: 'avatartoggle',
    toHighlight: ['avatartoggle'],
    content: {
      i18nText: 'OVERLAYS.PROFILE.TEXT',
    },
  },
  // { id: 'mainnav', origin: 'orgswitchbutton' },
];
