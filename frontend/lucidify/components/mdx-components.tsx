import dynamic from 'next/dynamic';
import { useMDXComponent } from "next-contentlayer/hooks";

// Dynamically import the Image component with SSR disabled
const DynamicImage = dynamic(() => import('next/image'), { ssr: false });

const components = {
  Image: DynamicImage,
};

interface MdxProps {
  code: string
}

export function Mdx({ code }: MdxProps) {
  const Component = useMDXComponent(code)

  return <Component components={components} />
}
