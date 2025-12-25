import Link from 'next/link';

export default function NotFound() {
  return (
    <div className="flex h-full items-center justify-center">
      <div className="text-center p-8">
        <h2 className="text-xl font-semibold text-gray-900 mb-2">
          Page not found
        </h2>
        <p className="text-gray-500 mb-4">
          The page you are looking for does not exist.
        </p>
        <Link
          href="/"
          className="px-6 py-3 bg-teal-400 text-white rounded-lg hover:bg-teal-500 transition-colors inline-block"
        >
          Go to Home
        </Link>
      </div>
    </div>
  );
}
