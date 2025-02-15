import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  transpilePackages: ['antd', '@ant-design/icons'],
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: (process.env.BACKEND || 'http://localhost:5000') + '/api/:path*',
      }
    ]
  },
};

export default nextConfig;