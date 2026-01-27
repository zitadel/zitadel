'use client';

import { useEffect, useState } from 'react';
import { useTheme } from 'next-themes';

export function InkeepSearch() {
  const { resolvedTheme } = useTheme();
  const [isLoaded, setIsLoaded] = useState(false);
  
  const config = {
    apiKey: process.env.NEXT_PUBLIC_INKEEP_API_KEY,
    integrationId: process.env.NEXT_PUBLIC_INKEEP_INTEGRATION_ID,
    organizationId: process.env.NEXT_PUBLIC_INKEEP_ORGANIZATION_ID,
    primaryBrandColor: '#26D6FF',
    organizationDisplayName: 'Zitadel',
    botName: 'Zitadel AI',
  };

  const handleOpen = () => {
    console.log("Inkeep: Setting isOpen = true");
    // Instead of calling .open(), we update the widget properties directly
    // @ts-ignore
    if (window.__INKEEP_UPDATE__) {
      // @ts-ignore
      window.__INKEEP_UPDATE__({ isOpen: true });
    } else {
      console.warn("Inkeep: Update function not ready.");
    }
  };

  useEffect(() => {
    // 1. Create Anchor
    let anchor = document.getElementById('inkeep-anchor');
    if (!anchor) {
      anchor = document.createElement('div');
      anchor.id = 'inkeep-anchor';
      document.body.appendChild(anchor);
    }
    // @ts-ignore
    window.__INKEEP_ANCHOR_ELEMENT__ = anchor;

    // 2. CSS Fix
    const style = document.createElement('style');
    style.innerHTML = `
      :root, body { --ikp-z-index-modal: 2147483647 !important; }
      .ikp-modal-container { z-index: 2147483647 !important; }
    `;
    document.head.appendChild(style);

    // 3. Inject Script
    if (document.getElementById('inkeep-script-injection')) {
      // @ts-ignore
      if (window.__INKEEP_INSTANCE__) setIsLoaded(true);
      return;
    }

    const script = document.createElement('script');
    script.id = 'inkeep-script-injection';
    script.type = 'module';
    script.crossOrigin = "anonymous";
    
    // Initial Config (Closed by default)
    const baseConfig = JSON.stringify({
      componentType: 'CustomTrigger', // Works best for manual control
      isOpen: false, // Start closed
      properties: {
        baseSettings: {
          apiKey: config.apiKey,
          integrationId: config.integrationId,
          organizationId: config.organizationId,
          organizationDisplayName: config.organizationDisplayName,
          primaryBrandColor: config.primaryBrandColor,
          colorMode: { forcedColorMode: resolvedTheme },
        },
        aiChatSettings: {
          botName: config.botName,
          aiAssistantAvatar: '/img/logo.svg',
        },
      },
    });

    script.textContent = `
      import { Inkeep } from 'https://unpkg.com/@inkeep/widgets-embed@0.2.292/dist/embed.js';
      
      const initialConfig = ${baseConfig};
      const anchorElement = window.__INKEEP_ANCHOR_ELEMENT__;

      if (anchorElement) {
        initialConfig.targetElement = anchorElement; 
        
        // Define the onClose callback to update state
        initialConfig.properties.onClose = () => {
           console.log("Inkeep: Closing...");
           // Re-render with isOpen: false
           if (window.__INKEEP_UPDATE__) window.__INKEEP_UPDATE__({ isOpen: false });
        };

        try {
          const inkeep = Inkeep();
          const instance = inkeep.embed(initialConfig);
          
          window.__INKEEP_INSTANCE__ = instance;
          
          // CRITICAL: We create a global update function.
          // This allows React to "push" new props (like isOpen: true) into the widget.
          window.__INKEEP_UPDATE__ = (newProps) => {
             // Merge new props into existing config
             const updatedConfig = { 
               ...initialConfig, 
               ...newProps 
             };
             // Re-render the widget with new props
             instance.render(updatedConfig);
          };

          window.dispatchEvent(new Event('inkeep-ready'));
          
        } catch(e) {
          console.error("Inkeep: CRASHED.", e);
        }
      }
    `;
    document.body.appendChild(script);

    // 4. Ready Listener
    const checkReady = () => {
       // @ts-ignore
       if (window.__INKEEP_INSTANCE__) setIsLoaded(true);
    };
    window.addEventListener('inkeep-ready', checkReady);
    const interval = setInterval(checkReady, 500);

    return () => {
      window.removeEventListener('inkeep-ready', checkReady);
      clearInterval(interval);
    };
  }, []);

  // Sync Theme
  useEffect(() => {
    // @ts-ignore
    if (window.__INKEEP_UPDATE__) {
       // @ts-ignore
       window.__INKEEP_UPDATE__({
         properties: {
           baseSettings: { colorMode: { forcedColorMode: resolvedTheme } }
         }
       });
    }
  }, [resolvedTheme]);

  // Keyboard Shortcut
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        handleOpen();
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, []);

  return (
    <button
      onClick={handleOpen}
      type="button"
      className="flex w-full items-center gap-2 rounded-lg border bg-secondary/50 px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
    >
      {!isLoaded ? (
         <span className="flex items-center gap-2">
            <svg className="animate-spin size-4" viewBox="0 0 24 24" fill="none">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Loading...
         </span>
      ) : (
        <>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="size-4"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.3-4.3" />
          </svg>
          <span className="flex-1 text-left">Search</span>
          <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
            <span className="text-xs">⌘</span>K
          </kbd>
        </>
      )}
    </button>
  );
}