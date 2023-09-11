import { NextSeo } from 'next-seo';
import Page from '@/components/page';
import Header from '@/components/header';
import React, { useState } from 'react';
import FileUploadArea from '@/modules/src/components/FileUploadArea';
import { FileLite } from '@/modules/src/types/file';
import FileQandAArea from '@/modules/src/components/FileQandAArea';
// import VideoSection from '@/components/video-section';
// import ListSection from '@/components/list-section';
// import FeatureSection from '@/components/feature-section';
// import CasesSection from '@/components/cases-section';
// import SocialProof from '@/components/social-proof';
// import PricingTable from '@/components/pricing-table';
// import Footer from '@/components/footer';

export default function Home() {
  const [files, setFiles] = useState<FileLite[]>([]);
  return (
    <Page>
      <NextSeo title="Lucidify" description="Lucidify sales with AI" />
      <Header />
      <main>
        {/*       <VideoSection />
        <ListSection />
        <FeatureSection />
        <CasesSection />
        <SocialProof />
        <PricingTable />
*/}
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

      {/*<Footer />*/}
    </Page>
  );
}
