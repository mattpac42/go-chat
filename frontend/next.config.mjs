/** @type {import('next').NextConfig} */
const nextConfig = {
  // Force dynamic rendering for all pages (no static export)
  output: 'standalone',
};

export default nextConfig;
