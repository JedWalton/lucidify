import { NextSeo } from 'next-seo';
import Page from '@/components/page';
import Header from '@/components/header';
import React, { useState } from 'react';
import FileUploadArea from '@/modules/openai-cookbook/components/FileUploadArea';
import { FileLite } from '@/modules/openai-cookbook/types/file';
import FileQandAArea from '@/modules/openai-cookbook/components/FileQandAArea';
import ChatbotWidget from '@/modules/chatbot-widget/ChatbotWidget';
// import VideoSection from '@/components/video-section';
// import ListSection from '@/components/list-section';
// import FeatureSection from '@/components/feature-section';
// import CasesSection from '@/components/cases-section';
// import SocialProof from '@/components/social-proof';
// import PricingTable from '@/components/pricing-table';
// import Footer from '@/components/footer';

        // <VideoSection />
        // <ListSection />
        // <FeatureSection />
        // <CasesSection />
        // <SocialProof />
        // <PricingTable />
export default function Home() {
  const [files, setFiles] = useState<FileLite[]>([]);
  return (
    <Page>
      <NextSeo title="Lucidify" description="Lucidify sales with AI" />
      <Header />
      <main>
      <div className="max-w-3xl mx-auto m-8 space-y-8 text-gray-800">
        <h1 className="text-4xl">File Q&A</h1>

        <div className="">
          To search for answers from the content in your files, upload them here
          and we will use OpenAI embeddings and GPT to find answers from the
          relevant documents.
        </div>

        <FileUploadArea
          handleSetFiles={setFiles}
          maxNumFiles={75}
          maxFileSizeMB={30}
        />

        <FileQandAArea files={files} />
      </div>
      </main>
      <ChatbotWidget />

      {/*<Footer />*/}
    </Page>
  );
}
