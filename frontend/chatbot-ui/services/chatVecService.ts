import { LocalStorage } from "@/types/storage";

import { ChatBody, Message } from '@/types/chat';

type ServerResponse = {
  success: boolean;
  data?: any;  // You can be more specific with this type if you know the expected structure
  message?: string;
};

export const chatVecService = {
  async performVectorSearchOnChatThread(): Promise<string| null> {
    try {
      // const url = `${process.env.PUBLIC_BACKEND_API_URL}/api/sync/localstorage/?key=${encodeURIComponent(key as string)}`;
      const url = `http://localhost:8080/api/chat`;

      const sessionToken = localStorage.getItem('sessionToken'); // Adjust this line to wherever your session token is stored
      console.log('sessionToken:', sessionToken)

      const headers: HeadersInit = {
        'Content-Type': 'application/json'
      };

      if (sessionToken) {
        headers['Authorization'] = `Bearer ${sessionToken}`;
      }

      let body 

      const options: RequestInit = {
        method: 'GET',
        headers: headers,
        body: body,
        mode: 'cors',
        credentials: 'include',
      };

      const response = await fetch(url, options);

      if (!response.ok) {
        // You can first attempt to decode the response as JSON, and then fall back to text if it fails.
        let errorMessage = 'Server responded with an error';
        try {
          const errorBody = await response.json();
          errorMessage = errorBody.message || `Server responded with ${response.status}`;
        } catch (jsonError) {
          errorMessage = await response.text(); // If response is not in JSON format
        }

        throw new Error(errorMessage);
      }

      // If the response is OK, we decode it from JSON
      const data: ServerResponse = await response.json();
      return data.data;

    } catch (error) {
      console.error(`Request failed: ${error}`);
      return null
    }
  },
}
