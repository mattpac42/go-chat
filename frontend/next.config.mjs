import { execSync } from 'child_process';
import { readFileSync } from 'fs';
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

// Read version from package.json
const packageJson = JSON.parse(readFileSync(join(__dirname, 'package.json'), 'utf8'));
const appVersion = packageJson.version || '0.0.0';

// Get git commit hash
let gitHash = 'unknown';
try {
  gitHash = execSync('git rev-parse --short HEAD', { encoding: 'utf8' }).trim();
} catch {
  // Git not available or not a git repository
  gitHash = 'dev';
}

/** @type {import('next').NextConfig} */
const nextConfig = {
  // Force dynamic rendering for all pages (no static export)
  output: 'standalone',
  env: {
    NEXT_PUBLIC_APP_VERSION: appVersion,
    NEXT_PUBLIC_GIT_HASH: gitHash,
  },
};

export default nextConfig;
