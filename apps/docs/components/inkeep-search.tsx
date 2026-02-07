"use client";

import type { SharedProps } from "fumadocs-ui/components/dialog/search";
import {
  InkeepModalSearchAndChat,
  type InkeepModalSearchAndChatProps,
} from "@inkeep/cxkit-react";
import { useEffect, useState } from "react";

export default function CustomDialog(props: SharedProps) {
  const [syncTarget, setSyncTarget] = useState<HTMLElement | null>(null);
  const { open, onOpenChange } = props;
  // We do this because document is not available in the server
  useEffect(() => {
    setSyncTarget(document.documentElement);
  }, []);

  const config: InkeepModalSearchAndChatProps = {
    baseSettings: {
      apiKey: process.env.NEXT_PUBLIC_INKEEP_API_KEY,
      primaryBrandColor: '#f25543',
      organizationDisplayName: "Zitadel",
      colorMode: {
        sync: {
          target: syncTarget,
          attributes: ["class"],
          isDarkMode: (attributes) => !!attributes.class?.includes("dark"),
        },
      },
    },
    modalSettings: {
      isOpen: open,
      onOpenChange,
      // optional settings
    },
    searchSettings: {
      // optional settings
    },
    aiChatSettings: {
      aiAssistantName: 'Zitadel AI Assistant',
      chatSubjectName: 'Zitadel',
      placeholder: 'How can I enable MFA?',
      introMessage: 'Hey! I’m the Zitadel AI assistant 🤖 — throw your questions my way!',
      aiAssistantAvatar: '/icons/favicon-32x32.png',
      userAvatar: '/icons/user.png',
      exampleQuestionsLabel: 'Example Questions:',
      exampleQuestions: [
        'What is an Instance?',
        'What is an Organization?',
        'How to create an Action?',
        'How to integrate Zitadel to my React app?',
      ],
      isFirstExampleQuestionHighlighted: true,
      shouldOpenLinksInNewTab: true,
      isCopyChatButtonVisible: true,
      disclaimerSettings: {
        isEnabled: true,
        label: 'AI Assistant',
        tooltip: 'Responses are AI-generated and may require verification.',
      },
      toolbarButtonLabels: {
        clear: 'Start Over',
        share: 'Share Chat',
        getHelp: 'Get Help',
        stop: 'Stop',
        copyChat: 'Copy Chat',
      },
    },
  };
  return <InkeepModalSearchAndChat {...config} />;
}