import { NextSeo } from 'next-seo';
import Page from '@/components/page';
import Header from '@/components/header';
import React, { useState } from 'react';
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
  return (
    <Page>
      <NextSeo title="Lucidify" description="Lucidify sales with AI" />
      <Header />
      <main>
      </main>
    </Page>
  );
}
