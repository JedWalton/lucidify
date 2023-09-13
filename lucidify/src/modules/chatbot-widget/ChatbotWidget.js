import React, { useState, useRef, useEffect } from 'react';

import './ChatbotWidget.module.css';

function ChatbotWidget() {
    const [isMaximized, setIsMaximized] = useState(true);
    const [messages, setMessages] = useState([]);
    const [userInput, setUserInput] = useState('');
    const [isSending, setIsSending] = useState(false);

    const messagesEndRef = useRef(null);

    useEffect(() => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollTop = messagesEndRef.current.scrollHeight;
        }
    }, [messages]);


    const toggleChat = () => {
        setIsMaximized(!isMaximized);
    };

    const sendMessage = async (e) => {
        e.preventDefault(); // Prevent default form submission

        if (isSending) return;
        setIsSending(true);

        // Add the user's message to the state
        setMessages(prevMessages => [...prevMessages, { type: 'user', text: userInput }]);

        try {
            const response = await fetch('http://localhost:8080/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ message: userInput })
            });

            const data = await response.json();

            // Add the bot's response to the state
            setMessages(prevMessages => [...prevMessages, { type: 'bot', text: data.response }]);
            setUserInput('');
        } catch (error) {
            console.error("Error sending message:", error);
        } finally {
            setIsSending(false); // Set back to false after request completes
        }
    };

    return (
        <div id="chatbox" className={isMaximized ? 'maximized' : ''}>
            <div id="chatHeader">
                <span id="minimizeButton" onClick={toggleChat}>Lucidify</span>
                <span id="minimizeButton" onClick={toggleChat}>
                    {isMaximized ? '−' : '+'}
                </span>
            </div>
            <div id="messages" ref={messagesEndRef}>
                {messages.map((message, index) => (
                    <div key={index} className={`message-container ${message.type}`}>
                        <div className={`${message.type}-message`}>
                            {message.text}
                        </div>
                    </div>
                ))}
            </div>
            <form id="inputArea" onSubmit={sendMessage}>
            <input
                type="text"
                id="userInput"
                placeholder="Type a message..."
                value={userInput}
                onChange={e => setUserInput(e.target.value)}
            />
            <button type="submit" disabled={isSending}>
                {isSending ? '...' : 'Send'}
            </button>
        </form>
        </div>
    );
}

export default ChatbotWidget;
