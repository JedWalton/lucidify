import React, { useState, ChangeEvent, FormEvent } from 'react';
import { UserButton } from '@clerk/nextjs';

export default function Example() {
  const [fileName, setFileName] = useState<string>('');
  const [fileContent, setFileContent] = useState<string>('');

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setFileName(file.name);
      const reader = new FileReader();
      reader.onload = function (e: ProgressEvent<FileReader>) {
        // Ensure that 'result' is a string before trying to set the state
        if (typeof e.target?.result === 'string') {
          setFileContent(e.target.result);
        }
      };
      reader.readAsText(file);
    }
  };

  const handleUpload = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const formData = new FormData();
    formData.append('file', new Blob([fileContent], { type: 'text/plain' }));
    formData.append('document_name', fileName);
    formData.append('content', fileContent);

    const response = await fetch('http://localhost:8080/documents/upload', {
      method: 'POST',
      body: formData,
      mode: 'cors',
      credentials: 'include',
    });

    if (response.ok) {
      console.log('File uploaded successfully');
    } else {
      console.error('File upload failed', response.statusText);
    }
  };

  return (
    <>
      <header>
        <UserButton afterSignOutUrl="/" />
      </header>
      <div>Your page's content can go here.</div>
      <form onSubmit={handleUpload}>
        <input
          type="file"
          name="file"
          accept=".txt"
          onChange={handleFileChange}
          required
        />
        <button type="submit">Upload</button>
      </form>
    </>
  );
}

