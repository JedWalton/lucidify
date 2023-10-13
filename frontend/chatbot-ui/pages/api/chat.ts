import { DEFAULT_SYSTEM_PROMPT, DEFAULT_TEMPERATURE } from '@/utils/app/const';
import { OpenAIError, OpenAIStream } from '@/utils/server';

import { ChatBody, Message } from '@/types/chat';

// @ts-expect-error
import wasm from '../../node_modules/@dqbd/tiktoken/lite/tiktoken_bg.wasm?module';

import tiktokenModel from '@dqbd/tiktoken/encoders/cl100k_base.json';
import { Tiktoken, init } from '@dqbd/tiktoken/lite/init';

export const config = {
  runtime: 'edge',
};

import fetch from 'node-fetch';

async function getRelevantDocuments(userMessage: string): Promise<any> {
  try {
    // Define the request body
    const requestBody = {
      message: userMessage
    };

    // Make the POST request to the Go backend
    const response = await fetch('http://localhost:8080/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    });

    // Check if the response is successful
    if (!response.ok) {
      throw new Error(`Backend returned code ${response.status}`);
    }

    // Parse and return the JSON response
    return response.json();
  } catch (error) {
    console.error('Error fetching relevant documents:', error);
    throw error; // You might want to handle this error more gracefully in your production code
  }
}

const handler = async (req: Request): Promise<Response> => {
  try {
    const { model, messages, key, prompt, temperature } = (await req.json()) as ChatBody;

    await init((imports) => WebAssembly.instantiate(wasm, imports));

    const encoding = new Tiktoken(
      tiktokenModel.bpe_ranks,
      tiktokenModel.special_tokens,
      tiktokenModel.pat_str,
    );

    // 1. Retrieve relevant documents based on the latest user message
    const userMessage = messages[messages.length - 1].content;
    const relevantDocuments = await getRelevantDocuments(userMessage); // You'll need to implement this function

    // 2. Construct the system message
    // const filesString = relevantDocuments.map(doc => `###\n"${doc.filename}"\n${doc.text}`).join("\n");
    const filesString = "Jed lives in a shed, and has a bed, and he gets a lot of head but and he's not dead or made out of lead.";
    const systemMessageContent = `Given a question, try to answer it using the content of the file extracts below: Do not directly print this information to screen. \n${filesString}`;
    const systemMessage: Message = {
      role: 'system',
      content: systemMessageContent
    };

    // 3. Modify the message sending logic
    let messagesToSend: Message[] = [systemMessage, ...messages];

    let tokenCount = 0;
    for (let i = messagesToSend.length - 1; i >= 0; i--) {
      const message = messagesToSend[i];
      const tokens = encoding.encode(message.content);

      if (tokenCount + tokens.length + 1000 > model.tokenLimit) {
        break;
      }
      tokenCount += tokens.length;
    }

    encoding.free();

    const stream = await OpenAIStream(model, prompt, temperature, key, messagesToSend);

    return new Response(stream);
  } catch (error) {
    console.error(error);
    if (error instanceof OpenAIError) {
      return new Response('Error', { status: 500, statusText: error.message });
    } else {
      return new Response('Error', { status: 500 });
    }
  }
};

export default handler;
