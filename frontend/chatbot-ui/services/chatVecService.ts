import { LocalStorage } from "@/types/storage";

import { ChatBody, Message } from '@/types/chat';

type ServerResponse = {
  success: boolean;
  data?: any;  // You can be more specific with this type if you know the expected structure
  message?: string;
};

export const chatVecService = {
  async performVectorSearchOnChatThread(messages: Message[]): Promise<string> {
    try {
      // const url = `${process.env.PUBLIC_BACKEND_API_URL}/api/sync/localstorage/?key=${encodeURIComponent(key as string)}`;
      const url = `http://localhost:8080/api/chat/vector-search`;

      const sessionToken = localStorage.getItem('sessionToken'); // Adjust this line to wherever your session token is stored
      console.log('sessionToken:', sessionToken)

      const headers: HeadersInit = {
        'Content-Type': 'application/json'
      };

      if (sessionToken) {
        headers['Authorization'] = `Bearer ${sessionToken}`;
      }

      let body

      body = JSON.stringify({ messages: messages });

      const options: RequestInit = {
        method: 'POST',
        headers: headers,
        body: body,
        mode: 'cors',
        credentials: 'include',
      };

      const response = await fetch(url, options);

      if (!response.ok) {
        let errorMessage = 'Server responded with an error';
        try {
          const errorBody = await response.json();
          errorMessage = errorBody.message || `Server responded with ${response.status}`;
        } catch (jsonError) {
          try {
            errorMessage = await response.text(); // If response is not in JSON format
          } catch (textError) {
            // handle the case where neither json nor text are readable
            errorMessage = `Server responded with ${response.status}, but the response was not readable.`;
          }
        }

        throw new Error(errorMessage);
      }

      // If the response is OK, we decode it from JSON
      const data = await response.json();
      console.log('data:', data)
      return data.data;

    } catch (error) {
      console.error(`Request failed: ${error}`);
      return "error"
    }
  },
}
