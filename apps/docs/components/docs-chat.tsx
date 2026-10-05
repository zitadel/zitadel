// components/docs-chat.tsx
'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { mixpanelClient } from '@/utils/mixpanel';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { cn } from '@/utils/cn';
import {
  ChevronDown,
  Send,
  RefreshCw,
  ThumbsUp,
  ThumbsDown,
} from 'lucide-react';

export type ChatMessage = {
  role: 'user' | 'assistant' | 'system';
  content: string;
};

export function DocsChat() {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const chatContainerRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Chat State
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [lastExecutionId, setLastExecutionId] = useState<string | null>(null);

  // Transcript Feedback State
  const [transcriptRating, setTranscriptRating] = useState<'up' | 'down' | null>(null);
  const [feedbackText, setFeedbackText] = useState('');
  const [isFeedbackSubmitted, setIsFeedbackSubmitted] = useState(false);

  const getIntroMessage = useCallback((): ChatMessage => {
    return {
      role: 'assistant',
      content: `Hey there, I’m the Zitadel AI assistant 🤖 — throw your questions my way!`,
    };
  }, []);

  // Initialize intro message
  useEffect(() => {
    if (messages.length === 0) {
      setMessages([getIntroMessage()]);
    }
  }, [messages.length, getIntroMessage]);

  useEffect(() => {
    if (!isLoading && isOpen) {
      setTimeout(() => inputRef.current?.focus(), 10);
    }
  }, [isLoading, isOpen]);

  useEffect(() => {
    if (isOpen) messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, isOpen, isLoading, transcriptRating]);

  // Handle click outside & escape
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setIsOpen(false);
    };
    const handleClickOutside = (e: MouseEvent) => {
      if (chatContainerRef.current && !chatContainerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

  const handleClearChat = () => {
    setMessages([getIntroMessage()]);
    setTranscriptRating(null);
    setFeedbackText('');
    setIsFeedbackSubmitted(false);
    setLastExecutionId(null);
  };

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim() || isLoading) return;

    const newMessages: ChatMessage[] = [...messages, { role: 'user', content: inputValue.trim() }];
    setMessages(newMessages);
    setInputValue('');
    setIsLoading(true);

    try {
      const response = await fetch('/docs/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ messages: newMessages }),
      });

      const data = await response.json();

      if (data.success && data.reply) {
        setMessages([...newMessages, { role: 'assistant', content: data.reply }]);
        if (data.executionId) {
          setLastExecutionId(data.executionId);
        }
      } else {
        setMessages([
          ...newMessages,
          { role: 'system', content: data.error || '⚠️ Failed to connect to the assistant.' },
        ]);
      }
    } catch {
      setMessages([...newMessages, { role: 'system', content: '⚠️ An unexpected error occurred.' }]);
    } finally {
      setIsLoading(false);
    }
  };

  const submitTranscriptFeedback = (rating: 'up' | 'down', reason: string = '') => {
    mixpanelClient.track('submitted_chat_feedback', {
      flow: 'docs_support',
      rating,
      reason,
      execution_id: lastExecutionId,
      message_count: messages.length,
    });
    setIsFeedbackSubmitted(true);
  };

  const handleInitialRating = (rating: 'up' | 'down') => {
    setTranscriptRating(rating);
    if (rating === 'up') {
      submitTranscriptFeedback('up');
    }
  };

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end font-sans">
      <div
        ref={chatContainerRef}
        className={cn(
          "absolute bottom-16 right-0 mb-2 flex flex-col w-[360px] sm:w-[400px] h-[600px] max-h-[80vh]",
          "bg-fd-background border border-fd-border rounded-2xl shadow-2xl overflow-hidden",
          "transition-all duration-300 ease-in-out origin-bottom-right",
          isOpen
            ? "opacity-100 scale-100 translate-y-0 pointer-events-auto"
            : "opacity-0 scale-[0.85] translate-y-8 pointer-events-none"
        )}
      >
        {/* Header */}
        <div className="bg-fd-accent text-fd-accent-foreground p-4 flex justify-between items-center border-b border-fd-border">
          <div className="flex items-center gap-2.5">
            <img
              src="/docs/zitadel-192x192.png"
              alt="Zitadel Logo"
              className="w-6 h-6 object-contain"
            />
            <span className="font-semibold text-lg tracking-tight text-black dark:text-white">Zitadel AI Assistant</span>
          </div>

          <div className="flex gap-1 items-center">
            <button
              onClick={handleClearChat}
              title="Start fresh"
              className="p-1.5 hover:bg-fd-background/20 rounded-md transition-colors"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
            <button
              onClick={() => setIsOpen(false)}
              title="Minimize chat"
              className="p-1.5 hover:bg-fd-background/20 rounded-md transition-colors"
            >
              <ChevronDown className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Message Container */}
        <div className="flex-1 p-4 overflow-y-auto bg-fd-background space-y-4">
          {messages.map((msg, idx) => (
            <div key={idx} className={cn("flex", msg.role === 'user' ? "justify-end" : "justify-start")}>
              <div
                className={cn(
                  "max-w-[85%] break-words overflow-hidden rounded-2xl px-4 py-3 text-sm shadow-sm",
                  msg.role === 'user'
                    ? "bg-[#187aff] text-white rounded-br-sm"
                    : msg.role === 'system'
                      ? "bg-red-500/10 text-red-600 w-full text-center border border-red-500/20 rounded-xl"
                      : "bg-fd-accent/10 text-fd-foreground border border-fd-border rounded-bl-sm"
                )}
              >
                {msg.role === 'user' || msg.role === 'system' ? (
                  msg.content
                ) : (
                  <div className="prose prose-sm dark:prose-invert max-w-none prose-p:leading-relaxed prose-pre:bg-fd-muted prose-pre:p-3 prose-pre:rounded-lg">
                    <ReactMarkdown
                      remarkPlugins={[remarkGfm]}
                      components={{
                        img: () => null,
                        a({ _node, children, ...props }: any) {
                          return (
                            <a
                              {...props}
                              className="text-[#187aff] dark:text-[#3f90ff] underline underline-offset-2 hover:opacity-80 transition-opacity font-medium"
                              target="_blank"
                              rel="noopener noreferrer"
                            >
                              {children}
                            </a>
                          );
                        },
                        code({ _node, inline, children, ...props }: any) {
                          return inline ? (
                            <code className="bg-fd-muted rounded px-1.5 py-0.5 font-mono text-xs text-fd-foreground" {...props}>
                              {children}
                            </code>
                          ) : (
                            <code className="block text-xs font-mono text-fd-foreground" {...props}>
                              {children}
                            </code>
                          );
                        },
                      }}
                    >
                      {msg.content}
                    </ReactMarkdown>
                  </div>
                )}
              </div>
            </div>
          ))}

          {isLoading && (
            <div className="flex justify-start">
              <div className="bg-fd-accent/10 border border-fd-border rounded-2xl rounded-bl-sm px-4 py-4 shadow-sm flex gap-1.5">
                <div className="w-1.5 h-1.5 bg-fd-muted-foreground rounded-full animate-bounce"></div>
                <div className="w-1.5 h-1.5 bg-fd-muted-foreground rounded-full animate-bounce" style={{ animationDelay: '0.15s' }}></div>
                <div className="w-1.5 h-1.5 bg-fd-muted-foreground rounded-full animate-bounce" style={{ animationDelay: '0.3s' }}></div>
              </div>
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>

        {/* Feedback Section */}
        {messages.length > 1 && (
          <div className="bg-fd-background border-t border-fd-border px-4 py-3 shrink-0">
            {isFeedbackSubmitted ? (
              <p className="text-sm text-center font-medium text-green-600 dark:text-green-400">
                Thank you for your feedback!
              </p>
            ) : transcriptRating === 'down' ? (
              <div className="flex flex-col gap-2 animate-in fade-in slide-in-from-bottom-2 duration-200">
                <label className="text-sm text-fd-foreground font-medium">What went wrong?</label>
                <textarea
                  maxLength={500}
                  value={feedbackText}
                  onChange={(e) => setFeedbackText(e.target.value)}
                  placeholder="The assistant didn't answer my question..."
                  className="w-full p-2 bg-transparent border border-fd-border rounded-lg text-sm outline-none focus:border-[#187aff] resize-none h-20 transition-colors text-fd-foreground"
                />
                <div className="flex justify-between items-center text-xs">
                  <span className={cn("font-medium", feedbackText.length === 500 ? 'text-red-500' : 'text-fd-muted-foreground')}>
                    {feedbackText.length}/500
                  </span>
                  <div className="flex gap-2">
                    <button
                      onClick={() => setTranscriptRating(null)}
                      className="px-3 py-1.5 text-fd-muted-foreground hover:bg-fd-accent/20 rounded-md transition-colors"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={() => submitTranscriptFeedback('down', feedbackText)}
                      disabled={!feedbackText.trim()}
                      className="px-3 py-1.5 bg-[#187aff] text-white rounded-md hover:bg-[#0f6ee6] disabled:opacity-50 transition-colors font-medium"
                    >
                      Submit
                    </button>
                  </div>
                </div>
              </div>
            ) : (
              <div className="flex items-center justify-between">
                <span className="text-sm text-fd-muted-foreground font-medium">How is this conversation going?</span>
                <div className="flex gap-1">
                  <button
                    onClick={() => handleInitialRating('up')}
                    aria-label="Good rating"
                    className="p-1.5 text-fd-muted-foreground hover:text-green-600 hover:bg-green-500/10 rounded-md transition-colors"
                  >
                    <ThumbsUp className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleInitialRating('down')}
                    aria-label="Bad rating"
                    className="p-1.5 text-fd-muted-foreground hover:text-red-600 hover:bg-red-500/10 rounded-md transition-colors"
                  >
                    <ThumbsDown className="w-4 h-4" />
                  </button>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Input area */}
        <form onSubmit={handleSend} className="p-3 bg-fd-background border-t border-fd-border">
          <div className="relative flex items-center">
            <input
              ref={inputRef}
              type="text"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              placeholder="Ask about Zitadel..."
              className="w-full pl-4 pr-12 py-3 bg-fd-accent/5 focus:bg-transparent border border-transparent focus:border-[#187aff] rounded-xl outline-none text-sm text-fd-foreground transition-colors"
              disabled={isLoading}
            />
            <button
              type="submit"
              disabled={isLoading || !inputValue.trim()}
              className="absolute right-2 p-1.5 bg-[#187aff] text-white rounded-lg hover:bg-[#0f6ee6] disabled:opacity-50 transition-colors"
            >
              <Send className="w-4 h-4" />
            </button>
          </div>
        </form>
      </div>

      <button
        onClick={() => setIsOpen((prev) => !prev)}
        className={cn(
          "flex items-center gap-1.5 px-4 h-12 rounded-full text-white bg-[#187aff] hover:bg-[#0f6ee6]",
          "shadow-[0_4px_14px_0_rgba(24,122,255,0.39)] transition-all duration-300 z-10",
          isOpen ? "scale-0 opacity-0 pointer-events-none" : "scale-100 opacity-100 hover:scale-105"
        )}
      >
        <img
          src="/docs/zitadel-192x192.png"
          alt="Zitadel Logo"
          className="w-5 h-5 object-contain -ml-1"
        />
        <span className="font-semibold text-sm whitespace-nowrap tracking-wide hidden sm:inline">Need Help?</span>
      </button>
    </div>
  );
}