import { CnslOverlay } from './overlay-workflow.service';

export const IntroWorkflowOverlays: CnslOverlay[] = [
  {
    id: 'orgswitcher',
    origin: 'orgswitchbutton',
    toHighlight: ['orgswitchbutton', 'orglink'],
    content: {
      number: 1,
      count: 4,
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
      number: 2,
      count: 4,
      i18nText: 'OVERLAYS.SYSTEM.TEXT',
    },
    requirements: {
      media: '(min-width: 600px)',
      permission: ['iam.read'],
    },
  },
  {
    id: 'profilebutton',
    origin: 'avatartoggle',
    toHighlight: ['avatartoggle'],
    content: {
      number: 3,
      count: 4,
      i18nText: 'OVERLAYS.PROFILE.TEXT',
    },
  },
  {
    id: 'mainnav',
    origin: 'mainnav',
    toHighlight: ['mainnav'],
    content: {
      number: 4,
      count: 4,
      i18nText: 'OVERLAYS.NAV.TEXT',
    },
    requirements: {
      permission: ['org.read'],
    },
  },
];

export const OrgContextChangedWorkflowOverlays: CnslOverlay[] = [
  {
    id: 'orgswitcher',
    origin: 'orgswitchbutton',
    toHighlight: ['orgswitchbutton', 'orglink'],
    content: {
      i18nText: 'OVERLAYS.CONTEXTCHANGED.TEXT',
    },
    requirements: {
      permission: ['org.read'],
    },
  },
];
