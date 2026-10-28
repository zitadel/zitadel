// apps/docs/app/api/chat/route.ts
import { NextResponse } from 'next/server';

export type ChatMessage = {
  role: 'user' | 'assistant' | 'system';
  content: string;
};

const rateLimitMap = new Map<string, { timestamps: number[] }>();
const RATE_LIMIT_MAX = 10;
const RATE_LIMIT_WINDOW_MS = 60000;
export const maxDuration = 30;

export async function POST(req: Request) {
  try {
    const apiKey = process.env.CHAT_AGENT_API_KEY;
    const endpointUrl = process.env.CHAT_AGENT_URL;

    if (!apiKey || !endpointUrl) {
      throw new Error('Chat agent configuration is missing');
    }

    const rawIp = req.headers.get('x-real-ip') || req.headers.get('x-forwarded-for') || 'anonymous';
    const ip = rawIp.split(',')[0].trim();

    const now = Date.now();

    const userRateData = rateLimitMap.get(ip) || { timestamps: [] };
    const validTimestamps = userRateData.timestamps.filter((ts) => now - ts < RATE_LIMIT_WINDOW_MS);

    if (validTimestamps.length >= RATE_LIMIT_MAX) {
      return NextResponse.json(
        { success: false, error: 'Rate limit exceeded. Please try again later.' },
        { status: 429 }
      );
    }

    validTimestamps.push(now);
    rateLimitMap.set(ip, { timestamps: validTimestamps });

    if (Math.random() < 0.01) {
      for (const [key, data] of rateLimitMap.entries()) {
        const valid = data.timestamps.filter((ts) => now - ts < RATE_LIMIT_WINDOW_MS);
        if (valid.length === 0) rateLimitMap.delete(key);
      }
    }

    const { messages } = await req.json();

    if (!Array.isArray(messages)) throw new Error('Invalid payload structure');

    const sanitizedMessages = messages
      .filter((m: ChatMessage) => m.role === 'user' || m.role === 'assistant')
      .map((m: ChatMessage) => ({
        role: m.role,
        content: String(m.content).substring(0, 1500),
      }));

    if (sanitizedMessages.length === 0) throw new Error('No valid messages provided');

    const prompt = sanitizedMessages[sanitizedMessages.length - 1].content;
    const recentHistory = sanitizedMessages.slice(-11, -1);
    const transcript =
      recentHistory.length > 0
        ? recentHistory.map((m) => `${m.role.toUpperCase()}: ${m.content}`).join('\n\n')
        : undefined;

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 25000);

    const response = await fetch(endpointUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${apiKey}`,
        'x-end-user-ip': ip,
      },
      body: JSON.stringify({
        sopId: 'docs_agent',
        properties: {
          user_question: prompt,
          ...(transcript && { transcript }),
        },
        verbose: true,
        background: false,
      }),
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    if (response.status === 429) {
      return NextResponse.json(
        { success: false, error: 'Rate limit exceeded. Please try again later.' },
        { status: 429 }
      );
    }

    const contentType = response.headers.get('content-type') || '';
    if (!contentType.includes('application/json')) {
      throw new Error(`Upstream returned non-JSON response: HTTP ${response.status}`);
    }

    if (!response.ok) {
      throw new Error(`AI Agent service returned HTTP ${response.status}`);
    }

    const jsonResponse = await response.json();

    if (!jsonResponse.success) {
      throw new Error(jsonResponse.message || 'Agent execution failed.');
    }

    let replyText = null;
    const responseData = jsonResponse.data || {};

    for (const [key, value] of Object.entries(responseData)) {
      if (key.endsWith('.chat_response')) {
        replyText = value as string;
        break;
      }
      if (value && typeof value === 'object' && 'chat_response' in value) {
        replyText = (value as any).chat_response as string;
        break;
      }
    }

    const finalReply = replyText || 'I finished thinking, but no response text was found in the agent output.';
    const executionId = jsonResponse.execution?.executionId;

    return NextResponse.json({ success: true, reply: finalReply, executionId });
  } catch (error: any) {
    if (error.name === 'AbortError') {
      return NextResponse.json(
        { success: false, error: 'The assistant took too long to respond. Please try again.' },
        { status: 504 }
      );
    }

    console.error('Docs Chat API Error:', error);
    return NextResponse.json(
      { success: false, error: 'Failed to process message' },
      { status: 500 }
    );
  }
}