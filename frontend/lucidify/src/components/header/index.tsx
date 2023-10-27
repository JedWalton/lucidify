import { tw, css } from 'twind/css';
import React, { useState } from 'react';

import Button from '@/components/button';
import Netlify from '@/constants/svg/netlify.svg';
import Nike from '@/constants/svg/nike.svg';
import Figma from '@/constants/svg/figma.svg';
import Aws from '@/constants/svg/aws.svg';
import Link from 'next/link';

const headerStyle = css`
  background-color: #ffffff;
  min-height: calc(100vh - 6rem);
`;

const Header = () => {
  const [email, setEmail] = useState(``); // State to manage the email input

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    // Discord webhook URL
    const webhookURL = `https://discord.com/api/webhooks/1149411228289093683/nG66YEpse8FyHauLtsDm79-SKV23j03B3IidFLJa7qWQapnxrlmbAjqdijQHRyVSrtAj`;

    try {
      const response = await fetch(webhookURL, {
        method: `POST`,
        headers: {
          'Content-Type': `application/json`,
        },
        body: JSON.stringify({
          content: `New email submitted: ${email}`,
        }),
      });

      if (!response.ok) {
        throw new Error(`Network response was not ok`);
      }

      console.log(`Email sent to Discord successfully!`);
    } catch (error) {
      if (error instanceof Error) {
        console.error(`There was a problem with the fetch operation:`, error.message);
      } else {
        console.error(`There was a problem with the fetch operation:`, error);
      }
    }
  };

  return (
    <header className={tw(headerStyle)}>
      <div className={tw(`max-w-4xl mx-auto py-16 px-14 sm:px-6 lg:px-8`)}>
        <h1
          className={tw(`font-sans font-bold text-4xl md:text-5xl lg:text-8xl text-center leading-snug text-gray-800`)}
        >
          Lucidify your business.
        </h1>
        <div className={tw(`max-w-xl mx-auto`)}>
          <p className={tw(`mt-10 text-gray-500 text-center text-xl lg:text-3xl`)}>Chat securely with your private ChatGPT.</p>
        </div>

        <div className={tw(`mt-10 flex justify-center items-center w-full mx-auto`)}>
          <Link href="/example">
            <Button primary type="submit">
              Get Started.
            </Button>
          </Link>
        </div>
      </div>
      <div className={tw(`flex justify-center w-full`)}>
        <div className={tw(`mt-4 w-full`)}>
          <p className={tw(`font-mono uppercase text-center font-medium text-sm text-gray-600`)}>These folks get it</p>
          <div className={tw(`flex items-center justify-center mx-auto flex-wrap`)}>
            <Aws className={tw(`m-12 mb-8`)} width={120} />
            <Netlify className={tw(`m-12`)} width={140} />
            <Nike className={tw(`m-12`)} width={140} />
            <Figma className={tw(`m-12`)} width={140} />
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
