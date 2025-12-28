import { NextRequest, NextResponse } from 'next/server';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081';

export async function POST(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const formData = await request.formData();

    // Forward the request to the backend
    const response = await fetch(`${API_BASE_URL}/api/projects/${params.id}/upload`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const error = await response.text();
      return NextResponse.json(
        { error: error || 'Upload failed' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Upload proxy error:', error);
    return NextResponse.json(
      { error: 'Failed to upload file' },
      { status: 500 }
    );
  }
}
