import React, { useState } from 'react';
import './ChatbotWidget.css';

function ChatbotWidget() {
    const [isMaximized, setIsMaximized] = useState(true);
    const [messages, setMessages] = useState([]);
    const [userInput, setUserInput] = useState('');

    const toggleChat = () => {
        setIsMaximized(!isMaximized);
    };

    const sendMessage = async (e) => {
        e.preventDefault(); // Prevent default form submission

        // Add the user's message to the state
        setMessages(prevMessages => [...prevMessages, { type: 'user', text: userInput }]);

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
    };

    return (
        <div id="chatbox" className={isMaximized ? 'maximized' : ''}>
            <div id="chatHeader">
                <span>Powered by Lucidify.xyz</span>
                <span id="minimizeButton" onClick={toggleChat}>
                    {isMaximized ? '−' : '+'}
                </span>
            </div>
            <div id="messages">
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
                <button type="submit">Send</button>
            </form>
        </div>
    );
}

export default ChatbotWidget;

